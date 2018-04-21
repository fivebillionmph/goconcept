package goconcept

import (
	"errors"
	"time"
)

const DBConcept__table string = "base_concepts"

type DBConcept struct {
	F_id int	`json:"-"`
	F_timestamp int `json:"-"`
	F_active int8	`json:"-"`
	F_type string	`json:"type"`
	F_name string	`json:"name"`

	Data *[]DBConceptData	`json:"data"`
	Relationships *[]DBConcept__Relationship `json:"relationships"`
}

type DBConcept__Relationship struct {
	Concept *DBConcept	`json:"item"`
	String1 string `json:"string1"`
	String2 string `json:"string2"`
	Reverse bool `json:"reverse"`
}

/* create */

func DBConcept__create(cxn *Connection, type_name string, name string) (*DBConcept, error) {
	row := cxn.DB.QueryRow("select count(*) from " + DBConcept__table + " where type = ? and name = ?", type_name, name)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("already exists")
	}

	timestamp := time.Now().Unix()
	stmt, err := cxn.DB.Prepare("insert into " + DBConcept__table + " values(NULL, ?, 1, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(timestamp, type_name, name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return DBConcept__getByID(cxn, int(id))
}

/* read */

func (d *DBConcept) readRow(row SQLRowInterface) error {
	err := row.Scan(
		&d.F_id,
		&d.F_timestamp,
		&d.F_active,
		&d.F_type,
		&d.F_name,
	)
	return err
}

func DBConcept__getByID(cxn *Connection, id int) (*DBConcept, error) {
	row := cxn.DB.QueryRow("select * from " + DBConcept__table + " where id=?", id)

	concept := DBConcept{}
	err := concept.readRow(row)
	if err != nil {
		return nil, errors.New("could not find concept")
	}

	concept.LoadData(cxn)
	return &concept, nil
}

func DBConcept__getAll(cxn *Connection, offset int, count int) (*[]DBConcept, error) {
	rows, err := cxn.DB.Query("select * from " + DBConcept__table + " limit ?, ?", offset, count)
	if err != nil {
		return nil, err
	}

	var concepts []DBConcept
	for rows.Next() {
		concept := DBConcept{}
		err := concept.readRow(rows)
		if err == nil {
			concept.LoadData(cxn)
			concepts = append(concepts, concept)
		}
	}

	return &concepts, nil
}

func DBConcept__getByTypeName(cxn *Connection, type_name string, name string) (*DBConcept, error) {
	row := cxn.DB.QueryRow("select * from " + DBConcept__table + " where type=? and name=?", type_name, name)

	concept := DBConcept{}
	err := concept.readRow(row)
	if err != nil {
		return nil, err
	}

	concept.LoadData(cxn)
	return &concept, nil
}

func DBConcept__getByType(cxn *Connection, type_name string, offset int, count int) (*[]DBConcept, error) {
	rows, err := cxn.DB.Query("select * from " + DBConcept__table + " where type=? limit ?, ?", type_name, offset, count)
	if err != nil {
		return nil, err
	}

	var concepts []DBConcept
	for rows.Next() {
		concept := DBConcept{}
		err := concept.readRow(rows)
		if err == nil {
			concept.LoadData(cxn)
			concepts = append(concepts, concept)
		}
	}

	return &concepts, nil
}

func DBConcept__getBySearchName(cxn *Connection, type_name string, name_search string, offset int, count int) (*[]DBConcept, error) {
	rows, err := cxn.DB.Query("select * from " + DBConcept__table + " where type = ? and name like ? limit ?, ?", type_name, "%" + name_search + "%", offset, count)
	if err != nil {
		return nil, err
	}

	var concepts []DBConcept
	for rows.Next() {
		concept := DBConcept{}
		err := concept.readRow(rows)
		if err == nil {
			concept.LoadData(cxn)
			concepts = append(concepts, concept)
		}
	}

	return &concepts, nil
}

/* update */

/* delete */

func DBConcept__delete(cxn *Connection, concept *DBConcept) error {
	if concept == nil {
		return errors.New("nil concept")
	}

	var err error

	concept.LoadData(cxn)
	for _, d := range *concept.Data {
		err = DBConceptData__delete(cxn, &d)
		if err != nil {
			return err
		}
	}

	relationships, err := DBConceptRelationship__getByConceptID(cxn, concept.F_id)
	if err != nil {
		return err
	}
	for _, r := range *relationships {
		err = DBConceptRelationship__delete(cxn, &r)
		if err != nil {
			return err
		}
	}

	stmt, err := cxn.DB.Prepare("delete from " + DBConcept__table + " where id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(concept.F_id)
	if err != nil {
		return err
	}

	concept.F_id = 0

	return nil
}

func (d *DBConcept) LoadData(cxn *Connection) {
	if d.Data != nil {
		return
	}

	data, err := DBConceptData__getByConceptID(cxn, d.F_id)
	if err != nil {
		return
	}

	d.Data = data
}

func DBConcept__getCountByType(cxn *Connection, type_name string) (int, error ) {
	row := cxn.DB.QueryRow("select count(*) from " + DBConcept__table + " where type=?", type_name)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *DBConcept) LoadRelationships(cxn *Connection) {
	if d.Relationships != nil {
		return
	}

	relationships, err := DBConceptRelationship__getByConceptID(cxn, d.F_id)
	if err != nil {
		return
	}

	var final_relationships []DBConcept__Relationship
	for _, rel := range *relationships {
		var other_concept *DBConcept
		var reverse bool
		if rel.F_id1 == d.F_id {
			rel.LoadConcept(cxn, 2)
			other_concept = rel.concept2
			reverse = false
		} else {
			rel.LoadConcept(cxn, 1)
			other_concept = rel.concept1
			reverse = true
		}
		if other_concept == nil {
			continue
		}
		this_relationship := DBConcept__Relationship{other_concept, rel.F_string1, rel.F_string2, reverse}
		final_relationships = append(final_relationships, this_relationship)
	}

	d.Relationships = &final_relationships
}
