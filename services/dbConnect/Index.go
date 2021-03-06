package dbConnect

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"../../exceptions"
	"../../config"
	"../../integrate/logger"
	"os"
	"fmt"
	"strings"
	"reflect"
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

type sqlStr struct {
	sqlType, tableName, orderBy string
	fields []string
	andConditions []string
	pageCondition string
}

func Select(tableName string) *sqlStr {
	return &sqlStr{
		sqlType: "SELECT",
		tableName: tableName,
	}
}

func (this *sqlStr) Fields(fields ...string) *sqlStr {
	this.fields = fields
	return this
}

func (this *sqlStr) AND(conditions ...string) *sqlStr {
	this.andConditions = conditions
	return this
}

func (this *sqlStr) Page(begin, limit int) *sqlStr {
	this.pageCondition = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, begin)
	return this
}

func (this *sqlStr) OrderBy(str string) *sqlStr {
	this.orderBy = str
	return this
}

func (this *sqlStr) Preview() string {
	baseSQL := fmt.Sprintf("%s %s FROM %s", this.sqlType, strings.Join(this.fields, ", "), this.tableName)
	if 0 != len(this.andConditions) {
		baseSQL = fmt.Sprintf("%s WHERE %s", baseSQL, strings.Join(this.andConditions, " AND "))
	}
	if "" != this.orderBy {
		baseSQL += fmt.Sprintf(" ORDER BY %s ", this.orderBy)
	}
	if "" != this.pageCondition {
		baseSQL += this.pageCondition
	}
	return baseSQL
}

func (this *sqlStr) Query(args ...interface{}) ([]map[string]interface{}, error) {
	conn := getConnection()
	defer conn.Close()
	logger.Logger("db-conn", this.Preview())
	rows, err := conn.Query(this.Preview(), args...)
	if nil != err {
		return nil, &exceptions.Error{Msg: "execute query fail. please check", Code: 400}
	}
	return ToMap(rows)
}

func ToMap(rows *sql.Rows) ([]map[string]interface{}, error){
	cols, _ := rows.Columns()
	result := make([]map[string]interface{}, 0)
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	for rows.Next() {
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, &exceptions.Error{Msg: err.Error(), Code: 500}
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			t := reflect.TypeOf(*val)
			if nil != *val && "[]uint8" == t.String() {
				m[colName] = string((*val).([]uint8))
			} else {
				m[colName] = *val
			}
		}
		result = append(result, m)
	}
	return result, nil
}

func ping() {
	conn, err := sql.Open(dbType, connStr)
	defer conn.Close()
	if nil != err {
		logger.Logger("db-sdk", "get Connection failed")
		os.Exit(100)
	}
	logger.Logger("db-sdk", "connection is ready")
}

func getConnection() *sql.DB {
	db, err := sql.Open(dbType, connStr)
	if nil != err {
		logger.Logger("db-sqlite", "get Connection failed")
		logger.Error("db-sqlite", err.Error())
	}
	return db
}

func WithQuery(sql string, callback func(rows *sql.Rows) (interface{}, error) , args ...interface{}) (interface{}, error) {
	conn := getConnection()
	defer conn.Close()
	stmt, err := conn.Prepare(sql)
	if nil != err {
		logger.Error("db-sdk", sql)
		return nil, &exceptions.Error{Msg: err.Error(), Code: 400}
	}
	logger.Logger("db-sdk", sql)
	rows, err := stmt.Query(args...)
	if nil != err {
		return nil, &exceptions.Error{Msg: err.Error(), Code: 400}
	}
	defer stmt.Close()
	return callback(rows)
}

func WithPrepare(sql string, callback func(*sql.Stmt) (interface{}, error)) (interface{}, error){
	conn := getConnection()
	defer conn.Close()
	stmt, _ := conn.Prepare(sql)
	defer stmt.Close()
	return callback(stmt)
}

func WithTransaction(callback func(*sql.Tx) (interface{}, error)) (interface{}, error) {
	conn := getConnection()
	defer conn.Close()
	tx, err := conn.Begin()
	if nil != err {
		return nil, &exceptions.Error{Msg: "Begin Transaction failed. please check code."}
	}
	defer tx.Commit()
	return callback(tx)
}