package goconcept

import (
	"errors"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"encoding/gob"
	"log"
	"strconv"
)

type Server struct {
	http_server *http.ServerMux
	Logger *log.Logger
	cookie_wrapper *CookieWrapper
	concept_types map[string]*ConceptType
	concept_relationship_types []*ConceptRelationshipType
	connection *Connection
	router *mux.Router
}

func NewServer(cookie_key string) (*Server, error) {
	gob.Register(&SessionUser{})

	logger = log.New(os.Stdout, "--- ", log.Ldate | log.Ltime | log.Lshortfile)
	http_server := http.NewServerMux()
	router := mux.NewRouter()
	cookie_wrapper := &CookieWrapper{sessions.NewCookieStore([]byte(cookie_key))}
	concept_types := make(map[string]*ConceptType)
	concept_relationship_types := []*ConceptRelationshipType{}
	connection, err := newConnection(dbhost, dbuser, dbpassword, dbdatabase)
	if err != nil {
		return nil, err
	}
	return &Server{http_server, logger, cookie_wrapper, concept_types, concept_relationship_types, connection, router}
}

fun (s *Server) Run(static_path string, static_dir string, html_file string, admin_path string, admin_html_file string, allow_user_create bool) error {
	s.addAdminRoutes()
	s.addUserRoutes(allow_user_create)
	s.addStaticRouterPath(static_path, static_dir)
	s.addHTMLRouterPath(admin_path, admin_html_file)
	s.addHTMLRouterPath("/" html_file)
}

func (s *Server) Start() error {
	s.Logger.Println("starting server")
}

func (s *Server) AddConceptType(concept_type *ConceptType) error {
	if concept_type == nil {
		return errors.New("nil concept type")
	}
	for _, ct := range s.concept_types {
		if ct.type_name == concept_type.Type_name {
			return errors.New("duplicate concept name")
		}
		if ct.pathname != "" && ct.pathname == concept_type.Pathname {
			return errors.New("duplicate pathname")
		}
	}
	s.concept_types[concept_type.Type_name] = concept_type

	if concept_type.Api_available {
		single_endpoint_handler := func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			vars := mux.Vars(r)
			name := vars["name"]
			concept, err := DBConcept__getByTypeName(cxn, concept_type.Type_name, name)
			if err != nil {
				s.Logger(err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			concept.LoadRelationships()

			s.SendJSONResponse(w, r, concept)
		}
		s.AddRouterPath("/api/v1/c/" + concept_type.Pathname + "/{name}", "GET", single_endpoint_handler)

		many_endpoint_handler := func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			query_vars := r.URL.Query()
			count := Util__queryToInt(query_vars, "count", 0, 20, false, false, 20)
			offset := Util__queryToInt(query_vars, "offset", 0, 0, false, true, 0)
			concepts, err := DBConcept__getByType(cxn, concept_type.Type_name, offset, count)
			if err != nil {
				s.Logger(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			s.SendJSONResponse(w, r, concepts)
		}
		s.AddRouterPath("/api/v1/c/" + concept_type.Pathname), "GET", many_endpoint_handler)
	}

	return nil
}

func (s *Server) AddConceptRelationshipType(concept_relationship_type *ConceptRelationshipType) error {
	if concept_relationship_type == nil {
		return errors.New("nil concept relationship type")
	}
	for _, crt := range s.concept_relationship_types {
		if crt.type1 == concept_relationship_type.Type1 &&
			crt.type2 == concept_relationship_type.Type2 &&
			crt.string1 == concept_relationship_type.String1 &&
			crt.string2 == concept_relationship_type.String2 {
			return errors.New("duplicate concept relationship type")
		}
	}
	type1_found := false
	type2_found := false
	for _, ct := range s.concept_types {
		if ct.type_name == concept_relationship_type.Type1 {
			type1 = true
		} else if ct.type_name == concept_relationship_type.Type2 {
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
