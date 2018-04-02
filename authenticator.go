package goconcept

import (
	"net/http"
)

func AuthenticateRequest(r *http.Request, cxn *Connection, cw *CookieWrapper, auth_type string) bool {
	switch auth_type {
		case "admin":
			return authenticate__admin(r, cxn, cw)
		default:
			return false
	}
}

func authenticate__admin(r *http.Request, cxn *Connection, cw *CookieWrapper) bool {
	session, _ := cw.Get(r, "base")
	logged_in_user, _ := Util__getUserFromSession(session)
	if logged_in_user != nil {
		if logged_in_user.Level > 1 {
			return true
		}
	}

	api_key_str := r.Header.Get(DBAPIKey__header_name)
	if api_key_str != "" {
		api_key, _ := DBAPIKey__getByKey(cxn, api_key_str)
		if api_key != nil {
			user, _ := DBUser__getByID(cxn, api_key.F_user_id)
			if user != nil && user.F_level > 1 {
				return true
			}
		}
	}

	return false
}
