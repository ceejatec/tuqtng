//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/couchbaselabs/tuqtng/ast"
	"github.com/couchbaselabs/tuqtng/network"
	"github.com/gorilla/mux"
)

type HttpResponse struct {
	w       http.ResponseWriter
	results ast.ValueChannel
}

func (this *HttpResponse) SendError(err error) {
	showError(this.w, fmt.Sprintf("%v", err), 500)
	close(this.results)
}

type HttpEndpoint struct {
	queryChannel network.QueryChannel
}

func NewHttpEndpoint(address string) *HttpEndpoint {
	rv := &HttpEndpoint{}

	r := mux.NewRouter()

	r.HandleFunc("/", welcome).Methods("GET")
	r.Handle("/query", rv).Methods("GET", "POST")

	go func() {
		err := http.ListenAndServe(address, r)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	return rv
}

func (this *HttpEndpoint) SendQueriesTo(queryChannel network.QueryChannel) {
	this.queryChannel = queryChannel
}

func welcome(w http.ResponseWriter, r *http.Request) {
	mustEncode(w, map[string]interface{}{
		"tuqtng":  "where no query has relaxed before",
		"version": "tuqtng 0.0",
	})
}

func (this *HttpEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	queryString := r.FormValue("q")
	if queryString == "" && r.Method == "POST" {
		queryStringBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			queryString = string(queryStringBytes)
		}
	}

	if queryString == "" {
		showError(w, "Missing required query string", 500)
		return
	} else {
		log.Printf("Query String: %v", queryString)
	}

	response := HttpResponse{w: w, results: make(ast.ValueChannel)}
	query := network.Query{
		Request:  network.UNQLStringQueryRequest{QueryString: queryString},
		Response: &response,
	}

	this.queryChannel <- query
	for val := range response.results {
		log.Printf("got a result val %v", val)
		// FIXME this where we serialize a result row and stream back to client
	}
}

func mustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}

func showError(w http.ResponseWriter, msg string, code int) {
	log.Printf("Reporting error %v/%v", code, msg)
	http.Error(w, msg, code)
}
