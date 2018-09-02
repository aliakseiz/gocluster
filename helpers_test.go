package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

////Helpers
type simplePoint struct {
	Lon, Lat float64
}

func (sp simplePoint) GetCoordinates() GeoCoordinates {
	return GeoCoordinates{sp.Lon, sp.Lat}
}

type TestPoint struct {
	Type       string
	Properties struct {
		//we don't need other data
		Name       string
		PointCount int `json:"point_count"`
	}
	Geometry struct {
		Coordinates []float64
	}
}

type GeoJSONResultFeature struct {
	Geometry [][]float64
	Tags     struct {
		PointCount int `json:"point_count"`
	}
}

func (tp *TestPoint) GetCoordinates() GeoCoordinates {
	return GeoCoordinates{
		Lon: tp.Geometry.Coordinates[0],
		Lat: tp.Geometry.Coordinates[1],
	}
}

func importData(filename string) []*TestPoint {
	var points = struct {
		Type     string
		Features []*TestPoint
	}{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	json.Unmarshal(raw, &points)
	//fmt.Printf("Get data: %+v\n",points)
	return points.Features
}

func importPoints(filename string) []Point {
	var result []Point
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	json.Unmarshal(raw, &result)
	return result

}

func importGeoJSONResultFeature(filename string) []GeoJSONResultFeature {
	var points = struct {
		Features []GeoJSONResultFeature
	}{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	json.Unmarshal(raw, &points)
	return points.Features
}
