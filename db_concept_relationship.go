package goconcept

import (
	"errors"
	"time"
)

const DBConceptRelationship__table string = "base_concept_relationships"

type DBConceptRelationship struct {
	F_id int
	F_timestamp int
	F_id1 int
	F_id2 int
	F_string1 string
	F_string2 string

	concept1 *DBConcept
	concept2 *DBConcept
}

func DBConceptRelationship__create(cxn *Connection, id1 int, id2 int, string1 string, string2 string) (*DBConceptRelationship, error) {
	row := cxn.DB.QueryRow("select count(*) from " + DBConceptRelationship__table + " where id1 = ? and id2 = ? and string1 = ? and string2 = ?", id1, id2, string1, string2)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("already exists")
	}

	timestamp := time.Now().Unix()
	stmt, err := cxn.DB.Prepare("insert into " + DBConceptRelationship__table + " values(NULL, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(timestamp, id1, id2, string1, string2)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return DBConceptRelationship__getByID(cxn, int(id))
}

func (d *DBConceptRelationship) readRow(row SQLRowInterface) error {
	err := row.Scan(
		&d.F_id,
		&d.F_timestamp,
		&d.F_id1,
		&d.F_id2,
		&d.F_string1,
		&d.F_string2,
	)
	return err
}

func DBConceptRelationship__getByID(cxn *Connection, id int) (*DBConceptRelationship, error) {
	row := cxn.DB.QueryRow("select * from " + DBConceptRelationship__table + " where id = ?", id)
	rel := DBConceptRelationship{}
	err := rel.readRow(row)
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func DBConceptRelationship__getByConceptID(cxn *Connection, id int) (*[]DBConceptRelationship, error) {
	rows, err := cxn.DB.Query("select * from " + DBConceptRelationship__table + " where id1 = ? or id2 = ?", id, id)
	if err != nil {
		return nil, err
	}

	var rels []DBConceptRelationship
	for rows.Next() {
		rel := DBConceptRelationship{}
		err := rel.readRow(rows)
		if err == nil {
			rels = append(rels, rel)
		}
	}

	return &rels, nil
}

func DBConceptRelationship__getByIDsStrings(cxn *Connection, id1 int, id2 int, string1 string, string2 string) (*DBConceptRelationship, error) {
	row := cxn.DB.QueryRow("select * from " + DBConceptRelationship__table + " where id1 = ? and id2 = ? and string1 = ? and string2 = ?", id1, id2, string1, string2)
	rel := DBConceptRelationship{}
	err := rel.readRow(row)
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func (d *DBConceptRelationship) LoadConcept(cxn *Connection, concept_id int) {
	if concept_id != 1 && concept_id != 2 {
		return
	}

	var id int
	if concept_id == 1 {
		if d.concept1 != nil {
			return
		}
		id = d.F_id1
	} else {
		if d.concept1 != nil {
			return
		}
		id = d.F_id2
	}

	concept, err := DBConcept__getByID(cxn, id)
	if err != nil {
		return
	}

	if concept_id == 1 {
		d.concept1 = concept
	} else {
		d.concept2 = concept
	}
}

func DBConceptRelationship__delete(cxn *Connection, relationship *DBConceptRelationship) error {
	if relationship == nil {
		return errors.New("nil relationship")
	}

	stmt, err := cxn.DB.Prepare("delete from " + DBConceptRelationship__table + " where id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(relationship.F_id)
	if err != nil {
		return err
	}

	relationship.F_id = 0

	return nil
}
