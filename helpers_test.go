package cluster_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	cluster "github.com/aliakseiz/gocluster"
)

type simplePoint struct {
	ID       int64
	Lon, Lat float64
}

func (sp simplePoint) GetID() int64 {
	return sp.ID
}

func (sp simplePoint) GetCoordinates() *cluster.GeoCoordinates {
	return &cluster.GeoCoordinates{Lng: sp.Lon, Lat: sp.Lat}
}

// TestPoint structure to import GeoJSON test data.
type TestPoint struct {
	ID         int64
	Type       string
	Properties properties
	Geometry   geometry
}

type properties struct {
	Name       string
	PointCount int `json:"point_count"`
}

type geometry struct {
	Coordinates []float64
}

func (tp *TestPoint) GetID() int64 {
	return tp.ID
}

func (tp *TestPoint) GetCoordinates() *cluster.GeoCoordinates {
	if tp.Geometry.Coordinates == nil {
		return nil
	}

	return &cluster.GeoCoordinates{
		Lng: tp.Geometry.Coordinates[0],
		Lat: tp.Geometry.Coordinates[1],
	}
}

type GeoJSONResultFeature struct {
	Geometry [][]float64
	Tags     struct {
		PointCount int `json:"point_count"`
	}
}

func importData(filename string) []*TestPoint {
	points := struct {
		Type     string
		Features []*TestPoint
	}{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())

		return nil
	}
	err = json.Unmarshal(raw, &points)
	if err != nil {
		fmt.Println(err.Error())

		return nil
	}

	return points.Features
}

func importPoints(filename string) []cluster.Point {
	var result []cluster.Point

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())

		return nil
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		fmt.Println(err.Error())

		return nil
	}

	return result
}

func importGeoJSONResultFeature(filename string) []GeoJSONResultFeature {
	points := struct {
		Features []GeoJSONResultFeature
	}{}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())

		return nil
	}

	err = json.Unmarshal(raw, &points)
	if err != nil {
		fmt.Println(err.Error())

		return nil
	}

	return points.Features
}

const epsilon = 0.0000000001

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}

	return false
}
