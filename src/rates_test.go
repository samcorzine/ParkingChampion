package main

import (
	"testing"
  "github.com/stretchr/testify/assert"
  "os"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "time"
)

func TestGetRate(t *testing.T) {
  rateFile, err := os.Open("/testing.json")
  if err != nil {
    fmt.Println(err)
  }
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
  t8, _ := time.Parse(time.RFC3339, "2015-07-01T19:00:01-05:00")
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
