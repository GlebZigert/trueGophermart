package dblayer

import (
	"log"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	//dbURL := "postgres://pg:pass@localhost:5432/crud"

	db, err := gorm.Open(postgres.Open(config.DatabaseDSN), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&users.User{})

	return db
}
