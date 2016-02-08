package gountries

import "math"

func Deg2Rad(deg float64) float64 {
	return deg * math.Pi / 180
}

func CalculatePythagorasEquirectangular(lat1, lon1, lat2, lon2 float64) (d float64) {
	lat1 = Deg2Rad(lat1)
	lat2 = Deg2Rad(lat2)
	lon1 = Deg2Rad(lon1)
	lon2 = Deg2Rad(lon2)
	var R = 6371.0 // km
	var x = (lon2 - lon1) * math.Cos((lat1+lat2)/2)
	var y = (lat2 - lat1)
	d = math.Sqrt(x*x+y*y) * R

	// Return Distance in Kilometers
	return
}

func CalculateHaversine(lat1, lon1, lat2, lon2 float64) (d float64) {
	var R = 6372.8 // Earth Radius in Kilometers

	var dLat = Deg2Rad(lat2 - lat1)
	var dLon = Deg2Rad(lon2 - lon1)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(Deg2Rad(lat1))*math.Cos(Deg2Rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d = R * c

	// Return Distance in Kilometers
	return
}
