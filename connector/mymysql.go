package connector

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type MysqlConnector struct {
	*Collection
}

func (c MysqlConnector) ValidateCollection() error {
	db, err := c.fetchDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (c MysqlConnector) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := c.fetchDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return rows, nil
}

func (c MysqlConnector) ConnectJoin() string {
	var collectionStr string
	if c.IsSSH {
		client, _ := DialWithPasswd(fmt.Sprintf("%v:%v", c.SSHHost, c.SSHPort), c.SSHUser, c.SSHPassword)

		procfStr := md5Str(fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.User, c.Password, c.Host, c.Port, c.DbName)+"|"+
			fmt.Sprintf("%v:%v@tcp(%v:%v)", c.SSHHost, c.SSHPort, c.SSHUser, c.SSHPassword)) + "+ssh"
		// 注册ssh代理
		mysql.RegisterDialContext(procfStr, (&ViaSSHDialer{client.client, nil}).Dial)
		collectionStr = fmt.Sprintf("%v:%v@%v(%v:%v)/%v", c.User, c.Password, procfStr, c.Host, c.Port, c.DbName)
	} else {
		collectionStr = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.User, c.Password, c.Host, c.Port, c.DbName)
	}
	return collectionStr
}

func (c MysqlConnector) fetchDB() (*sql.DB, error) {
	collectionStr := c.ConnectJoin()
	db, err := sql.Open("mysql", collectionStr)
	if err != nil {
		return nil, errors.New("该数据库连接失败，错误：" + err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.New("该数据库连接失败，错误：" + err.Error())
	}
	return db, nil
}

func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
