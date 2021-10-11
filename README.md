# Cluster

[![CircleCI](https://circleci.com/gh/aliakseiz/gocluster/tree/master.svg?style=svg)](https://circleci.com/gh/aliakseiz/gocluster/tree/master)
[![Go Report](https://goreportcard.com/badge/github.com/aliakseiz/gocluster)](https://goreportcard.com/report/github.com/aliakseiz/gocluster)
[![Go Doc](https://godoc.org/github.com/aliakseiz/gocluster?status.svg)](http://godoc.org/github.com/aliakseiz/gocluster)

`gocluster` is a very fast Golang library for geospatial point clustering.

The origin of this library is in [GoCluster](https://github.com/MadAppGang/gocluster).

This fork has few additional features.

	- Method to obtain expansion zoom
	- Google maps example
	- Refactored implementation
    - Improved test coverage

![clusters2](https://cloud.githubusercontent.com/assets/25395/11857351/43407b46-a40c-11e5-8662-e99ab1cd2cb7.gif)

The cluster uses hierarchical greedy clustering approach. The same approach used by Dave Leaver in his
Leaflet.markercluster plugin.

So this approach is extremely fast, the only drawback is that all clustered points are stored in memory.

This library is deeply inspired by MapBox's superclaster JS library and blog
post: https://www.mapbox.com/blog/supercluster/

Easy to use:

```go
package main

import (
	"fmt"
	"github.com/aliakseiz/gocluster"
)

func main() {
	var points []*cluster.Point
	// Convert slice of your objects to slice of GeoPoint (interface) objects
	geoPoints := make([]cluster.GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	// Create new cluster (this will build index)
	c, err := cluster.New(geoPoints, cluster.WithinZoom(0, 21))
	if err != nil {
		fmt.Println(err)

		return
	}

	// Get tour tile with mercator coordinate projections to display directly on the map
	result := c.GetTile(0, 0, 0)
	// or get all clusters for zoom 10 without points count limit
	results := c.AllClusters(10, -1)

	fmt.Printf("%v\n", result)
	fmt.Printf("%v\n", results)
}
```

All IDs of `Point` that you have as result are the index of initial array of Geopoint, so you could get you point by
this index.

## Init cluster index

To init index, you need to prepare your data. All your points should implement `GeoPoint` interface:

```go
type GeoPoint interface {
	GetCoordinates() GeoCoordinates
}

type GeoCoordinates struct {
	Lng float64
	Lat float64
}
```

The `Cluster` could be tweaked:

|parameter | default value | description |
|---|---|---|
|MinZoom | 0 | Minimum zoom level at which clusters are generated |
|MaxZoom | 16 | Minimum zoom level at which clusters are generated |
|PointSize | 40 | Cluster radius, in pixels |
|TileSize | 512 | Tile extent. Radius is calculated relative to this value |
|NodeSize | 64 | NodeSize is size of the KD-tree node. Higher means faster indexing but slower search, and vise versa. |

Available option functions:

```go
WithPointSize(size int) Option
WithTileSize(size int) Option
WithinZoom(min, max int) Option
WithNodeSize(size int) Option

// Creating new cluster
New(points []GeoPoint, opts ...Option) (*Cluster, error)
```

## Search point in boundary box

To search all points inside the box, that are limited by the box, formed by north-west point and east-south points. You
need to provide Z index as well.

```go

northWest := simplePoint{71.36718750000001, -83.79204408779539}
southEast := simplePoint{-71.01562500000001, 83.7539108491127}
zoom := 2

results := c.GetClusters(northWest, southEast, zoom)
```

Returns the array of 'ClusterPoint' for zoom level. Each point has the following coordinates:

* X coordinate of returned object is Longitude
* Y coordinate of returned object is Latitude
* if the object is cluster of points (NumPoints > 1), the ID is generated started from ClusterIdxSeed (ID>
  ClusterIdxSeed)
* if the object represents only one point, it's id is the index of initial GeoPoints array

## Search points for tile

OSM and Google
maps [uses tiles system](https://developers.google.com/maps/documentation/javascript/maptypes#TileCoordinates) to
optimize map loading. So you could get all points for the tile with tileX, tileY and zoom:

```go
c := NewCluster(geoPoints)
tileX := 0
tileY := 1
zoom := 4

results := c.GetTile(tileX, tileY, zoom)
```

In this case all coordinates are returned in pixels for that tile. If you want to return objects with Lat, Lng,
use `GetTileWithLatLng` method.

## Test data

Testdata in `testdata` directory is based on [GeoJSON](https://en.wikipedia.org/wiki/GeoJSON) format.
