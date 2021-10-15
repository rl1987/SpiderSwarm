package spsw

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func getCountForTable(sqliteFilePath string, tableName string) int {
	var output string

	dbConn, _ := sql.Open("sqlite3", sqliteFilePath+"?cache=shared&mode=rwc")
	defer dbConn.Close()

	query, _ := dbConn.Prepare(fmt.Sprintf("SELECT COUNT(*) FROM %s;", tableName))

	defer query.Close()

	query.QueryRow().Scan(&output)

	count, _ := strconv.Atoi(output)

	return count
}

func TestSQLiteSpiderBusBackendScheduledTaskE2E(t *testing.T) {
	taskPromise := &TaskPromise{UUID: "D412D565-B2A8-4BE3-B3CB-B37008FDA099"}
	taskTemplate := &TaskTemplate{TaskName: "testTask"}

	workflowName := "testWorkflow"
	workflowVersion := "2.0"
	jobUUID := "5369DD61-E98E-465E-9619-4641D06728FB"

	scheduledTask := NewScheduledTask(taskPromise, taskTemplate, workflowName, workflowVersion, jobUUID)

	assert.NotNil(t, scheduledTask)

	backend := NewSQLiteSpiderBusBackend("")
	defer func() {
		backend.Close()
		os.Remove(backend.sqliteFilePath)
	}()

	assert.Equal(t, 0, getCountForTable(backend.sqliteFilePath, "scheduledTasks"))
	assert.NotNil(t, backend)

	gotScheduledTask := backend.ReceiveScheduledTask()
	assert.Nil(t, gotScheduledTask)

	err := backend.SendScheduledTask(scheduledTask)
	assert.Nil(t, err)
	assert.Equal(t, 1, getCountForTable(backend.sqliteFilePath, "scheduledTasks"))

	gotScheduledTask2 := backend.ReceiveScheduledTask()
	assert.Equal(t, scheduledTask, gotScheduledTask2)

	gotScheduledTask = backend.ReceiveScheduledTask()
	assert.Nil(t, gotScheduledTask)
	assert.Equal(t, 0, getCountForTable(backend.sqliteFilePath, "scheduledTasks"))
}

func TestSQLiteSpiderBusBackendTaskPromiseE2E(t *testing.T) {
	taskPromise := &TaskPromise{UUID: "215B5E28-56AA-48DE-ADFB-8641E0547161"}

	backend := NewSQLiteSpiderBusBackend("")
	defer func() {
		backend.Close()
		os.Remove(backend.sqliteFilePath)
	}()

	assert.Equal(t, 0, getCountForTable(backend.sqliteFilePath, "taskPromises"))
	assert.NotNil(t, backend)

	gotTaskPromise := backend.ReceiveTaskPromise()
	assert.Nil(t, gotTaskPromise)

	err := backend.SendTaskPromise(taskPromise)
	assert.Nil(t, err)
	assert.Equal(t, 1, getCountForTable(backend.sqliteFilePath, "taskPromises"))

	gotTaskPromise2 := backend.ReceiveTaskPromise()
	assert.Equal(t, taskPromise, gotTaskPromise2)

	gotTaskPromise = backend.ReceiveTaskPromise()
	assert.Nil(t, gotTaskPromise)
	assert.Equal(t, 0, getCountForTable(backend.sqliteFilePath, "taskPromises"))
}

func TestSQLiteSpiderBusBackendItemE2E(t *testing.T) {
	item := &Item{UUID: "3350F665-F1EB-48A9-8FD8-704BDCCA4941"}

	backend := NewSQLiteSpiderBusBackend("")
	defer func() {
		backend.Close()
		os.Remove(backend.sqliteFilePath)
	}()

	assert.NotNil(t, backend)
	assert.Equal(t, 0, getCountForTable(backend.sqliteFilePath, "items"))

	gotItem := backend.ReceiveItem()
	assert.Nil(t, gotItem)

	err := backend.SendItem(item)
	assert.Nil(t, err)
	assert.Equal(t, 1, getCountForTable(backend.sqliteFilePath, "items"))

	gotItem = backend.ReceiveItem()
	assert.Equal(t, item, gotItem)

	gotItem = backend.ReceiveItem()
	assert.Nil(t, gotItem)
	assert.Equal(t, 0, getCountForTable(backend.sqliteFilePath, "items"))
}
