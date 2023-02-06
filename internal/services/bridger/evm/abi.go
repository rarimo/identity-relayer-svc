package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

func mustABIType(evmType string) abi.Type {
	abiType, err := abi.NewType(evmType, "", nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to register ABI type"))
	}

	return abiType
}
