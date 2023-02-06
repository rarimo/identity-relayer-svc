package core

import (
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

type TransferDetails struct {
	Transfer     rarimocore.Transfer
	Token        tokenmanager.Info
	TokenDetails tokenmanager.Item
	Signature    string
	Origin       string
	MerklePath   [][32]byte
}
