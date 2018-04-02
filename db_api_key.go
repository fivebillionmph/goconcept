package goconcept

import (
	"crypto/rand"
	"time"
	"encoding/base64"
)

const DBAPIKey__table string = "base_api_keys"
const DBAPIKey__keylen int = 32
const DBAPIKey__user_max int = 2
const DBAPIKey__header_name string = "X-api-key"

type DBAPIKey struct {
	F_id int	`json:"-"`
	F_user_id int	`json:"-"`
	F_timestamp int	`json:"-"`
	F_active int	`json:"active"`
	F_key string	`json:"key"`
}

func DBAPIKey__create(cxn *Connection, user *DBUser) (*DBAPIKey, error) {
	crypto_bytes := make([]byte, DBAPIKey__keylen)
	_, err := rand.Read(crypto_bytes)
	if err != nil {
		return nil, err
	}
	crypto_string := base64.URLEncoding.EncodeToString(crypto_bytes)
	new_key := crypto_string[0:DBAPIKey__keylen]

	timestamp := int(time.Now().Unix())

	stmt, err := cxn.DB.Prepare("insert into " + DBAPIKey__table + " values(NULL, ?, ?, 1, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.F_id, timestamp, new_key)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return DBAPIKey__getByID(cxn, int(id))
}

func (d *DBAPIKey) readRow(row SQLRowInterface) error {
	err := row.Scan(
		&d.F_id,
		&d.F_user_id,
		&d.F_timestamp,
		&d.F_active,
		&d.F_key,
	)
	return err
}

func DBAPIKey__getByID(cxn *Connection, id int) (*DBAPIKey, error) {
	row := cxn.DB.QueryRow("select * from " + DBAPIKey__table + " where id = ?", id)
	key := DBAPIKey{}
	err := key.readRow(row)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func DBAPIKey__getByKey(cxn *Connection, key_str string) (*DBAPIKey, error) {
	row := cxn.DB.QueryRow("select * from " + DBAPIKey__table + " where BINARY key = ?", key_str)
	key := DBAPIKey{}
	err := key.readRow(row)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func DBAPIKey__getByUserID(cxn *Connection, user_id int) (*[]DBAPIKey, error) {
	rows, err := cxn.DB.Query("select * from " + DBAPIKey__table + " where user_id = ?", user_id)
	if err != nil {
		return nil, err
	}

	var keys []DBAPIKey;
	for rows.Next() {
		key := DBAPIKey{}
		err := key.readRow(rows)
		if err == nil {
			keys = append(keys, key)
		}
	}

	return &keys, nil
}

func DBAPIKey__getCountByUserID(cxn *Connection, user_id int, active_only bool) (int, error) {
	var active_only_str string
	if active_only {
		active_only_str = " and active = 1"
	} else {
		active_only_str = ""
	}
	row := cxn.DB.QueryRow("select count(*) from " + DBAPIKey__table + " where user_id = ?" + active_only_str, user_id)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *DBAPIKey) deactive(cxn *Connection) error {
	stmt, err := cxn.DB.Prepare("update " + DBAPIKey__table + " set active = 0 where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(d.F_id)
	if err != nil {
		return err
	}
	d.F_active = 0

	return nil
}
