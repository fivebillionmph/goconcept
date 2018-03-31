package goconcept

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type CookieWrapper struct {
	store *sessions.CookieStore
}

func (cw *CookieWrapper) Get(r *http.Request, session_name string) (*sessions.Session, error) {
	return cw.store.Get(r, session_name)
}
