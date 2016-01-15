# gountries

The shameless port of the countries gem to golang (well kinda).

This go package provides country information via a simple API.

*This is currently a work in progress, so things may change. More stuff will be added*

# Examples


```go

  import (
    "github.com/pariz/gountries/gountries"
  )

  gountries := NewGountries()

  // Find sweden
  sweden, _ := gountries.FindCountryByName("sweden")

  // Get the bordering countries of sweden
  for _, country := range sweden.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

  // Output:
  // Finland
  // Norway


```

The complete country struct can be found here:

```go

type Country struct {
	Name        CountryName
	Code        string
	TLD         []string
	CCA3        string
	CIOC        string
	Currency    []string
	CallingCode []string

	Capital      string
	AltSpellings []string

	Region    string
	SubRegion string

	Languages    map[string]string
	Translations map[string]BaseLang

	LatLng []float64

	Demonym string

	LandLocked bool

	Borders []string

	Area int64
}

```

# Testing

Has pretty much full test coverage.

# Todo

* Province/County selection
* GeoJSON information
* Suggestions?
