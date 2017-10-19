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
	// 1430.5502701671583
	// 2667.2283097795016
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
	// 1430.1937864547901
	// 2641.26145088825

}
