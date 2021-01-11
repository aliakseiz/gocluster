package cluster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCluster(t *testing.T) {
	points := importData("./testdata/places.json")
	if len(points) == 0 {
		t.Error("Getting empty test data")
	} else {
		t.Logf("Getting %v points to test\n", len(points))
	}
	c, _ := New([]GeoPoint{})
	assert.Equal(t, c.MinZoom, 0, "they should be equal")
	assert.Equal(t, c.MaxZoom, 21, "they should be equal")
	assert.Equal(t, c.PointSize, 40, "they should be equal")
	assert.Equal(t, c.TileSize, 512, "they should be equal")
	assert.Equal(t, c.NodeSize, 64, "they should be equal")
}

func TestAllClusters(t *testing.T) {
	var point GeoPoint = simplePoint{-1, 71.36718750000001, -83.79204408779539}
	c, _ := New([]GeoPoint{point})
	p := c.AllClusters(21)[0]
	assert.InDelta(t, p.X, 71.36718750000001, 0.000001)
	assert.InDelta(t, p.Y, -83.79204408779539, 0.000001)
}

func TestCluster_GetClusters(t *testing.T) {
	points := importData("./testdata/places.json")
	if len(points) == 0 {
		t.Error("Getting empty test data")
	} else {
		t.Logf("Getting %v points to test\n", len(points))
	}
	geoPoints := make([]GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c, _ := New(geoPoints,
		WithinZoom(0, 17),
		WithPointSize(40),
		WithTileSize(512),
		WithNodeSize(64))
	southEast := simplePoint{-1, 71.36718750000001, -83.79204408779539}
	northWest := simplePoint{-1, -71.01562500000001, 83.7539108491127}
	result := c.GetClusters(northWest, southEast, 2)
	assert.NotEmpty(t, result)
	expectedPoints := importData("./testdata/cluster.json")
	assert.Equal(t, len(result), len(expectedPoints))
	for i := range result {
		rp := result[i]
		ep := expectedPoints[i]
		assert.True(t, floatEquals(rp.X, ep.Geometry.Coordinates[0]))
		assert.True(t, floatEquals(rp.Y, ep.Geometry.Coordinates[1]))
		if rp.NumPoints > 1 {
			assert.Equal(t, rp.NumPoints, ep.Properties.PointCount)
		}
	}
}

func TestCluster_AllClusters(t *testing.T) {
	points := importData("./testdata/places.json")
	if len(points) == 0 {
		t.Error("Getting empty test data")
	} else {
		t.Logf("Getting %v points to test\n", len(points))
	}
	geoPoints := make([]GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c, _ := New(geoPoints,
		WithinZoom(0, 17),
		WithPointSize(40),
		WithTileSize(512),
		WithNodeSize(64))
	result := c.AllClusters(2)
	assert.NotEmpty(t, result)
	assert.Equal(t, 100, len(result))

}

func ExampleCluster_GetClusters() {
	points := importData("./testdata/places.json")
	geoPoints := make([]GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c, _ := New(geoPoints)
	northWest := simplePoint{-1, -71.01562500000001, 83.7539108491127}
	southEast := simplePoint{-1, 71.36718750000001, -83.79204408779539}
	result := c.GetClusters(northWest, southEast, 2)
	fmt.Printf("%+v", result[:3])
	// Output: [{X:-14.473194953510028 Y:26.157965399212813 zoom:1 ID:107 NumPoints:1} {X:-12.408741828510014 Y:58.16339752811905 zoom:1 ID:159 NumPoints:1} {X:-9.269962828651519 Y:42.928736057812586 zoom:1 ID:127 NumPoints:1}]
}
