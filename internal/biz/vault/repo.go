package vault

import (
	"context"

	v1 "cobo-ucw-backend/api/ucw/v1"

	"gorm.io/gorm"
)

type Repo interface {
	Update(ctx context.Context, vault Vault) error
	UpdateByVaultID(ctx context.Context, data Vault, vaultID string) error
	Save(ctx context.Context, vault *Vault) (int64, error)
	GetByVaultID(ctx context.Context, vaultID string) (*Vault, error)
	Tx(db *gorm.DB) Repo
}

type GroupRepo interface {
	Update(ctx context.Context, data Group) error
	Save(ctx context.Context, data *Group) (int64, error)
	GetByGroupID(ctx context.Context, groupID string) (*Group, error)
	ListGroups(ctx context.Context, vaultID string, groupType v1.Group_GroupType) ([]*Group, error)
	Tx(db *gorm.DB) GroupRepo
}

type GroupNodeRepo interface {
	Update(ctx context.Context, data GroupNode) error
	Save(ctx context.Context, data *GroupNode) (int64, error)
	ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*GroupNode, error)
	ListNodeGroups(ctx context.Context, userID, nodeID string) ([]*GroupNode, error)
	BatchSave(ctx context.Context, data []*GroupNode) error
	Tx(db *gorm.DB) GroupNodeRepo
}

type TssRequestRepo interface {
	Tx(db *gorm.DB) TssRequestRepo
	Update(ctx context.Context, data TssRequest) error
	Save(ctx context.Context, data *TssRequest) (int64, error)
	GetByTssRequestID(ctx context.Context, tssRequestID string) (*TssRequest, error)
	ListGroupRelatedTssRequest(ctx context.Context, userID string, groupIDs []string, status int64) ([]*TssRequest, error)
	ListTssRequest(ctx context.Context, params ListTssRequestParams) ([]*TssRequest, error)
}

type ListTssRequestParams struct {
	VaultID       string
	TargetGroupID string
	SourceGroupID string
	Status        []v1.TssRequest_Status
	Limit         int
	LastID        int64
}
