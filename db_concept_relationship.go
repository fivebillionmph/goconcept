package goconcept

var DBConceptRelationship__table string = "base_concept_relationships"

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

func (d *DBConceptRelationship) readRow(row sqlRowInterface) error {
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

func (d *DBConceptRelationship) loadConcept(cxn *Connection, concept_id int) {
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
		logger.Println(err)
		return
	}

	if concept_id == 1 {
		d.concept1 = concept
	} else {
		d.concept2 = concept
	}
}

func DBConceptRelationship__delete(cxn *Connection, relationship *DBConceptRelationship) error {
	if relationships == nil {
		return errors.New("nil relationship")
	}

	stmt, err := cxn.Prepare("delete from " + DBConceptRelationship__table + " where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(relationship.F_id)
	if err != nil {
		return err
	}

	relationship.F_id = 0

	return nil
}
