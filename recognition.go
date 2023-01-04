package lpr

type Direction string

const (
	Approaching Direction = "Approaching"
	Leaving     Direction = "Leaving"
	Unknown     Direction = "Unknown"
)

type Recognition struct {
	LicencePlate string
	Confidence   int
	Nation       string
	Country      string
	Direction    Direction
}

func (recognition *Recognition) IsDirectionViolated() bool {
	switch recognition.Direction {
	case Leaving:
		return true
	default:
		return false
	}
}
