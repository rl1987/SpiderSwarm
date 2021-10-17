package spsw

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SQLiteSpiderBusBackend struct {
	SpiderBusBackend

	dbConn         *sql.DB
	sqliteFilePath string
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

	// https://github.com/mattn/go-sqlite3/issues/39#issuecomment-13469905
	dbConn, err := sql.Open("sqlite3", sqliteFilePath+"?cache=shared&mode=rwc")
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open DB connection: %v", err))
		return nil
	}

	tx, _ := dbConn.Begin()

	tx.Exec("CREATE TABLE IF NOT EXISTS scheduledTasks (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")
	tx.Exec("CREATE TABLE IF NOT EXISTS taskPromises (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")
	tx.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")

	tx.Commit()

	fmt.Printf("Created new SQLite DB at: %s\n", sqliteFilePath)

	return &SQLiteSpiderBusBackend{
		dbConn:         dbConn,
		sqliteFilePath: sqliteFilePath,
	}
}

func (ssbb *SQLiteSpiderBusBackend) maybePrintError(err error) {
	if err != nil && err != sql.ErrNoRows {
		spew.Dump(err)
	}
}

func (ssbb *SQLiteSpiderBusBackend) encodeEntry(entry interface{}) []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(entry)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}

func (ssbb *SQLiteSpiderBusBackend) decodeEntry(raw []byte, entry interface{}) interface{} {
	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(entry)
	if err != nil {
		spew.Dump(err)
	}

	return &entry
}

func (ssbb *SQLiteSpiderBusBackend) SendScheduledTask(scheduledTask *ScheduledTask) error {
	raw := ssbb.encodeEntry(scheduledTask)

	tmpDbConn, _ := sql.Open("sqlite3", ssbb.sqliteFilePath+"?cache=shared&mode=rwc")
	defer tmpDbConn.Close()
	tx, _ := tmpDbConn.Begin()

	tx.Exec("INSERT INTO scheduledTasks (raw) VALUES (?)", raw)

	tx.Commit()

	return nil
}

func (ssbb *SQLiteSpiderBusBackend) ReceiveScheduledTask() *ScheduledTask {
	tmpDbConn, _ := sql.Open("sqlite3", ssbb.sqliteFilePath+"?cache=shared&mode=rwc")
	defer tmpDbConn.Close()
	tx, _ := tmpDbConn.Begin()

	var row_id int
	var raw []byte

	row := tx.QueryRow("SELECT * FROM scheduledTasks ORDER BY id ASC LIMIT 1")

	err := row.Scan(&row_id, &raw)
	if err != nil {
		tx.Rollback()
		ssbb.maybePrintError(err)
		return nil
	}

	scheduledTask := &ScheduledTask{}

	ssbb.decodeEntry(raw, scheduledTask)

	tx.Exec(fmt.Sprintf("DELETE FROM scheduledTasks WHERE id=%d", row_id))
	tx.Commit()

	return scheduledTask
}

func (ssbb *SQLiteSpiderBusBackend) SendTaskPromise(taskPromise *TaskPromise) error {
	raw := ssbb.encodeEntry(taskPromise)

	tmpDbConn, _ := sql.Open("sqlite3", ssbb.sqliteFilePath+"?cache=shared&mode=rwc")
	defer tmpDbConn.Close()
	tx, _ := tmpDbConn.Begin()

	tx.Exec("INSERT INTO taskPromises (raw) VALUES (?)", raw)

	tx.Commit()

	return nil
}

func (ssbb *SQLiteSpiderBusBackend) ReceiveTaskPromise() *TaskPromise {
	tmpDbConn, _ := sql.Open("sqlite3", ssbb.sqliteFilePath+"?cache=shared&mode=rwc")
	defer tmpDbConn.Close()
	tx, _ := tmpDbConn.Begin()

	var row_id int
	var raw []byte

	row := tx.QueryRow("SELECT * FROM taskPromises ORDER BY id ASC LIMIT 1")

	err := row.Scan(&row_id, &raw)
	if err != nil {
		tx.Rollback()
		ssbb.maybePrintError(err)
		return nil
	}

	taskPromise := &TaskPromise{}

	ssbb.decodeEntry(raw, taskPromise)

	tx.Exec(fmt.Sprintf("DELETE FROM taskPromises WHERE id=%d", row_id))
	tx.Commit()

	return taskPromise
}

func (ssbb *SQLiteSpiderBusBackend) SendItem(item *Item) error {
	raw := ssbb.encodeEntry(item)

	tmpDbConn, _ := sql.Open("sqlite3", ssbb.sqliteFilePath+"?cache=shared&mode=rwc")
	defer tmpDbConn.Close()
	tx, _ := tmpDbConn.Begin()

	tx.Exec("INSERT INTO items (raw) VALUES (?)", raw)

	tx.Commit()

	return nil
}

func (ssbb *SQLiteSpiderBusBackend) ReceiveItem() *Item {
	tmpDbConn, _ := sql.Open("sqlite3", ssbb.sqliteFilePath+"?cache=shared&mode=rwc")
	defer tmpDbConn.Close()
	tx, _ := tmpDbConn.Begin()

	var row_id int
	var raw []byte

	row := tx.QueryRow("SELECT * FROM items ORDER BY id ASC LIMIT 1")

	err := row.Scan(&row_id, &raw)
	if err != nil {
		tx.Rollback()
		ssbb.maybePrintError(err)
		return nil
	}

	item := &Item{}

	ssbb.decodeEntry(raw, item)

	tx.Exec(fmt.Sprintf("DELETE FROM items WHERE id=%d", row_id))
	tx.Commit()

	return item
}

func (ssbb *SQLiteSpiderBusBackend) getCountForTable(tableName string) int {
	var output string

	query, _ := ssbb.dbConn.Prepare(fmt.Sprintf("SELECT COUNT(*) FROM %s;", tableName))

	defer query.Close()

	query.QueryRow().Scan(&output)

	count, _ := strconv.Atoi(output)

	return count
}

func (ssbb *SQLiteSpiderBusBackend) Close() {
	ssbb.dbConn.Close()
}
