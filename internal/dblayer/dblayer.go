package dblayer

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/logger"
)

const (
	qTimeout = 1500
)

var dbl *DBLayer

type Fields map[string]interface{}

type DBLayer struct {
	db      *sql.DB
	timeout time.Duration
}

type QUD struct { // Query-Update-Delete
	DBLayer
	table string
	//db      *sql.DB
	//tx      *sql.Tx
	cond   string // WHERE a = ? AND b = ?
	group  string
	order  string
	params []interface{} // params for condition expression
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
	dbl = &DBLayer{}
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

func Table(t string) *QUD {
	// TODO: maybe db: &dbl.DB ?
	//return &QUD{table: t, db: dbl.db}
	return &QUD{table: t, DBLayer: *dbl}
}

func (qud *QUD) Seek(args ...interface{}) *QUD {
	var where, list string
	var count int
	if len(args) == 0 {
		return qud
	}

	for pos, arg := range args {
		logger.Log.Info("arg: ", zap.Any("", arg))
		switch v := arg.(type) {
		case string:
			if 0 == pos {
				count = strings.Count(v, "?")
				where = v
			} else {
				qud.params = append(qud.params, arg)
			}
		case int, int64:
			if 0 == pos || pos > count {
				list = " = ?"
			}
			qud.params = append(qud.params, arg)

		case []int64: // TODO: use strings.Repeat()
			//list = "IN(" + JoinSlice(v) + ")"
			if len(v) > 0 {
				list = "IN(?" + strings.Repeat(", ?", len(v)-1) + ")"
				for i := range v {
					qud.params = append(qud.params, v[i])
				}
			}
		case []string: // TODO: use strings.Repeat()
			//list = "IN(" + JoinSlice(v) + ")"
			if len(v) > 0 {
				list = "IN(?" + strings.Repeat(", ?", len(v)-1) + ")"
				for i := range v {
					qud.params = append(qud.params, v[i])
				}
			}

		default:
			qud.params = append(qud.params, arg)
		}
	}
	if "" != list && "" == where {
		where = "id"
	}
	if "" != list || "" != where {
		qud.cond = " WHERE " + where + " " + list
	}

	return qud
}

func (qud *QUD) Get(tx *sql.Tx, mymap Fields, limits ...int64) (*sql.Rows, []interface{}, error) {
	return qud.get(tx, "", mymap, limits...)
}

func (qud *QUD) get(tx *sql.Tx, hint string, mymap Fields, limits ...int64) (res *sql.Rows, values []interface{}, err error) {

	//   fmt.Println("(qud *QUD) get hint ",hint)
	keys := ""
	values = make([]interface{}, len(mymap))

	i := 0
	for k, v := range mymap {
		if "" != keys {
			keys += ", "
		}
		if strings.ContainsRune(k, '(') || strings.ContainsRune(k, ' ') {
			keys += k
		} else if strings.ContainsRune(k, '.') {
			keys += strings.Replace(k, ".", ".`", 1) + "`"
		} else {
			keys += "" + k + ""
		}

		values[i] = v
		i++
	}
	q := "SELECT " + hint + " " + keys + " FROM " + qud.table + " " + qud.cond
	if "" != qud.order {
		q += " ORDER BY " + qud.order
	}
	if "" != qud.group {
		q += " GROUP BY " + qud.group
	}
	if len(limits) > 0 {
		q += " LIMIT ?" // + strconv.FormatInt(limits[0], 10)
		qud.params = append(qud.params, limits[0])
	}
	if len(limits) > 1 {
		q += ", ?"
		qud.params = append(qud.params, limits[1])
	}

	//  fmt.Println("get: ",q)
	logger.Log.Info("query: ", zap.String("", q))

	if nil == tx {
		//   fmt.Println("get qud.db.QueryContext ",q)
		ctx, _ := context.WithTimeout(context.TODO(), qud.timeout)
		res, err = qud.db.QueryContext(ctx, q, qud.params...)
	} else {
		//    fmt.Println("get tx.Query ",q)
		res, err = tx.Query(q, qud.params...)
	}

	qud.reset()

	return res, values, err
}

func (qud *QUD) Save(tx *sql.Tx, fields Fields) (err error) {
	var pId *int64
	tmp, ok := fields["id"]
	if ok {
		pId = tmp.(*int64)
	}
	if !ok { // new WITHOUT id field
		_, err = qud.Insert(tx, fields)
	} else if 0 == *pId { // new WITH id
		delete(fields, "id")
		*pId, err = qud.Insert(tx, fields)
	}
	//qud.reset()
	return
}

func inc(v *int) string {
	r := *v
	*v = *v + 1
	return strconv.Itoa(r)
}

func (qud *QUD) Insert(tx *sql.Tx, fields Fields) (id int64, err error) {
	//keys, values := fieldsMap(fields)
	cnt := 1
	keys := ""
	values := make([]interface{}, len(fields))

	i := 0
	for k, v := range fields {
		if "" != keys {
			keys += ", "
		}
		keys += k
		logger.Log.Info("field: ", zap.Any("", v))
		values[i] = v
		i++
	}

	q := "INSERT INTO " + qud.table + " (" + keys + ") VALUES ($" + inc(&cnt) + strings.Repeat(", $"+inc(&cnt), len(values)-1) + ")"

	logger.Log.Info("query: ", zap.String("", q))
	//var res sql.Result
	qud.params = values
	_, err = qud.execQuery(tx, q)

	/*
		if nil == err {
			id, err = res.LastInsertId()
		}
	*/

	return
}

func (qud *QUD) execQuery(tx *sql.Tx, q string) (sql.Result, error) {
	if nil == tx {
		ctx, _ := context.WithTimeout(context.TODO(), qud.timeout)
		return qud.db.ExecContext(ctx, q, qud.params...)
	} else {
		return tx.Exec(q, qud.params...)
	}
}

// all calls shoud be chained .Table().Seek().Get(nil, )
// but someone can reuse .Table() multiple times
// so we need to clean conditional part of object (cond & params)
// for just in case
func (qud *QUD) reset(args ...interface{}) {
	qud.cond = ""
	qud.group = ""
	qud.order = ""
	qud.params = []interface{}{}
}
