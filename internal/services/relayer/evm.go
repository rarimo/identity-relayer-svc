package relayer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	evmbind "gitlab.com/rarimo/evm-bridge/gobind"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/helpers"
)

var bytes32SliceType = mustABIType("bytes32[]")
var bytesType = mustABIType("bytes")
var uint256Type = mustABIType("uint256")
var addressType = mustABIType("address")
var stringType = mustABIType("string")

var proofABI = abi.Arguments{{Type: bytes32SliceType}, {Type: bytesType}}
var nativeABI = abi.Arguments{{Type: uint256Type}}
var erc20ABI = abi.Arguments{{Type: addressType}, {Type: uint256Type}}
var erc721ABI = abi.Arguments{{Type: addressType}, {Type: uint256Type}, {Type: stringType}}
var erc1155ABI = abi.Arguments{{Type: addressType}, {Type: uint256Type}, {Type: stringType}, {Type: uint256Type}}

func (c *relayerConsumer) processEVMTransfer(
	task data.RelayTask,
	transfer rarimocore.Transfer,
	token tokenmanager.ChainInfo,
	tokenDetails tokenmanager.Item,
	network tokenmanager.Params,
	chain *config.EVMChain,
) error {
	log := c.log.WithField("op_id", task.OperationIndex)

	bridge, err := evmbind.NewBridge(chain.BridgeAddress, chain.RPC)
	if err != nil {
		return errors.Wrap(err, "failed to make an instance of ethereum bridge")
	}

	amount, err := getAmountOrDefault(transfer.Amount, big.NewInt(1))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid amount: %s", transfer.Amount))
	}
	receiver := common.HexToAddress(transfer.Receiver)
	origin := helpers.ToByte32(hexutil.MustDecode(task.Origin))
	/**
	Tweak the V value to make it compatible with the OpenZeppelin ECDSA implementation
	https://github.com/OpenZeppelin/openzeppelin-contracts/blob/a1948250ab8c441f6d327a65754cb20d2b1b4554/contracts/utils/cryptography/ECDSA.sol#L143
	*/
	signature := hexutil.MustDecode(task.Signature)
	signature[64] += 27

	proof, err := proofABI.Pack(task.MustParseMerklePath(), signature)
	if err != nil {
		return errors.Wrap(err, "failed to ABI encode the proof")
	}
	bundle, err := getBundleData(transfer)
	if err != nil {
		return errors.Wrap(err, "failed to parse bundle data")
	}

	opts := chain.TransactorOpts()
	nonce, err := chain.RPC.PendingNonceAt(context.TODO(), chain.SubmitterAddress)
	if err != nil {
		return errors.Wrap(err, "failed to fetch a nonce")
	}
	opts.Nonce = big.NewInt(int64(nonce))
	gasPrice, err := chain.RPC.SuggestGasPrice(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to get suggested gas price")
	}
	opts.GasPrice = gasPrice
	opts.GasLimit = uint64(300000)

	var call func() (*types.Transaction, error)
	switch tokenDetails.TokenType {
	case tokenmanager.Type_NATIVE:
		tokenData, err := nativeABI.Pack(amount)
		if err != nil {
			return errors.Wrap(err, "failed to ABI encode token data")
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
		tokenData, err := erc20ABI.Pack(common.HexToAddress(token.TokenAddress), amount)
		if err != nil {
			return errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawERC20(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
				tokenDetails.Wrapped,
			)
		}
	case tokenmanager.Type_ERC721:
		tokenID, err := parseTokenID(token.TokenId)
		if err != nil {
			return errors.Wrap(err, "failed to parse the tokenID")
		}

		tokenData, err := erc721ABI.Pack(common.HexToAddress(token.TokenAddress), tokenID, tokenDetails.Uri)
		if err != nil {
			return errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawERC721(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
				tokenDetails.Wrapped,
			)
		}
	case tokenmanager.Type_ERC1155:
		tokenID, err := parseTokenID(token.TokenId)
		if err != nil {
			return errors.Wrap(err, "failed to parse the tokenID")
		}
		tokenData, err := erc1155ABI.Pack(common.HexToAddress(token.TokenAddress), tokenID, tokenDetails.Uri, amount)
		if err != nil {
			return errors.Wrap(err, "failed to ABI encode token data")
		}

		call = func() (*types.Transaction, error) {
			return bridge.BridgeTransactor.WithdrawERC1155(
				opts,
				tokenData,
				bundle,
				origin,
				receiver,
				proof,
				tokenDetails.Wrapped,
			)
		}
	default:
		return errors.Errorf("token type %d is not supported", tokenDetails.TokenType)
	}

	tx, err := call()
	if err != nil {
		return errors.Wrap(err, "failed to call the withdraw method")
	}
	log.WithField("tx_id", tx.Hash()).Info("submitted transaction")

	receipt, err := bind.WaitMined(context.TODO(), chain.RPC, tx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for the transaction to be mined")
	}
	if receipt.Status == 0 {
		log.WithField("receipt", prettify(receipt)).Errorf("%s transaction failed", transfer.ToChain)
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

func mustABIType(evmType string) abi.Type {
	abiType, err := abi.NewType(evmType, "", nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to register ABI type"))
	}

	return abiType
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
	result.Salt = helpers.ToByte32(salt)

	return result, nil
}

func parseTokenID(rawTokenID string) (*big.Int, error) {
	rawBytes, err := hexutil.Decode(rawTokenID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the tokenID: %s", logan.Field("token_id", rawTokenID))
	}

	return big.NewInt(0).SetBytes(rawBytes), nil
}
