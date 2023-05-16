package lpr

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func parsePullAddress(body string) (string, error) {
	address, err := stringInBetween(body, "<wsa:Address>", "</wsa:Address>")
	if err != nil {
		return "", errors.New("<wsa:Address> tag is not found in response")
	}
	return address, nil
}

func parseRecognition(body string) *Recognition {
	if !strings.Contains(body, "PlateNumber") {
		return nil
	}

	data, err := stringInBetween(body, "<tt:Data>", "</tt:Data>")
	if err != nil {
		return nil
	}

	split := strings.SplitAfter(data, ">")
	split = split[:len(split)-1]

	now := time.Now()

	rec := &Recognition{
		Timestamp: now,
	}

	for _, line := range split {
		name, value, err := parseItem(line)
		if err != nil {
			return nil
		}

		switch name {
		case "PlateNumber":
			rec.LicencePlate = value
		case "Likelihood":
			intValue, err := strconv.Atoi(value)
			if err == nil {
				if intValue > 100 {
					rec.Confidence = intValue / 10
				} else {
					rec.Confidence = intValue
				}
			}
		case "Nation":
			rec.Nation = value
		case "Country":
			rec.Country = value
		case "VehicleDirection":
			switch value {
			case "reverse":
				rec.Direction = Leaving
			case "forward":
				rec.Direction = Approaching
			default:
				rec.Direction = Unknown
			}
		}
	}

	if rec.LicencePlate == "" {
		return nil
	}

	return rec
}

func parseItem(line string) (string, string, error) {
	startIndex := strings.Index(line, "<tt:SimpleItem Name=\"")
	if startIndex == -1 {
		return "", "", errors.New("invalid input")
	}

	endIndex := strings.Index(line[startIndex:], "\" Value=\"")
	if endIndex == -1 {
		return "", "", errors.New("invalid input")
	}

	name := line[startIndex+len("<tt:SimpleItem Name=\"") : startIndex+endIndex]

	startIndex = startIndex + endIndex + len("\" Value=\"")
	endIndex = strings.Index(line[startIndex:], "\"/>")
	if endIndex == -1 {
		return "", "", errors.New("invalid input")
	}

	value := line[startIndex : startIndex+endIndex]
	return name, value, nil
}

func stringInBetween(source, start, end string) (string, error) {
	startIndex := strings.Index(source, start)
	if startIndex == -1 {
		return "", errors.New("no string found")
	}

	endIndex := strings.Index(source[startIndex:], end)
	if endIndex == -1 {
		return "", errors.New("no string found")
	}

	return source[startIndex+len(start) : startIndex+endIndex], nil
}
