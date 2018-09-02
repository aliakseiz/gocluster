package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kr/pretty"

	"github.com/electrious/cluster"
	"github.com/electrious/cluster/examples/googlemaps/spherand"
)

type testPoint struct {
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

func (tp testPoint) GetCoordinates() cluster.GeoCoordinates {
	return cluster.GeoCoordinates{
		Lon: tp.Geometry.Coordinates[0],
		Lat: tp.Geometry.Coordinates[1],
	}
}

type latlng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (tp latlng) GetCoordinates() cluster.GeoCoordinates {
	return cluster.GeoCoordinates{Lon: tp.Lng, Lat: tp.Lat}
}

type boundingBox struct {
	NW latlng `json:"nw"`
	SE latlng `json:"se"`
}

type payload struct {
	Zoom        int         `json:"zoom"`
	BoundingBox boundingBox `json:"bb"`
}

var c *cluster.Cluster

func main() {
	fmt.Printf("creating random samples\n")
	createSamplesRandomly(5000000)
	fmt.Printf("samples created\n")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/clusters", parseGhPost)
	fmt.Printf("listening to 8080\n")
	http.ListenAndServe(":8080", nil)
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	log.Println("received request")
	log.Println(request.URL.String())
	decoder := json.NewDecoder(request.Body)

	var t payload
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	pretty.Println(t.BoundingBox.NW, t.BoundingBox.SE, t.Zoom)
	points := c.GetClusters(t.BoundingBox.NW, t.BoundingBox.SE, t.Zoom)
	// points := c.AllClusters(t.Zoom)
	data, err := json.Marshal(points)
	if err != nil {
		panic(err)
	}
	rw.Write(data)
}

// lat, lng := spherand.Geographical()

func createSamplesFromJSON() {
	points := importData("../../testdata/places.json")
	c = cluster.NewCluster()
	geoPoints := make([]cluster.GeoPoint, len(points))
	for i := range points {
		geoPoints[i] = points[i]
	}
	c.ClusterPoints(geoPoints)
}

func createSamplesRandomly(num int) {
	c = cluster.NewCluster()
	log.Printf("generating ranom lat/lngs")
	latlngs := make([]cluster.GeoPoint, num)
	for i := range latlngs {
		lat, lng := spherand.Geographical()
		latlngs[i] = latlng{lat, lng}
	}
	geoPoints := make([]cluster.GeoPoint, len(latlngs))
	for i := range latlngs {
		geoPoints[i] = latlngs[i]
	}

	log.Printf("starting clustring")
	c.ClusterPoints(geoPoints)
	log.Printf("clustering done")
}

func importData(filename string) []testPoint {
	var points = struct {
		Type     string
		Features []testPoint
	}{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	json.Unmarshal(raw, &points)
	//fmt.Printf("Gett data: %+v\n",points)
	return points.Features
}
