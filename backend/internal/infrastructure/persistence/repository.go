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
	dsn := cfg.GetString("database.dsn")
	//	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
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

func newRepositoryTask(repository *Repository) *RepositoryTask {
	return &RepositoryTask{repository: repository}
}

func (r RepositoryTask) OnStart(ctx context.Context) error {
	return nil
}

func (r RepositoryTask) OnStop(ctx context.Context) error {
	return r.repository.Close()
}
