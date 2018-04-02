package goconcept

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type SQLRowInterface interface {
	Scan(dest ...interface{}) error
}

type Connection struct {
	DB *sql.DB
}

func newConnection(host string, user string, password string, db string) (*Connection, error) {
	conn, err := sql.Open("mysql", user + ":" + password + "@tcp(" + host + ":3306)/" + db + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil
}

func (c *Connection) Close() {
	c.DB.Close()
}
