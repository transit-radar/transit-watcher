package ttgt

import (
	"context"
	"net/http"
	"net/url"
)

const (
	DefaultHost   = "api.metrohcm.ttgt.vn"
	DefaultScheme = "https"

	PathTransitVehicles = "/transit/vehicles"
)

var (
	DefaultURL = url.URL{
		Host:   DefaultHost,
		Scheme: DefaultScheme,
	}
)

type Client interface {
	ListTransitVehicles(ctx context.Context, params ListTransitVehiclesParams) (*TransitVehicleResponse, error)
}

type client struct {
	baseURL    url.URL
	httpClient *http.Client
}

type option func(c *client) error

func NewClient(options ...option) Client {
	c := &client{
		baseURL:    DefaultURL,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithHTTPClient(httpClient *http.Client) option {
	return func(c *client) error {
		c.httpClient = httpClient
		return nil
	}
}

func WithDomain(domain string) option {
	return func(c *client) error {
		baseURL, err := url.Parse(domain)
		if err != nil {
			return err
		}

		c.baseURL = *baseURL
		return nil
	}
}
