package goconcept

import (
	"errors"
)

type ConceptType struct {
	Type_name string	`json:"type"`
	Concept_data []ConceptData	`json:"concept_data"`
	Pathname string	`json:"pathname"`
	Api_available bool	`json:"api_available"`
}

type ConceptData struct {
	Type_name string	`json:"type"`
	Approve_func func(string) bool	`json:"-"`
	Single bool	`json:"single"`
}

type ConceptRelationshipType struct {
	Type1 string	`json:"type1"`
	Type2 string	`json:"type2"`
	String1 string	`json:"string1"`
	String2 string	`json:"string2"`
}

func NewConceptType(name string, concept_data []ConceptData, api_available bool, pathname string) (*ConceptType, error) {
	var real_pathname string
	if api_available {
		if pathname != "" {
			return nil, errors.New("invalid pathname")
		}
		real_pathname = pathname
	} else {
		real_pathname = ""
	}

	if name == "" {
		return nil, errors.New("invalid name")
	}

	used_names := make(map[string]bool)
	for _, cd := range concept_data {
		_, ok := used_names[cd.Type_name]
		if ok {
			return nil, errors.New("name used twice")
		}
		used_names[cd.Type_name] = true
	}

	return &ConceptType{name, concept_data, real_pathname, api_available}, nil
}

func NewConceptData(name string, approve_func *func(str string) bool, single bool) (*ConceptData, error) {
	if name == "" {
		return nil, errors.New("invalid name")
	}

	var real_approve_func func(str string) bool
	if approve_func == nil {
		real_approve_func = func(str string) bool {
			return true
		}
	} else {
		real_approve_func = *approve_func
	}

	return &ConceptData{name, real_approve_func, single}, nil
}

func NewConceptRelationshipType(type1 string, type2 string, string1 string, string2 string) (*ConceptRelationshipType, error) {
	if type1 == "" || type2 == "" || string1 == "" || string2 == "" {
		return nil, errors.New("invalid concept relationship type")
	}

	return &ConceptRelationshipType{type1, type2, string1, string2}, nil
}
