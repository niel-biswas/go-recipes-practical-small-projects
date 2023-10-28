// What is the maximal ride speed in rides.json and by which car?
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const format = "2006-01-02T15:04"

func parseTime(ts string) (time.Time, error) {
	return time.Parse(format, ts)
}

type ride struct {
	StartTime string  `json:"start"`
	EndTime   string  `json:"end"`
	ID        string  `json:"id"`
	Distance  float64 `json:"distance"`
}

func NewRide(st, et, id string, dt float64) *ride {
	r := ride{
		StartTime: st,
		EndTime:   et,
		ID:        id,
		Distance:  dt,
	}
	return &r
}

func (r *ride) getSpeed(startTime, endTime time.Time) (float64, error) {
	dt := endTime.Sub(startTime)
	dtHour := float64(dt) / float64(time.Hour)
	if dtHour == 0 {
		return 0, errors.New("Time cannot be zero for a given distance to calculate speed")
	}
	speed := r.Distance / dtHour
	return speed, nil
}

func maxRideSpeed(r io.Reader) (string, float64, error) {
	var maxSpeedCarID string
	var maxSpeed float64
	dec := json.NewDecoder(r)
	for {
		rideInfo := NewRide("", "", "", 0)
		err := dec.Decode(&rideInfo)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", 0, err
		}
		st, err := parseTime(rideInfo.StartTime)
		et, err := parseTime(rideInfo.EndTime)
		if Speed, err := rideInfo.getSpeed(st, et); err == nil {
			if Speed > maxSpeed {
				maxSpeed = Speed
				maxSpeedCarID = rideInfo.ID
			}
		}
	}
	return maxSpeedCarID, maxSpeed, nil
}

func main() {
	file, err := os.Open("rides.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	id, speed, err := maxRideSpeed(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s drives with a maximum speed of %f", id, speed)
}
