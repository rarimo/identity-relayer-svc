package core

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	merkle "gitlab.com/rarimo/go-merkle"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimo "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	token "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

var (
	ErrUnsupportedContent = errors.New("unsupported content")
)

func (c *core) getContents(operations ...rarimo.Operation) ([]merkle.Content, error) {
	contents := make([]merkle.Content, 0, len(operations))

	for _, op := range operations {
		switch op.OperationType {
		case rarimo.OpType_TRANSFER:
			content, err := c.getTransferContent(op)
			if err != nil {
				return nil, err
			}

			if content != nil {
				contents = append(contents, content)
			}
		case rarimo.OpType_CHANGE_PARTIES:
			return nil, ErrUnsupportedContent
		case rarimo.OpType_FEE_TOKEN_MANAGEMENT:
			content, err := c.getFeeManagementContent(op)
			if err != nil {
				return nil, err
			}

			if content != nil {
				contents = append(contents, content)
			}
		case rarimo.OpType_CONTRACT_UPGRADE:
			content, err := c.getContractUpgradeContent(op)
			if err != nil {
				return nil, err
			}

			if content != nil {
				contents = append(contents, content)
			}
		case rarimo.OpType_IDENTITY_DEFAULT_TRANSFER:
			content, err := c.getIdentityDefaultTransferContent(op)
			if err != nil {
				return nil, err
			}

			if content != nil {
				contents = append(contents, content)
			}
		default:
			return nil, ErrUnsupportedContent
		}
	}

	return contents, nil
}

func (c *core) getTransferContent(op rarimo.Operation) (merkle.Content, error) {
	transfer, err := pkg.GetTransfer(op)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing operation details")
	}

	collectionDataResp, err := c.tm.CollectionData(context.TODO(), &token.QueryGetCollectionDataRequest{Chain: transfer.To.Chain, Address: transfer.To.Address})
	if err != nil {
		return nil, errors.Wrap(err, "error getting collection data entry")
	}

	collectionResp, err := c.tm.Collection(context.TODO(), &token.QueryGetCollectionRequest{Index: collectionDataResp.Data.Collection})
	if err != nil {
		return nil, errors.Wrap(err, "error getting collection data entry")
	}

	onChainItemResp, err := c.tm.OnChainItem(context.TODO(), &token.QueryGetOnChainItemRequest{Chain: transfer.To.Chain, Address: transfer.To.Address, TokenID: transfer.To.TokenID})
	if err != nil {
		return nil, errors.Wrap(err, "error getting on chain item entry")
	}

	itemResp, err := c.tm.Item(context.TODO(), &token.QueryGetItemRequest{Index: onChainItemResp.Item.Item})
	if err != nil {
		return nil, errors.Wrap(err, "error getting item entry")
	}

	networkResp, err := c.tm.NetworkParams(context.TODO(), &token.QueryNetworkParamsRequest{Name: transfer.To.Chain})
	if err != nil {
		return nil, errors.Wrap(err, "error getting network param entry")
	}

	bridgeparams := networkResp.Params.GetBridgeParams()
	if err != nil {
		return nil, errors.New("bridge params not found")
	}

	content, err := pkg.GetTransferContent(collectionResp.Collection, collectionDataResp.Data, itemResp.Item, bridgeparams, transfer)
	return *content, errors.Wrap(err, "error creating content")
}

func (c *core) getFeeManagementContent(op rarimo.Operation) (merkle.Content, error) {
	manage, err := pkg.GetFeeTokenManagement(op)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing operation details")
	}

	networkResp, err := c.tm.NetworkParams(context.TODO(), &token.QueryNetworkParamsRequest{Name: manage.Chain})
	if err != nil {
		return nil, errors.Wrap(err, "error getting network param entry")
	}

	feeparams := networkResp.Params.GetFeeParams()
	if err != nil {
		return nil, errors.New("bridge params not found")
	}

	content, err := pkg.GetFeeTokenManagementContent(feeparams, manage)
	return *content, errors.Wrap(err, "error creating content")
}

func (c *core) getContractUpgradeContent(op rarimo.Operation) (merkle.Content, error) {
	upgrade, err := pkg.GetContractUpgrade(op)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing operation details")
	}

	networkResp, err := c.tm.NetworkParams(context.TODO(), &token.QueryNetworkParamsRequest{Name: upgrade.Chain})
	if err != nil {
		return nil, errors.Wrap(err, "error getting network param entry")
	}

	content, err := pkg.GetContractUpgradeContent(networkResp.Params, upgrade)
	return *content, errors.Wrap(err, "error creating content")
}

func (c *core) getIdentityDefaultTransferContent(op rarimo.Operation) (merkle.Content, error) {
	transfer, err := pkg.GetIdentityDefaultTransfer(op)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing operation details")
	}

	content, err := pkg.GetIdentityDefaultTransferContent(transfer)
	return *content, errors.Wrap(err, "error creating content")
}
