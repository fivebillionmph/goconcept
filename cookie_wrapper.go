package goconcept

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type SessionUser struct {
	Id int	`json:"-"`
	Username string `json:"username"`
	Level int	`json:"level"`
}

type CookieWrapper struct {
	store *sessions.CookieStore
}

func (cw *CookieWrapper) Get(r *http.Request, session_name string) (*sessions.CookieSession, error) {
	return cw.store.Get(r, session_name)
}
