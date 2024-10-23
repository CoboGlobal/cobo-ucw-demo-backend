package dao

import (
	"context"
	"strings"

	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewTssRequest(db *database.Data) *TssRequest {
	return &TssRequest{
		db,
	}
}

type TssRequest struct {
	*database.Data
}

func (d *TssRequest) Update(ctx context.Context, tx vault.TssRequest) error {
	return d.DB.WithContext(ctx).Model(&tx).Updates(tx).Error
}

func (d *TssRequest) Save(ctx context.Context, tx *vault.TssRequest) (int64, error) {
	return int64(tx.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(tx).Error
}

func (d *TssRequest) GetByTssRequestID(ctx context.Context, tssRequestID string) (*vault.TssRequest, error) {
	var res *model.TssRequest
	err := d.DB.WithContext(ctx).Where("request_id = ? ", tssRequestID).First(&res).Error
	return &vault.TssRequest{TssRequest: res}, err
}

func (d *TssRequest) ListGroupRelatedTssRequest(ctx context.Context, userID string, groupIDs []string, status int64) ([]*vault.TssRequest, error) {
	var res []*model.TssRequest

	if len(groupIDs) == 0 {
		return []*vault.TssRequest{}, nil
	}

	req := d.DB.WithContext(ctx).
		Model(&model.TssRequest{}).
		Where("user_id = ?", userID).
		Where("target_group_id in ? or source_group_id in ?", groupIDs, groupIDs)

	if status != 0 {
		req = req.Where("status = ?", status)
	}

	err := req.Find(&res).Error
	var list []*vault.TssRequest

	for _, each := range res {
		list = append(list, &vault.TssRequest{
			TssRequest: each,
		})
	}
	return list, err
}

func (d *TssRequest) Tx(tx *gorm.DB) vault.TssRequestRepo {
	return &TssRequest{Data: &database.Data{DB: tx}}
}

func (d *TssRequest) ListTssRequest(ctx context.Context, params vault.ListTssRequestParams) ([]*vault.TssRequest, error) {
	var res []*model.TssRequest

	var query []string
	var args []interface{}

	if params.VaultID != "" {
		query = append(query, " vault_id = ? ")
		args = append(args, params.VaultID)
	}

	if len(params.Status) != 0 {
		query = append(query, " status in (?) ")
		args = append(args, params.Status)
	}

	if params.LastID != 0 {
		query = append(query, " id > ? ")
		args = append(args, params.LastID)
	}

	prepare := d.DB.WithContext(ctx).Where(strings.Join(query, " and "), args...)
	var err error
	if params.Limit != 0 {
		err = prepare.Limit(params.Limit).Find(&res).Error
	} else {
		err = prepare.Find(&res).Error
	}
	if err != nil {
		return nil, err
	}

	var list []*vault.TssRequest

	for _, each := range res {
		list = append(list, &vault.TssRequest{TssRequest: each})
	}
	return list, err
}
