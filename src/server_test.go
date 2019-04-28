package main

import (
	"testing"
  "github.com/stretchr/testify/assert"
  "net/http"
  "github.com/gorilla/mux"
  "net/http/httptest"
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
}
