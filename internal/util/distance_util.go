package util

import "math"

const earthRadius = 6371.0

func HaversineBounds(lat, long, distance float64) (minLat, maxLat, minLong, maxLong float64) {
	latRange := distance / earthRadius
	latRadians := lat * (math.Pi / 180.0)
	longRadians := distance / (earthRadius * math.Cos(latRadians))

	minLat = lat - latRange*(180.0/math.Pi)
	maxLat = lat + latRange*(180.0/math.Pi)
	minLong = long - longRadians*(180.0/math.Pi)
	maxLong = long + longRadians*(180.0/math.Pi)
	return
}
