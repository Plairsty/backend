package persistence

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/postgres"
	"plairsty/backend/domain/entity"
)

type Repositories struct {
	db *gorm.DB
}

func NewRepository(
	DbDriver,
	DbUser,
	DbPassword,
	DbPort,
	DbHost,
	DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost,
		DbPort,
		DbUser,
		DbName,
		DbPassword,
	)
	db, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	return &Repositories{
		db: db,
	}, nil
}

// Close closes the database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

func (s *Repositories) GetDB() *gorm.DB {
	return s.db
}

// AutoMigrate migrates all tables
func (s *Repositories) AutoMigrate() error {
	return s.db.AutoMigrate(&entity.User{}).Error
}
