// What is the maximal ride speed in rides.json and by which car with total time taken (in hours) and distance (in kms)?
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

func (r *ride) getSpeedHour(startTime, endTime time.Time) (float64, float64, error) {
	dt := endTime.Sub(startTime)
	dtHour := float64(dt) / float64(time.Hour)
	if dtHour == 0 {
		return 0, 0, errors.New("Time cannot be zero for a given distance to calculate speed")
	}
	speed := r.Distance / dtHour
	return speed, dtHour, nil
}

func maxRideSpeedHourDist(r io.Reader) (string, float64, float64, float64, error) {
	var maxSpeedCarID string
	var maxSpeed float64
	var totalTime float64
	var totalDistance float64
	dec := json.NewDecoder(r)
	for {
		rideInfo := NewRide("", "", "", 0)
		err := dec.Decode(&rideInfo)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", 0, 0, 0, err
		}
		st, err := parseTime(rideInfo.StartTime)
		et, err := parseTime(rideInfo.EndTime)
		if Speed, time2cover, err := rideInfo.getSpeedHour(st, et); err == nil {
			if Speed > maxSpeed {
				maxSpeed = Speed
				maxSpeedCarID = rideInfo.ID
				totalTime = time2cover
				totalDistance = rideInfo.Distance
			}
		}
	}
	return maxSpeedCarID, maxSpeed, totalTime, totalDistance, nil
}

func main() {
	file, err := os.Open("rides.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	id, speed, time, dist, err := maxRideSpeedHourDist(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s drives with a maximum speed of %f covering a distance of %f km(s) in %f hour(s)", id, speed, dist, time)
}
