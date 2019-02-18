package test

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"../util"
	"testing"
	"fmt"
	"time"
)

const (
	dbType = "sqlite3"
	connStr = "/Users/afterloe/Afterloe/go/upm.dll"
	insertSQL = "INSERT INTO file(name, savePath, contentType, key, uploadTime, size, status, modifyTime, rev) values(?, ?, ?, ? ,? ,? ,? ,? ,?)"
)

func getConnection() (*sql.DB, error) {
	return sql.Open(dbType, connStr)
}

func Test_InsertFile(t *testing.T) {
	conn, _ := getConnection()
	defer conn.Close()
	tx, err := conn.Begin()
	if nil != err {
		t.Error(err)
		return
	}
	stmt, err := tx.Prepare(insertSQL)
	for i := 0; i < 18000; i++ {
		result, err := stmt.Exec("75239521d2744fe295123e0b7d7ae430_th.jpg", "/tmp/filesystem", "image/jpeg", util.GeneratorUUID(), time.Now().Unix(), 305810, true, 0, "")
		if nil != err {
			t.Error(err)
			return
		}
		lastId, err := result.LastInsertId()
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(fmt.Sprintf("insert id is %d", lastId))
	}
	tx.Commit()
}