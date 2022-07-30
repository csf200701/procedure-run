package connector

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnector struct {
	*Collection
}

func (c MysqlConnector) ValidateCollection() error {
	collectionStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.User, c.Password, c.Host, c.Port, c.DbName)
	info, err := validateCollection(collectionStr)
	if err != nil {
		return errors.New(info + "，错误：" + err.Error())
	}
	return nil
}

func (c MysqlConnector) Query(query string, args ...interface{}) (*sql.Rows, error) {
	collectionStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.User, c.Password, c.Host, c.Port, c.DbName)
	info, err := validateCollection(collectionStr)
	if err != nil {
		return nil, errors.New(info + "，错误：" + err.Error())
	}
	db, err := sql.Open("mysql", collectionStr)
	if err != nil {
		return nil, errors.New("该数据库连接失败，错误：" + err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return nil, errors.New("该数据库连接失败，错误：" + err.Error())
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	//defer rows.Close()
	return rows, nil
}

func (c MysqlConnector) ConnectJoin() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.User, c.Password, c.Host, c.Port, c.DbName)
}

func validateCollection(c string) (string, error) {
	db, err := sql.Open("mysql", c)
	if err != nil {
		return "该数据库连接失败", err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return "该数据库连接失败", err
	}
	return "", nil
}
