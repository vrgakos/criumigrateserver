package main

import (
	"net/http"

	"github.com/vrgakos/criumigrateserver/handlers"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},

	Route{
		"BaseScp",
		"POST",
		"/base/scp",
		handlers.BaseScp,
	},
	Route{
		"BaseClean",
		"POST",
		"/base/clean",
		handlers.BaseClean,
	},

	Route{
		"CriuPreDump",
		"POST",
		"/criu/pre-dump",
		handlers.CriuPreDump,
	},
	Route{
		"CriuDump",
		"POST",
		"/criu/dump",
		handlers.CriuDump,
	},
	Route{
		"CriuRestore",
		"POST",
		"/criu/restore",
		handlers.CriuRestore,
	},
	Route{
		"CriuLazyPages",
		"POST",
		"/criu/lazy-pages",
		handlers.CriuLazyPages,
	},
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}