package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/go-redis/redis/v8"
	evmbind "gitlab.com/rarimo/evm-bridge/gobind"

	"gitlab.com/distributed_lab/logan/v3"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger"
	chains "gitlab.com/rarimo/relayer-svc/internal/types"
	"gitlab.com/rarimo/relayer-svc/internal/utils"

	uniswap "gitlab.com/rarimo/relayer-svc/pkg/uniswap"
)

const evmEstimateTTLMinutes = 30

type evmBridger struct {
	log          *logan.Entry
	rarimocore   rarimocore.QueryClient
	tokenmanager tokenmanager.QueryClient
	evm          *config.EVM
	redis        *redis.Client
}

func NewEVMBridger(cfg config.Config) bridger.Bridger {
	return &evmBridger{
		log:          cfg.Log().WithField("service", "evm_bridger"),
		rarimocore:   rarimocore.NewQueryClient(cfg.Cosmos()),
		tokenmanager: tokenmanager.NewQueryClient(cfg.Cosmos()),
		evm:          cfg.EVM(),
		redis:        cfg.Redis().Client(),
	}
}

func (b *evmBridger) makeWitdrawTx(
	ctx context.Context,
	transfer core.TransferDetails,
	simulation bool,
) (*config.EVMChain, *types.Transaction, error) {
	chain := b.mustGetChain(transfer.Transfer.ToChain)
	bridge, err := evmbind.NewBridge(chain.BridgeAddress, chain.RPC)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to make an instance of ethereum bridge")
	}

	withdrawn, err := bridge.BridgeCaller.UsedHashes(
		&bind.CallOpts{Pending: false, Context: ctx},
		utils.ToByte32(hexutil.MustDecode(transfer.Origin)),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to check if the transfer was already withdrawn")
	}
	if withdrawn {
		return nil, nil, bridger.ErrAlreadyWithdrawn
	}

	amount, err := utils.GetAmountOrDefault(transfer.Transfer.Amount, big.NewInt(1))
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("invalid amount: %s", transfer.Transfer.Amount))
	}
	receiver := common.HexToAddress(transfer.Transfer.Receiver)
	origin := utils.ToByte32(hexutil.MustDecode(transfer.Origin))
	/**
	Tweak the V value to make it compatible with the OpenZeppelin ECDSA implementation
	https://github.com/OpenZeppelin/openzeppelin-contracts/blob/a1948250ab8c441f6d327a65754cb20d2b1b4554/contracts/utils/cryptography/ECDSA.sol#L143
	*/
	signature := hexutil.MustDecode(transfer.Signature)
	signature[64] += 27

	proof, err := proofABI.Pack(transfer.MerklePath, signature)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to ABI encode the proof")
	}
	bundle, err := getBundleData(transfer.Transfer)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse bundle data")
	}

	opts := chain.TransactorOpts()
	opts.NoSend = simulation
	nonce, err := chain.RPC.PendingNonceAt(context.TODO(), chain.SubmitterAddress)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch a nonce")
	}
	opts.Nonce = big.NewInt(int64(nonce))
	gasPrice, err := chain.RPC.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get suggested gas price")
	}
	opts.GasPrice = gasPrice
	opts.GasLimit = uint64(300000)

	var call func() (*types.Transaction, error)
	switch transfer.TokenDetails.TokenType {
	case tokenmanager.Type_NATIVE:
		tokenData, err := nativeABI.Pack(amount)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawNative(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
			)
		}
	case tokenmanager.Type_ERC20:
		tokenData, err := erc20ABI.Pack(common.HexToAddress(transfer.TokenDetails.TokenAddress), amount)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawERC20(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
				transfer.TokenDetails.Wrapped,
			)
		}
	case tokenmanager.Type_ERC721:
		tokenID, err := parseTokenID(transfer.TokenDetails.TokenId)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to parse the tokenID")
		}

		tokenData, err := erc721ABI.Pack(
			common.HexToAddress(transfer.TokenDetails.TokenAddress),
			tokenID,
			transfer.TokenDetails.Uri,
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawERC721(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
				transfer.TokenDetails.Wrapped,
			)
		}
	case tokenmanager.Type_ERC1155:
		tokenID, err := parseTokenID(transfer.TokenDetails.TokenId)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to parse the tokenID")
		}
		tokenData, err := erc1155ABI.Pack(
			common.HexToAddress(transfer.TokenDetails.TokenAddress),
			tokenID,
			transfer.TokenDetails.Uri,
			amount,
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawERC1155(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
				transfer.TokenDetails.Wrapped,
			)
		}
	default:
		return nil, nil, errors.Errorf("token type %d is not supported", transfer.TokenDetails.TokenType)
	}

	tx, err := call()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to call the withdraw method")
	}

	return chain, tx, nil
}

func (b *evmBridger) getSavedEstimate(ctx context.Context, transferID string) (*bridger.FeeEstimate, error) {
	resp := b.redis.Get(ctx, transferID)
	if resp.Err() == redis.Nil {
		return nil, nil
	} else if resp.Err() != nil {
		return nil, errors.Wrap(resp.Err(), "failed to get the estimate")
	}

	raw, err := resp.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read estimate from redis")
	}

	var estimate bridger.FeeEstimate
	if err := json.Unmarshal(raw, &estimate); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal estimate")
	}

	return &estimate, nil
}

func (b *evmBridger) saveEstimate(ctx context.Context, estimate *bridger.FeeEstimate) error {
	raw, err := json.Marshal(estimate)
	if err != nil {
		return errors.Wrap(err, "failed to marshal estimate")
	}

	if resp := b.redis.Set(ctx, estimate.TransferID, raw, evmEstimateTTLMinutes*time.Minute); resp.Err() != nil {
		return errors.Wrap(resp.Err(), "failed to set the cursor")
	}

	return nil
}

func (b *evmBridger) Withdraw(
	ctx context.Context,
	transfer core.TransferDetails,
) error {
	log := b.log.WithField("op_id", transfer.Origin)

	chain, tx, err := b.makeWitdrawTx(ctx, transfer, false)
	if err != nil {
		return errors.Wrap(err, "failed to call the withdraw method")
	}

	log.WithField("tx_id", tx.Hash()).Info("submitted transaction")

	receipt, err := bind.WaitMined(context.TODO(), chain.RPC, tx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for the transaction to be mined")
	}
	if receipt.Status == 0 {
		log.WithField("receipt", utils.Prettify(receipt)).Errorf("%s transaction failed", transfer.Transfer.ToChain)
		return errors.New("transaction failed")
	}

	log.
		WithFields(logan.F{
			"tx_id":        tx.Hash(),
			"tx_index":     receipt.TransactionIndex,
			"block_number": receipt.BlockNumber,
			"gas_used":     receipt.GasUsed,
		}).
		Info("evm transaction confirmed")

	return nil
}

func (b *evmBridger) EstimateRelayFee(
	ctx context.Context,
	transfer core.TransferDetails,
) (bridger.FeeEstimate, error) {
	estimate, err := b.getSavedEstimate(ctx, transfer.Origin)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to get saved estimate")
	}
	if estimate != nil {
		return *estimate, nil
	}

	chain, tx, err := b.makeWitdrawTx(ctx, transfer, true)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to call the withdraw method")
	}

	chainParams, err := getChainParams(transfer.Transfer.ToChain)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to get the chain params")
	}

	bn, err := chain.RPC.BlockNumber(ctx)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to get the block number")
	}

	bignumBn := big.NewInt(0).SetUint64(bn)
	block, err := chain.RPC.BlockByNumber(ctx, bignumBn)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to get the block")
	}

	var baseFee *big.Int
	switch transfer.Transfer.ToChain {
	case chains.MaticMainnet, chains.Mumbai:
		baseFee = calculatePolygonBaseFee(chainParams, block.Header())
	case chains.BSCMainnet, chains.Chapel:
		baseFee = nil
	case chains.AvalancheMainnet, chains.Fuji:
		baseFee, err = chain.AvalancheRPC().EstimateBaseFee(ctx)
		if err != nil {
			return bridger.FeeEstimate{}, errors.Wrap(err, "failed to get the avalanche base fee")
		}
	default:
		baseFee = misc.CalcBaseFee(chainParams, block.Header())
	}

	msg, err := tx.AsMessage(types.LatestSigner(chainParams), baseFee)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to get the message")
	}

	callMsg := ethereum.CallMsg{
		From:      msg.From(),
		To:        msg.To(),
		Gas:       msg.Gas(),
		GasPrice:  msg.GasPrice(),
		GasFeeCap: msg.GasFeeCap(),
		GasTipCap: msg.GasTipCap(),
		Data:      msg.Data(),
	}

	gasEstimate, err := chain.RPC.EstimateGas(ctx, callMsg)
	if err != nil {
		errors.Wrap(err, "failed to get the gas estimate")
	}

	gasPrice, err := chain.RPC.SuggestGasPrice(ctx)
	if err != nil {
		errors.Wrap(err, "failed to get the gas price")
	}

	totalGasPrice := big.NewInt(0).Mul(gasPrice, big.NewInt(0).SetUint64(gasEstimate))
	fee, err := b.calcFee(ctx, chain, totalGasPrice)
	if err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to calculate the fee")
	}

	createdAt := time.Now()
	feeEstimate := bridger.FeeEstimate{
		TransferID:      transfer.Transfer.Origin,
		GasEstimate:     totalGasPrice,
		GasToken:        chain.GasToken,
		FeeAmount:       fee,
		FeeToken:        chain.RelayFeeToken,
		FeeTokenAddress: chain.RelayFeeTokenAddress.Hex(),
		FromChain:       transfer.Transfer.FromChain,
		ToChain:         transfer.Transfer.ToChain,
		CreatedAt:       createdAt,
		ExpiresAt:       createdAt.Add(evmEstimateTTLMinutes * time.Minute),
	}

	if err := b.saveEstimate(ctx, &feeEstimate); err != nil {
		return bridger.FeeEstimate{}, errors.Wrap(err, "failed to save the estimate")
	}

	return feeEstimate, nil
}

func getBundleData(transfer rarimocore.Transfer) (evmbind.IBundlerBundle, error) {
	if len(transfer.BundleData) == 0 {
		return evmbind.IBundlerBundle{}, nil
	}

	result := evmbind.IBundlerBundle{}
	bundle, err := hexutil.Decode(transfer.BundleData)
	if err != nil {
		return result, errors.Wrap(err, "failed to parse bundle data")
	}
	salt, err := hexutil.Decode(transfer.BundleSalt)
	if err != nil {
		return result, errors.Wrap(err, "failed to parse bundle salt")
	}

	result.Bundle = bundle
	result.Salt = utils.ToByte32(salt)

	return result, nil
}

func parseTokenID(rawTokenID string) (*big.Int, error) {
	rawBytes, err := hexutil.Decode(rawTokenID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the tokenID: %s", logan.Field("token_id", rawTokenID))
	}

	return big.NewInt(0).SetBytes(rawBytes), nil
}

func (b *evmBridger) mustGetChain(chainName string) *config.EVMChain {
	chain, ok := b.evm.GetChainByName(chainName)
	if !ok {
		panic(errors.Errorf("unknown EVM chain: %s", chainName))
	}

	return chain
}

func getChainParams(chain string) (*params.ChainConfig, error) {
	var chainParams *params.ChainConfig

	switch chain {
	case chains.EthereumMainnet:
		chainParams = params.MainnetChainConfig
	case chains.Goerli:
		chainParams = params.GoerliChainConfig
	case chains.Sepolia:
		chainParams = params.SepoliaChainConfig
	case chains.MaticMainnet:
		chainParams = &params.ChainConfig{
			ChainID:             big.NewInt(137),
			HomesteadBlock:      big.NewInt(0),
			DAOForkBlock:        nil,
			DAOForkSupport:      true,
			EIP150Hash:          common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(3395000),
			MuirGlacierBlock:    big.NewInt(3395000),
			BerlinBlock:         big.NewInt(14750000),
			LondonBlock:         big.NewInt(23850000),
		}
	case chains.Mumbai:
		chainParams = &params.ChainConfig{
			ChainID:             big.NewInt(80001),
			HomesteadBlock:      big.NewInt(0),
			DAOForkBlock:        nil,
			DAOForkSupport:      true,
			EIP150Hash:          common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(2722000),
			MuirGlacierBlock:    big.NewInt(2722000),
			BerlinBlock:         big.NewInt(13996000),
			LondonBlock:         big.NewInt(22640000),
		}
	case chains.BSCMainnet:
		chainParams = &params.ChainConfig{
			ChainID:             big.NewInt(56),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
		}
	case chains.Chapel:
		chainParams = &params.ChainConfig{
			ChainID:             big.NewInt(97),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		}
	case chains.AvalancheMainnet:
		chainParams = &params.ChainConfig{
			ChainID:             big.NewInt(43114),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		}
	case chains.Fuji:
		chainParams = &params.ChainConfig{
			ChainID:             big.NewInt(43113),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP150Hash:          common.Hash{},
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		}
	default:
		return nil, errors.Errorf("chain %s is not supported", chain)
	}

	return chainParams, nil

}

func (b *evmBridger) calcFee(
	ctx context.Context,
	chain *config.EVMChain,
	gasEstimate *big.Int,
) (*big.Int, error) {
	var pool uniswap.Pool
	var err error

	if chain.UseUniswapV2 {
		pool, err = uniswap.NewPoolV2(
			b.log.WithField("chain", chain.Name),
			chain.RPC,
			chain.UniswapPoolAddress,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the uniswap pool")
		}
	} else {
		pool, err = uniswap.NewPoolV3(
			b.log.WithField("chain", chain.Name),
			chain.RPC,
			chain.UniswapPoolAddress,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the uniswap pool")
		}
	}

	converted, err := pool.Convert(ctx, chain.RelayFeeTokenAddress, chain.GasTokenAddress, gasEstimate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert gas estimate to relay fee token")
	}

	// sanity check
	if converted.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("converted gas estimate is zero")
	}

	return converted, nil
}
