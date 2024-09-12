package dblayer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/logger"
)

const (
	qTimeout = 1500
)

var dbl DBLayer

type DBLayer struct {
	db      *sql.DB
	timeout time.Duration
}

func Init() (err error) {
	logger.Log.Info("config.DatabaseDSN: ", zap.String("config.DatabaseDSN", config.DatabaseDSN))
	database, err := sql.Open("pgx", config.DatabaseDSN)
	if nil != err {
		logger.Log.Error("sql.Open: ", zap.String("", err.Error()))

		return
	}

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err = database.PingContext(ctx)

	if err != nil {
		logger.Log.Error("database.PingContext: ", zap.String("", err.Error()))
		database.Close()
		return
	}

	dbl.Bind(database, qTimeout)

	err = dbl.MakeTables(tables, true)
	if nil != err {
		logger.Log.Error("dbl.MakeTables: ", zap.String("", err.Error()))
		dbl.Close()
		return
	}
	return
}

func (dbl *DBLayer) Bind(db *sql.DB, timeout int) {
	dbl.db = db
	dbl.timeout = time.Duration(timeout) * time.Millisecond
}

func (dbl *DBLayer) Close() (err error) {
	return dbl.db.Close()
}

func (dbl *DBLayer) MakeTables(tables []string, strict bool) (err error) {
	for i := 0; i < len(tables) && (nil == err || !strict); i++ {
		//log.Println(tables[i])
		ctx, _ := context.WithTimeout(context.TODO(), dbl.timeout)
		_, err = dbl.db.ExecContext(ctx, tables[i])
	}
	return
}
