// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ILightweightStateGistRootData is an auto generated low-level Go binding around an user-defined struct.
type ILightweightStateGistRootData struct {
	Root               *big.Int
	CreatedAtTimestamp *big.Int
}

// ILightweightStateIdentitiesStatesRootData is an auto generated low-level Go binding around an user-defined struct.
type ILightweightStateIdentitiesStatesRootData struct {
	Root         [32]byte
	SetTimestamp *big.Int
}

// ILightweightStateStatesMerkleData is an auto generated low-level Go binding around an user-defined struct.
type ILightweightStateStatesMerkleData struct {
	IssuerId           *big.Int
	IssuerState        *big.Int
	CreatedAtTimestamp *big.Int
	MerkleProof        [][32]byte
}

// LightweightStateV3MetaData contains all meta data concerning the LightweightStateV3 contract.
var LightweightStateV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newGistRoot\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newIdentitesStatesRoot\",\"type\":\"bytes32\"}],\"name\":\"SignedStateTransited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"P\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sourceStateContract_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sourceChainName_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"chainName_\",\"type\":\"string\"}],\"name\":\"__LightweightState_init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"facade_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"chainName_\",\"type\":\"string\"}],\"name\":\"__Signers_init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"newSignerPubKey_\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature_\",\"type\":\"bytes\"}],\"name\":\"changeSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newSourceStateContract_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature_\",\"type\":\"bytes\"}],\"name\":\"changeSourceStateContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"methodId_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"contractAddress_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"signHash_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature_\",\"type\":\"bytes\"}],\"name\":\"checkSignatureAndIncrementNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"facade\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"root_\",\"type\":\"uint256\"}],\"name\":\"geGISTRootData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structILightweightState.GistRootData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentGISTRootInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structILightweightState.GistRootData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGISTRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root_\",\"type\":\"bytes32\"}],\"name\":\"getIdentitiesStatesRootData\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"setTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structILightweightState.IdentitiesStatesRootData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"methodId_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"contractAddress_\",\"type\":\"address\"}],\"name\":\"getSigComponents\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"chainName_\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"nonce_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structILightweightState.GistRootData\",\"name\":\"gistData_\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"identitiesStatesRoot_\",\"type\":\"bytes32\"}],\"name\":\"getSignHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"identitiesStatesRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root_\",\"type\":\"bytes32\"}],\"name\":\"isIdentitiesStatesRootExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newIdentitiesStatesRoot_\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structILightweightState.GistRootData\",\"name\":\"gistData_\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proof_\",\"type\":\"bytes\"}],\"name\":\"signedTransitState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"signer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sourceChainName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sourceStateContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature_\",\"type\":\"bytes\"}],\"name\":\"upgradeToWithSig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"methodId_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"contractAddress_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newAddress_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature_\",\"type\":\"bytes\"}],\"name\":\"validateChangeAddressSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"issuerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issuerState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"internalType\":\"structILightweightState.StatesMerkleData\",\"name\":\"statesMerkleData_\",\"type\":\"tuple\"}],\"name\":\"verifyStatesMerkleData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// LightweightStateV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use LightweightStateV3MetaData.ABI instead.
var LightweightStateV3ABI = LightweightStateV3MetaData.ABI

// LightweightStateV3 is an auto generated Go binding around an Ethereum contract.
type LightweightStateV3 struct {
	LightweightStateV3Caller     // Read-only binding to the contract
	LightweightStateV3Transactor // Write-only binding to the contract
	LightweightStateV3Filterer   // Log filterer for contract events
}

// LightweightStateV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type LightweightStateV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LightweightStateV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type LightweightStateV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LightweightStateV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LightweightStateV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LightweightStateV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LightweightStateV3Session struct {
	Contract     *LightweightStateV3 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LightweightStateV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LightweightStateV3CallerSession struct {
	Contract *LightweightStateV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// LightweightStateV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LightweightStateV3TransactorSession struct {
	Contract     *LightweightStateV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LightweightStateV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type LightweightStateV3Raw struct {
	Contract *LightweightStateV3 // Generic contract binding to access the raw methods on
}

// LightweightStateV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LightweightStateV3CallerRaw struct {
	Contract *LightweightStateV3Caller // Generic read-only contract binding to access the raw methods on
}

// LightweightStateV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LightweightStateV3TransactorRaw struct {
	Contract *LightweightStateV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewLightweightStateV3 creates a new instance of LightweightStateV3, bound to a specific deployed contract.
func NewLightweightStateV3(address common.Address, backend bind.ContractBackend) (*LightweightStateV3, error) {
	contract, err := bindLightweightStateV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3{LightweightStateV3Caller: LightweightStateV3Caller{contract: contract}, LightweightStateV3Transactor: LightweightStateV3Transactor{contract: contract}, LightweightStateV3Filterer: LightweightStateV3Filterer{contract: contract}}, nil
}

// NewLightweightStateV3Caller creates a new read-only instance of LightweightStateV3, bound to a specific deployed contract.
func NewLightweightStateV3Caller(address common.Address, caller bind.ContractCaller) (*LightweightStateV3Caller, error) {
	contract, err := bindLightweightStateV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3Caller{contract: contract}, nil
}

// NewLightweightStateV3Transactor creates a new write-only instance of LightweightStateV3, bound to a specific deployed contract.
func NewLightweightStateV3Transactor(address common.Address, transactor bind.ContractTransactor) (*LightweightStateV3Transactor, error) {
	contract, err := bindLightweightStateV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3Transactor{contract: contract}, nil
}

// NewLightweightStateV3Filterer creates a new log filterer instance of LightweightStateV3, bound to a specific deployed contract.
func NewLightweightStateV3Filterer(address common.Address, filterer bind.ContractFilterer) (*LightweightStateV3Filterer, error) {
	contract, err := bindLightweightStateV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3Filterer{contract: contract}, nil
}

// bindLightweightStateV3 binds a generic wrapper to an already deployed contract.
func bindLightweightStateV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LightweightStateV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LightweightStateV3 *LightweightStateV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LightweightStateV3.Contract.LightweightStateV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LightweightStateV3 *LightweightStateV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.LightweightStateV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LightweightStateV3 *LightweightStateV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.LightweightStateV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LightweightStateV3 *LightweightStateV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LightweightStateV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LightweightStateV3 *LightweightStateV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LightweightStateV3 *LightweightStateV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.contract.Transact(opts, method, params...)
}

// P is a free data retrieval call binding the contract method 0x8b8fbd92.
//
// Solidity: function P() view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3Caller) P(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "P")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// P is a free data retrieval call binding the contract method 0x8b8fbd92.
//
// Solidity: function P() view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3Session) P() (*big.Int, error) {
	return _LightweightStateV3.Contract.P(&_LightweightStateV3.CallOpts)
}

// P is a free data retrieval call binding the contract method 0x8b8fbd92.
//
// Solidity: function P() view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3CallerSession) P() (*big.Int, error) {
	return _LightweightStateV3.Contract.P(&_LightweightStateV3.CallOpts)
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() view returns(string)
func (_LightweightStateV3 *LightweightStateV3Caller) ChainName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "chainName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() view returns(string)
func (_LightweightStateV3 *LightweightStateV3Session) ChainName() (string, error) {
	return _LightweightStateV3.Contract.ChainName(&_LightweightStateV3.CallOpts)
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() view returns(string)
func (_LightweightStateV3 *LightweightStateV3CallerSession) ChainName() (string, error) {
	return _LightweightStateV3.Contract.ChainName(&_LightweightStateV3.CallOpts)
}

// Facade is a free data retrieval call binding the contract method 0x5014a0fb.
//
// Solidity: function facade() view returns(address)
func (_LightweightStateV3 *LightweightStateV3Caller) Facade(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "facade")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Facade is a free data retrieval call binding the contract method 0x5014a0fb.
//
// Solidity: function facade() view returns(address)
func (_LightweightStateV3 *LightweightStateV3Session) Facade() (common.Address, error) {
	return _LightweightStateV3.Contract.Facade(&_LightweightStateV3.CallOpts)
}

// Facade is a free data retrieval call binding the contract method 0x5014a0fb.
//
// Solidity: function facade() view returns(address)
func (_LightweightStateV3 *LightweightStateV3CallerSession) Facade() (common.Address, error) {
	return _LightweightStateV3.Contract.Facade(&_LightweightStateV3.CallOpts)
}

// GeGISTRootData is a free data retrieval call binding the contract method 0x0dd93b5d.
//
// Solidity: function geGISTRootData(uint256 root_) view returns((uint256,uint256))
func (_LightweightStateV3 *LightweightStateV3Caller) GeGISTRootData(opts *bind.CallOpts, root_ *big.Int) (ILightweightStateGistRootData, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "geGISTRootData", root_)

	if err != nil {
		return *new(ILightweightStateGistRootData), err
	}

	out0 := *abi.ConvertType(out[0], new(ILightweightStateGistRootData)).(*ILightweightStateGistRootData)

	return out0, err

}

// GeGISTRootData is a free data retrieval call binding the contract method 0x0dd93b5d.
//
// Solidity: function geGISTRootData(uint256 root_) view returns((uint256,uint256))
func (_LightweightStateV3 *LightweightStateV3Session) GeGISTRootData(root_ *big.Int) (ILightweightStateGistRootData, error) {
	return _LightweightStateV3.Contract.GeGISTRootData(&_LightweightStateV3.CallOpts, root_)
}

// GeGISTRootData is a free data retrieval call binding the contract method 0x0dd93b5d.
//
// Solidity: function geGISTRootData(uint256 root_) view returns((uint256,uint256))
func (_LightweightStateV3 *LightweightStateV3CallerSession) GeGISTRootData(root_ *big.Int) (ILightweightStateGistRootData, error) {
	return _LightweightStateV3.Contract.GeGISTRootData(&_LightweightStateV3.CallOpts, root_)
}

// GetCurrentGISTRootInfo is a free data retrieval call binding the contract method 0xaf7a3f59.
//
// Solidity: function getCurrentGISTRootInfo() view returns((uint256,uint256))
func (_LightweightStateV3 *LightweightStateV3Caller) GetCurrentGISTRootInfo(opts *bind.CallOpts) (ILightweightStateGistRootData, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "getCurrentGISTRootInfo")

	if err != nil {
		return *new(ILightweightStateGistRootData), err
	}

	out0 := *abi.ConvertType(out[0], new(ILightweightStateGistRootData)).(*ILightweightStateGistRootData)

	return out0, err

}

// GetCurrentGISTRootInfo is a free data retrieval call binding the contract method 0xaf7a3f59.
//
// Solidity: function getCurrentGISTRootInfo() view returns((uint256,uint256))
func (_LightweightStateV3 *LightweightStateV3Session) GetCurrentGISTRootInfo() (ILightweightStateGistRootData, error) {
	return _LightweightStateV3.Contract.GetCurrentGISTRootInfo(&_LightweightStateV3.CallOpts)
}

// GetCurrentGISTRootInfo is a free data retrieval call binding the contract method 0xaf7a3f59.
//
// Solidity: function getCurrentGISTRootInfo() view returns((uint256,uint256))
func (_LightweightStateV3 *LightweightStateV3CallerSession) GetCurrentGISTRootInfo() (ILightweightStateGistRootData, error) {
	return _LightweightStateV3.Contract.GetCurrentGISTRootInfo(&_LightweightStateV3.CallOpts)
}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3Caller) GetGISTRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "getGISTRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3Session) GetGISTRoot() (*big.Int, error) {
	return _LightweightStateV3.Contract.GetGISTRoot(&_LightweightStateV3.CallOpts)
}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3CallerSession) GetGISTRoot() (*big.Int, error) {
	return _LightweightStateV3.Contract.GetGISTRoot(&_LightweightStateV3.CallOpts)
}

// GetIdentitiesStatesRootData is a free data retrieval call binding the contract method 0xa055a692.
//
// Solidity: function getIdentitiesStatesRootData(bytes32 root_) view returns((bytes32,uint256))
func (_LightweightStateV3 *LightweightStateV3Caller) GetIdentitiesStatesRootData(opts *bind.CallOpts, root_ [32]byte) (ILightweightStateIdentitiesStatesRootData, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "getIdentitiesStatesRootData", root_)

	if err != nil {
		return *new(ILightweightStateIdentitiesStatesRootData), err
	}

	out0 := *abi.ConvertType(out[0], new(ILightweightStateIdentitiesStatesRootData)).(*ILightweightStateIdentitiesStatesRootData)

	return out0, err

}

// GetIdentitiesStatesRootData is a free data retrieval call binding the contract method 0xa055a692.
//
// Solidity: function getIdentitiesStatesRootData(bytes32 root_) view returns((bytes32,uint256))
func (_LightweightStateV3 *LightweightStateV3Session) GetIdentitiesStatesRootData(root_ [32]byte) (ILightweightStateIdentitiesStatesRootData, error) {
	return _LightweightStateV3.Contract.GetIdentitiesStatesRootData(&_LightweightStateV3.CallOpts, root_)
}

// GetIdentitiesStatesRootData is a free data retrieval call binding the contract method 0xa055a692.
//
// Solidity: function getIdentitiesStatesRootData(bytes32 root_) view returns((bytes32,uint256))
func (_LightweightStateV3 *LightweightStateV3CallerSession) GetIdentitiesStatesRootData(root_ [32]byte) (ILightweightStateIdentitiesStatesRootData, error) {
	return _LightweightStateV3.Contract.GetIdentitiesStatesRootData(&_LightweightStateV3.CallOpts, root_)
}

// GetSigComponents is a free data retrieval call binding the contract method 0x827e099e.
//
// Solidity: function getSigComponents(uint8 methodId_, address contractAddress_) view returns(string chainName_, uint256 nonce_)
func (_LightweightStateV3 *LightweightStateV3Caller) GetSigComponents(opts *bind.CallOpts, methodId_ uint8, contractAddress_ common.Address) (struct {
	ChainName string
	Nonce     *big.Int
}, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "getSigComponents", methodId_, contractAddress_)

	outstruct := new(struct {
		ChainName string
		Nonce     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainName = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Nonce = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSigComponents is a free data retrieval call binding the contract method 0x827e099e.
//
// Solidity: function getSigComponents(uint8 methodId_, address contractAddress_) view returns(string chainName_, uint256 nonce_)
func (_LightweightStateV3 *LightweightStateV3Session) GetSigComponents(methodId_ uint8, contractAddress_ common.Address) (struct {
	ChainName string
	Nonce     *big.Int
}, error) {
	return _LightweightStateV3.Contract.GetSigComponents(&_LightweightStateV3.CallOpts, methodId_, contractAddress_)
}

// GetSigComponents is a free data retrieval call binding the contract method 0x827e099e.
//
// Solidity: function getSigComponents(uint8 methodId_, address contractAddress_) view returns(string chainName_, uint256 nonce_)
func (_LightweightStateV3 *LightweightStateV3CallerSession) GetSigComponents(methodId_ uint8, contractAddress_ common.Address) (struct {
	ChainName string
	Nonce     *big.Int
}, error) {
	return _LightweightStateV3.Contract.GetSigComponents(&_LightweightStateV3.CallOpts, methodId_, contractAddress_)
}

// GetSignHash is a free data retrieval call binding the contract method 0x62f8f8a4.
//
// Solidity: function getSignHash((uint256,uint256) gistData_, bytes32 identitiesStatesRoot_) view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3Caller) GetSignHash(opts *bind.CallOpts, gistData_ ILightweightStateGistRootData, identitiesStatesRoot_ [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "getSignHash", gistData_, identitiesStatesRoot_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetSignHash is a free data retrieval call binding the contract method 0x62f8f8a4.
//
// Solidity: function getSignHash((uint256,uint256) gistData_, bytes32 identitiesStatesRoot_) view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3Session) GetSignHash(gistData_ ILightweightStateGistRootData, identitiesStatesRoot_ [32]byte) ([32]byte, error) {
	return _LightweightStateV3.Contract.GetSignHash(&_LightweightStateV3.CallOpts, gistData_, identitiesStatesRoot_)
}

// GetSignHash is a free data retrieval call binding the contract method 0x62f8f8a4.
//
// Solidity: function getSignHash((uint256,uint256) gistData_, bytes32 identitiesStatesRoot_) view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3CallerSession) GetSignHash(gistData_ ILightweightStateGistRootData, identitiesStatesRoot_ [32]byte) ([32]byte, error) {
	return _LightweightStateV3.Contract.GetSignHash(&_LightweightStateV3.CallOpts, gistData_, identitiesStatesRoot_)
}

// IdentitiesStatesRoot is a free data retrieval call binding the contract method 0xe08e70bb.
//
// Solidity: function identitiesStatesRoot() view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3Caller) IdentitiesStatesRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "identitiesStatesRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// IdentitiesStatesRoot is a free data retrieval call binding the contract method 0xe08e70bb.
//
// Solidity: function identitiesStatesRoot() view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3Session) IdentitiesStatesRoot() ([32]byte, error) {
	return _LightweightStateV3.Contract.IdentitiesStatesRoot(&_LightweightStateV3.CallOpts)
}

// IdentitiesStatesRoot is a free data retrieval call binding the contract method 0xe08e70bb.
//
// Solidity: function identitiesStatesRoot() view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3CallerSession) IdentitiesStatesRoot() ([32]byte, error) {
	return _LightweightStateV3.Contract.IdentitiesStatesRoot(&_LightweightStateV3.CallOpts)
}

// IsIdentitiesStatesRootExists is a free data retrieval call binding the contract method 0xbfd73455.
//
// Solidity: function isIdentitiesStatesRootExists(bytes32 root_) view returns(bool)
func (_LightweightStateV3 *LightweightStateV3Caller) IsIdentitiesStatesRootExists(opts *bind.CallOpts, root_ [32]byte) (bool, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "isIdentitiesStatesRootExists", root_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsIdentitiesStatesRootExists is a free data retrieval call binding the contract method 0xbfd73455.
//
// Solidity: function isIdentitiesStatesRootExists(bytes32 root_) view returns(bool)
func (_LightweightStateV3 *LightweightStateV3Session) IsIdentitiesStatesRootExists(root_ [32]byte) (bool, error) {
	return _LightweightStateV3.Contract.IsIdentitiesStatesRootExists(&_LightweightStateV3.CallOpts, root_)
}

// IsIdentitiesStatesRootExists is a free data retrieval call binding the contract method 0xbfd73455.
//
// Solidity: function isIdentitiesStatesRootExists(bytes32 root_) view returns(bool)
func (_LightweightStateV3 *LightweightStateV3CallerSession) IsIdentitiesStatesRootExists(root_ [32]byte) (bool, error) {
	return _LightweightStateV3.Contract.IsIdentitiesStatesRootExists(&_LightweightStateV3.CallOpts, root_)
}

// Nonces is a free data retrieval call binding the contract method 0xed3218a2.
//
// Solidity: function nonces(address , uint8 ) view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3Caller) Nonces(opts *bind.CallOpts, arg0 common.Address, arg1 uint8) (*big.Int, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "nonces", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0xed3218a2.
//
// Solidity: function nonces(address , uint8 ) view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3Session) Nonces(arg0 common.Address, arg1 uint8) (*big.Int, error) {
	return _LightweightStateV3.Contract.Nonces(&_LightweightStateV3.CallOpts, arg0, arg1)
}

// Nonces is a free data retrieval call binding the contract method 0xed3218a2.
//
// Solidity: function nonces(address , uint8 ) view returns(uint256)
func (_LightweightStateV3 *LightweightStateV3CallerSession) Nonces(arg0 common.Address, arg1 uint8) (*big.Int, error) {
	return _LightweightStateV3.Contract.Nonces(&_LightweightStateV3.CallOpts, arg0, arg1)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3Session) ProxiableUUID() ([32]byte, error) {
	return _LightweightStateV3.Contract.ProxiableUUID(&_LightweightStateV3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LightweightStateV3 *LightweightStateV3CallerSession) ProxiableUUID() ([32]byte, error) {
	return _LightweightStateV3.Contract.ProxiableUUID(&_LightweightStateV3.CallOpts)
}

// Signer is a free data retrieval call binding the contract method 0x238ac933.
//
// Solidity: function signer() view returns(address)
func (_LightweightStateV3 *LightweightStateV3Caller) Signer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "signer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Signer is a free data retrieval call binding the contract method 0x238ac933.
//
// Solidity: function signer() view returns(address)
func (_LightweightStateV3 *LightweightStateV3Session) Signer() (common.Address, error) {
	return _LightweightStateV3.Contract.Signer(&_LightweightStateV3.CallOpts)
}

// Signer is a free data retrieval call binding the contract method 0x238ac933.
//
// Solidity: function signer() view returns(address)
func (_LightweightStateV3 *LightweightStateV3CallerSession) Signer() (common.Address, error) {
	return _LightweightStateV3.Contract.Signer(&_LightweightStateV3.CallOpts)
}

// SourceChainName is a free data retrieval call binding the contract method 0xe4ffd04a.
//
// Solidity: function sourceChainName() view returns(string)
func (_LightweightStateV3 *LightweightStateV3Caller) SourceChainName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "sourceChainName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// SourceChainName is a free data retrieval call binding the contract method 0xe4ffd04a.
//
// Solidity: function sourceChainName() view returns(string)
func (_LightweightStateV3 *LightweightStateV3Session) SourceChainName() (string, error) {
	return _LightweightStateV3.Contract.SourceChainName(&_LightweightStateV3.CallOpts)
}

// SourceChainName is a free data retrieval call binding the contract method 0xe4ffd04a.
//
// Solidity: function sourceChainName() view returns(string)
func (_LightweightStateV3 *LightweightStateV3CallerSession) SourceChainName() (string, error) {
	return _LightweightStateV3.Contract.SourceChainName(&_LightweightStateV3.CallOpts)
}

// SourceStateContract is a free data retrieval call binding the contract method 0xfc228319.
//
// Solidity: function sourceStateContract() view returns(address)
func (_LightweightStateV3 *LightweightStateV3Caller) SourceStateContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "sourceStateContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SourceStateContract is a free data retrieval call binding the contract method 0xfc228319.
//
// Solidity: function sourceStateContract() view returns(address)
func (_LightweightStateV3 *LightweightStateV3Session) SourceStateContract() (common.Address, error) {
	return _LightweightStateV3.Contract.SourceStateContract(&_LightweightStateV3.CallOpts)
}

// SourceStateContract is a free data retrieval call binding the contract method 0xfc228319.
//
// Solidity: function sourceStateContract() view returns(address)
func (_LightweightStateV3 *LightweightStateV3CallerSession) SourceStateContract() (common.Address, error) {
	return _LightweightStateV3.Contract.SourceStateContract(&_LightweightStateV3.CallOpts)
}

// VerifyStatesMerkleData is a free data retrieval call binding the contract method 0xb0d46f2c.
//
// Solidity: function verifyStatesMerkleData((uint256,uint256,uint256,bytes32[]) statesMerkleData_) view returns(bool, bytes32)
func (_LightweightStateV3 *LightweightStateV3Caller) VerifyStatesMerkleData(opts *bind.CallOpts, statesMerkleData_ ILightweightStateStatesMerkleData) (bool, [32]byte, error) {
	var out []interface{}
	err := _LightweightStateV3.contract.Call(opts, &out, "verifyStatesMerkleData", statesMerkleData_)

	if err != nil {
		return *new(bool), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// VerifyStatesMerkleData is a free data retrieval call binding the contract method 0xb0d46f2c.
//
// Solidity: function verifyStatesMerkleData((uint256,uint256,uint256,bytes32[]) statesMerkleData_) view returns(bool, bytes32)
func (_LightweightStateV3 *LightweightStateV3Session) VerifyStatesMerkleData(statesMerkleData_ ILightweightStateStatesMerkleData) (bool, [32]byte, error) {
	return _LightweightStateV3.Contract.VerifyStatesMerkleData(&_LightweightStateV3.CallOpts, statesMerkleData_)
}

// VerifyStatesMerkleData is a free data retrieval call binding the contract method 0xb0d46f2c.
//
// Solidity: function verifyStatesMerkleData((uint256,uint256,uint256,bytes32[]) statesMerkleData_) view returns(bool, bytes32)
func (_LightweightStateV3 *LightweightStateV3CallerSession) VerifyStatesMerkleData(statesMerkleData_ ILightweightStateStatesMerkleData) (bool, [32]byte, error) {
	return _LightweightStateV3.Contract.VerifyStatesMerkleData(&_LightweightStateV3.CallOpts, statesMerkleData_)
}

// LightweightStateInit is a paid mutator transaction binding the contract method 0x27e25ecf.
//
// Solidity: function __LightweightState_init(address signer_, address sourceStateContract_, string sourceChainName_, string chainName_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) LightweightStateInit(opts *bind.TransactOpts, signer_ common.Address, sourceStateContract_ common.Address, sourceChainName_ string, chainName_ string) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "__LightweightState_init", signer_, sourceStateContract_, sourceChainName_, chainName_)
}

// LightweightStateInit is a paid mutator transaction binding the contract method 0x27e25ecf.
//
// Solidity: function __LightweightState_init(address signer_, address sourceStateContract_, string sourceChainName_, string chainName_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) LightweightStateInit(signer_ common.Address, sourceStateContract_ common.Address, sourceChainName_ string, chainName_ string) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.LightweightStateInit(&_LightweightStateV3.TransactOpts, signer_, sourceStateContract_, sourceChainName_, chainName_)
}

// LightweightStateInit is a paid mutator transaction binding the contract method 0x27e25ecf.
//
// Solidity: function __LightweightState_init(address signer_, address sourceStateContract_, string sourceChainName_, string chainName_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) LightweightStateInit(signer_ common.Address, sourceStateContract_ common.Address, sourceChainName_ string, chainName_ string) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.LightweightStateInit(&_LightweightStateV3.TransactOpts, signer_, sourceStateContract_, sourceChainName_, chainName_)
}

// SignersInit is a paid mutator transaction binding the contract method 0x3baa7892.
//
// Solidity: function __Signers_init(address signer_, address facade_, string chainName_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) SignersInit(opts *bind.TransactOpts, signer_ common.Address, facade_ common.Address, chainName_ string) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "__Signers_init", signer_, facade_, chainName_)
}

// SignersInit is a paid mutator transaction binding the contract method 0x3baa7892.
//
// Solidity: function __Signers_init(address signer_, address facade_, string chainName_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) SignersInit(signer_ common.Address, facade_ common.Address, chainName_ string) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.SignersInit(&_LightweightStateV3.TransactOpts, signer_, facade_, chainName_)
}

// SignersInit is a paid mutator transaction binding the contract method 0x3baa7892.
//
// Solidity: function __Signers_init(address signer_, address facade_, string chainName_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) SignersInit(signer_ common.Address, facade_ common.Address, chainName_ string) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.SignersInit(&_LightweightStateV3.TransactOpts, signer_, facade_, chainName_)
}

// ChangeSigner is a paid mutator transaction binding the contract method 0x497f6959.
//
// Solidity: function changeSigner(bytes newSignerPubKey_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) ChangeSigner(opts *bind.TransactOpts, newSignerPubKey_ []byte, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "changeSigner", newSignerPubKey_, signature_)
}

// ChangeSigner is a paid mutator transaction binding the contract method 0x497f6959.
//
// Solidity: function changeSigner(bytes newSignerPubKey_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) ChangeSigner(newSignerPubKey_ []byte, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.ChangeSigner(&_LightweightStateV3.TransactOpts, newSignerPubKey_, signature_)
}

// ChangeSigner is a paid mutator transaction binding the contract method 0x497f6959.
//
// Solidity: function changeSigner(bytes newSignerPubKey_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) ChangeSigner(newSignerPubKey_ []byte, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.ChangeSigner(&_LightweightStateV3.TransactOpts, newSignerPubKey_, signature_)
}

// ChangeSourceStateContract is a paid mutator transaction binding the contract method 0x89aeb0f5.
//
// Solidity: function changeSourceStateContract(address newSourceStateContract_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) ChangeSourceStateContract(opts *bind.TransactOpts, newSourceStateContract_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "changeSourceStateContract", newSourceStateContract_, signature_)
}

// ChangeSourceStateContract is a paid mutator transaction binding the contract method 0x89aeb0f5.
//
// Solidity: function changeSourceStateContract(address newSourceStateContract_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) ChangeSourceStateContract(newSourceStateContract_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.ChangeSourceStateContract(&_LightweightStateV3.TransactOpts, newSourceStateContract_, signature_)
}

// ChangeSourceStateContract is a paid mutator transaction binding the contract method 0x89aeb0f5.
//
// Solidity: function changeSourceStateContract(address newSourceStateContract_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) ChangeSourceStateContract(newSourceStateContract_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.ChangeSourceStateContract(&_LightweightStateV3.TransactOpts, newSourceStateContract_, signature_)
}

// CheckSignatureAndIncrementNonce is a paid mutator transaction binding the contract method 0xe3754f90.
//
// Solidity: function checkSignatureAndIncrementNonce(uint8 methodId_, address contractAddress_, bytes32 signHash_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) CheckSignatureAndIncrementNonce(opts *bind.TransactOpts, methodId_ uint8, contractAddress_ common.Address, signHash_ [32]byte, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "checkSignatureAndIncrementNonce", methodId_, contractAddress_, signHash_, signature_)
}

// CheckSignatureAndIncrementNonce is a paid mutator transaction binding the contract method 0xe3754f90.
//
// Solidity: function checkSignatureAndIncrementNonce(uint8 methodId_, address contractAddress_, bytes32 signHash_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) CheckSignatureAndIncrementNonce(methodId_ uint8, contractAddress_ common.Address, signHash_ [32]byte, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.CheckSignatureAndIncrementNonce(&_LightweightStateV3.TransactOpts, methodId_, contractAddress_, signHash_, signature_)
}

// CheckSignatureAndIncrementNonce is a paid mutator transaction binding the contract method 0xe3754f90.
//
// Solidity: function checkSignatureAndIncrementNonce(uint8 methodId_, address contractAddress_, bytes32 signHash_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) CheckSignatureAndIncrementNonce(methodId_ uint8, contractAddress_ common.Address, signHash_ [32]byte, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.CheckSignatureAndIncrementNonce(&_LightweightStateV3.TransactOpts, methodId_, contractAddress_, signHash_, signature_)
}

// SignedTransitState is a paid mutator transaction binding the contract method 0x189a5073.
//
// Solidity: function signedTransitState(bytes32 newIdentitiesStatesRoot_, (uint256,uint256) gistData_, bytes proof_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) SignedTransitState(opts *bind.TransactOpts, newIdentitiesStatesRoot_ [32]byte, gistData_ ILightweightStateGistRootData, proof_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "signedTransitState", newIdentitiesStatesRoot_, gistData_, proof_)
}

// SignedTransitState is a paid mutator transaction binding the contract method 0x189a5073.
//
// Solidity: function signedTransitState(bytes32 newIdentitiesStatesRoot_, (uint256,uint256) gistData_, bytes proof_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) SignedTransitState(newIdentitiesStatesRoot_ [32]byte, gistData_ ILightweightStateGistRootData, proof_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.SignedTransitState(&_LightweightStateV3.TransactOpts, newIdentitiesStatesRoot_, gistData_, proof_)
}

// SignedTransitState is a paid mutator transaction binding the contract method 0x189a5073.
//
// Solidity: function signedTransitState(bytes32 newIdentitiesStatesRoot_, (uint256,uint256) gistData_, bytes proof_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) SignedTransitState(newIdentitiesStatesRoot_ [32]byte, gistData_ ILightweightStateGistRootData, proof_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.SignedTransitState(&_LightweightStateV3.TransactOpts, newIdentitiesStatesRoot_, gistData_, proof_)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_LightweightStateV3 *LightweightStateV3Session) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.UpgradeTo(&_LightweightStateV3.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.UpgradeTo(&_LightweightStateV3.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LightweightStateV3 *LightweightStateV3Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.UpgradeToAndCall(&_LightweightStateV3.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.UpgradeToAndCall(&_LightweightStateV3.TransactOpts, newImplementation, data)
}

// UpgradeToWithSig is a paid mutator transaction binding the contract method 0x52d04661.
//
// Solidity: function upgradeToWithSig(address newImplementation_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) UpgradeToWithSig(opts *bind.TransactOpts, newImplementation_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "upgradeToWithSig", newImplementation_, signature_)
}

// UpgradeToWithSig is a paid mutator transaction binding the contract method 0x52d04661.
//
// Solidity: function upgradeToWithSig(address newImplementation_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) UpgradeToWithSig(newImplementation_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.UpgradeToWithSig(&_LightweightStateV3.TransactOpts, newImplementation_, signature_)
}

// UpgradeToWithSig is a paid mutator transaction binding the contract method 0x52d04661.
//
// Solidity: function upgradeToWithSig(address newImplementation_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) UpgradeToWithSig(newImplementation_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.UpgradeToWithSig(&_LightweightStateV3.TransactOpts, newImplementation_, signature_)
}

// ValidateChangeAddressSignature is a paid mutator transaction binding the contract method 0x7d1e764b.
//
// Solidity: function validateChangeAddressSignature(uint8 methodId_, address contractAddress_, address newAddress_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Transactor) ValidateChangeAddressSignature(opts *bind.TransactOpts, methodId_ uint8, contractAddress_ common.Address, newAddress_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.contract.Transact(opts, "validateChangeAddressSignature", methodId_, contractAddress_, newAddress_, signature_)
}

// ValidateChangeAddressSignature is a paid mutator transaction binding the contract method 0x7d1e764b.
//
// Solidity: function validateChangeAddressSignature(uint8 methodId_, address contractAddress_, address newAddress_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3Session) ValidateChangeAddressSignature(methodId_ uint8, contractAddress_ common.Address, newAddress_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.ValidateChangeAddressSignature(&_LightweightStateV3.TransactOpts, methodId_, contractAddress_, newAddress_, signature_)
}

// ValidateChangeAddressSignature is a paid mutator transaction binding the contract method 0x7d1e764b.
//
// Solidity: function validateChangeAddressSignature(uint8 methodId_, address contractAddress_, address newAddress_, bytes signature_) returns()
func (_LightweightStateV3 *LightweightStateV3TransactorSession) ValidateChangeAddressSignature(methodId_ uint8, contractAddress_ common.Address, newAddress_ common.Address, signature_ []byte) (*types.Transaction, error) {
	return _LightweightStateV3.Contract.ValidateChangeAddressSignature(&_LightweightStateV3.TransactOpts, methodId_, contractAddress_, newAddress_, signature_)
}

// LightweightStateV3AdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the LightweightStateV3 contract.
type LightweightStateV3AdminChangedIterator struct {
	Event *LightweightStateV3AdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LightweightStateV3AdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LightweightStateV3AdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LightweightStateV3AdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LightweightStateV3AdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LightweightStateV3AdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LightweightStateV3AdminChanged represents a AdminChanged event raised by the LightweightStateV3 contract.
type LightweightStateV3AdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_LightweightStateV3 *LightweightStateV3Filterer) FilterAdminChanged(opts *bind.FilterOpts) (*LightweightStateV3AdminChangedIterator, error) {

	logs, sub, err := _LightweightStateV3.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3AdminChangedIterator{contract: _LightweightStateV3.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_LightweightStateV3 *LightweightStateV3Filterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *LightweightStateV3AdminChanged) (event.Subscription, error) {

	logs, sub, err := _LightweightStateV3.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LightweightStateV3AdminChanged)
				if err := _LightweightStateV3.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_LightweightStateV3 *LightweightStateV3Filterer) ParseAdminChanged(log types.Log) (*LightweightStateV3AdminChanged, error) {
	event := new(LightweightStateV3AdminChanged)
	if err := _LightweightStateV3.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LightweightStateV3BeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the LightweightStateV3 contract.
type LightweightStateV3BeaconUpgradedIterator struct {
	Event *LightweightStateV3BeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LightweightStateV3BeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LightweightStateV3BeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LightweightStateV3BeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LightweightStateV3BeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LightweightStateV3BeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LightweightStateV3BeaconUpgraded represents a BeaconUpgraded event raised by the LightweightStateV3 contract.
type LightweightStateV3BeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_LightweightStateV3 *LightweightStateV3Filterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*LightweightStateV3BeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _LightweightStateV3.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3BeaconUpgradedIterator{contract: _LightweightStateV3.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_LightweightStateV3 *LightweightStateV3Filterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *LightweightStateV3BeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _LightweightStateV3.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LightweightStateV3BeaconUpgraded)
				if err := _LightweightStateV3.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_LightweightStateV3 *LightweightStateV3Filterer) ParseBeaconUpgraded(log types.Log) (*LightweightStateV3BeaconUpgraded, error) {
	event := new(LightweightStateV3BeaconUpgraded)
	if err := _LightweightStateV3.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LightweightStateV3InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the LightweightStateV3 contract.
type LightweightStateV3InitializedIterator struct {
	Event *LightweightStateV3Initialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LightweightStateV3InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LightweightStateV3Initialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LightweightStateV3Initialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LightweightStateV3InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LightweightStateV3InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LightweightStateV3Initialized represents a Initialized event raised by the LightweightStateV3 contract.
type LightweightStateV3Initialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_LightweightStateV3 *LightweightStateV3Filterer) FilterInitialized(opts *bind.FilterOpts) (*LightweightStateV3InitializedIterator, error) {

	logs, sub, err := _LightweightStateV3.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3InitializedIterator{contract: _LightweightStateV3.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_LightweightStateV3 *LightweightStateV3Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *LightweightStateV3Initialized) (event.Subscription, error) {

	logs, sub, err := _LightweightStateV3.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LightweightStateV3Initialized)
				if err := _LightweightStateV3.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_LightweightStateV3 *LightweightStateV3Filterer) ParseInitialized(log types.Log) (*LightweightStateV3Initialized, error) {
	event := new(LightweightStateV3Initialized)
	if err := _LightweightStateV3.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LightweightStateV3SignedStateTransitedIterator is returned from FilterSignedStateTransited and is used to iterate over the raw logs and unpacked data for SignedStateTransited events raised by the LightweightStateV3 contract.
type LightweightStateV3SignedStateTransitedIterator struct {
	Event *LightweightStateV3SignedStateTransited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LightweightStateV3SignedStateTransitedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LightweightStateV3SignedStateTransited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LightweightStateV3SignedStateTransited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LightweightStateV3SignedStateTransitedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LightweightStateV3SignedStateTransitedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LightweightStateV3SignedStateTransited represents a SignedStateTransited event raised by the LightweightStateV3 contract.
type LightweightStateV3SignedStateTransited struct {
	NewGistRoot            *big.Int
	NewIdentitesStatesRoot [32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSignedStateTransited is a free log retrieval operation binding the contract event 0x8e8ff16553fbf4a457c08a5e07cc27c8aac14b9a1a8e1f546a6c1b9366304a56.
//
// Solidity: event SignedStateTransited(uint256 newGistRoot, bytes32 newIdentitesStatesRoot)
func (_LightweightStateV3 *LightweightStateV3Filterer) FilterSignedStateTransited(opts *bind.FilterOpts) (*LightweightStateV3SignedStateTransitedIterator, error) {

	logs, sub, err := _LightweightStateV3.contract.FilterLogs(opts, "SignedStateTransited")
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3SignedStateTransitedIterator{contract: _LightweightStateV3.contract, event: "SignedStateTransited", logs: logs, sub: sub}, nil
}

// WatchSignedStateTransited is a free log subscription operation binding the contract event 0x8e8ff16553fbf4a457c08a5e07cc27c8aac14b9a1a8e1f546a6c1b9366304a56.
//
// Solidity: event SignedStateTransited(uint256 newGistRoot, bytes32 newIdentitesStatesRoot)
func (_LightweightStateV3 *LightweightStateV3Filterer) WatchSignedStateTransited(opts *bind.WatchOpts, sink chan<- *LightweightStateV3SignedStateTransited) (event.Subscription, error) {

	logs, sub, err := _LightweightStateV3.contract.WatchLogs(opts, "SignedStateTransited")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LightweightStateV3SignedStateTransited)
				if err := _LightweightStateV3.contract.UnpackLog(event, "SignedStateTransited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignedStateTransited is a log parse operation binding the contract event 0x8e8ff16553fbf4a457c08a5e07cc27c8aac14b9a1a8e1f546a6c1b9366304a56.
//
// Solidity: event SignedStateTransited(uint256 newGistRoot, bytes32 newIdentitesStatesRoot)
func (_LightweightStateV3 *LightweightStateV3Filterer) ParseSignedStateTransited(log types.Log) (*LightweightStateV3SignedStateTransited, error) {
	event := new(LightweightStateV3SignedStateTransited)
	if err := _LightweightStateV3.contract.UnpackLog(event, "SignedStateTransited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LightweightStateV3UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the LightweightStateV3 contract.
type LightweightStateV3UpgradedIterator struct {
	Event *LightweightStateV3Upgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LightweightStateV3UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LightweightStateV3Upgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LightweightStateV3Upgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LightweightStateV3UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LightweightStateV3UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LightweightStateV3Upgraded represents a Upgraded event raised by the LightweightStateV3 contract.
type LightweightStateV3Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LightweightStateV3 *LightweightStateV3Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*LightweightStateV3UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _LightweightStateV3.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &LightweightStateV3UpgradedIterator{contract: _LightweightStateV3.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LightweightStateV3 *LightweightStateV3Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *LightweightStateV3Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _LightweightStateV3.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LightweightStateV3Upgraded)
				if err := _LightweightStateV3.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LightweightStateV3 *LightweightStateV3Filterer) ParseUpgraded(log types.Log) (*LightweightStateV3Upgraded, error) {
	event := new(LightweightStateV3Upgraded)
	if err := _LightweightStateV3.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
