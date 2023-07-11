package core

import (
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

type IdentityTransferDetails struct {
	OpIndex    string
	Proof      []byte
	Signature  string
	MerklePath [][32]byte
}

type TransferDetails struct {
	Transfer      rarimocore.Transfer
	DstCollection tokenmanager.CollectionData
	Item          tokenmanager.Item
	Signature     string
	Origin        string
	MerklePath    [][32]byte
}
