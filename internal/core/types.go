package core

import rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"

type IdentityStateTransferDetails struct {
	Operation *rarimocore.IdentityStateTransfer
	Proof     []byte
}

type IdentityGISTTransferDetails struct {
	Operation *rarimocore.IdentityGISTTransfer
	Proof     []byte
}
