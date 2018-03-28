package goconcept

import (
	"errors"
)

type ConceptType struct {
	type_name string
	concept_data []ConceptData
	pathname string
	api_available bool
}

type ConceptData struct {
	type_name string
	approve_func func(str string) bool
	single bool
}

type ConceptRelationshipType struct {
	type1 string
	type2 string
	string1 string
	string2 string
}

func NewConceptType(name string, concept_data []ConceptData, api_available bool, pathname string) (*ConceptType, error) {
	var real_pathname string
	if api_available {
		if pathname != "" {
			return (nil, errors.New("invalid pathname"))
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
		_, ok := used_names[cd.type_name]
		if ok {
			return (nil, errors.New("name used twice"))
		}
		used_names[cd.type_name] = true
	}

	return &ConceptType{name, concept_data, api_available, pathname}, nil
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
		return (nil, errors.New("invalid concept relationship type"))
	}

	return &ConceptRelationshipType{type1, type2, string1, string2}, nil
}
