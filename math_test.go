package gountries

import "fmt"

func ExamplePythagorasEquirectangular() {

	se, _ := query.FindCountryByAlpha("SWE")
	de, _ := query.FindCountryByAlpha("DEU")

	distance := MeasureDistancePythagoras(se, de)

	fmt.Println(distance)

	// Output:
	// 1430.5502701671583

}

func ExampleHaversine() {

	se, _ := query.FindCountryByAlpha("SWE")
	de, _ := query.FindCountryByAlpha("DEU")

	distance := MeasureDistanceHaversine(se, de)

	fmt.Println(distance)

	// Output:
	// 1430.1937864547901

}
