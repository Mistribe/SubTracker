package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg config.Configuration) *Repository {
	dsn := cfg.GetString("DATABASE_DSN")
	//	dsn := "host=localhost user=postgres password=postgres dbname=app port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	return &Repository{
		db: db,
	}
}

func (r Repository) GetDB() *gorm.DB {
	return r.db
}

func (r Repository) Close() error {
	if r.db != nil {
		sqlDB, err := r.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (r Repository) Ping() error {
	if r.db != nil {
		sqlDB, err := r.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}
	return errors.New("database connection is nil")
}

type RepositoryTask struct {
	repository *Repository
}

func (r RepositoryTask) Priority() int {
	return -1
}

func newRepositoryTask(repository *Repository) *RepositoryTask {
	return &RepositoryTask{repository: repository}
}

func (r RepositoryTask) OnStart(_ context.Context) error {
	if err := r.repository.db.AutoMigrate(
		&subscriptionModel{},
		&labelModel{},
		&familyMemberModel{},
		&subscriptionPaymentModel{},
		&subscriptionFamilyMemberModel{},
		&subscriptionLabelModel{},
	); err != nil {
		return err
	}

	return nil
}

func (r RepositoryTask) OnStop(_ context.Context) error {
	return r.repository.Close()
}
