package gobus

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	AllDataUrl   = "https://api.gobus.vn/transit/data/getAllData"
	StopsDataUrl = "https://api.gobus.vn/transit/stops/geojson"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client}
}

func (c *Client) GetStops(ctx context.Context) ([]Stop, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, StopsDataUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("status is not ok, cannot get data from URL")
	}

	body, err := io.ReadAll(res.Body)

	stops := &StopResponse{}

	err = json.Unmarshal(body, stops)
	if err != nil {
		return nil, err
	}

	return stops.Stops, nil
}

func (c *Client) GetRoutes(ctx context.Context) ([]Route, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, AllDataUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("status is not ok, cannot get data from URL")
	}

	body, err := io.ReadAll(res.Body)

	routes := []Route{}

	err = json.Unmarshal(body, &routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}
