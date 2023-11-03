package core

import rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"

type IdentityTransferDetails struct {
	Operation *rarimocore.IdentityDefaultTransfer
	Proof     []byte
}

type IdentityGISTTransferDetails struct {
	Operation *rarimocore.IdentityGISTTransfer
	Proof     []byte
}
