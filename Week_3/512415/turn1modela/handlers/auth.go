package handlers

import (
	"net/http"
)

var users = map[string]string{
	"user": "password",
}

func Authenticate(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()

	if ok && password == users[username] {
		return true
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	return false
}
