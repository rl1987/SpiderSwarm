package spiderswarm

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SQLiteSpiderBusBackend struct {
	SpiderBusBackend

	dbConn *sql.DB
}

func NewSQLiteSpiderBusBackend(sqliteFilePath string) *SQLiteSpiderBusBackend {
	if sqliteFilePath == "" {
		sqliteDirPath, err := ioutil.TempDir(os.TempDir(), "spiderbus_")
		if err != nil {
			log.Error(fmt.Sprintf("Failed to create temp dir for SQLiteSpiderBusBackend: %v", err))
			return nil
		}

		sqliteFilePath = sqliteDirPath + "/" + "spiderbus.db"
	}

	dbConn, err := sql.Open("sqlite3", sqliteFilePath)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open DB connection: %v", err))
		return nil
	}

	tx, _ := dbConn.Begin()

	tx.Exec("CREATE TABLE IF NOT EXISTS scheduledTasks (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")
	tx.Exec("CREATE TABLE IF NOT EXISTS taskPromises (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")
	tx.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")

	tx.Commit()

	return &SQLiteSpiderBusBackend{
		dbConn: dbConn,
	}
}
