package goconcept

import (
	"net/http"
	"github.com/gorilla/mux"
	"errors"
	"log"
	"encoding/json"
)

func (s *Server) AddRouterPath(path string, method string, handler func(http.ResponseWriter, *http.Request, *Connection, *CookieWrapper)) error {
	if method != "GET" && method != "POST" {
		return errors.New("method must be GET or POST")
	}

	s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, s.connection, s.cookie_cookie_wrapper)
	}).Methods(method)

	return nil
}

func (s *Server) addStaticRouterPath(path string, dir string) error {
	s.router.PathPrefix("/" + path).Handler(http.StripPrefix("/" + path + "/", http.FileServer(http.Dir("./" + dir + "/"))))
	return nil
}

func (s *Server) addHTMLRouterPath(path_prefix string, html_file string) {
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
			s.Logger("unauthorized admin request")
			return
		}

		types := []string
		for type_name, _ := range s.concept_types {
			types = append(types, *type_name)
		}

		s.SendJSONResponse(w, r, &types)
	})

	s.AddRouterPath("/api/v1/ca/concept/types/{typename}", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
			return
		}

		vars := mux.Vars(r)
		typename := vars["typename"]
		concept_type, ok := concept_types[typename]
		if !ok {
			s.Logger.Printf("%+v\n", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		s.SendJSONRequest(w,r, concept_type)
	})

	s.AddRouterPath("/api/v1/ca/concept/relationships", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
			return
		}

		s.SendJSONResponse(w, r, &s.concept_relationship_types)
	})

	s.AddRouterPath("/api/v1/ca/concept/add", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
			return
		}

		body_data := struct {
			Type_name string	`json:"type"`
			Name string	`json:"name"`
		}{}
		err := Util__requestJSONDecode(r, &body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
			return
		}

		s.SendJSONRequest(new_concept)
	})

	s.AddRouterPath("/api/v1/ca/concept/delete", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
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
		s.SendJSONRequest(w, r, true)
	})

	s.AddRouterPath("/api/v1/ca/concept/add/data", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
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
				server_concept_data = cd
				break
			}
		}
		if server_concept_data == nil {
			http.Error(w, "invalid data key", http.StatusBadRequest)
			return
		}

		concept_type := DBConcept__getByTypeName(cxn, body_data.Type_name, body_data.Name)

		if server_concept_data.Single {
			for _, data := range concept_type.Data {
				if data.F_key == server_concept_data.Type_name {
					http.Error(w, "data key already exists", http.StatusBadRequest)
					return
				}
			}
		}

		if !server_concept_data.Approve_func(Data_value) {
			http.Error(w, "data is not valid", http.StatusBadRequest)
			return
		}

		new_concept_data, err := DBConceptData__create(cxn, concept_type.F_id, body_data.Key, body_data.Value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.SendJSONRequest(w, r, new_concept_data)
	})

	s.AddRouterPath("/api/v1/ca/concept/delete/data", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
			return
		}

		body_data := struct {
			Type_name string `json:"type"`
			Name string `json:"name"`
			Data_key string `json:"data_key"`
			Data_value string `json:"data_value"`
		}{}
		err := Util__requestJSONDecode(&body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		concept_type, err := DBConcept__getByTypeName(cxn, body_data.Type_name, body_data.Name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var concept_data DBConceptData
		for _, cd := range concept_type.Data {
			if cd.F_key == body_data.Data_key && cd.F_value == body_data.Data_value {
				concept_data == cd
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
		s.SendJSONRequest(w, r, true)
	})

	s.AddRouterPath("/api/v1/ca/concept/add/rel", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
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
		err := Util__requestJSONDecode(&body_data)
		if err != nil {
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
		if server_relationshp == nil {
			http.Error(w, "invalid type and string combination", http.StatusBadRequest)
			return
		}

		db_concept1, _ := DBConcept__getByTypeName(cxn, body_data.Type1, body_data.Name1)
		if db_concept1 == nil {
			http.Error(w, "invalid type1", http.StatusBadRequest)
			return
		}
		db_concept2, _ := DBConcept__getByTypeName(cxn, body_data.Type2, body_data.Name2)
		if db_concept2 == nil {
			http.Error(w, "invalid type2", http.StatusBadRequest)
			return
		}

		new_relationship, err := DBConceptRelationship__create(cxn, db_concept1.F_id, db_concept2.F_id, body_data.string1, body_data_string2)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.SendJSONRequest(w, r, new_relationship)
	})

	s.AddRouterPath("/api/v1/ca/concept/delete/rel", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		if(!AuthenticateRequest(r, cxn, cw, "admin")) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			s.Logger("unauthorized admin request")
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
		err := Util__requestJSONDecode(&body_data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		db_concept1, _ := DBConcept__getByTypeName(cxn, body_data.Type1, body_data.Name1)
		if db_concept1 == nil {
			http.Error(w, "invalid type1", http.StatusBadRequest)
			return
		}
		db_concept2, _ := DBConcept__getByTypeName(cxn, body_data.Type2, body_data.Name2)
		if db_concept2 == nil {
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
		s.SendJSONRequest(w, r, true)
	})
}

func (s *Server) addUserRoutes(allow_user_create bool) {
	if allow_user_create {
		s.AddRouterPath("/api/v1/user/register", "POST", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
			session, _ := cw.Get("base")
			logged_in_user, _ := Util__getUserFromSession(session)
			if logged_in_user != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			body_data := struct {
				Email string `json:"email"`
				Username string `json:"username"`
				Password string `json:"password"`
			}
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
		session, _ := cw.Get("base")
		logged_in_user, _ := Util__getUserFromSession(session)
		if logged_in_user != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		body_data := struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}
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
		session, _ := cw.Get("base")
		logged_in_user, _ := Util__getUserFromSession(session)
		if logged_in_user == nil {
			http.Error(w, "not logged in", http.StatusBadRequest)
			return
		}

		delete(session.Values, "user")
		session.Save(r, w)
		s.SendJSONRespose(w, r, true)
	})

	s.AddRouterPath("/api/v1/user/info", "GET", func(w http.ResponseWriter, r *http.Request, cxn *Connection, cw *CookieWrapper) {
		session, _ := cw.Get("base")
		logged_in_user, _ := Util__getUserFromSession(session)
		if logged_in_user == nil {
			http.Error(w, "not logged in", http.StatusBadRequest)
			return
		}

		s.SendJSONResponse(w, r, logged_in_user)
	})
}
