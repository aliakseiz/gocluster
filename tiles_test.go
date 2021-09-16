package cluster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_GetTile00(t *testing.T) {
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
		WithinZoom(0, 3),
		WithPointSize(60),
		WithTileSize(256),
		WithNodeSize(64))
	result := c.GetTile(0, 0, 0)
	assert.NotEmpty(t, result)
	expectedPoints := importPoints("./testdata/expect_tile0_0_0.json")
	for i := range result {
		result[i].Included = nil
	}
	assert.Equal(t, result, expectedPoints)
}

// validate original result from JS library.
func TestCluster_GetTileDefault(t *testing.T) {
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
	c, _ := New(geoPoints)
	result := c.GetTile(0, 0, 0)
	assert.NotEmpty(t, result)
	expectedPoints := importGeoJSONResultFeature("./testdata/places-z0-0-0.json")
	assert.Equal(t, len(result), len(expectedPoints))
	for i := range result {
		rp := result[i]
		ep := expectedPoints[i]
		assert.Equal(t, rp.X, ep.Geometry[0][0])
		assert.Equal(t, rp.Y, ep.Geometry[0][1])
		if rp.NumPoints > 1 {
			assert.Equal(t, rp.NumPoints, ep.Tags.PointCount)
		}
	}
}

func ExampleCluster_GetTile() {
	points := importData("./testdata/places.json")
	geoPoints := make([]GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c, _ := New(geoPoints,
		WithinZoom(0, 3),
		WithPointSize(60),
		WithTileSize(256),
		WithNodeSize(64))
	result := c.GetTile(0, 0, 4)
	fmt.Printf("%+v", result)
	// Output: [{X:-2418 Y:165 zoom:0 ID:62 NumPoints:1} {X:-3350 Y:253 zoom:0 ID:22 NumPoints:1}]
}
