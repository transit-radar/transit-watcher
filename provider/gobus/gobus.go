package gobus

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/leebenson/conform"
)

const (
	DefaultHost   = "api.gobus.vn"
	DefaultScheme = "https"

	PathAllData = "/transit/data/getAllData"
	PathStops   = "/transit/stops/geojson"
)

var (
	DefaultURL = url.URL{
		Host:   DefaultHost,
		Scheme: DefaultScheme,
	}
)

type Client interface {
	ListStops(context.Context) ([]Stop, error)
	ListRoutes(context.Context) ([]Route, error)
}

type client struct {
	httpClient *http.Client
}

type option func(c *client)

func NewClient(options ...option) Client {
	c := &client{
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithHTTPClient(httpClient *http.Client) option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func (c *client) ListStops(ctx context.Context) ([]Stop, error) {
	url := DefaultURL
	url.Path = PathStops

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

	stops := &StopResponse{}

	err = json.Unmarshal(body, stops)
	if err != nil {
		return nil, err
	}

	err = conform.Strings(stops)
	if err != nil {
		return nil, err
	}

	return stops.Stops, nil
}

func (c *client) ListRoutes(ctx context.Context) ([]Route, error) {
	url := DefaultURL
	url.Path = PathAllData

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

	routes := []Route{}

	err = json.Unmarshal(body, &routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}
