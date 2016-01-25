# gountries

The shameless port of the countries gem to golang (well kinda).

This go package provides country information via a simple API.

*This is currently a work in progress, so things may change. More stuff will be added*

# Examples


```go

  import (
    "github.com/pariz/gountries"
  )

  query := gountries.New()

  // Find sweden
  sweden, _ := query.FindCountryByName("sweden")

  // Get the bordering countries of sweden
  for _, country := range sweden.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

  // Output:
  // Finland
  // Norway


```

# Testing

Has pretty much full test coverage.

# Todo

* (in progress) Province/County selection
* GeoJSON information
* Suggestions?
