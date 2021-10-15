package spsw

import ()

type MySQLSpiderBusBackend struct {
	SpiderBusBackend

	dbConn *sql.DB
}
