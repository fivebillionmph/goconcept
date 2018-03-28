package goconcept

import (
	"errors"
)

var DBConceptData__table string = "base_concept_data"

type DBConceptData struct {
	F_id int	`json:"-"`
	F_timestamp int	`json:"-"`
	F_concept_id int	`json:"-"`
	F_active int8	`json:"-"`
	F_key string	`json:"key"`
	F_value string	`json:"value"`
}

func (d *DBConceptData) readRow(row sqlRowInterface) error {
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

	stmt, err := cxn.Prepare("delete from " + DBConceptData__table + " where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(concept_data.F_id)
	if err != nil {
		return err
	}

	concept_data.F_id = 0

	return nil
}
