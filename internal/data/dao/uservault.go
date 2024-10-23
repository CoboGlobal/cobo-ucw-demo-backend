package dao

import (
	"context"

	"cobo-ucw-backend/internal/biz/user"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewUserVault(db *database.Data) *UserVault {
	return &UserVault{db}
}

type UserVault struct {
	*database.Data
}

func (d *UserVault) Update(ctx context.Context, data user.UserVault) error {
	return d.DB.WithContext(ctx).Model(&data.UserVault).Updates(data.UserVault).Error
}

func (d *UserVault) Save(ctx context.Context, data *user.UserVault) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.UserVault).Error
}

func (d *UserVault) GetByUserID(ctx context.Context, userID string) ([]*user.UserVault, error) {
	var res []*model.UserVault
	err := d.DB.WithContext(ctx).Where("user_id = ? ", userID).Find(&res).Error

	var list []*user.UserVault
	for _, each := range res {
		list = append(list, &user.UserVault{UserVault: each})
	}
	return list, err
}

func (d *UserVault) Tx(tx *gorm.DB) *UserVault {
	return &UserVault{Data: &database.Data{DB: tx}}
}
