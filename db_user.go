package goconcept

import (
	"time"
	"github.com/asaskevich/govalidator"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"unicode/utf8"
)

var DBUser__table string = "base_users"

type DBUser struct {
	F_id int `json:"-"`
	F_timestamp int `json:"-"`
	F_email string `json:"-"`
	F_password string `json:"-"`
	F_username string `json:"username"`
	F_level int8 `json:"level"`
	F_active int8 `json:"-"`
}

/* create */
func DBUser__create(cxn *Connection, email string, password_plaintext string, username string, level uint8) (*DBUser, error) {
	time := int(time.Now().Unix())

	if !govalidator.IsEmail(email) {
		return nil, errors.New("invalid email address")
	}

	if utf8.RuneCountInString(password_plaintext) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	password_ba, err := bcrypt.GenerateFromPassword([]byte(password_plaintext), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("invalid password")
	}
	password := string(password_ba)

	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 16 {
		return nil, errors.New("username must be between 3 and 16 characters")
	}

	if level != 1 && level != 2 {
		level = 1
	}

	active := 1

	stmt, err := cxn.DB.Prepare("INSERT INTO " +  DBUser_table + " values(NULL, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, errors.New("db error")
	}
	defer stmt.Close()

	res, err := stmt.Exec(time, email, password, username, level, active)
	if err != nil {
		return nil, errors.New("db error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.New("db error")
	}

	return DBUser__getByID(cxn, int(id))
}

/* read */
func DBUser__getByID(cxn *Connection, id int) (*DBUser, error) {
	row := cxn.DB.QueryRow("select * from " + DBUser_table + " where id=?", id)

	user := DBUser{}
	err := row.Scan(
		&user.F_id,
		&user.F_timestamp,
		&user.F_email,
		&user.F_password,
		&user.F_username,
		&user.F_level,
		&user.F_active,
	)
	if err != nil {
		return nil, errors.New("could not find user")
	}

	return &user, nil
}

func DBUser__getByPasswordChallenge(cxn *Connection, email string, password_plaintext string) (*DBUser, error) {
	row := cxn.DB.QueryRow("select id, password from " + DBUser__table + " where email=?", email)

	var id int
	var password string
	err := row.Scan(&id, &password)
	if err != nil {
		return nil, err
	}

	password_plaintext_ba := []byte(password_plaintext)
	password_ba := []byte(password)
	err = bcrypt.CompareHashAndPassword(password_ba, password_plaintext_ba)
	if err != nil {
		return nil, err
	}

	return DBUser__getByID(cxn, id)
}

/* update */

/* delete */


func (u DBUser) level() int8 {
	return u.F_level
}
