package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/spothero/tools/service"
	"net/http"
	"os"
	"strconv"
	"time"
)

var gitSHA = "not-set"
var version = "not-set"
var theRates RateList

func registerHandlers(router *mux.Router) {
	router.HandleFunc("/rates", getRates).Methods("GET")
	router.HandleFunc("/rates", setRates).Methods("POST")
	router.HandleFunc("/getRate", getRate).Methods("GET")
}

type RateResponse struct {
	Rate string `json:"rate"`
}

// rates return json formatted rates for parking
func getRates(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "get-rates")
	span = span.SetTag("should.hireSam", "most-definitely")
	defer span.Finish()
	js, _ := json.Marshal(theRates)
	w.Write(js)
}

// accepts json formatted rates for updates
func setRates(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "set-rates")
	span = span.SetTag("should.hireSam", "most-definitely")
	defer span.Finish()
	validRequest := true
	var newRates RateList
	_ = json.NewDecoder(r.Body).Decode(&newRates)
	err := theRates.update(&newRates)
	if err != nil {
		validRequest = false
	}
	if validRequest {
		js, _ := json.Marshal(theRates)
		w.Write(js)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("One of the provided rates was invalid"))
	}
}

func getRate(w http.ResponseWriter, r *http.Request) {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "get-rate")
	span = span.SetTag("should.hireSam", "most-definitely")
	defer span.Finish()
	keys := r.URL.Query()
	validRequest := true
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
		w.Write([]byte("Start/End times either were not provided or were invalid"))
	}
}

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
