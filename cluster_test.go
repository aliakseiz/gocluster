package cluster_test

import (
	"fmt"
	"testing"

	cluster "github.com/aliakseiz/gocluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCluster(t *testing.T) {
	points := importData("./testdata/places.json")
	assert.NotEmptyf(t, points, "no points for clustering")

	c, _ := cluster.New([]cluster.GeoPoint{})
	assert.Equal(t, c.MinZoom, 0, "they should be equal")
	assert.Equal(t, c.MaxZoom, 21, "they should be equal")
	assert.Equal(t, c.PointSize, 40, "they should be equal")
	assert.Equal(t, c.TileSize, 512, "they should be equal")
	assert.Equal(t, c.NodeSize, 64, "they should be equal")
}

func TestAllClusters(t *testing.T) {
	var point cluster.GeoPoint = simplePoint{-1, 71.36718750000001, -83.79204408779539}

	c, _ := cluster.New([]cluster.GeoPoint{point})

	p := c.AllClusters(21, -1)[0]
	assert.InDelta(t, p.X, 71.36718750000001, 0.000001)
	assert.InDelta(t, p.Y, -83.79204408779539, 0.000001)
}

func TestCluster_GetClusters(t *testing.T) {
	points := importData("./testdata/places.json")
	assert.NotEmptyf(t, points, "no points for clustering")

	geoPoints := make([]cluster.GeoPoint, len(points))

	for i := range points {
		geoPoints[i] = points[i]
	}

	c, _ := cluster.New(geoPoints,
		cluster.WithinZoom(0, 17),
		cluster.WithPointSize(40),
		cluster.WithTileSize(512),
		cluster.WithNodeSize(64))
	southEast := simplePoint{-1, 71.36718750000001, -83.79204408779539}
	northWest := simplePoint{-1, -71.01562500000001, 83.7539108491127}

	result := c.GetClusters(northWest, southEast, 2, -1)
	assert.NotEmpty(t, result)

	expectedPoints := importData("./testdata/cluster.json")
	require.Equal(t, len(expectedPoints), len(result))

	for i := range result {
		rp := result[i]
		ep := expectedPoints[i]
		assert.True(t, floatEquals(rp.X, ep.GetCoordinates().Lng))
		assert.True(t, floatEquals(rp.Y, ep.GetCoordinates().Lat))
		// Verify points count for clusters only
		if rp.IsCluster(c) {
			assert.Equal(t, rp.NumPoints, ep.Properties.PointCount)
		}
		// Included field is tested separately
	}
}

func TestCluster_CrossingNotCrossing(t *testing.T) {
	points := []*TestPoint{
		{ID: 1, Geometry: geometry{[]float64{-178.989, 0}}},
		{ID: 2, Geometry: geometry{[]float64{-178.990, 0}}},
		{ID: 3, Geometry: geometry{[]float64{-178.9991, 0}}},
		{ID: 4, Geometry: geometry{[]float64{-178.992, 0}}},
	}

	geoPoints := make([]cluster.GeoPoint, len(points))

	for i := range points {
		geoPoints[i] = points[i]
	}

	c, _ := cluster.New(geoPoints,
		cluster.WithinZoom(0, 17),
		cluster.WithPointSize(40),
		cluster.WithTileSize(512),
		cluster.WithNodeSize(64))
	southEast := simplePoint{-1, -177, -10}
	northWest := simplePoint{-1, -179, 10}

	nonCrossing := c.GetClusters(northWest, southEast, 1, -1)
	assert.NotEmpty(t, nonCrossing)

	southEast = simplePoint{-1, -177, -10}
	northWest = simplePoint{-1, 179, 10}

	crossing := c.GetClusters(northWest, southEast, 1, -1)
	assert.NotEmpty(t, crossing)

	assert.EqualValues(t, nonCrossing, crossing)
}

func TestCluster_GetClusters_Included(t *testing.T) {
	northWest := simplePoint{-1, -15.8, 72.8}
	southEast := simplePoint{-1, 46.3, 4.7}
	zoom := 5
	limit := -1

	tests := []struct {
		name     string
		input    []*cluster.Point
		expected []*cluster.Point
	}{
		{
			name: "two points same location, one cluster",
			input: []*cluster.Point{
				{ID: 0, NumPoints: 1, X: 20.8, Y: 52.2},
				{ID: 1, NumPoints: 1, X: 20.8, Y: 52.2},
			},
			expected: []*cluster.Point{
				{ID: 0, NumPoints: 2, X: 20.8, Y: 52.2, Included: []int64{0, 1}},
			},
		},
		{
			name: "two points different location, one cluster",
			input: []*cluster.Point{
				{ID: 0, NumPoints: 1, X: 20.81, Y: 52.21},
				{ID: 1, NumPoints: 1, X: 20.83, Y: 52.23},
			},
			expected: []*cluster.Point{
				{ID: 0, NumPoints: 2, X: 20.82, Y: 52.2200011258, Included: []int64{0, 1}},
			},
		},
		{
			name: "three points different location, one cluster",
			input: []*cluster.Point{
				{ID: 0, NumPoints: 1, X: 20.81, Y: 52.21},
				{ID: 1, NumPoints: 1, X: 20.83, Y: 52.23},
				{ID: 2, NumPoints: 1, X: 20.85, Y: 52.25},
			},
			expected: []*cluster.Point{
				{ID: 0, NumPoints: 3, X: 20.83, Y: 52.23000300333, Included: []int64{0, 1, 2}},
			},
		},
		{
			name: "three points different location, two clusters",
			input: []*cluster.Point{
				{ID: 0, NumPoints: 1, X: 20.81, Y: 52.21},
				{ID: 1, NumPoints: 1, X: 20.83, Y: 52.23},
				{ID: 2, NumPoints: 1, X: 22.00, Y: 54.00},
			},
			expected: []*cluster.Point{
				{ID: 0, NumPoints: 2, X: 20.82, Y: 52.2200011258, Included: []int64{0, 1}},
				{ID: 1, NumPoints: 1, X: 22.00, Y: 54.00, Included: []int64{2}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geoPoints := make([]cluster.GeoPoint, len(tt.input))
			for i := range tt.input {
				geoPoints[i] = tt.input[i]
			}

			c, _ := cluster.New(geoPoints,
				cluster.WithinZoom(0, 17),
				cluster.WithPointSize(60),
				cluster.WithTileSize(512),
				cluster.WithNodeSize(64))
			got := c.GetClusters(northWest, southEast, zoom, limit)

			assert.NotEmptyf(t, got, "no clusters created")
			require.Equalf(t, len(tt.expected), len(got), "expected and result arrays length must be equal")

			for i := range got {
				rp := got[i]
				ep := tt.expected[i]
				assert.Truef(t, floatEquals(ep.X, rp.X), "X coordinates don't match: %2.11f and %2.11f", ep.X, rp.X)
				assert.Truef(t, floatEquals(ep.Y, rp.Y), "Y coordinates don't match: %2.11f and %2.11f", ep.Y, rp.Y)
				assert.Equalf(t, ep.NumPoints, rp.NumPoints, "points count doesn't match")
				assert.Equalf(t, ep.Included, rp.Included, "included points don't match")
			}
		})
	}
}

func TestCluster_AllClusters(t *testing.T) {
	points := importData("./testdata/places.json")
	assert.NotEmptyf(t, points, "no points for clustering")

	geoPoints := make([]cluster.GeoPoint, len(points))

	for i := range points {
		geoPoints[i] = points[i]
	}

	c, _ := cluster.New(geoPoints,
		cluster.WithinZoom(0, 17),
		cluster.WithPointSize(40),
		cluster.WithTileSize(512),
		cluster.WithNodeSize(64))

	result := c.AllClusters(2, -1)
	assert.NotEmpty(t, result)
	assert.Equal(t, 100, len(result))
}

func ExampleCluster_GetClusters() {
	points := importData("./testdata/places.json")
	// var points []*TestPoint

	geoPoints := make([]cluster.GeoPoint, len(points))

	for i := range points {
		geoPoints[i] = points[i]
	}

	c, _ := cluster.New(geoPoints)
	northWest := simplePoint{-1, -71.01562500000001, 83.7539108491127}
	southEast := simplePoint{-1, 71.36718750000001, -83.79204408779539}
	result := c.GetClusters(northWest, southEast, 2, -1)

	fmt.Printf("%+v", result[:3])
	// Output: [{X:-14.473194953510028 Y:26.157965399212813 zoom:1 ID:107 NumPoints:1 Included:[0]} {X:-12.408741828510014 Y:58.16339752811905 zoom:1 ID:159 NumPoints:1 Included:[0]} {X:-9.269962828651519 Y:42.928736057812586 zoom:1 ID:127 NumPoints:1 Included:[0]}]
}
