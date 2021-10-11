package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MercatorProjection(t *testing.T) {
	c := GeoCoordinates{
		Lng: -79.04411780507252, // 0.2804330060970208
		Lat: 43.08771393436908,  // 0.36711590445377973
	}
	x, y := MercatorProjection(c)
	assert.Equal(t, x, 0.2804330060970208)
	assert.Equal(t, y, 0.36711590445377973)
	c = GeoCoordinates{
		Lng: -62.06181800038502,
		Lat: 5.686896063275327,
	}
	x, y = MercatorProjection(c)
	assert.Equal(t, x, 0.32760606111004165)
	assert.Equal(t, y, 0.4841770650015434)
}

func Test_MercatorReversedProjection(t *testing.T) {
	c := GeoCoordinates{
		Lng: -79.044117805,
		Lat: 43.0877139344,
	}
	reversed := ReverseMercatorProjection(MercatorProjection(c))
	assert.Equal(t, reversed, c)
}
