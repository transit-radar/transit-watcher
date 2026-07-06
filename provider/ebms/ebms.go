package ebms

import (
	"context"
	"net/http"
	"net/url"
)

const (
	DefaultHost   = "apicms.ebms.vn"
	DefaultScheme = "https"

	PathGetAllRoute         = "/businfo/getallroute"
	PathGetRouteByID        = "/businfo/getroutebyid/%s"
	PathGetVarsByRoute      = "/businfo/getvarsbyroute/%s"
	PathGetTimetableByRoute = "/businfo/gettimetablebyroute/%s"
	PathGetTripsByTimetable = "/businfo/gettripsbytimetable/%s/%s"
	PathGetStopsByVar       = "/businfo/getstopsbyvar/%s/%s"
	PathGetPathsByVar       = "/businfo/getpathsbyvar/%s/%s"
)

var (
	DefaultURL = url.URL{
		Host:   DefaultHost,
		Scheme: DefaultScheme,
	}
)

type Client interface {
	ListRoutes(ctx context.Context) ([]RouteListItem, error)
	GetRouteByID(ctx context.Context, routeID string) (Route, error)
	ListVariants(ctx context.Context, routeID string) ([]Variant, error)
	ListTimetables(ctx context.Context, routeID string) ([]Timetable, error)
	ListTrips(ctx context.Context, routeID, timetableID string) ([]Trip, error)
	ListStops(ctx context.Context, routeID, variantID string) ([]Stop, error)
	ListPaths(ctx context.Context, routeID, variantID string) ([]Path, error)
}

type client struct {
	baseURL    url.URL
	httpClient *http.Client
}

type option func(c *client) error

func NewClient(options ...option) (Client, error) {
	c := &client{
		baseURL:    DefaultURL,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
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
