package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/electrious-go/cluster"
	"github.com/electrious-go/cluster/examples/googlemaps/spherand"
	_ "github.com/electrious-go/cluster/examples/googlemaps/static"
	"github.com/rakyll/statik/fs"
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
		Lng: tp.Geometry.Coordinates[0],
		Lat: tp.Geometry.Coordinates[1],
	}
}

type latlng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (tp latlng) GetCoordinates() cluster.GeoCoordinates {
	return cluster.GeoCoordinates{Lng: tp.Lng, Lat: tp.Lat}
}

type boundingBox struct {
	NW latlng `json:"nw"`
	SE latlng `json:"se"`
}

type payload struct {
	Zoom        int         `json:"zoom"`
	BoundingBox boundingBox `json:"bb"`
}

type payload2 struct {
	ClusterID int `json:"clusterID"`
}

var c *cluster.Cluster

func main() {
	log.Printf("creating random samples\n")
	c = createCluster(1000000)
	log.Printf("samples created\n")
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(statikFS))
	http.HandleFunc("/clusters", clustersEndpoint)
	http.HandleFunc("/zoom", zoomEndpoint)
	log.Printf("listening to 8080\n")
	http.ListenAndServe(":8080", nil)
}

func clustersEndpoint(rw http.ResponseWriter, request *http.Request) {
	log.Println("received request")
	log.Println(request.URL.String())
	decoder := json.NewDecoder(request.Body)
	var t payload
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	points := c.GetClusters(t.BoundingBox.NW, t.BoundingBox.SE, t.Zoom)
	data, err := json.Marshal(points)
	if err != nil {
		panic(err)
	}
	rw.Write(data)
}

func zoomEndpoint(rw http.ResponseWriter, request *http.Request) {
	log.Println("received request")
	log.Println(request.URL.String())
	decoder := json.NewDecoder(request.Body)
	var t payload2
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	zoom := c.GetClusterExpansionZoom(t.ClusterID)
	l := []byte(fmt.Sprintf(`{"zoom": %d}`, zoom))
	rw.Write(l)
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
		cluster.WithPointSize(120))
	if err != nil {
		panic(err)
	}
	log.Printf("clustering done")
	return c
}
