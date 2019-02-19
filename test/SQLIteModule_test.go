package test

import (
	"testing"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"os"
)

func Test_SQLiteDemo(t *testing.T) {
	os.Remove("./foo.db")
	db, err := sql.Open("sqlite3", "./foo.db")
	if nil != err {
		t.Error(err)
	}
	defer db.Close()

	sqlStmt := `
		create table fol (id integer not null primary key, name text);
		delete from fol;
	`
	_, err = db.Exec(sqlStmt)
	if nil != err {
		t.Error(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}
	stmt, err := tx.Prepare("insert into fol(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			t.Error(err)
		}
	}
	tx.Commit()
}
