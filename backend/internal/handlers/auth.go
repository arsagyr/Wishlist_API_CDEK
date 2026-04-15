package handlers

import (
	"net/http"
)

type AuthHandlerFunc func(w http.ResponseWriter, r *http.Request)

func (f AuthHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}
