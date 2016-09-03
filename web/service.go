package web

import (
	"fmt"
	"frankinstore/store"
	"net/http"
)

// REVU: don't think we need to return the Service
func StartService(part int, db store.Store) error {
	if db == nil {
		return fmt.Errorf("arg 'db' is nil")
	}

	http.HandleFunc("/set", getSetHandler(db))
	http.HandleFunc("/get/", getGetHandler(db))

	return nil
}

func getSetHandler(db store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}

func getGetHandler(db store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}
