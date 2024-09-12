package dblayer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GlebZigert/gophermart/internal/config"
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

	fmt.Println("config.DatabaseDSN: ", config.DatabaseDSN)
	database, err := sql.Open("pgx", config.DatabaseDSN)
	if nil != err {
		return
	}

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err = database.PingContext(ctx)

	if nil != err {
		database.Close()
		return
	}

	dbl.Bind(database, qTimeout)

	err = dbl.MakeTables(tables, true)
	if nil != err {
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
