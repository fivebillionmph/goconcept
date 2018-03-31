package goconcept

import (
	"net/http"
)

func AuthenticateRequest(r *http.Request, cxn *Connection, cw *CookieWrapper, auth_type string) bool {
	switch auth_type {
		case "admin":
			return authenticate__admin(r *http.Request, cxn *Connection, cw *CookieWrapper)
		default:
			return false
	}
}

func authenticate__admin(r *http.Request, cxn *Connection, cw *CookieWrapper) bool {
	
}
