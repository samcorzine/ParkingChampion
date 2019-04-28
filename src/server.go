// Copyright 2019 SpotHero
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"os"
  "encoding/json"
  "time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/spothero/tools/service"
)

// These variables should be set during build with the Go link tool
// e.x.: when running go build, provide -ldflags="-X main.version=1.0.0"
var gitSHA = "not-set"
var version = "not-set"

var theRates RateList

// registerHandlers is a callback used to register HTTP endpoints to the default server
// NOTE: The HTTP server automatically registers /health and /metrics -- Have a look in your
// browser!
func registerHandlers(router *mux.Router) {
	router.HandleFunc("/rate", getRates).Methods("GET")
	router.HandleFunc("/rate", setRates).Methods("POST")
  router.HandleFunc("/getRate", getRate).Methods("GET")
}

type RateResponse struct {
	Rate         int `json:"rate"`
}

// rates return json formatted rates for parking
func getRates(w http.ResponseWriter, r *http.Request) {
	// NOTE: This is an example of an opentracing span
	span, _ := opentracing.StartSpanFromContext(r.Context(), "get-rates")
	span = span.SetTag("best.language", "golang")
	span = span.SetTag("best.mascot", "gopher")
	defer span.Finish()
  js, _ := json.Marshal(theRates)
	w.Write(js)
}

// rates return json formatted rates for parking
func setRates(w http.ResponseWriter, r *http.Request) {
	// NOTE: This is an example of an opentracing span
	span, _ := opentracing.StartSpanFromContext(r.Context(), "set-rate")
	span = span.SetTag("best.language", "golang")
	span = span.SetTag("best.mascot", "gopher")
	defer span.Finish()
  _ = json.NewDecoder(r.Body).Decode(&theRates)
  js, _ := json.Marshal(theRates)
	w.Write(js)
}

// getRate returns a rate for any valid
func getRate(w http.ResponseWriter, r *http.Request) {
	// NOTE: This is an example of an opentracing span
	span, _ := opentracing.StartSpanFromContext(r.Context(), "get-rate")
	span = span.SetTag("best.language", "golang")
	span = span.SetTag("best.mascot", "gopher")
	defer span.Finish()
  keys := r.URL.Query()
  startString := keys.Get("start")
  endString := keys.Get("end")
  startT, _ := time.Parse(time.RFC3339, startString)
  endT, _ := time.Parse(time.RFC3339, endString)
  rq := RateQuery{startTime: startT, endTime: endT}
  rate := theRates.getRate(rq)
  rateResponse := RateResponse{Rate:rate}
  js, _ := json.Marshal(rateResponse)
  w.Write(js)
}

// This is the main entrypoint of the program. Here we create our root command and then execute it.
func main() {
	serverCmd := service.HTTPConfig{
		Config: service.Config{
			Name:        "parking_rates",
			Version:     version,
			GitSHA:      gitSHA,
			Environment: "local",
		},
		RegisterHandlers: registerHandlers,
	}
	if err := serverCmd.ServerCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
