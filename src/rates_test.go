package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGetRate(t *testing.T) {
	rateFile, _ := os.Open("/testfiles/testing.json")
	defer rateFile.Close()
	byteValue, _ := ioutil.ReadAll(rateFile)
	var theRates RateList
	json.Unmarshal(byteValue, &theRates)
	// 2015-07-01T07:00:00-05:00 to 2015-07-01T12:00:00-05:00 should yield 1750
	t1, _ := time.Parse(time.RFC3339, "2015-07-01T07:00:00-05:00")
	t2, _ := time.Parse(time.RFC3339, "2015-07-01T12:00:00-05:00")
	rq := RateQuery{startTime: t1, endTime: t2}
	assert.Equal(t, 1750, theRates.getRate(rq))
	// 2015-07-04T15:00:00+00:00 to 2015-07-04T20:00:00+00:00 should yield 2000
	t3, _ := time.Parse(time.RFC3339, "2015-07-04T15:00:00+00:00")
	t4, _ := time.Parse(time.RFC3339, "2015-07-04T20:00:00+00:00")
	rq2 := RateQuery{startTime: t3, endTime: t4}
	assert.Equal(t, 2000, theRates.getRate(rq2))
	// 2015-07-04T07:00:00+05:00 to 2015-07-04T20:00:00+05:00 should yield unavailable (-1)
	t5, _ := time.Parse(time.RFC3339, "2015-07-04T07:00:00+05:00")
	t6, _ := time.Parse(time.RFC3339, "2015-07-04T20:00:00+05:00")
	rq3 := RateQuery{startTime: t5, endTime: t6}
	assert.Equal(t, -1, theRates.getRate(rq3))
	// test some second stuff
	t7, _ := time.Parse(time.RFC3339, "2015-07-01T07:00:00-05:00")
	t8, _ := time.Parse(time.RFC3339, "2015-07-01T18:00:01-05:00")
	rq4 := RateQuery{startTime: t7, endTime: t8}
	assert.Equal(t, -1, theRates.getRate(rq4))

}

func TestQueryValidation(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2015-07-01T07:00:00-05:00")
	t2, _ := time.Parse(time.RFC3339, "2015-07-01T12:00:00-05:00")
	rq1 := RateQuery{startTime: t1, endTime: t2}
	assert.True(t, rq1.validate())
	t3, _ := time.Parse(time.RFC3339, "2015-07-02T07:00:00-05:00")
	t4, _ := time.Parse(time.RFC3339, "2015-07-01T12:00:00-05:00")
	rq2 := RateQuery{startTime: t3, endTime: t4}
	assert.False(t, rq2.validate())
}

func TestRateValidation(t *testing.T) {
	r1 := Rate{Days: "wed", Times: "0900-1200", TimeZone: "America/Chicago", Price: 500}
	v1 := r1.validate()
	assert.True(t, v1)
	r2 := Rate{Days: "wed", Times: "0900-1200", TimeZone: "America/Chigo", Price: 500}
	v2 := r2.validate()
	assert.False(t, v2)
	r3 := Rate{Days: "wed", Times: "0900-1200", TimeZone: "America/Chicago", Price: -1}
	v3 := r3.validate()
	assert.False(t, v3)
	r4 := Rate{Days: "tuesday", Times: "0900-1200", TimeZone: "America/Chicago", Price: 500}
	v4 := r4.validate()
	assert.False(t, v4)
}

func TestRateListUpdate(t *testing.T) {
	var theRates RateList
	rateFile, _ := os.Open("/testfiles/testing.json")
	defer rateFile.Close()
	byteValue, _ := ioutil.ReadAll(rateFile)
	var theNewRates RateList
	json.Unmarshal(byteValue, &theNewRates)
	theRates.update(&theNewRates)
	t1, _ := time.Parse(time.RFC3339, "2015-07-01T07:00:00-05:00")
	t2, _ := time.Parse(time.RFC3339, "2015-07-01T12:00:00-05:00")
	rq := RateQuery{startTime: t1, endTime: t2}
	assert.Equal(t, 1750, theRates.getRate(rq))
}
