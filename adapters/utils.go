package adapters

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type TimeSpan struct {
	Value int
	Unit  string
}

var durationPattern = regexp.MustCompile(`^(\d+) (minute|hour|day|week|month)s?$`)

func NewTimeSpan(input string) (*TimeSpan, error) {
	isValid, err := ValidateDurationString(input)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, fmt.Errorf("invalid duration string: %s", input)
	}

	matches := durationPattern.FindStringSubmatch(input)
	value, _ := strconv.Atoi(matches[1])
	unit := matches[2]

	return &TimeSpan{
		Value: value,
		Unit:  unit,
	}, nil
}

func (t TimeSpan) ModifyDate(date time.Time, add bool) time.Time {
	multiplier := 1
	if !add {
		multiplier = -1
	}
	switch t.Unit {
	case "minute":
		return date.Add(time.Minute * time.Duration(multiplier*t.Value))
	case "hour":
		return date.Add(time.Hour * time.Duration(multiplier*t.Value))
	case "day":
		return date.AddDate(0, 0, multiplier*t.Value)
	case "week":
		return date.AddDate(0, 0, multiplier*t.Value*7)
	case "month":
		return date.AddDate(0, multiplier*t.Value, 0)
	default:
		return date
	}
}

func (t TimeSpan) String() string {
	return fmt.Sprintf("%d %s", t.Value, t.Unit)
}

func ValidateDurationString(input string) (bool, error) {
	return durationPattern.MatchString(input), nil
}
