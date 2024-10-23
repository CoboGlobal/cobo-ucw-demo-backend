package dao

import (
	"context"

	"cobo-ucw-backend/internal/biz/user"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm/clause"
)

func NewUserNode(db *database.Data) *UserNode {
	return &UserNode{db}
}

type UserNode struct {
	*database.Data
}

func (d *UserNode) Update(ctx context.Context, data user.UserNode) error {
	return d.DB.WithContext(ctx).Model(&data.UserNode).Updates(data.UserNode).Error
}

func (d *UserNode) Save(ctx context.Context, data *user.UserNode) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.UserNode).Error
}

func (d *UserNode) GetByNodeID(ctx context.Context, nodeID string) (*user.UserNode, error) {
	var res *model.UserNode
	err := d.DB.WithContext(ctx).Where("node_id = ? ", nodeID).First(&res).Error
	return &user.UserNode{UserNode: res}, err
}

func (d *UserNode) GetByUserID(ctx context.Context, userID string) ([]*user.UserNode, error) {
	var res []*model.UserNode
	err := d.DB.WithContext(ctx).Where("user_id = ? ", userID).Find(&res).Error

	var list []*user.UserNode
	for _, each := range res {
		list = append(list, &user.UserNode{UserNode: each})
	}
	return list, err
}

func (d *UserNode) GetNodesByIDs(ctx context.Context, nodeIDs []string) ([]*user.UserNode, error) {
	var res []*model.UserNode
	err := d.DB.WithContext(ctx).Where("node_id in (?) ", nodeIDs).Find(&res).Error
	var list []*user.UserNode
	for _, each := range res {
		list = append(list, &user.UserNode{UserNode: each})
	}
	return list, err
}
