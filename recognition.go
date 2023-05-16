package lpr

import "time"

type Direction string

const (
	Approaching Direction = "Approaching"
	Leaving     Direction = "Leaving"
	Unknown     Direction = "Unknown"
)

type Recognition struct {
	Timestamp    time.Time `json:"timestamp"`
	LicencePlate string    `json:"licence_plate"`
	Direction    Direction `json:"direction"`
	Confidence   int       `json:"confidence"`
	Nation       string    `json:"nation"`
	Country      string    `json:"country"`
}
