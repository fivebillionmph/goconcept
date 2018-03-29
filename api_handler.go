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
		handler(w, r, s.connection, s.cookie_store)
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
