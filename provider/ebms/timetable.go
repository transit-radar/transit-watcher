package ebms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *client) ListTimetables(ctx context.Context, routeID string) ([]Timetable, error) {
	url := c.baseURL
	url.Path = fmt.Sprintf(PathGetTimetableByRoute, routeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("status is not ok, cannot get data from URL")
	}

	body, err := io.ReadAll(res.Body)

	timetables := []Timetable{}

	err = json.Unmarshal(body, &timetables)
	if err != nil {
		return nil, err
	}

	return timetables, nil
}

func (c *client) ListTrips(ctx context.Context, routeID, timetableID string) ([]Trip, error) {
	url := c.baseURL
	url.Path = fmt.Sprintf(PathGetTripsByTimetable, routeID, timetableID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("status is not ok, cannot get data from URL")
	}

	body, err := io.ReadAll(res.Body)

	trips := []Trip{}

	err = json.Unmarshal(body, &trips)
	if err != nil {
		return nil, err
	}

	return trips, nil
}
