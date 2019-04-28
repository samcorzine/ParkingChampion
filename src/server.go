package main

import (
	"net/http"
	"os"
  "encoding/json"
  "time"
	"strconv"

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
	router.HandleFunc("/rates", getRates).Methods("GET")
	router.HandleFunc("/rates", setRates).Methods("POST")
  router.HandleFunc("/getRate", getRate).Methods("GET")
}

type RateResponse struct {
	Rate         string `json:"rate"`
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
	span, _ := opentracing.StartSpanFromContext(r.Context(), "set-rates")
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
	validRequest  := true
  startString := keys.Get("start")
  endString := keys.Get("end")
  startT, startErr := time.Parse(time.RFC3339, startString)
	if startErr != nil {
		validRequest = false
	}
  endT, endErr := time.Parse(time.RFC3339, endString)
	if endErr != nil {
		validRequest = false
	}
	if validRequest {
	  rq := RateQuery{startTime: startT, endTime: endT}
	  rate := theRates.getRate(rq)
		var resp RateResponse
		if rate == -1 {
			resp.Rate = "unavailable"
		} else {
			resp.Rate = strconv.Itoa(rate)
		}
	  js, _ := json.Marshal(resp)
	  w.Write(js)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Start/End times either were not provided or were invalid "))
	}
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
