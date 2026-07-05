package multigo

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Geolocation struct {
	Id               json.Number `json:"Id"`
	Degree           json.Number `json:"Deg"`
	Latitude         json.Number `json:"Lat"`
	Longitude        json.Number `json:"Lng"`
	Speed            json.Number `json:"Speed"`
	LicensePlate     string      `json:"VehicleNumber"`
	RouteId          json.Number `json:"RouteId"`
	CurrentStationId json.Number `json:"currentStationId"`
	NextStationId    json.Number `json:"nextStationId"`
	Timestamp        time.Time   `json:"lastUpdateTime"`
}

type ListGeolocationParams struct {
	RegionCode *string
	RouteID    *string
	Direction  *int
}

func (c *client) ListGeolocations(ctx context.Context, params ListGeolocationParams) ([]Geolocation, error) {
	q := url.Values{
		"regionCode": []string{"hcm"},
	}

	if params.RegionCode != nil {
		q.Set("regionCode", *params.RegionCode)
	}

	if params.RouteID != nil {
		q.Set("routeId", *params.RouteID)
	}

	if params.Direction != nil {
		q.Set("direction", strconv.Itoa(*params.Direction))
	}

	url := DefaultURL
	url.Path = PathGeolocation
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	err = injectCredentials(&req.Header)
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

	geolocations := []Geolocation{}

	err = json.Unmarshal(body, &geolocations)
	if err != nil {
		return nil, err
	}

	return geolocations, nil
}
