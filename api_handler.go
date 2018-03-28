package goconcept

import (
	"net/http"
	"github.com/gorilla/mux"
	"errors"
)

type APIHandler struct {
	router *mux.Router
	connection *Connection
	cookie_store *CookieWrapper
}

func (h *APIHandler) AddPath(path string, method string, handler func(http.ResponseWriter, *http.Request, *Connection, *CookieWrapper)) error {
	if method != "GET" && method != "POST" {
		return errors.New("method must be GET or POST")
	}

	h.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, h.connection, h.cookie_store)
	}).Methods(method)

	return nil
}

func (h *APIHandler) addStaticPath(path string, dir string) error {
	s.router.PathPrefix("/" + path).Handler(http.StripPrefix("/" + path + "/", http.FileServer(http.Dir("./" + dir + "/"))))
	return nil
}

func (h *APIHandler) addHTMLPath(path_prefix string, html_file string) {
	h.router.PathPrefix(path_prefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, html_file)
	})
}
