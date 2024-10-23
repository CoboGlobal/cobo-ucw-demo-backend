package database

import (
	"time"

	"cobo-ucw-backend/internal/conf"
	model2 "cobo-ucw-backend/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Data struct {
	DB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	var dialector gorm.Dialector
	switch c.Database.Driver {
	case "mysql":
		dialector = mysql.Open(c.Database.Source)
	case "postgres":
		dialector = postgres.Open(c.Database.Source)
	case "sqlite":
		dialector = sqlite.Open(c.Database.Source)
	default:
		panic("unsupported database driver")
	}
	schema.RegisterSerializer("json", schema.JSONSerializer{})

	gormLog := GormLoggerConfig{
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  log.LevelDebug,
	}
	if c.Log != nil {
		gormLog = GormLoggerConfig{
			SlowThreshold:             c.Log.SlowThreshold.AsDuration(),
			IgnoreRecordNotFoundError: true,
			LogLevel:                  log.ParseLevel(c.Log.Level),
		}
	}
	db, err := gorm.Open(dialector,
		&gorm.Config{
			PrepareStmt: true,
			Logger:      GormLogger(gormLog, log.NewHelper(logger))})
	if err != nil {
		log.NewHelper(logger).Fatalf("failed to open db driver: %s, err: %v", c.Database.Driver, err)
	}
	err = db.AutoMigrate(
		&model2.User{},
		&model2.Vault{},
		&model2.Transaction{},
		&model2.UserVault{},
		&model2.UserNode{},
		&model2.GroupNode{},
		&model2.Wallet{},
		&model2.Group{},
		&model2.TssRequest{},
		&model2.Address{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate table model: %v", err)
	}

	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{DB: db}, cleanup, nil
}
