package main

import (
	"fmt"
	"time"
)

// tsConvert convert time stamp in "YYYY-MM-DDTHH:MM" format from one time zone to another
func tsConvert(ts, from, to string) (string, error) {

	const format = "2006-01-02T15:04"
	// Parsing input timestamp as per given format
	t, errt := time.Parse(format, ts)
	if errt != nil {
		fmt.Printf("error: %s\n", errt)
	}

	// Loading from location timezone
	fromtz, err := time.LoadLocation(from)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	// Deriving time.Date based on fromTZ
	fromtzDateTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), fromtz)

	// Loading to location timezone
	totz, err2 := time.LoadLocation(to)
	if err2 != nil {
		fmt.Printf("err2: %s\n", err2)
	}

	// Deriving totzDateTime based on fromtzDateTime
	totzDateTime := fromtzDateTime.In(totz)
	return totzDateTime.Format(format), nil
}

func main() {
	ts := "2021-03-08T19:12"
	out, err := tsConvert(ts, "America/Los_Angeles", "Asia/Jerusalem")
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}

	fmt.Println(out) // 2021-03-09T05:12
}
