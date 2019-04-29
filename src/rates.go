package main


import (
  "time"
  "strings"
  "strconv"
  "sync"
  "errors"
)
// valid formats for day strings, as well as their associated integer representations in golang
var validDays = [...]string {"sun", "mon", "tues", "wed", "thurs", "fri", "sat"}
var dayIntMap = map[string]int{
   "sun"  : 0,
   "mon"  : 1,
   "tues" : 2,
   "wed"  : 3,
   "thurs": 4,
   "fri"  : 5,
   "sat"  : 6,
 }

type Rate struct {
	Days         string `json:"days"`
	Times        string `json:"times"`
	TimeZone     string `json:"tz"`
	Price        int    `json:"price"`
}

type RateList struct {
  Rates       []Rate `json:"rates"`
  sync.Mutex
}

type RateQuery struct {
  startTime   time.Time
  endTime     time.Time
}

// checks if a day string is one of the accepted representations
func dayInValidDays(day string) bool {
    for _, b := range validDays {
        if b == day {
            return true
        }
    }
    return false
}

// Updates a ratelist to the rates in the argument ratelist
func (rateList *RateList) update(newRates *RateList) error {
  for _, x := range newRates.Rates {
    if !x.validate() {
      return errors.New("New rates contains an invalid rate")
    }
  }
  rateList.Rates = newRates.Rates
  return nil
}

// validates that a rate has acceptable values
// Timezone is valid in golang
// Days are valid to this api's conventions
// Price is positive
func (rate Rate) validate() bool {
  _, err := time.LoadLocation(rate.TimeZone)
  if err != nil {
    return false
  }
  for _, day := range strings.Split(rate.Days, ",") {
    if !dayInValidDays(day) {
      return false
    }
  }
  if rate.Price < 0 {
    return false
  }
  return true
}

// Takes the times formatted in "military time" from the rate ranges and translates them to time structs on the query day
func timeFromMil(queryTime time.Time, milTime string, tz string) time.Time {
  y, m, d := queryTime.Date()
  milHour, _ := strconv.Atoi(milTime[:2])
  milMinute, _ := strconv.Atoi(milTime[2:])
  location, _ := time.LoadLocation(tz)
  theTime := time.Date(y, m, d, milHour, milMinute, 0, 0, location)
  return theTime
}

// Checks if start and end times fall inside a rate's range
func (rate *Rate) timeMatch(rq RateQuery) bool {
  start := rq.startTime
  end := rq.endTime
  days := strings.Split(rate.Days, ",")
  var daysAsInts []int
  for _, day:= range days {
    if dayInValidDays(day) {
      daysAsInts = append(daysAsInts, dayIntMap[day])
    }
  }
  for _, x := range daysAsInts {
    if int(start.Weekday()) == x {
      rateTimeRange := rate.Times
      rateTimes := strings.Split(rateTimeRange, "-")
      beginTimeRange := timeFromMil(start, rateTimes[0], rate.TimeZone)
      endTimeRange := timeFromMil(start, rateTimes[1], rate.TimeZone)
      if start.After(beginTimeRange) && end.Before(endTimeRange) {
        return true
      }
    }
  }
  return false
}

// Get's the applicable rate for a time, returns -1 no rates apply
func (rateList *RateList) getRate(rq RateQuery) int {
  if !rq.validate() {
    return -1
  }
  for _, x := range rateList.Rates {
    if x.timeMatch(rq) {
      return x.Price
    }
  }
  return -1
}

// Validates that a RateQuery's start and end are the same date
func (rq *RateQuery) validate() bool {
  y0, m0, d0 := rq.startTime.Date()
  y1, m1, d1 := rq.endTime.Date()
  if y0 != y1 {
    return false
  }
  if m0 != m1 {
    return false
  }
  if d0 != d1 {
    return false
  }
  return true
}
