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
	cookie_store *CookieWrapper
	concept_types map[string]*ConceptType
	concept_relationship_types []*ConceptRelationshipType
	connection *Connection
	router *mux.Router
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
	api_handler := &APIHandler{router, cxn, cookie_store, logger}
	return &Server{http_server, logger, cookie_store, concept_types, concept_relationship_types, connection, router}
}

fun (s *Server) Run(static_path string, static_dir, string, html_file string, admin_path string, admin_html_file string) error {
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
		if ct.type_name == concept_type.type_name {
			return errors.New("duplicate concept name")
		}
		if ct.pathname != "" && ct.pathname == concept_type.pathname {
			return errors.New("duplicate pathname")
		}
	}
	s.concept_types[concept_type.type_name] = concept_type

	if concept_type.api_available {
		single_endpoint_handler := func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			vars := mux.Vars(r)
			name := vars["name"]
			concept, err := DBConcept__getByTypeName(cxn, concept_type.type_name, name)
			if err != nil {
				s.Logger(err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			concept.LoadRelationships()

			s.SendJSONResponse(w, r, concept)
		}
		s.AddRouterPath("/api/v1/c/" + concept_type.pathname + "/{name}", "GET", single_endpoint_handler)

		many_endpoint_handler := func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			query_vars := r.URL.Query()
			count := 20
			offset := 0

			req_count, ok := query_vars["count"]
			if ok {
				if len(req_count) > 0 {
					req_count_int, err := strconv.Atoi(req_count[0])
					if err == nil && req_count_int < 100 && req_count_int > 0 {
						count = req_count_int
					}
				}
			}

			req_offset, ok := query_vars["offset"]
			if ok {
				if len(req_offset) > 0 {
					req_offset_int, err := strconv.Atoi(req_offset[0])
					if err == nil && req_offset_int >= 0 {
						offset = req_offset_int
					}
				}
			}

			concepts, err := DBConcept__getByType(cxn, concept_type.type_name, offset, count)
			if err != nil {
				s.Logger(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			s.SendJSONResponse(w, r, concepts)
		}
		s.AddRouterPath("/api/v1/c/" + concept_type.pathname), "GET", many_endpoint_handler)
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
