package gountries

import "fmt"

func ExampleCalculatePythagorasEquirectangular() {

	se, _ := query.FindCountryByAlpha("SWE")
	de, _ := query.FindCountryByAlpha("DEU")

	distance := MeasureDistancePythagoras(se, de)

	fmt.Println(distance)

	distance = CalculatePythagorasEquirectangular(
		se.Coordinates.MaxLatitude, se.Coordinates.MaxLongitude,
		de.Coordinates.MinLatitude, de.Coordinates.MinLongitude)

	fmt.Println(distance)

	// Output:
	// 1326.748584888358
	// 2650.281173643059
}

func ExampleCalculateHaversine() {

	se, _ := query.FindCountryByAlpha("SWE")
	de, _ := query.FindCountryByAlpha("DEU")

	distance := MeasureDistanceHaversine(se, de)

	fmt.Println(distance)

	distance = CalculateHaversine(
		se.Coordinates.MaxLatitude, se.Coordinates.MaxLongitude,
		de.Coordinates.MinLatitude, de.Coordinates.MinLongitude)

	fmt.Println(distance)

	// Output:
	// 1326.3589326391889
	// 2626.379314009457

}
