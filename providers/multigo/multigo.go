package multigo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	GeolocationUrl = "https://multipass-api.golabs.vn/v2/public/busmap/route_bus_gps?regionCode=hcm&routeId=%s&direction=%d"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client}
}

func (c *Client) GetGeolocations(ctx context.Context, routeID string, direction int) ([]Geolocation, error) {
	url := fmt.Sprintf(GeolocationUrl, routeID, direction)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	err = injectCredentials(&req.Header)
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

	geolocations := []Geolocation{}

	err = json.Unmarshal(body, &geolocations)
	if err != nil {
		return nil, err
	}

	return geolocations, nil
}

func injectCredentials(header *http.Header) error {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return err
	}

	epoch := time.Now().UnixMilli()
	epochStr := strconv.FormatInt(epoch, 10)
	secret := epochStr + ":Bus2019M@p_"
	proof := md5.Sum([]byte(secret))

	header.Add("device-id", uuid.String())
	header.Add("epoch", epochStr)
	header.Add("proof", hex.EncodeToString(proof[:]))
	header.Add("client-version", "ios|22")

	return nil
}
