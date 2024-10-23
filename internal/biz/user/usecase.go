package user

import (
	"context"

	"cobo-ucw-backend/integration/portal"
	"cobo-ucw-backend/internal/data/model"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usecase struct {
	client    *portal.Client
	repo      Repo
	userNode  UserNodeRepo
	userVault UserVaultRepo
	logger    *log.Helper
}

func NewUsecase(client *portal.Client,
	repo Repo,
	userNode UserNodeRepo,
	userVault UserVaultRepo,
	logger log.Logger,
) *Usecase {
	return &Usecase{
		client:    client,
		repo:      repo,
		userNode:  userNode,
		userVault: userVault,
		logger:    log.NewHelper(logger),
	}
}

func (u *Usecase) Login(ctx context.Context, email string) (*User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Register(ctx, email)
	} else {
		uv, err := u.userVault.GetByUserID(ctx, user.UserID)
		if err != nil {
			return nil, err
		}

		if len(uv) > 0 {
			user.UserVault = uv[0]
		}
	}

	return user, nil
}

func (u *Usecase) Register(ctx context.Context, email string) (*User, error) {
	user := &User{
		User: &model.User{
			UserID: uuid.NewString(),
			Email:  email,
		},
	}
	_, err := u.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Usecase) GetUserNodes(ctx context.Context, userID string) (Nodes, error) {
	nodes, err := u.userNode.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return nodes, err
}
func (u *Usecase) BindUserNode(ctx context.Context, userID, nodeID string) (*UserNode, error) {
	un := &UserNode{
		UserNode: &model.UserNode{
			UserID: userID,
			NodeID: nodeID,
			Status: 0,
		},
	}
	_, err := u.userNode.Save(ctx, un)
	if err != nil {
		return nil, err
	}

	return un, err
}

func (u *Usecase) GetUserInfo(ctx context.Context, userID string) (*User, error) {
	user, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	uv, err := u.userVault.GetByUserID(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	if len(uv) > 0 {
		user.UserVault = uv[0]
	}
	return user, nil
}

func (u *Usecase) GetUserNodeByNodeIDs(ctx context.Context, nodeIDs []string) ([]*UserNode, error) {
	userNodes, err := u.userNode.GetNodesByIDs(ctx, nodeIDs)
	if err != nil {
		return nil, err
	}

	//userIDs := make([]string, 0, len(userNodes))
	//nodeUserSet := make(map[string][]*UserNode)
	//for _, each := range userNodes {
	//	nodes, ok := nodeUserSet[each.UserID]
	//	if ok {
	//		nodes = append(nodes, each)
	//	} else {
	//		nodes = []*UserNode{each}
	//		userIDs = append(userIDs, each.UserID)
	//	}
	//	nodeUserSet[each.UserID] = nodes
	//}
	//users, err := u.repo.GetByUserIDs(ctx, userIDs)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, each := range users {
	//	nodes := nodeUserSet[each.UserID]
	//	for _, node := range nodes {
	//		node.User = *each
	//	}
	//}
	return userNodes, nil
}

func (u *Usecase) BindUserVault(ctx context.Context, userID, vaultID string) error {
	_, err := u.userVault.Save(ctx, &UserVault{
		UserVault: &model.UserVault{
			UserID:  userID,
			VaultID: vaultID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
