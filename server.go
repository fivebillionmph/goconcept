package goconcept

import (
	"errors"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"encoding/gob"
	"log"
)

type Server struct {
	http_server *http.ServerMux
	Logger *log.Logger
	cookie_store *CookieWrapper
	concept_types map[string]*ConceptType
	concept_relationship_types []*ConceptRelationshipType
	connection *Connection
	API_handler *APIHandler
}

func NewServer(cookie_key string) (*Server, error) {
	logger = log.New(os.Stdout, "--- ", log.Ldate | log.Ltime | log.Lshortfile)
	http_server := http.NewServerMux()
	router := mux.NewRouter()
	cookie_store := &CookieWrapper{sessions.NewCookieStore([]byte(cookie_key))}
	concept_types := make(map[string]*ConceptType)
	concept_relationship_types := []*ConceptRelationshipType{}
	connection, err := newConnection(dbhost, dbuser, dbpassword, dbdatabase)
	if err != nil {
		return nil, err
	}
	api_handler := &apiHandler{router, cxn, cookie_store}
	return &Server{http_server, logger, cookie_store, concept_types, concept_relationship_types, connection, api_handler}
}

func (s *Server) Start() error {
	s.Logger.Println("starting server")
}

func (s *Server) AddConceptType(concept_type *ConceptType) error {
	if concept_type == nil {
		return errors.New("nil concept type")
	}
	for _, ct := range s.concept_types {
		if ct.type_name == concept_type.type_name {
			return errors.New("duplicate concept name")
		}
		if ct.pathname != "" && ct.pathname == concept_type.pathname {
			return errors.New("duplicate pathname")
		}
	}
	s.concept_types[concept_type.type_name] = concept_type

	if concept_type.api_available {
		s.API_handler.AddPath("/api/v1/c/" + concept_type.pathname + "/{concept}")
	}

	return nil
}

func (s *Server) AddConceptRelationshipType(concept_relationship_type *ConceptRelationshipType) error {
	if concept_relationship_type == nil {
		return errors.New("nil concept relationship type")
	}
	for _, crt := range s.concept_relationship_types {
		if crt.type1 == concept_relationship_type.type1 &&
			crt.type2 == concept_relationship_type.type2 &&
			crt.string1 == concept_relationship_type.string1 &&
			crt.string2 == concept_relationship_type.string2 {
			return errors.New("duplicate concept relationship type")
		}
	}
	type1_found := false
	type2_found := false
	for _, ct := range s.concept_types {
		if ct.type_name == concept_relationship_type.type1 {
			type1 = true
		} else if ct.type_name == concept_relationship_type.type2 {
			type2 = true
		}
		if type1 && type2 {
			break
		}
	}
	if !type1 || !type2 {
		return errors.New("invalid relationship type")
	}
	s.concept_relationship_types = append(s.concept_relationship_types, concept_relationship_type)
	return nil
}
