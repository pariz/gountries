package gountries

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var query *Query

func TestMain(m *testing.M) {
	query = New()
	flag.Parse()
	os.Exit(m.Run())
}

func ExampleCountry_BorderingCountries() {

	se, _ := query.FindCountryByAlpha("SWE")
	for _, country := range se.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

	// Output:
	// Finland
	// Norway

}

func ExampleCountry_translations() {

	se, _ := query.FindCountryByAlpha("SWE")
	fmt.Println(se.Translations["DEU"].Common)

	// Output:
	// Schweden
}
