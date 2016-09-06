//    Copyright © 2016 Joubin Houshyar. All rights reserved.
//
//    This file is part of Frankinstore.
//
//    Frankinstore is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as
//    published by the Free Software Foundation, either version 3 of
//    the License, or (at your option) any later version.
//
//    Frankinstore is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public
//    License along with Frankinstore.  If not, see <http://www.gnu.org/licenses/>.

// package defines the web-api and supporting services.
package web

import (
	"encoding/hex"
	"fmt"
	"frankinstore/store"
	"io/ioutil"
	"net/http"
	"path"
)

/// services //////////////////////////////////////////////////////////////////

const DefaultPort = 5722

// starts frankinstore webservices on specified port 'port'
// and delegating to the provided backend store 'db'
func RunService(port int, db store.Store) error {
	if db == nil {
		return fmt.Errorf("arg 'db' is nil")
	}

	http.HandleFunc("/info", getInfoHandler(db))
	http.HandleFunc("/set", getSetHandler(db))
	http.HandleFunc("/get/", getGetHandler(db))

	addr := fmt.Sprintf(":%d", port)

	return http.ListenAndServe(addr, nil)
}

// convenince error response function
func onError(w http.ResponseWriter, code int, fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	http.Error(w, msg, code)
}

/// handlers //////////////////////////////////////////////////////////////////

// returns a new http request handler function for Set semantics
//
// The returned handler will service POST method requests, with request
// body (binary blob) uses as 'value' to store. Successful addtions to store
// will result in return of (hex encoded) key or error as returned by the db.
func getSetHandler(db store.Store) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		/* assert constraints */
		if req.Method != "POST" {
			onError(w, http.StatusBadRequest, "expect POST method - have %s", req.Method)
			return
		}
		if req.ContentLength < 1 {
			onError(w, http.StatusBadRequest, "value data not provided")
			return
		}

		// get post data
		blob, e := ioutil.ReadAll(req.Body)
		if e != nil {
			onError(w, http.StatusInternalServerError, e.Error())
			return
		}

		// process request
		key, e := db.Put(blob)
		if e != nil {
			// TODO: need to distinguish top level errors e.g. NotFouund
			// REVU: ok for now
			onError(w, http.StatusBadRequest, e.Error())
			return
		}

		// post response - note binary key is hex encoded
		encoded := []byte(hex.EncodeToString(key[:]))
		w.Write(encoded)

		return
	}
}

// returns a new http request handler function for Get semantics
func getGetHandler(db store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		/* assert constraints */
		if req.Method != "GET" {
			onError(w, http.StatusBadRequest, "expect GET method - have %s", req.Method)
			return
		}

		// service api is assumed as ../get/<sha-hexstring>
		_, keystr := path.Split(req.URL.Path)
		if keystr == "" {
			onError(w, http.StatusBadRequest, "key not provided")
			return
		}
		b, e := hex.DecodeString(keystr)
		if e != nil {
			onError(w, http.StatusBadRequest, e.Error())
			return
		}

		// process request
		var key store.Key
		copy(key[:], b)
		val, e := db.Get(key)
		if e != nil {
			onError(w, http.StatusBadRequest, e.Error())
			return
		}
		// post response - note value is returned in binary form as original
		w.Write(val)
	}
}

func getInfoHandler(db store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		/* assert constraints */
		if req.Method != "GET" {
			onError(w, http.StatusBadRequest, "expect GET method - have %s", req.Method)
			return
		}

		// process request
		info, e := db.Info()
		if e != nil {
			onError(w, http.StatusBadRequest, e.Error())
			return
		}
		// post response - note value is returned in binary form as original
		w.Write(info)
	}
}
