package goconcept

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type SQLRowInterface interface {
	Scan(dest ...interface{}) error
}

type Connection struct {
	db *sql.db
}

func newConnection(host string, user string, password string, db string) (*Connection, error) {
	db, err := sql.Open("mysql", user + ":" + password + "@tcp(" + host + ":3306)/" + db + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return &Connection{db}
}

func (c *Connection) close() {
	c.db.close()
}
