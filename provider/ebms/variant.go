package ebms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *client) ListVariants(ctx context.Context, routeID string) ([]Variant, error) {
	url := c.baseURL
	url.Path = fmt.Sprintf(PathGetVarsByRoute, routeID)

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

	variants := []Variant{}

	err = json.Unmarshal(body, &variants)
	if err != nil {
		return nil, err
	}

	return variants, nil
}
