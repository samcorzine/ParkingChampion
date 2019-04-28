package main

import (
	"testing"
  "github.com/stretchr/testify/assert"
  "net/http"
  "github.com/gorilla/mux"
  "net/http/httptest"
  "os"
)

func TestGetRateEndpoint(t *testing.T) {
  router := mux.NewRouter()
  router.HandleFunc("/getRate", getRate).Methods("GET")
  req, err := http.NewRequest("GET", "/getRate", nil)
	assert.NoError(t, err)
  httpRec := httptest.NewRecorder()
  router.ServeHTTP(httpRec, req)
  resp := httpRec.Result()
  assert.Equal(t, 400, resp.StatusCode)
  req2, err2 := http.NewRequest("GET", "getRate?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00", nil)
	assert.NoError(t, err2)
  httpRec2 := httptest.NewRecorder()
  router.ServeHTTP(httpRec2, req2)
  resp2 := httpRec2.Result()
  assert.Equal(t, 301, resp2.StatusCode)
}

func TestSetRatesEndpoint(t *testing.T) {
  router := mux.NewRouter()
  router.HandleFunc("/rates", setRates).Methods("POST")

  rateFile, _ := os.Open("/testfiles/testing.json")
  defer rateFile.Close()
  req, err := http.NewRequest("POST", "/rates", rateFile)
  assert.NoError(t, err)
  httpRec := httptest.NewRecorder()
  router.ServeHTTP(httpRec, req)
  resp := httpRec.Result()
  assert.Equal(t, 200, resp.StatusCode)

  rateFile2, _ := os.Open("/testfiles/testing_invalid.json")
  defer rateFile2.Close()
  req2, err2 := http.NewRequest("POST", "/rates", rateFile2)
  assert.NoError(t, err2)
  httpRec2 := httptest.NewRecorder()
  router.ServeHTTP(httpRec2, req2)
  resp2 := httpRec2.Result()
  assert.Equal(t, 400, resp2.StatusCode)
}
