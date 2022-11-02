# Cluster

[![CircleCI](https://circleci.com/gh/aliakseiz/gocluster/tree/master.svg?style=svg)](https://circleci.com/gh/aliakseiz/gocluster/tree/master)
[![Go Report](https://goreportcard.com/badge/github.com/aliakseiz/gocluster)](https://goreportcard.com/report/github.com/aliakseiz/gocluster)
[![Go Doc](https://godoc.org/github.com/aliakseiz/gocluster?status.svg)](http://godoc.org/github.com/aliakseiz/gocluster)

`gocluster` is a very fast Golang library for geospatial point clustering. Benchmarks available below.

This is a fork of [GoCluster](https://github.com/MadAppGang/gocluster). Which is basically a port of Mapbox [supercluster](https://github.com/mapbox/supercluster/). 

Additional features comparing to the origin:
    
    - Correct grouping when the view covers both hemispheres
	- Method to obtain a cluster expansion zoom
	- Refactored implementation
    - Improved test coverage
	- Google maps example

![clusters2](https://cloud.githubusercontent.com/assets/25395/11857351/43407b46-a40c-11e5-8662-e99ab1cd2cb7.gif)

The cluster uses hierarchical greedy clustering approach.

Usage example:

```go
package main

import (
  "github.com/aliakseiz/gocluster"
  "log"
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
    log.Fatal(err)
  }

  // Get tour tile with mercator coordinate projections to display directly on the map
  result := c.GetTile(0, 0, 0)
  // or get all clusters for zoom 10 without points count limit
  results := c.AllClusters(10, -1)

  log.Printf("%v\n", result)
  log.Printf("%v\n", results)
}
```

All IDs of `Point` that you have as result are the index of initial array of Geopoint, so you could get you point by
this index.

## Init cluster index

To init the index, points should be prepared first. All points should implement `GeoPoint` interface:

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

To search all points inside the box, that are limited by the box, formed by north-west point and east-south points. 
Z index should be provided as well.

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

OSM and Google maps [uses tiles system](https://developers.google.com/maps/documentation/javascript/maptypes#TileCoordinates) to
optimize map loading. So it is possible to get all points for the tile with tileX, tileY and zoom:

```go
c := NewCluster(geoPoints)
tileX := 0
tileY := 1
zoom := 4

results := c.GetTile(tileX, tileY, zoom)
```

In this case all coordinates are returned in pixels for that tile. To retrieve objects with Lat, Lng,
`GetTileWithLatLng` method should be used.

## Test data

Testdata in `testdata` directory is based on [GeoJSON](https://en.wikipedia.org/wiki/GeoJSON) format.

## Benchmarks

