package goconcept

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

func NewConceptType(name string, pathname string, api_available bool, concept_data []ConceptData) ConceptType {
	var real_pathname string
	if api_available {
		if pathname == "" {
			real_pathname = name
		}
		real_pathname = pathname
	} else {
		real_pathname = ""
	}

	return ConceptType{name, concept_data, real_pathname, api_available}
}

func NewConceptData(name string, approve_func *func(str string) bool, single bool) ConceptData {
	var real_approve_func func(str string) bool
	if approve_func == nil {
		real_approve_func = func(str string) bool {
			return true
		}
	} else {
		real_approve_func = *approve_func
	}

	return ConceptData{name, real_approve_func, single}
}

func NewConceptRelationshipType(type1 string, type2 string, string1 string, string2 string) ConceptRelationshipType {
	return ConceptRelationshipType{type1, type2, string1, string2}
}
