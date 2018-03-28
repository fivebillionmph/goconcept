package goconcept

import (
	"errors"
)

var DBConcept__table string = "base_concepts"

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
	Reltype string `json:"reltype"`
}

const DBConcept__TYPE_PAGE string = "page"
const DBConcept__TYPE_PROGRAMMING_LANGUAGE string = "programming-language"
const DBConcept__TYPE_TOOL string = "tool"
const DBConcept__TYPE_ENGINE string = "engine"

/* create */

/* read */

func (d *DBConcept) readRow(row sqlRowInterface) error {
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

	concept.loadData(cxn)
	return &concept, nil
}

func DBConcept__getByTypeName(cxn *Connection, type_name string, name string) (*DBConcept, error) {
	row := cxn.DB.QueryRow("select * from " + DBConcept__table + " where type=? and name=?", type_name, name)

	concept := DBConcept{}
	err := concept.readRow(row)
	if err != nil {
		return nil, err
	}

	concept.loadData(cxn)
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
			concept.loadData(cxn)
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

	concept.loadData(cxn)
	for _, d := range concept.Data {
		err = DBConceptData__delete(d)
		if err != nil {
			return err
		}
	}

	concept.loadRelationships(cxn)
	for _, r := range concept.Relationships {
		err = DBConceptRelationship__delete(r)
		if err != nil {
			return err
		}
	}

	stmt, err := cxn.Prepare("delete from " + DBConcept__table + " where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(concept.F_id)
	if err != nil {
		return err
	}

	concept.F_id = 0

	return nil
}

func (d *DBConcept) loadData(cxn *Connection) {
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

func (d *DBConcept) loadRelationships(cxn *Connection) {
	if d.Relationships != nil {
		return
	}

	relationships, err := DBConceptRelationship__getByConceptID(cxn, d.F_id)
	if err != nil {
		logger.Println(err)
		return
	}

	var final_relationships []DBConcept__Relationship
	for _, rel := range *relationships {
		var other_concept *DBConcept
		var reltype string
		if rel.F_id1 == d.F_id {
			rel.loadConcept(cxn, 2)
			other_concept = rel.concept2
			reltype = rel.F_string1
		} else {
			rel.loadConcept(cxn, 1)
			other_concept = rel.concept1
			reltype = rel.F_string2
		}
		if other_concept == nil {
			continue
		}
		this_relationship := DBConcept__Relationship{other_concept, reltype}
		final_relationships = append(final_relationships, this_relationship)
	}

	d.Relationships = &final_relationships
}
