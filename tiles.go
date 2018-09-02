package cluster

import "github.com/electrious/kdbush"

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
