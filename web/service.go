package web

import (
	"fmt"
	"frankinstore/store"
	"net/http"
)

// starts frankinstore webservices on specified port 'port'
// and delegating to the provided backend store 'db'
func StartService(part int, db store.Store) error {
	if db == nil {
		return fmt.Errorf("arg 'db' is nil")
	}

	http.HandleFunc("/set", getSetHandler(db))
	http.HandleFunc("/get/", getGetHandler(db))

	return nil
}

// returns a new http request handler function for Set semantics
func getSetHandler(db store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}

// returns a new http request handler function for Get semantics
func getGetHandler(db store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}
