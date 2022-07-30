package connector

import (
	"database/sql"
)

const (
	MYSQL      string = "Mysql"
	POSTGRESQL string = "PostgreSQL"
	ORACLE     string = "Oracle"
)

type Collection struct {
	Host     string `survey:"host"`
	Port     string `survey:"port"`
	User     string `survey:"user"`
	Password string `survey:"password"`
	DbName   string `survey:"db"`
	DbType   string `survey:"type"`
	IsSSH    bool
}

type Connector interface {
	ValidateCollection() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	ConnectJoin() string
}

func Database(c *Collection) Connector {
	var connector Connector
	if c.DbType == MYSQL {
		connector = Mysql(c)
	}

	return connector
}

func Mysql(c *Collection) *MysqlConnector {
	return &MysqlConnector{Collection: c}
}