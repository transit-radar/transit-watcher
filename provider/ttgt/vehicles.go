package ttgt

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type ListTransitVehiclesParams struct {
	RouteID   string
	VariantID string
}

func (c *client) ListTransitVehicles(ctx context.Context, params ListTransitVehiclesParams) (*TransitVehicleResponse, error) {
	q := url.Values{
		"routeId": []string{params.RouteID},
		"varId":   []string{params.VariantID},
	}

	url := c.baseURL
	url.Path = PathTransitVehicles
	url.RawQuery = q.Encode()

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

	var response TransitVehicleResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
