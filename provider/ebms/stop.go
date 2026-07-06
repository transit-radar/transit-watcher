package ebms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *client) ListStops(ctx context.Context, routeID, variantID string) ([]Stop, error) {
	url := c.baseURL
	url.Path = fmt.Sprintf(PathGetStopsByVar, routeID, variantID)

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

	stops := []Stop{}

	err = json.Unmarshal(body, &stops)
	if err != nil {
		return nil, err
	}

	return stops, nil
}

func (c *client) ListPaths(ctx context.Context, routeID, variantID string) ([]Path, error) {
	url := c.baseURL
	url.Path = fmt.Sprintf(PathGetPathsByVar, routeID, variantID)

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

	paths := []Path{}

	err = json.Unmarshal(body, &paths)
	if err != nil {
		return nil, err
	}

	return paths, nil
}
