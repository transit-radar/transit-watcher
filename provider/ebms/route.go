package ebms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *client) ListRoutes(ctx context.Context) ([]RouteListItem, error) {
	url := c.baseURL
	url.Path = PathGetAllRoute

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

	routes := []RouteListItem{}

	err = json.Unmarshal(body, &routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (c *client) GetRouteByID(ctx context.Context, routeID string) (Route, error) {
	url := c.baseURL
	url.Path = fmt.Sprintf(PathGetRouteByID, routeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return Route{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Route{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Route{}, errors.New("status is not ok, cannot get data from URL")
	}

	body, err := io.ReadAll(res.Body)

	var route Route

	err = json.Unmarshal(body, &route)
	if err != nil {
		return Route{}, err
	}

	return route, nil
}
