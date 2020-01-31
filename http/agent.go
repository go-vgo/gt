package http

import (
	"math/rand"
	"time"
)

var (
	userAgent = [...]string{
		"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
		"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
		"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GetRandomUserAgent get random UserAgent
func GetRandomUserAgent(args ...[]string) string {
	if len(args) > 0 {
		return userAgent[r.Intn(len(args[0]))]
	}

	return userAgent[r.Intn(len(userAgent))]
}
