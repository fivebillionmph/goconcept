package goconcept

import (
	"errors"
	"net/url"
	"net/http"
	"encoding/json"
	"github.com/gorilla/sessions"
	"strconv"
)

func Util__queryToInt(vars url.Values, var_name string, min int, max int, min_inf bool, max_inf bool, default_value int) int {
	val_str, ok := vars[var_name]
	if !ok {
		return default_value
	}
	if len(val_str) < 1 {
		return default_value
	}
	val_int, err := strconv.Atoi(val_str[0])
	if err != nil {
		return default_value
	}
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

func Util__getUserFromSession(session *sessions.Session) (*SessionUser, error) {
	user_interface, ok := session.Values["user"]
	if !ok {
		return nil, errors.New("no user")
	}

	user, ok := user_interface.(*SessionUser)
	if !ok {
		return nil, errors.New("invalid stored type")
	}
	return user, nil
}

func Util__saveUserToSession(w http.ResponseWriter, r *http.Request, session *sessions.Session, user *DBUser) *SessionUser {
	session_user := &SessionUser{user.F_id, user.F_username, int(user.F_level)}
	session.Values["user"] = session_user
	session.Save(r, w)
	return session_user
}
