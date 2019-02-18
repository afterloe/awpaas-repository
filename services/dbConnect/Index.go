package dbConnect

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"../../exceptions"
	"../../config"
	"../../integrate/logger"
	"os"
)

var (
	dbType, connStr string
)

func init() {
	db := config.GetByTarget(config.Get("services"), "db")
	dbType = config.GetByTarget(db, "dbType").(string)
	connStr = config.GetByTarget(db, "database").(string)
	// 如果开启检测
	flg := config.GetByTarget(db, "ping")
	if nil != flg {
		if flg.(bool) {
			ping()
		}
	}
}

func ping() {
	conn, err := getConnection()
	defer conn.Close()
	if nil != err {
		logger.Logger("db-sdk", string("get Connection failed"))
		os.Exit(100)
	}
}

func getConnection() (*sql.DB, error) {
	return sql.Open(dbType, connStr)
}

func WithPrepare(sql string, callback func(*sql.Stmt) (map[string]interface{}, error)) (map[string]interface{}, error){
	conn, err := getConnection()
	if nil != err {
		return nil, &exceptions.Error{Msg: "get Connection failed. please check."}
	}
	defer conn.Close()
	stmt, err := conn.Prepare(sql)
	stmt.Close()
	return callback(stmt)
}

func WithTransaction(callback func(*sql.Tx) (map[string]interface{}, error)) (map[string]interface{}, error) {
	conn, err := getConnection()
	if nil != err {
		return nil, &exceptions.Error{Msg: "get Connection failed. please check."}
	}
	defer conn.Close()
	tx, err := conn.Begin()
	if nil != err {
		return nil, &exceptions.Error{Msg: "Begin Transaction failed. please check code."}
	}
	defer tx.Commit()
	return callback(tx)
}