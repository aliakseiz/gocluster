package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kr/pretty"

	"github.com/electrious/cluster"
)

type testPoint struct {
	Type       string
	Properties struct {
		//we don't need other data
		Name string
	}
	Geometry struct {
		Coordinates []float64
	}
}

func (tp *testPoint) GetCoordinates() cluster.GeoCoordinates {
	return cluster.GeoCoordinates{
		Lon: tp.Geometry.Coordinates[0],
		Lat: tp.Geometry.Coordinates[1],
	}
}

//type MercatorPoint struct {
//	Cluster cluster.ClusterPoint
//	MercatorX int
//	MercatorY int
//}
//
//func mercator(p cluster.ClusterPoint) MercatorPoint {
//	mp := MercatorPoint{}
//	mp.Cluster = p
//	mp.MercatorX =
//
//}

func importData(filename string) []*testPoint {
	var points = struct {
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
	Lon, Lat float64
}

func (sp simplePoint) GetCoordinates() cluster.GeoCoordinates {
	return cluster.GeoCoordinates{Lon: sp.Lon, Lat: sp.Lat}
}

func main() {
	points := importData("../../testdata/places.json")

	c := cluster.NewCluster()
	geoPoints := make([]cluster.GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c.PointSize = 60
	c.MaxZoom = 3
	c.TileSize = 256
	//c.NodeSize = 64
	southEast := simplePoint{71.36718750000001, -83.79204408779539}
	northWest := simplePoint{-71.01562500000001, 83.7539108491127}
	c.WithPoints(geoPoints)

	result := c.GetClusters(northWest, southEast, 0)
	pretty.Println(c.GetClusterExpansionZoom(32001))
	//result = c.GetTile(0,0,0)
	fmt.Printf("Getting points: %+v\n length %v \n", result, len(result))

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(resultJSON))

}
