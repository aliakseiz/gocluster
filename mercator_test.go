package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MercatorProjection(t *testing.T) {
	coords := GeoCoordinates{
		Lng: -79.04411780507252, //0.2804330060970208
		Lat: 43.08771393436908,  //0.36711590445377973
	}
	x, y := MercatorProjection(coords)
	assert.Equal(t, x, 0.2804330060970208)
	assert.Equal(t, y, 0.36711590445377973)
	coords = GeoCoordinates{
		Lng: -62.06181800038502,
		Lat: 5.686896063275327,
	}
	x, y = MercatorProjection(coords)
	assert.Equal(t, x, 0.32760606111004165)
	assert.Equal(t, y, 0.4841770650015434)
}
