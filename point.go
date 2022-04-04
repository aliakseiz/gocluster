package cluster

import (
	"github.com/electrious-go/kdbush"
)

// Point struct that implements clustered points
// could have only one point or set of points.
type Point struct {
	X, Y      float64
	zoom      int
	ID        int // Index for pint, Id for cluster
	NumPoints int
	Included  []int64
}

// GetID to be compatible with interface.
func (cp *Point) GetID() int64 {
	return int64(cp.ID)
}

// Coordinates to be compatible with interface.
func (cp *Point) Coordinates() (float64, float64) {
	return cp.X, cp.Y
}

// GetCoordinates to be compatible with interface.
func (cp *Point) GetCoordinates() *GeoCoordinates {
	return &GeoCoordinates{
		Lng: cp.X,
		Lat: cp.Y,
	}
}

// IsCluster tells you if this point is cluster or rather regular point.
func (cp *Point) IsCluster(c *Cluster) bool {
	return cp.ID >= c.clusterIdxSeed
}

// GeoCoordinates represent position in the Earth.
type GeoCoordinates struct {
	Lng float64
	Lat float64
}

// GeoPoint interface returning lat/lng coordinates.
// All objects, that you want to cluster should implement this interface.
type GeoPoint interface {
	GetID() int64
	GetCoordinates() *GeoCoordinates
}

// translate geopoints to Points with projection coordinates.
func translateGeoPointsToPoints(points []GeoPoint) []*Point {
	result := make([]*Point, 0, len(points))
	for i, p := range points {
		geoPoint := p.GetCoordinates()
		if geoPoint == nil { // Skip points without coordinates
			continue
		}

		cp := Point{}
		cp.zoom = InfinityZoomLevel
		cp.X, cp.Y = MercatorProjection(*geoPoint) // nil check is above
		result = append(result, &cp)
		cp.NumPoints = 1
		cp.ID = i
		cp.Included = []int64{p.GetID()}
	}
	return result
}

func clustersToPoints(points []*Point) []kdbush.Point {
	result := make([]kdbush.Point, len(points))
	for i, v := range points {
		result[i] = v
	}
	return result
}
