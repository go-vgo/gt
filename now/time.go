package now

import (
	"time"

	"github.com/jinzhu/now"
)

// GoTime the different string formats for go dates
const (
	// DefaultFormat       = "2006-01-02 15:04:05"
	Format              = "2006-01-02 15:04:05"
	GoFormat            = "2006-01-02 15:04:05.999999999"
	DateFormat          = "2006-01-02"
	FormattedDateFormat = "Jan 2, 2006"
	TimeFormat          = "15:04:05"
	HourMinuteFormat    = "15:04"
	HourFormat          = "15"
	DayDateTimeFormat   = "Mon, Aug 2, 2006 3:04 PM"
	CookieFormat        = "Monday, 02-Jan-2006 15:04:05 MST"
	RFC822Format        = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC1036Format       = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC2822Format       = "Mon, 02 Jan 2006 15:04:05 -0700"
	RFC3339Format       = "2006-01-02T15:04:05-07:00"
	RSSFormat           = "Mon, 02 Jan 2006 15:04:05 -0700"
)

// Add add a tm time to the begin time
func Add(begin time.Time, tm string) (end time.Time) {
	day, _ := time.ParseDuration(tm)
	end = begin.Add(day)
	return
}

// Timestamp get the timestamp of now
func Timestamp() int64 {
	return time.Now().Unix()
}

// TimestampNano get the nano timestamp of now
func TimestampNano() int64 {
	return time.Now().UnixNano()
}

// GetMonTS get the timestamp of the monday
func GetMonTS(m ...bool) (int64, error) {
	var tm string
	if len(m) > 0 {
		tm = now.BeginningOfMonth().Format(Format)
	} else {
		tm = now.BeginningOfWeek().Format(Format)
	}

	theTime, err := time.ParseInLocation(Format, tm, time.Local)
	if err != nil {
		return 0, err
	}

	day, _ := time.ParseDuration("24h")
	theTime = theTime.Add(day)

	timestamp := theTime.Unix()
	return timestamp, nil
}

// GetMonthTS returns the timestamp of the first day of the month
func GetMonthTS() (int64, error) {
	return GetMonTS(true)
}
