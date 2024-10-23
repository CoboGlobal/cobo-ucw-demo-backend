package user

import "context"

type Repo interface {
	Update(ctx context.Context, user User) error
	Save(ctx context.Context, user *User) (int64, error)
	GetByUserID(ctx context.Context, userID string) (*User, error)
	GetByUserIDs(ctx context.Context, userIDs []string) ([]*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserNodeRepo interface {
	Update(ctx context.Context, node UserNode) error
	Save(ctx context.Context, node *UserNode) (int64, error)
	GetByNodeID(ctx context.Context, nodeID string) (*UserNode, error)
	GetByUserID(ctx context.Context, userID string) ([]*UserNode, error)
	GetNodesByIDs(ctx context.Context, nodeIDs []string) ([]*UserNode, error)
}

type UserVaultRepo interface {
	Update(ctx context.Context, data UserVault) error
	Save(ctx context.Context, data *UserVault) (int64, error)
	GetByUserID(ctx context.Context, userID string) ([]*UserVault, error)
}
