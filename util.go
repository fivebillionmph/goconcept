package goconcept

import (
	"errors"
	"net/url"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gorilla/sessions"
)

func Util__queryToInt(vars url.Values, var_name string, min int, max int, min_inf bool, max_inf bool, default_value int) int {
	val_str, ok := vars[var_name]
	if !ok {
		return default_value
	}
	if len(val_str) < 1 {
		return default_value
	}
	val_int, err := strcon.Atoi(val_str[0])
	if !min_inf && val_int < min {
		return default_value
	}
	if !max_inf && val_int > max {
		return default_value
	}
	return val_int
}

func Util__requestJSONDecode(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(s)
	return err
}

func Util__getUserFromSession(session *sessions.CookieStore) (*SessionUser, error) {
	user, ok := session.Values["user"]
	if !ok {
		return nil errors.New("no user")
	}
	return user, nil
}

func Util__saveUserToSession(w http.ResponseWriter, r *http.Request, session *sessions.CookieStore, user *DBUser) (*SessionUser) {
	session_user := &SessionUser{user.F_id, user.F_username, user.F_level}
	session["user"] = session_user
	session.Save(r, w)
}
