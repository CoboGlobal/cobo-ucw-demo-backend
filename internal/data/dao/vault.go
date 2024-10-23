package dao

import (
	"context"

	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewVault(db *database.Data) *Vault {
	return &Vault{db}
}

type Vault struct {
	*database.Data
}

func (d *Vault) Update(ctx context.Context, data vault.Vault) error {
	return d.DB.WithContext(ctx).Model(&data.Vault).Updates(data.Vault).Error
}

func (d *Vault) UpdateByVaultID(ctx context.Context, data vault.Vault, vaultID string) error {
	return d.DB.WithContext(ctx).Model(&data.Vault).Where("vault_id = ? ", vaultID).Updates(data).Error
}

func (d *Vault) Save(ctx context.Context, data *vault.Vault) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.Vault).Error
}

func (d *Vault) GetByVaultID(ctx context.Context, vaultID string) (*vault.Vault, error) {
	var res *model.Vault
	err := d.DB.WithContext(ctx).Where("vault_id = ? ", vaultID).First(&res).Error
	return &vault.Vault{Vault: res}, err
}

func (d *Vault) GetByUserIDProjectID(ctx context.Context, userID, projectID string) (*vault.Vault, error) {
	var res *model.Vault
	err := d.DB.WithContext(ctx).Where("user_id = ? and project_id = ? ", userID, projectID).First(&res).Error
	return &vault.Vault{Vault: res}, err
}

func (d *Vault) Tx(tx *gorm.DB) vault.Repo {
	return &Vault{Data: &database.Data{DB: tx}}
}
