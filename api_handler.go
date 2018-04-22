package goconcept

import (
	"net/http"
	"github.com/gorilla/mux"
	"errors"
	"encoding/json"
)

func (s *Server) AddRouterPath(path string, method string, handler func(http.ResponseWriter, *http.Request, *Connection, *CookieWrapper)) error {
	if method != "GET" && method != "POST" {
		return errors.New("method must be GET or POST")
	}

	s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, s.connection, s.cookie_wrapper)
	}).Methods(method)

	return nil
}

func (s *Server) AddStaticRouterPathPrefix(path string, dir string) error {
	s.router.PathPrefix(path).Handler(http.StripPrefix(path + "/", http.FileServer(http.Dir(dir))))
	return nil
}

func (s *Server) AddSingleFileServePathPrefix(path_prefix string, html_file string) {
	s.router.PathPrefix(path_prefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, html_file)
	})
}

func (s *Server) SendJSONResponse(w http.ResponseWriter, r *http.Request, obj interface{}) {
	json_response, err := json.Marshal(obj)
	if err != nil {
		s.Logger.Printf("%+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(json_response)
}

func (s *Server) addAdminRoutes() {
	s.AddRouterPath("/api/v1/ca/concept/types", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		types := []*ConceptType{}
		for _, t := range s.concept_types {
			types = append(types, t)
		}

		s.SendJSONResponse(w, r, &types)
	})

	s.AddRouterPath("/api/v1/ca/concept/data", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if !AuthenticateRequest(r, cxn, cw, "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		query_vars := r.URL.Query()
		count := Util__queryToInt(query_vars, "count", 0, 20, false, false, 20)
		offset := Util__queryToInt(query_vars, "offset", 0, 0, false, true, 0)

		concepts, err := DBConcept__getAll(cxn, offset, count)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.SendJSONResponse(w, r, concepts)
	})

	s.AddRouterPath("/api/v1/ca/concept/data/{type}", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if !AuthenticateRequest(r, cxn, cw, "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		path_vars := mux.Vars(r)
		type_name := path_vars["type"]

		query_vars := r.URL.Query()
		count := Util__queryToInt(query_vars, "count", 0, 20, false, false, 20)
		offset := Util__queryToInt(query_vars, "offset", 0, 0, false, true, 0)

		concepts, err := DBConcept__getByType(cxn, type_name, offset, count)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.SendJSONResponse(w, r, concepts)
	})

	s.AddRouterPath("/api/v1/ca/concept/data/{type}/{name}", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if !AuthenticateRequest(r, cxn, cw, "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		path_vars := mux.Vars(r)
		type_name := path_vars["type"]
		name := path_vars["name"]

		concept, err := DBConcept__getByTypeName(cxn, type_name, name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		concept.LoadRelationships(cxn)

		s.SendJSONResponse(w, r, concept)
	})

	s.AddRouterPath("/api/v1/ca/concept/relationships", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if !AuthenticateRequest(r, cxn, cw, "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		s.SendJSONResponse(w, r, &s.concept_relationship_types)
	})

	s.AddRouterPath("/api/v1/ca/concept/add", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		body_data := struct {
			Type_name string	`json:"type"`
			Name string	`json:"name"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			s.Logger.Println(err)
			return
		}

		_, type_exists := s.concept_types[body_data.Type_name]
		if !type_exists {
			http.Error(w, "type does not exist", http.StatusBadRequest)
			return
		}

		existing_concept, _ := DBConcept__getByTypeName(cxn, body_data.Type_name, body_data.Name)
		if existing_concept != nil {
			http.Error(w, "type already exists", http.StatusBadRequest)
			return
		}

		new_concept, err := DBConcept__create(cxn, body_data.Type_name, body_data.Name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			s.Logger.Println(err)
			return
		}

		s.SendJSONResponse(w, r, new_concept)
	})

	s.AddRouterPath("/api/v1/ca/concept/delete", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		body_data := struct {
			Type_name string `json:"type"`
			Name string `json:"name"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		existing_concept, err := DBConcept__getByTypeName(cxn, body_data.Type_name, body_data.Name)
		if existing_concept == nil {
			http.Error(w, "could not find concept", http.StatusBadRequest)
			return
		}

		err = DBConcept__delete(cxn, existing_concept)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			s.Logger.Println("%+v\n", err)
			return
		}
		s.SendJSONResponse(w, r, true)
	})

	s.AddRouterPath("/api/v1/ca/concept/add/data", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		body_data := struct {
			Type_name string `json:"type"`
			Name string `json:"name"`
			Data_key string `json:"data_key"`
			Data_value string `json:"data_value"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		server_type, type_exists := s.concept_types[body_data.Type_name]
		if !type_exists {
			http.Error(w, "type does not exist", http.StatusBadRequest)
			return
		}

		var server_concept_data *ConceptData = nil
		for _, cd := range server_type.Concept_data {
			if cd.Type_name == body_data.Data_key {
				server_concept_data = &cd
				break
			}
		}
		if server_concept_data == nil {
			http.Error(w, "invalid data key", http.StatusBadRequest)
			return
		}

		concept_type, err := DBConcept__getByTypeName(cxn, body_data.Type_name, body_data.Name)
		if err != nil {
			http.Error(w, "type does not exist", http.StatusBadRequest)
			return
		}

		if server_concept_data.Single {
			for _, data := range *concept_type.Data {
				if data.F_key == server_concept_data.Type_name {
					http.Error(w, "data key already exists", http.StatusBadRequest)
					return
				}
			}
		}

		if !server_concept_data.Approve_func(body_data.Data_value) {
			http.Error(w, "data is not valid", http.StatusBadRequest)
			return
		}

		new_concept_data, err := DBConceptData__create(cxn, concept_type.F_id, body_data.Data_key, body_data.Data_value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.SendJSONResponse(w, r, new_concept_data)
	})

	s.AddRouterPath("/api/v1/ca/concept/delete/data", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		body_data := struct {
			Type_name string `json:"type"`
			Name string `json:"name"`
			Data_key string `json:"data_key"`
			Data_value string `json:"data_value"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		concept_type, err := DBConcept__getByTypeName(cxn, body_data.Type_name, body_data.Name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var concept_data *DBConceptData
		for _, cd := range *concept_type.Data {
			if cd.F_key == body_data.Data_key && cd.F_value == body_data.Data_value {
				concept_data = &cd
				break
			}
		}

		if concept_data == nil {
			http.Error(w, "data not found", http.StatusBadRequest)
			return
		}

		err = DBConceptData__delete(cxn, concept_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			s.Logger.Println("%+v\n", err)
			return
		}
		s.SendJSONResponse(w, r, true)
	})

	s.AddRouterPath("/api/v1/ca/concept/add/rel", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		body_data := struct {
			Type1 string `json:"type1"`
			Type2 string `json:"type2"`
			Name1 string `json:"name1"`
			Name2 string `json:"name2"`
			String1 string `json:"string1"`
			String2 string `json:"string2"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			s.Logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var server_relationship *ConceptRelationshipType
		for _, sr := range s.concept_relationship_types {
			if sr.Type1 == body_data.Type1 && sr.Type2 == body_data.Type2 && sr.String1 == body_data.String1 && sr.String2 == body_data.String2 {
				server_relationship = sr
				break
			}
		}
		if server_relationship == nil {
			http.Error(w, "invalid type and string combination", http.StatusBadRequest)
			return
		}

		db_concept1, err := DBConcept__getByTypeName(cxn, body_data.Type1, body_data.Name1)
		if db_concept1 == nil {
			s.Logger.Println(err)
			http.Error(w, "invalid type1", http.StatusBadRequest)
			return
		}
		db_concept2, err := DBConcept__getByTypeName(cxn, body_data.Type2, body_data.Name2)
		if db_concept2 == nil {
			s.Logger.Println(err)
			http.Error(w, "invalid type2", http.StatusBadRequest)
			return
		}

		new_relationship, err := DBConceptRelationship__create(cxn, db_concept1.F_id, db_concept2.F_id, body_data.String1, body_data.String2)
		if err != nil {
			s.Logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.SendJSONResponse(w, r, new_relationship)
	})

	s.AddRouterPath("/api/v1/ca/concept/delete/rel", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		body_data := struct {
			Type1 string `json:"type1"`
			Type2 string `json:"type2"`
			Name1 string `json:"name1"`
			Name2 string `json:"name2"`
			String1 string `json:"string1"`
			String2 string `json:"string2"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		db_concept1, err := DBConcept__getByTypeName(cxn, body_data.Type1, body_data.Name1)
		if db_concept1 == nil {
			s.Logger.Println(err)
			http.Error(w, "invalid type1", http.StatusBadRequest)
			return
		}
		db_concept2, err := DBConcept__getByTypeName(cxn, body_data.Type2, body_data.Name2)
		if db_concept2 == nil {
			s.Logger.Println(err)
			http.Error(w, "invalid type2", http.StatusBadRequest)
			return
		}

		relationship, err := DBConceptRelationship__getByIDsStrings(cxn, db_concept1.F_id, db_concept2.F_id, body_data.String1, body_data.String2)
		if err != nil {
			http.Error(w, "could not find relationship", http.StatusBadRequest)
			return
		}

		err = DBConceptRelationship__delete(cxn, relationship)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			s.Logger.Println("%+v\n", err)
			return
		}
		s.SendJSONResponse(w, r, true)
	})

	s.AddRouterPath("/api/v1/ca/concept/data-search/{type}", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if !AuthenticateRequest(r, cxn, cw, "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger.Println("unauthorized admin request")
			return
		}

		path_vars := mux.Vars(r)
		type_name := path_vars["type"]

		query_vars := r.URL.Query()
		q, ok := query_vars["q"]
		if !ok || len(q) == 0 {
			http.Error(w, "q query required", http.StatusBadRequest)
			return
		}
		count := Util__queryToInt(query_vars, "count", 0, 100, false, false, 100)
		offset := Util__queryToInt(query_vars, "offset", 0, 0, false, true, 0)

		concepts, err := DBConcept__getBySearchName(cxn, type_name, q[0], offset, count)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.SendJSONResponse(w, r, concepts)
	})
}

func (s *Server) addUserRoutes(allow_user_create bool) {
	if allow_user_create {
		s.AddRouterPath("/api/v1/user/register", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			session, _ := cw.Get(r, "base")
			logged_in_user, _ := Util__getUserFromSession(session)
			if logged_in_user != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			body_data := struct {
				Email string `json:"email"`
				Username string `json:"username"`
				Password string `json:"password"`
			}{}
			err := Util__requestJSONDecode(r, &body_data)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			user, err := DBUser__create(cxn, body_data.Email, body_data.Password, body_data.Username, 1)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			session_user := Util__saveUserToSession(w, r, session, user)
			s.SendJSONResponse(w, r, session_user)
		})
	}

	s.AddRouterPath("/api/v1/user/login", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		session, _ := cw.Get(r, "base")
		logged_in_user, _ := Util__getUserFromSession(session)
		if logged_in_user != nil {
			http.Error(w, "Already logged in", http.StatusBadRequest)
			return
		}

		body_data := struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user, err := DBUser__getByPasswordChallenge(cxn, body_data.Email, body_data.Password)
		if err != nil {
			http.Error(w, "email password combination does not exist", http.StatusBadRequest)
			return
		}

		session_user := Util__saveUserToSession(w, r, session, user)
		s.SendJSONResponse(w, r, session_user)
	})

	s.AddRouterPath("/api/v1/user/logout", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		session, _ := cw.Get(r, "base")
		logged_in_user, _ := Util__getUserFromSession(session)
		if logged_in_user == nil {
			http.Error(w, "not logged in", http.StatusBadRequest)
			return
		}

		delete(session.Values, "user")
		session.Save(r, w)
		s.SendJSONResponse(w, r, true)
	})

	s.AddRouterPath("/api/v1/user/info", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		session, _ := cw.Get(r, "base")
		logged_in_user, _ := Util__getUserFromSession(session)
		if logged_in_user == nil {
			http.Error(w, "not logged in", http.StatusBadRequest)
			return
		}

		s.SendJSONResponse(w, r, logged_in_user)
	})

	s.AddRouterPath("/api/v1/user/keys", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		session, _ := cw.Get(r, "base")
		user, _ := Util__getUserFromSession(session)
		if user == nil {
			http.Error(w, "not logged in", http.StatusBadRequest)
			return
		}

		api_keys, err := DBAPIKey__getByUserID(cxn, user.Id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.SendJSONResponse(w, r, api_keys)
	})

	s.AddRouterPath("/api/v1/user/keys/add", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		session, _ := cw.Get(r, "base")
		user, _ := Util__getUserFromSession(session)
		if user == nil {
			http.Error(w, "not logged in", http.StatusBadRequest)
			return
		}

		existing_count, err := DBAPIKey__getCountByUserID(cxn, user.Id, true)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if existing_count > DBAPIKey__user_max {
			http.Error(w, "max number of keys reached", http.StatusBadRequest)
			return
		}

		db_user, err := DBUser__getByID(cxn, user.Id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		key, err := DBAPIKey__create(cxn, db_user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		s.SendJSONResponse(w, r, key)
	})
}
