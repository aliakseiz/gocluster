package cluster

// Point struct that implements clustered points
//could have only one point or set of points
type Point struct {
	X, Y      float64
	zoom      int
	ID        int //Index for pint, Id for cluster
	NumPoints int
	ParentID  int
	//IncludedPoints []int TODO: Implement inclusion of objects
}

// Coordinates to be compatible with interface
func (cp *Point) Coordinates() (float64, float64) {
	return cp.X, cp.Y
}
