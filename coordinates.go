package cluster

// GeoCoordinates represent position in the Earth
type GeoCoordinates struct {
	Lon float64
	Lat float64
}

// GeoPoint interface returning lat/lng coordinates.
// All object, that you want to cluster should implement this protocol
type GeoPoint interface {
	GetCoordinates() GeoCoordinates
}
