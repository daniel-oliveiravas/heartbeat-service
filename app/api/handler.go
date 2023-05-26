package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func New() http.Handler {
	router := httprouter.New()
	return router
}
