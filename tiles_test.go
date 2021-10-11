package cluster_test

import (
	"fmt"
	cluster "github.com/aliakseiz/gocluster"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_GetTile00(t *testing.T) {
	points := importData("./testdata/places.json")
	assert.NotEmptyf(t, points, "no points for clustering")

	geoPoints := make([]cluster.GeoPoint, len(points))

	for i := range points {
		geoPoints[i] = points[i]
	}

	c, _ := cluster.New(geoPoints,
		cluster.WithinZoom(0, 3),
		cluster.WithPointSize(60),
		cluster.WithTileSize(256),
		cluster.WithNodeSize(64))

	result := c.GetTile(0, 0, 0)
	assert.NotEmpty(t, result)

	expectedPoints := importPoints("./testdata/expect_tile0_0_0.json")
	// Included field value is tested separately
	for i := range result {
		result[i].Included = nil
	}
	assert.Equal(t, expectedPoints, result)
}

// validate original result from JS library.
func TestCluster_GetTileDefault(t *testing.T) {
	points := importData("./testdata/places.json")
	if len(points) == 0 {
		t.Error("Getting empty test data")
	} else {
		t.Logf("Getting %v points to test\n", len(points))
	}
	geoPoints := make([]cluster.GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c, _ := cluster.New(geoPoints)
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
	geoPoints := make([]cluster.GeoPoint, len(points))

	for i := range points {
		geoPoints[i] = points[i]
	}

	c, _ := cluster.New(geoPoints,
		cluster.WithinZoom(0, 3),
		cluster.WithPointSize(60),
		cluster.WithTileSize(256),
		cluster.WithNodeSize(64))
	result := c.GetTile(0, 0, 4)

	fmt.Printf("%+v", result)
	// Output: [{X:-3350 Y:253 zoom:0 ID:22 NumPoints:1 Included:[0]} {X:-2418 Y:165 zoom:0 ID:62 NumPoints:1 Included:[0]}]
}
