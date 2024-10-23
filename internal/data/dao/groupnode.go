package dao

import (
	"context"

	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewGroupNode(db *database.Data) *GroupNode {
	return &GroupNode{db}
}

type GroupNode struct {
	*database.Data
}

func (d *GroupNode) Update(ctx context.Context, data vault.GroupNode) error {
	return d.DB.WithContext(ctx).Model(&data).Updates(data).Error
}

func (d *GroupNode) Save(ctx context.Context, data *vault.GroupNode) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data).Error
}

func (d *GroupNode) BatchSave(ctx context.Context, data []*vault.GroupNode) error {
	return d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data).Error
}

func (d *GroupNode) ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*vault.GroupNode, error) {
	var res []*model.GroupNode
	err := d.DB.WithContext(ctx).Where("group_id in (?) ", groupIDs).Find(&res).Error
	if err != nil {
		return nil, err
	}
	var list []*vault.GroupNode

	for _, each := range res {
		list = append(list, &vault.GroupNode{GroupNode: each})
	}
	return list, err
}

func (d *GroupNode) ListNodeGroups(ctx context.Context, userID, nodeID string) ([]*vault.GroupNode, error) {
	var res []*model.GroupNode
	err := d.DB.WithContext(ctx).Where("user_id = ? and node_id = ? ", userID, nodeID).Find(&res).Error
	var list []*vault.GroupNode

	for _, each := range res {
		list = append(list, &vault.GroupNode{GroupNode: each})
	}
	return list, err
}

func (d *GroupNode) Tx(tx *gorm.DB) vault.GroupNodeRepo {
	return &GroupNode{Data: &database.Data{DB: tx}}
}
