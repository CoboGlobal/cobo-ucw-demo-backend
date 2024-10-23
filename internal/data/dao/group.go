package dao

import (
	"context"
	"strings"

	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewGroup(db *database.Data) *Group {
	return &Group{db}
}

type Group struct {
	*database.Data
}

func (d *Group) Update(ctx context.Context, data vault.Group) error {
	return d.DB.WithContext(ctx).Model(&data).Updates(data).Error
}

func (d *Group) Save(ctx context.Context, data *vault.Group) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.Group).Error
}

func (d *Group) GetByGroupID(ctx context.Context, groupID string) (*vault.Group, error) {
	var res *model.Group
	err := d.DB.WithContext(ctx).Where("group_id = ? ", groupID).First(&res).Error
	return &vault.Group{Group: res}, err
}

func (d *Group) ListGroups(ctx context.Context, vaultID string, groupType v1.Group_GroupType) ([]*vault.Group, error) {
	var res []*model.Group

	var query = []string{"vault_id = ?"}

	var args = []interface{}{vaultID}

	if groupType != 0 {
		query = append(query, "group_type = ? ")
		args = append(args, groupType)
	}

	err := d.DB.WithContext(ctx).Where(strings.Join(query, " and "), args...).Find(&res).Error
	if err != nil {
		return nil, err
	}
	var list []*vault.Group

	for _, each := range res {
		list = append(list, &vault.Group{Group: each})
	}
	return list, err
}

func (d *Group) Tx(tx *gorm.DB) vault.GroupRepo {
	return &Group{Data: &database.Data{DB: tx}}
}
