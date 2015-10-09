package parser

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"
)

var (
	unknownTime = time.Unix(0, 0)
)

func ParseFeed(source []byte, funcs ...func([]byte) (Feed, error)) (Feed, error) {
	var feed Feed
	var err error

	for _, f := range funcs {
		feed, err = f(source)
		if err != nil {
			if _, ok := err.(xml.UnmarshalError); !ok {
				return feed, err
			}
		} else {
			break
		}
	}

	return feed, err
}

func parseDate(date string) (time.Time, error) {
	formats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		http.TimeFormat,
		"Mon, 2 Jan 2006 15:04:05 GMT", // http.TimeFormat, but without leading zeros
	}

	date = strings.TrimSpace(date)
	var err error
	var t time.Time
	for _, f := range formats {
		t, err = time.Parse(f, date)
		if err == nil {
			return t, nil
		}
	}

	return t, err
}
