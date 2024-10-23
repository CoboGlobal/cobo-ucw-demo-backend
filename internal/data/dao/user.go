package dao

import (
	"context"

	"cobo-ucw-backend/internal/biz/user"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm/clause"
)

func NewUser(db *database.Data) *User {
	return &User{
		db,
	}
}

type User struct {
	*database.Data
}

func (d *User) Update(ctx context.Context, data user.User) error {
	return d.DB.WithContext(ctx).Model(&data.User).Updates(data.User).Error
}

func (d *User) Save(ctx context.Context, data *user.User) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.User).Error
}

func (d *User) GetByUserID(ctx context.Context, userID string) (*user.User, error) {
	var res *model.User
	err := d.DB.WithContext(ctx).Where("user_id = ? ", userID).First(&res).Error
	return &user.User{User: res}, err
}
func (d *User) GetByUserIDs(ctx context.Context, userIDs []string) ([]*user.User, error) {
	var res []*model.User
	err := d.DB.WithContext(ctx).Where("user_id in (?) ", userIDs).Find(&res).Error
	var list []*user.User
	for _, each := range res {
		list = append(list, &user.User{User: each})
	}
	return list, err
}

func (d *User) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var res *model.User
	err := d.DB.WithContext(ctx).Where("email = ? ", email).First(&res).Error
	return &user.User{User: res}, err
}
