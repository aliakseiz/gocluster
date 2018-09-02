package cluster

import (
	"math"

	"github.com/electrious/kdbush"
)

const (
	// InfinityZoomLevel indicate impossible large zoom level (Cluster's max is 21)
	InfinityZoomLevel = 100
	// NoParent means there is no higher cluster in tree
	NoParent = -1
)

// Cluster struct get a list or stream of geo objects
// and produce all levels of clusters
// MinZoom - minimum  zoom level to generate clusters
// MaxZoom - maximum zoom level to generate clusters
// Zoom range is limited by 0 to 21, and MinZoom could not be larger, then MaxZoom
// PointSize - pixel size of marker, affects clustering radius
// TileSize - size of tile in pixels, affects clustering radius
type Cluster struct {
	MinZoom   int
	MaxZoom   int
	PointSize int
	TileSize  int
	NodeSize  int
	Radius    int
	Extent    int
	Indexes   []*kdbush.KDBush
	Points    []GeoPoint

	ClusterIdxSeed int
}

// NewCluster create new Cluster instance with default parameters:
// MinZoom = 0
// MaxZoom = 16
// PointSize = 40
// TileSize = 512 (GMaps and OSM default)
// NodeSize is size of the KD-tree node, 64 by default. Higher means faster indexing but slower search, and vise versa.
func NewCluster() *Cluster {
	return &Cluster{
		MinZoom:   0,
		MaxZoom:   16,
		PointSize: 40, // 240
		TileSize:  512,
		NodeSize:  64, // 128
		Radius:    40,
		Extent:    512,
	}
}

// WithPoints get points and create multilevel clustered indexes
// All points should implement GeoPoint interface
// they are not copied, so you could not worry about memory efficiency
// And GetCoordinates called only once for each object, so you could calc it on the fly, if you need
func (c *Cluster) WithPoints(points []GeoPoint) error {
	//limit max Zoom
	if c.MaxZoom > 21 {
		c.MaxZoom = 21
	}
	//adding extra layer for infinite zoom (initial) layers data storage
	c.Indexes = make([]*kdbush.KDBush, c.MaxZoom-c.MinZoom+2)
	c.Points = points
	//get digits number, start from next exponent
	//if we have 78, all cluster will start from 100...
	//if we have 986 points, all clusters ids will start from 1000
	c.ClusterIdxSeed = int(math.Pow(10, float64(digitsCount(len(points)))))
	clusters := translateGeoPointsToPoints(points)
	for z := c.MaxZoom; z >= c.MinZoom; z-- {
		//create index from clusters from previous iteration
		c.Indexes[z+1] = kdbush.NewBush(clustersToPoints(clusters), c.NodeSize)
		//create clusters for level up using just created index
		clusters = c.clusterize(clusters, z)
	}
	//index topmost points
	c.Indexes[c.MinZoom] = kdbush.NewBush(clustersToPoints(clusters), c.NodeSize)
	return nil
}

// GetClusters returns the array of clusters for zoom level.
// The northWest and southEast points are boundary points of square, that should be returned.
// northWest is left topmost point.
// southEast is right bottom point.
// return the object for clustered points,
// X coordinate of returned object is Longitude and
// Y coordinate of returned object is Latitude
func (c *Cluster) GetClusters(northWest, southEast GeoPoint, zoom int) []Point {
	index := c.Indexes[c.limitZoom(zoom)]
	nwX, nwY := MercatorProjection(northWest.GetCoordinates())
	seX, seY := MercatorProjection(southEast.GetCoordinates())
	ids := index.Range(nwX, nwY, seX, seY)
	result := make([]Point, len(ids))
	for i := range ids {
		p := index.Points[ids[i]].(*Point)
		cp := *p
		coordinates := ReverseMercatorProjection(cp.X, cp.Y)
		cp.X = coordinates.Lon
		cp.Y = coordinates.Lat
		result[i] = cp
	}

	return result
}

// GetClustersPointsInRadius will return child points for specific cluster
// this is done with kdbush.Within method allowing fast search
func (c *Cluster) GetClustersPointsInRadius(clusterID int) []*Point {
	// if clusterID is smaller than initial seed
	// it means that it is original point from which
	// cluster(s) are made
	if clusterID < c.ClusterIdxSeed {
		return []*Point{}
	}
	originIndex := (clusterID >> 5) - c.ClusterIdxSeed
	originZoom := (clusterID % 32) - 1
	originTree := c.Indexes[originZoom]
	originPoint := originTree.Points[originIndex]
	r := float64(c.Radius) / (float64(c.Extent) * math.Pow(2.0, float64(originZoom)))
	// r := 200.
	treeBelow := c.Indexes[originZoom+1]
	ids := treeBelow.Within(originPoint, r)
	children := []*Point{}
	for _, i := range ids {
		c := treeBelow.Points[i].(*Point)
		if c.ParentID != clusterID {
			continue
		}
		children = append(children, c)
	}
	return children
}

// GetClusterExpansionZoom will return how much you need to zoom
// to get to a next cluster
func (c *Cluster) GetClusterExpansionZoom(clusterID int) *int {
	clusterZoom := (clusterID % 32) - 1
	id := clusterID
	for clusterZoom < int(c.MaxZoom) {
		children := c.GetClustersPointsInRadius(id)
		if len(children) == 0 {
			return nil
		}
		clusterZoom++

		// in case it's more then 1, then return current zoom
		if len(children) != 1 {
			break
		}
		id = children[0].ID
	}
	return &clusterZoom
}

// AllClusters returns all cluster points, array of Point,  for zoom on the map.
// X coordinate of returned object is Longitude and.
// Y coordinate of returned object is Latitude.
func (c *Cluster) AllClusters(zoom int) []Point {
	index := c.Indexes[c.limitZoom(zoom)]
	points := index.Points
	result := make([]Point, len(points))
	for i := range points {
		p := index.Points[i].(*Point)
		cp := *p
		coordinates := ReverseMercatorProjection(cp.X, cp.Y)
		cp.X = coordinates.Lon
		cp.Y = coordinates.Lat
		result[i] = cp
	}

	return result
}

// GetTile return points for  Tile with coordinates x and y and for zoom z
// return objects with pixel coordinates
func (c *Cluster) GetTile(x, y, z int) []Point {
	return c.getTile(x, y, z, false)
}

// GetTileWithLatLon return points for  Tile with coordinates x and y and for zoom z
// return objects with LatLon coordinates
func (c *Cluster) GetTileWithLatLon(x, y, z int) []Point {
	return c.getTile(x, y, z, true)
}

func (c *Cluster) getTile(x, y, z int, latlon bool) []Point {
	index := c.Indexes[c.limitZoom(z)]
	z2 := 1 << uint(z)
	z2f := float64(z2)
	extent := c.TileSize
	r := c.PointSize
	p := float64(r) / float64(extent)
	top := (float64(y) - p) / z2f
	bottom := (float64(y) + 1 + p) / z2f
	resultIds := index.Range(
		(float64(x)-p)/z2f,
		float64(top),
		(float64(x)+1+p)/z2f,
		bottom,
	)
	var result []Point
	if latlon == true {
		result = c.pointIDToLatLonPoint(resultIds, index.Points)
	} else {
		result = c.pointIDToMercatorPoint(resultIds, index.Points, float64(x), float64(y), z2f)
	}
	if x == 0 {
		minX1 := float64(1-p) / z2f
		minY1 := float64(top)
		maxX1 := 1.0
		maxY1 := float64(bottom)
		resultIds = index.Range(minX1, minY1, maxX1, maxY1)
		var sr1 []Point
		if latlon == true {
			sr1 = c.pointIDToLatLonPoint(resultIds, index.Points)
		} else {
			sr1 = c.pointIDToMercatorPoint(resultIds, index.Points, z2f, float64(y), z2f)
		}
		result = append(result, sr1...)

	}
	if x == (z2 - 1) {
		minX2 := 0.0
		minY2 := float64(top)
		maxX2 := float64(p) / z2f
		maxY2 := float64(bottom)
		resultIds = index.Range(minX2, minY2, maxX2, maxY2)
		var sr2 []Point
		if latlon == true {
			sr2 = c.pointIDToLatLonPoint(resultIds, index.Points)
		} else {
			sr2 = c.pointIDToMercatorPoint(resultIds, index.Points, -1, float64(y), z2f)
		}
		result = append(result, sr2...)
	}
	return result
}

// calc Point mercator projection regarding tile
func (c *Cluster) pointIDToMercatorPoint(ids []int, points []kdbush.Point, x, y, z2 float64) []Point {
	var result []Point
	for i := range ids {
		p := points[ids[i]].(*Point)
		cp := *p
		//translate our coordinate system to mercator
		cp.X = float64(round(float64(c.TileSize) * (p.X*z2 - x)))
		cp.Y = float64(round(float64(c.TileSize) * (p.Y*z2 - y)))
		cp.zoom = 0
		result = append(result, cp)
	}
	return result
}

func (c *Cluster) pointIDToLatLonPoint(ids []int, points []kdbush.Point) []Point {
	result := make([]Point, len(ids))
	for i := range ids {
		p := points[ids[i]].(*Point)
		cp := *p
		coordinates := ReverseMercatorProjection(cp.X, cp.Y)
		cp.X = coordinates.Lon
		cp.Y = coordinates.Lat
		result[i] = cp
	}
	return result
}

// clusterize points for zoom level
func (c *Cluster) clusterize(points []*Point, zoom int) []*Point {
	var result []*Point
	r := float64(c.PointSize) / float64(c.TileSize*(1<<uint(zoom)))
	// iterate all clusters
	for pi := range points {
		// skip points we have already clustered
		p := points[pi]
		if p.zoom <= zoom {
			continue
		}
		// mark this point as visited
		p.zoom = zoom
		// find all neighbours
		tree := c.Indexes[zoom+1]
		neighbourIds := tree.Within(&kdbush.SimplePoint{X: p.X, Y: p.Y}, r)
		nPoints := p.NumPoints
		wx := p.X * float64(nPoints)
		wy := p.Y * float64(nPoints)
		var foundNeighbours []*Point
		for j := range neighbourIds {
			b := points[neighbourIds[j]]
			// filter out neighbours, that are already processed (and processed point "p" as well)
			if zoom < b.zoom {
				wx += b.X * float64(b.NumPoints)
				wy += b.Y * float64(b.NumPoints)
				nPoints += b.NumPoints
				b.zoom = zoom //set the zoom to skip in other iterations
				foundNeighbours = append(foundNeighbours, b)
			}
		}
		newCluster := p
		// create new cluster
		if len(foundNeighbours) > 0 {
			newCluster = &Point{}
			newCluster.X = wx / float64(nPoints)
			newCluster.Y = wy / float64(nPoints)
			newCluster.NumPoints = nPoints
			newCluster.zoom = InfinityZoomLevel
			newCluster.ParentID = NoParent
			// create ID based on seed + index
			// this is then shifted to create space for zoom
			// this is useful when you need extract zoom from ID
			newCluster.ID = ((c.ClusterIdxSeed + pi) << 5) + zoom + 1
			for _, p := range foundNeighbours {
				p.ParentID = newCluster.ID
			}
		}
		result = append(result, newCluster)
	}
	return result
}

func (c *Cluster) limitZoom(zoom int) int {
	if zoom > c.MaxZoom+1 {
		zoom = c.MaxZoom + 1
	}
	if zoom < c.MinZoom {
		zoom = c.MinZoom
	}
	return zoom
}

////////// End of Cluster implementation

/////////////////////////////////
// private stuff
/////////////////////////////////

//translate geopoints to Points witrh projection coordinates
func translateGeoPointsToPoints(points []GeoPoint) []*Point {
	var result = make([]*Point, len(points))
	for i, p := range points {
		cp := Point{}
		cp.zoom = InfinityZoomLevel
		cp.X, cp.Y = MercatorProjection(p.GetCoordinates())
		result[i] = &cp
		cp.NumPoints = 1
		cp.ID = i
		cp.ParentID = NoParent
	}
	return result
}

func clustersToPoints(points []*Point) []kdbush.Point {
	result := make([]kdbush.Point, len(points))
	for i, v := range points {
		result[i] = v
	}
	return result
}
