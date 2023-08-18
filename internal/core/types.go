package core

import rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"

type IdentityTransferDetails struct {
	Operation *rarimocore.IdentityDefaultTransfer
	Proof     []byte
}
