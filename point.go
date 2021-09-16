package cluster

import (
	"github.com/electrious-go/kdbush"
)

// Point struct that implements clustered points
// could have only one point or set of points.
type Point struct {
	X, Y        float64
	zoom        int
	ID          int // Index for pint, Id for cluster
	NumPoints   int
	Included    []int64
	Descendants []int64
}

// Coordinates to be compatible with interface.
func (cp *Point) Coordinates() (float64, float64) {
	return cp.X, cp.Y
}

// IsCluster tells you if this point is cluster or
// rather regular point.
func (cp *Point) IsCluster(c *Cluster) bool {
	return cp.ID >= c.clusterIdxSeed
}

// GeoCoordinates represent position in the Earth.
type GeoCoordinates struct {
	Lng float64
	Lat float64
}

// GeoPoint interface returning lat/lng coordinates.
// All object, that you want to cluster should implement this protocol.
type GeoPoint interface {
	GetID() int64
	GetCoordinates() GeoCoordinates
}

// translate geopoints to Points witrh projection coordinates.
func translateGeoPointsToPoints(points []GeoPoint) []*Point {
	result := make([]*Point, len(points))
	for i, p := range points {
		cp := Point{}
		cp.zoom = InfinityZoomLevel
		cp.X, cp.Y = MercatorProjection(p.GetCoordinates())
		result[i] = &cp
		cp.NumPoints = 1
		cp.ID = i
		cp.Included = []int64{p.GetID()}
		cp.Descendants = []int64{p.GetID()}
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
