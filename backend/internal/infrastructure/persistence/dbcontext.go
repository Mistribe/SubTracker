package persistence

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Oleexo/config-go"
	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseContext struct {
	db *gorm.DB
}

func NewDatabaseContext(
	cfg config.Configuration,
	logger *slog.Logger) *DatabaseContext {
	dsn := cfg.GetString("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: sloggorm.New(
			sloggorm.WithLogger(logger),
			sloggorm.WithTraceAll(), // Log all SQL queries
		),
	})
	if err != nil {
		panic(err)
	}
	// Get the underlying SQL database connection
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxOpenConns(25)                  // Maximum number of open connections
	sqlDB.SetMaxIdleConns(5)                   // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum connection lifetime
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Maximum connection idle time

	return &DatabaseContext{
		db: db,
	}
}

func (r DatabaseContext) GetDB() *gorm.DB {
	return r.db
}

func (r DatabaseContext) Close() error {
	if r.db != nil {
		sqlDB, err := r.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (r DatabaseContext) Ping() error {
	if r.db != nil {
		sqlDB, err := r.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}
	return errors.New("database connection is nil")
}

type DatabaseContextTask struct {
	repository *DatabaseContext
}

func (r DatabaseContextTask) Priority() int {
	return -1
}

func newRepositoryTask(repository *DatabaseContext) *DatabaseContextTask {
	return &DatabaseContextTask{repository: repository}
}

func (r DatabaseContextTask) OnStart(_ context.Context) error {
	if err := r.repository.db.AutoMigrate(
		&subscriptionSqlModel{},
		&subscriptionServiceUserModel{},
		&labelSqlModel{},
		&familySqlModel{},
		&familyMemberSqlModel{},
		&providerSqlModel{},
		&providerLabelSqlModel{},
		&providerPlanSqlModel{},
		&providerPriceSqlModel{},
	); err != nil {
		return err
	}

	return nil
}

func (r DatabaseContextTask) OnStop(_ context.Context) error {
	return r.repository.Close()
}
