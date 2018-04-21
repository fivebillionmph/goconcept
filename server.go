package goconcept

import (
	"errors"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"encoding/gob"
	"log"
	"os"
)

type Server struct {
	http_server *http.ServeMux
	Logger *log.Logger
	cookie_wrapper *CookieWrapper
	concept_types map[string]*ConceptType
	concept_relationship_types []*ConceptRelationshipType
	connection *Connection
	router *mux.Router
}

func NewServer(cookie_key string, dbhost string, dbuser string, dbpassword string, dbdatabase string) (*Server, error) {
	gob.Register(&SessionUser{})

	logger := log.New(os.Stdout, "--- ", log.Ldate | log.Ltime | log.Lshortfile)
	http_server := http.NewServeMux()
	router := mux.NewRouter()
	cookie_wrapper := &CookieWrapper{sessions.NewCookieStore([]byte(cookie_key))}
	concept_types := make(map[string]*ConceptType)
	concept_relationship_types := []*ConceptRelationshipType{}
	connection, err := newConnection(dbhost, dbuser, dbpassword, dbdatabase)
	if err != nil {
		return nil, err
	}
	return &Server{http_server, logger, cookie_wrapper, concept_types, concept_relationship_types, connection, router}, nil
}

func (s *Server) Start(static_path string, static_dir string, html_file string, admin_path string, allow_user_create bool) {
	s.addAdminRoutes()
	s.addUserRoutes(allow_user_create)
	s.AddStaticRouterPath(static_path, static_dir)
	s.AddStaticRouterPath(admin_path, "./goconcept-files/admin-frontend")
	s.AddHTMLRouterPath("/", html_file)
	s.http_server.Handle("/", s.router)

	go func() {
		s.Logger.Println("starting server")
		_ = http.ListenAndServe(":8080", s.http_server)
	}()

	serverCommands(s)
}

func (s *Server) AddConceptType(concept_type ConceptType) error {
	/* validate concept type data */
	existing_data_types := make(map[string]bool)
	for _, cd := range concept_type.Concept_data {
		_, exists := existing_data_types[cd.Type_name]
		if exists {
			return errors.New("duplicate data type name: " + concept_type.Type_name + " - " + cd.Type_name)
		}
		existing_data_types[cd.Type_name] = true
	}

	/* validate not duplicate concept type */
	for _, ct := range s.concept_types {
		if ct.Type_name == concept_type.Type_name {
			return errors.New("duplicate concept name: " + ct.Type_name)
		}
		if ct.Pathname != "" && ct.Pathname == concept_type.Pathname {
			return errors.New("duplicate pathname: " + ct.Pathname)
		}
	}
	s.concept_types[concept_type.Type_name] = &concept_type

	/* add api routes */
	if concept_type.Api_available {
		single_endpoint_handler := func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			vars := mux.Vars(r)
			name := vars["name"]
			concept, err := DBConcept__getByTypeName(cxn, concept_type.Type_name, name)
			if err != nil {
				s.Logger.Println(err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			concept.LoadRelationships(s.connection)

			s.SendJSONResponse(w, r, concept)
		}
		s.AddRouterPath("/api/v1/c/" + concept_type.Pathname + "/{name}", "GET", single_endpoint_handler)

		many_endpoint_handler := func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			query_vars := r.URL.Query()
			count := Util__queryToInt(query_vars, "count", 0, 20, false, false, 20)
			offset := Util__queryToInt(query_vars, "offset", 0, 0, false, true, 0)
			concepts, err := DBConcept__getByType(cxn, concept_type.Type_name, offset, count)
			if err != nil {
				s.Logger.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			s.SendJSONResponse(w, r, concepts)
		}
		s.AddRouterPath("/api/v1/c/" + concept_type.Pathname, "GET", many_endpoint_handler)
	}

	return nil
}

func (s *Server) AddConceptRelationshipType(concept_relationship_type ConceptRelationshipType) error {
	for _, crt := range s.concept_relationship_types {
		if crt.Type1 == concept_relationship_type.Type1 &&
			crt.Type2 == concept_relationship_type.Type2 &&
			crt.String1 == concept_relationship_type.String1 &&
			crt.String2 == concept_relationship_type.String2 {
			return errors.New("duplicate concept relationship type")
		}
	}
	type1_found := false
	type2_found := false
	for _, ct := range s.concept_types {
		if ct.Type_name == concept_relationship_type.Type1 {
			type1_found = true
		}
		if ct.Type_name == concept_relationship_type.Type2 {
			type2_found = true
		}
		if type1_found && type2_found {
			break
		}
	}
	if !type1_found || !type2_found {
		return errors.New("invalid relationship type")
	}
	s.concept_relationship_types = append(s.concept_relationship_types, &concept_relationship_type)
	return nil
}
