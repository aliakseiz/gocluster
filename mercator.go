package cluster

import "math"

// MercatorProjection will convert lat,lng into spherical mercator range which is 0 to 1
func MercatorProjection(coordinates GeoCoordinates) (float64, float64) {
	x := coordinates.Lon/360.0 + 0.5
	sin := math.Sin(coordinates.Lat * math.Pi / 180.0)
	y := (0.5 - 0.25*math.Log((1+sin)/(1-sin))/math.Pi)
	if y < 0 {
		y = 0
	}
	if y > 1 {
		y = 1
	}
	return x, y
}

// ReverseMercatorProjection converts spherical mercator range to lat,lng
func ReverseMercatorProjection(x, y float64) GeoCoordinates {
	result := GeoCoordinates{}
	result.Lon = (x - 0.5) * 360
	y2 := (180 - y*360) * math.Pi / 180.0
	result.Lat = 360*math.Atan(math.Exp(y2))/math.Pi - 90
	return result
}
