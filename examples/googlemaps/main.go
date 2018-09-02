package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	log.Printf("creating random samples\n")
	c = createCluster(1000000)
	log.Printf("samples created\n")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/clusters", parseGhPost)
	log.Printf("listening to 8080\n")
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
	// pretty.Println(t.BoundingBox.NW, t.BoundingBox.SE, t.Zoom)
	points := c.GetClusters(t.BoundingBox.NW, t.BoundingBox.SE, t.Zoom)
	data, err := json.Marshal(points)
	if err != nil {
		panic(err)
	}
	rw.Write(data)
}

func createCluster(num int) *cluster.Cluster {
	log.Printf("generating random lat/lng")
	coords := make([]cluster.GeoPoint, num)
	for i := range coords {
		lat, lng := spherand.Geographical()
		coords[i] = latlng{lat, lng}
	}
	geoPoints := make([]cluster.GeoPoint, len(coords))
	for i := range coords {
		geoPoints[i] = coords[i]
	}
	log.Printf("starting clustering")
	c, err := cluster.New(coords,
		cluster.WithNodeSize(64),
		cluster.WithPointSize(240))
	if err != nil {
		panic(err)
	}
	log.Printf("clustering done")
	return c
}
