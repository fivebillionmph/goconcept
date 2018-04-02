package goconcept

import (
	"errors"
	"time"
)

const DBConceptData__table string = "base_concept_data"

type DBConceptData struct {
	F_id int	`json:"-"`
	F_timestamp int	`json:"-"`
	F_concept_id int	`json:"-"`
	F_active int8	`json:"-"`
	F_key string	`json:"key"`
	F_value string	`json:"value"`
}

func DBConceptData__create(cxn *Connection, concept_id int, key string, val string) (*DBConceptData, error) {
	timestamp := time.Now().Unix()
	stmt, err := cxn.DB.Prepare("insert into " + DBConceptData__table + " values(NULL, ?, ?, 1, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(timestamp, concept_id, key, val)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return DBConceptData__getByID(cxn, int(id))
}

func (d *DBConceptData) readRow(row SQLRowInterface) error {
	err := row.Scan(
		&d.F_id,
		&d.F_timestamp,
		&d.F_concept_id,
		&d.F_active,
		&d.F_key,
		&d.F_value,
	)
	return err
}

func DBConceptData__getByID(cxn *Connection, id int) (*DBConceptData, error) {
	row := cxn.DB.QueryRow("select * from " + DBConceptData__table + " where id = ?", id)

	concept_data := DBConceptData{}
	err := concept_data.readRow(row)
	if err != nil {
		return nil, errors.New("could not find concept data")
	}

	return &concept_data, nil
}

func DBConceptData__getByConceptID(cxn *Connection, concept_id int) (*[]DBConceptData, error) {
	rows, err := cxn.DB.Query("select * from " + DBConceptData__table + " where concept_id = ?", concept_id)
	if err != nil {
		return nil, err
	}

	var data []DBConceptData;
	for rows.Next() {
		datum := DBConceptData{}
		err := datum.readRow(rows)
		if err == nil {
			data = append(data, datum)
		}
	}

	return &data, nil
}

func DBConceptData__delete(cxn *Connection, concept_data *DBConceptData) error {
	if concept_data == nil {
		return errors.New("nil concept_data")
	}

	stmt, err := cxn.DB.Prepare("delete from " + DBConceptData__table + " where id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(concept_data.F_id)
	if err != nil {
		return err
	}

	concept_data.F_id = 0

	return nil
}
