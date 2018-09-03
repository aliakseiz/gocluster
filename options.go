package cluster

// Option allows to modify cluster properties or cluster itself
type Option func(*Cluster) error

// cluster := &Cluster{
// 	MinZoom:   0,
// 	MaxZoom:   16,
// 	PointSize: 40, // 240
// 	TileSize:  512,
// 	NodeSize:  64, // 128
// 	Radius:    40,
// 	Extent:    512,
// }

// WithPointSize will set point size
func WithPointSize(size int) Option {
	return func(c *Cluster) error {
		c.PointSize = size
		return nil
	}
}

// WithTileSize will set tile size.
// TileSize = 512 (GMaps and OSM default)
func WithTileSize(size int) Option {
	return func(c *Cluster) error {
		c.TileSize = size
		return nil
	}
}

// WithinZoom will set min/max zoom
func WithinZoom(min, max int) Option {
	return func(c *Cluster) error {
		c.MinZoom = min
		c.MaxZoom = max
		return nil
	}
}

// WithNodeSize will set node size.
func WithNodeSize(size int) Option {
	return func(c *Cluster) error {
		c.NodeSize = size
		return nil
	}
}
