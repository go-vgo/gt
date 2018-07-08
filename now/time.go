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

// Add add some time
func Add(begin time.Time, tm string) (end time.Time) {
	day, _ := time.ParseDuration(tm)
	end = begin.Add(day)
	return
}

// Timestamp get now timestamp
func Timestamp() int64 {
	return time.Now().Unix()
}

// TimestampNano get now timestamp
func TimestampNano() int64 {
	return time.Now().UnixNano()
}

// GetMonTS get monday timestamp
func GetMonTS() (int64, error) {
	weekTm := now.BeginningOfWeek().Format(Format)
	theTime, err := time.ParseInLocation(Format, weekTm, time.Local)
	if err != nil {
		return 0, err
	}

	day, _ := time.ParseDuration("24h")
	theTime = theTime.Add(day)

	timestamp := theTime.Unix()
	return timestamp, nil
}
