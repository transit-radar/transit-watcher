package ttgt

type TransitVehicleResponse struct {
	Vehicles []TransitVehicle `json:"vehicles"`
}

type TransitVehicle struct {
	Angle      float32    `json:"angle"`
	Coordinate [2]float64 `json:"coordinate"`
	Percent    float32    `json:"percent"`
}
