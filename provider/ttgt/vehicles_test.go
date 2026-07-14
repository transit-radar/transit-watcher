package ttgt

import (
	"testing"

	ttgttest "codeberg.org/transit-radar/transit-watcher/test/data/api.metrohcm.ttgt.vn"
	"github.com/stretchr/testify/assert"
)

func TestClient_ListTransitVehicles(t *testing.T) {
	ts := ttgttest.NewServer()
	t.Cleanup(ts.Close)

	tests := []struct {
		name             string
		routeID          string
		variantID        string
		expectedResponse *TransitVehicleResponse
		expectedErr      string
	}{{
		name:      "Line1NorthBound",
		routeID:   "384",
		variantID: "1",
		expectedResponse: &TransitVehicleResponse{Vehicles: []TransitVehicle{
			{Angle: 59, Coordinate: [2]float64{106.7918792141948, 10.860784816320976}, Percent: 0.8650349650349627},
			{Angle: 22, Coordinate: [2]float64{106.75568126626956, 10.809340662664953}, Percent: 0.5438461538461524},
			{Angle: 53, Coordinate: [2]float64{106.71600416500532, 10.796476768344395}, Percent: 0.23749999999999838},
		}},
	}, {
		name:      "Line1SouthBound",
		routeID:   "384",
		variantID: "1",
		expectedResponse: &TransitVehicleResponse{Vehicles: []TransitVehicle{
			{Angle: 59, Coordinate: [2]float64{106.7918792141948, 10.860784816320976}, Percent: 0.8650349650349627},
			{Angle: 22, Coordinate: [2]float64{106.75568126626956, 10.809340662664953}, Percent: 0.5438461538461524},
			{Angle: 53, Coordinate: [2]float64{106.71600416500532, 10.796476768344395}, Percent: 0.23749999999999838},
		}},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(WithDomain(ts.URL), WithHTTPClient(ts.Client()))
			response, err := client.ListTransitVehicles(t.Context(), ListTransitVehiclesParams{
				RouteID:   test.routeID,
				VariantID: test.variantID,
			})

			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResponse, response)
			} else {
				assert.EqualError(t, err, test.expectedErr)
			}
		})
	}
}
