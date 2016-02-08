package gountries

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var query *Query

func TestMain(m *testing.M) {
	var err error

	if err != nil {
		panic(err)
	}
	query = New()
	flag.Parse()
	os.Exit(m.Run())
}

func ExampleBorderingCountries() {

	se, _ := query.FindCountryByAlpha("SWE")
	for _, country := range se.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

	// Output:
	// Finland
	// Norway

}

func ExampleTranslations() {

	se, _ := query.FindCountryByAlpha("SWE")
	fmt.Println(se.Translations["DEU"].Common)

	// Output:
	// Schweden
}
