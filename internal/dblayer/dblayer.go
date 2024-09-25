package dblayer

import (
	"github.com/GlebZigert/trueGophermart/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(DatabaseDSN string) (db *gorm.DB, err error) {
	//dbURL := "postgres://pg:pass@localhost:5432/crud"

	if db, err = gorm.Open(postgres.Open(DatabaseDSN), &gorm.Config{}); err != nil {
		return
	}

	if err = db.AutoMigrate(&model.User{}); err != nil {
		return
	}
	if err = db.AutoMigrate(&model.Order{}); err != nil {
		return
	}

	return
}
