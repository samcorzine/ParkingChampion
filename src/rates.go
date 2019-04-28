package main


import (
  "time"
  "strings"
  "strconv"
  "sync"
  "errors"
)

var validDays = [...]string {"sun", "mon", "tues", "wed", "thurs", "fri", "sat"}

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


func dayInValidDays(day string) bool {
    for _, b := range validDays {
        if b == day {
            return true
        }
    }
    return false
}



func (rateList *RateList) update(newRates *RateList) error {
  for _, x := range newRates.Rates {
    if !x.validate(){
      return errors.New("New rates contains an invalid rate")
    }
  }
  rateList = newRates
  return nil
}

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



func timeFromMil(queryTime time.Time, milTime string, tz string) time.Time {
  y, m, d := queryTime.Date()
  milHour, _ := strconv.Atoi(milTime[:2])
  milMinute, _ := strconv.Atoi(milTime[2:])
  location, _ := time.LoadLocation(tz)
  theTime := time.Date(y, m, d, milHour, milMinute, 0, 0, location)
  return theTime
}

func (rate *Rate) timeMatch(rq RateQuery) bool {
  start := rq.startTime
  end := rq.endTime
  days := strings.Split(rate.Days, ",")
  var daysAsInts []int
  for _, x := range days {
    day := -1
    switch x {
      case "sun":
        day = 0
      case "mon":
        day = 1
      case "tues":
        day = 2
      case "wed":
        day = 3
      case "thurs":
        day = 4
      case "fri":
        day = 5
      case "sat":
        day = 6
    }
    if day != -1 {
      daysAsInts = append(daysAsInts, day)
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
