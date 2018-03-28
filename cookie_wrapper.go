package goconcept

import (
	"github.com/gorilla/sessions"
)

type CookieWrapper struct {
	store *sessions.CookieStore
}
