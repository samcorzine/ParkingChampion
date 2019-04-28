package main


import (
  "time"
  "strings"
  "strconv"
)

type Rate struct {
	Days         string `json:"days"`
	Times        string `json:"times"`
	TimeZone     string `json:"tz"`
	Price        int    `json:"price"`
}

type RateList struct {
  Rates       []Rate `json:"rates"`
}

type RateQuery struct {
  startTime   time.Time
  endTime     time.Time
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
      beginningR := rateTimes[0]
      endR := rateTimes[1]
      beginningHour, _ := strconv.Atoi(beginningR[:2])
      beginningMinute, _ := strconv.Atoi(beginningR[2:])
      endHour, _ := strconv.Atoi(endR[:2])
      endMinute, _ := strconv.Atoi(endR[2:])
      // TODO fix seconds logic
      if start.Hour() > beginningHour {
        if end.Hour() < endHour {
          return true
        }
        if end.Hour() == endHour && end.Minute() <= endMinute {
          return true
        }
      }
      if start.Hour() == beginningHour && start.Minute() >= beginningMinute {
        if end.Hour() < endHour {
          return true
        }
        if end.Hour() == endHour && end.Minute() <= endMinute {
          return true
        }
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
