# gountries

The shameless port of the countries gem to golang (well kinda).

This go package provides country information via a simple API.

All data is derived from the [pariz/countries](https://github.com/pariz/countries) repo.

*This is currently a work in progress, so things may change. More stuff will be added*

# Installation

```
go get github.com/pariz/gountries
```

# Examples


## Basic
```go


  import (
    "github.com/pariz/gountries"
    "fmt"
  )


  query := gountries.New()

  /////////////////
  // Find sweden //
  /////////////////

  sweden, _ := query.FindCountryByName("sweden")
  // sweden, _ := query.FindCountryByAlpha("SE")
  // sweden, _ := query.FindCountryByAlpha("SWE")

  fmt.Println(sweden.Name.Common) // Output: Sweden
  fmt.Println(sweden.Name.Official) // Output: Konungariket Sverige

  fmt.Println(sweden.Translations["deu"].Common) // Output: Schweden
  fmt.Println(sweden.Translations["deu"].Official) // Output: Königreich Schweden


```
## A bit more advanced
```go

  import (
    "github.com/pariz/gountries"
    "fmt"
  )

  query := gountries.New()

  ////////////////////////////////////////////
  // Find the bordering countries of Sweden //
  ////////////////////////////////////////////

  sweden, _ := query.FindCountryByAlpha("swe")

  // Get the bordering countries of sweden
  for _, country := range sweden.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

  // Output:
  // Finland
  // Norway

  ////////////////////////////////////
  // Find all subdivions for Sweden //
  ////////////////////////////////////

  subdivisions := sweden.SubDivisions()

	for _, subdivision := range subdivisions {
		fmt.Println(subdivision.Name)
	}

  // Output:
  // Västerbottens län
  // Uppsala län
  // Södermanlands län
  // Gotlands län
  // Dalarnas län
  // ...

  //////////////////////////////////////////////////////////
  // Find all countries bordering Germany and Switzerland //
  //////////////////////////////////////////////////////////

  countryQuery := Country{
		Borders: []string{
			"DEU",
			"CHE",
		},
	}

	countries := query.FindCountries(countryQuery)

	for _, c := range countries {
		fmt.Println(c.Name.Common)
	}

	// Output:
	// Austria
	// France

```

# Testing

Has a pretty solid test coverage but is constantly improving.

# Todo


* (partially complete) Province/County querying
* Measurement between coordinates
* GeoJSON information
* Suggestions?
