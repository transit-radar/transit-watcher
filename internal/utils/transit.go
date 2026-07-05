package utils

import (
	"errors"
	"strings"
	"time"

	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
)

const (
	OpTimeLayout = "15:04"
)

func FilterTransitRoutes(routes []gobus.Route) ([]gobus.Route, error) {
	filtered := []gobus.Route{}

	for _, route := range routes {
		from, to, err := ParseOperationTime(route.Info.OperationTime)
		if err != nil {
			return nil, err
		}

		// skip non-operating routes
		if time.Now().Before(from) || to.After(time.Now().Add(-2*time.Hour)) {
			continue
		}

		filtered = append(filtered, route)
	}

	return filtered, nil
}

func ParseOperationTime(operationTime string) (from, to time.Time, err error) {
	unparsedTimes := strings.Split(operationTime, " - ")

	if len(unparsedTimes) != 2 {
		// nothing to parse
		return time.Time{}, time.Time{}, errors.New("Operation time is not supported")
	}

	now := time.Now()

	from, err = time.Parse(OpTimeLayout, unparsedTimes[0])

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	from = time.Date(now.Year(), now.Month(), now.Day(), from.Hour(), from.Minute(), from.Second(), 0, now.Location())

	to, err = time.Parse(OpTimeLayout, unparsedTimes[1])

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	to = time.Date(now.Year(), now.Month(), now.Day(), from.Hour(), from.Minute(), from.Second(), 0, now.Location())

	return from, to, nil
}
