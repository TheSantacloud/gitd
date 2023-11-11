package adapters

import (
	"fmt"
	"regexp"
	"strconv"
)

type TimeSpan struct {
	Value int
	Unit  string
}

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

func (t TimeSpan) String() string {
	return fmt.Sprintf("%d %s", t.Value, t.Unit)
}

func ValidateDurationString(input string) (bool, error) {
	return durationPattern.MatchString(input), nil
}

var durationPattern = regexp.MustCompile(`^(\d+) (day|week|month)s?$`)
