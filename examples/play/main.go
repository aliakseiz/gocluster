package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aliakseiz/gocluster"
)

type testPoint struct {
	Type       string
	Properties struct {
		Name string
	}
	Geometry struct {
		Coordinates []float64
	}
}

func (tp *testPoint) GetCoordinates() *cluster.GeoCoordinates {
	if tp.Geometry.Coordinates == nil {
		return nil
	}

	return &cluster.GeoCoordinates{
		Lng: tp.Geometry.Coordinates[0],
		Lat: tp.Geometry.Coordinates[1],
	}
}

func (tp *testPoint) GetID() int64 {
	return 0
}

func importData(filename string) []*testPoint {
	points := struct {
		Type     string
		Features []*testPoint
	}{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	json.Unmarshal(raw, &points)
	return points.Features
}

type simplePoint struct {
	Lng, Lat float64
}

func (sp simplePoint) GetCoordinates() *cluster.GeoCoordinates {
	return &cluster.GeoCoordinates{Lng: sp.Lng, Lat: sp.Lat}
}

func (sp simplePoint) GetID() int64 {
	return 0
}

func main() {
	points := importData("../../testdata/places.json")
	geoPoints := make([]cluster.GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c, _ := cluster.New(geoPoints,
		cluster.WithinZoom(0, 3),
		cluster.WithPointSize(60),
		cluster.WithTileSize(256))
	southEast := simplePoint{71.36718750000001, -83.79204408779539}
	northWest := simplePoint{-71.01562500000001, 83.7539108491127}
	result := c.GetClusters(northWest, southEast, 0, -1)
	fmt.Printf("Getting points: %+v\n length %v \n", result, len(result))
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(resultJSON))
}
