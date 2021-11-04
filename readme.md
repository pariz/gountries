# gountries

[![wercker status](https://app.wercker.com/status/909d6a059d7d0b49b74ec8b658f97df4/s/master "wercker status")](https://app.wercker.com/project/bykey/909d6a059d7d0b49b74ec8b658f97df4) [![codecov.io](https://codecov.io/github/pariz/gountries/coverage.svg?branch=master)](https://codecov.io/github/pariz/gountries?branch=master) [![Go Report Card](https://goreportcard.com/badge/pariz/gountries)](https://goreportcard.com/report/pariz/gountries)

Inspired by the [countries](https://github.com/hexorx/countries) gem for ruby.

Countries (ISO-3166-1), Country Subdivisions(ISO-3166-2), Currencies (ISO 4217), Geo Coordinates(ISO-6709) as well as translations, country borders and other stuff exposed as struct data.

All data is derived from the [pariz/countries](https://github.com/pariz/countries) repo.

*This is currently a work in progress, so things may change. More stuff will be added*

# Installation

```
go get github.com/pariz/gountries
```

# Examples


## Basic
```go
package main

import (
	"fmt"
	"github.com/pariz/gountries"
)

func main() {
	query := gountries.New()

	/////////////////
	// Find sweden //
	/////////////////

	sweden, _ := query.FindCountryByName("sweden")
	// sweden, _ := query.FindCountryByAlpha("SE")
	// sweden, _ := query.FindCountryByAlpha("SWE")

	fmt.Println(sweden.Name.Common)   // Output: Sweden
	fmt.Println(sweden.Name.Official) // Output: Konungariket Sverige

	fmt.Println(sweden.Translations["DEU"].Common)   // Output: Schweden
	fmt.Println(sweden.Translations["DEU"].Official) // Output: Königreich Schweden
}

```
## A bit more advanced
```go
package main

import (
	"fmt"
	"github.com/pariz/gountries"
)

func main() {

	query := gountries.New()

	////////////////////////////////////////////
	// Find the bordering countries of Sweden //
	////////////////////////////////////////////

	sweden, _ := query.FindCountryByAlpha("SWE") // "SE" also works..

	// Get the bordering countries of sweden
	for _, country := range sweden.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

	// Output:
	// Finland
	// Norway

	////////////////////////////////////
	// Find all subdivisons for Sweden //
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

	countryQuery := gountries.Country{
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

	///////////////////////////////////////////////////////////////////
	// Calculate distance between Sweden and Germany (in Kilometers) //
	///////////////////////////////////////////////////////////////////

	se, _ := query.FindCountryByAlpha("SWE")
	de, _ := query.FindCountryByAlpha("DEU")

	distance := gountries.MeasureDistanceHaversine(se, de)
	//distance := MeasureDistancePythagoras(se, de)

	fmt.Println(distance)

	// Output:
	// 1430.1937864547901

	distance = gountries.CalculateHaversine(
		se.Coordinates.MaxLatitude, se.Coordinates.MaxLongitude,
		de.Coordinates.MinLatitude, de.Coordinates.MinLongitude)

	fmt.Println(distance)

	// Output:
	// 2641.26145088825
}

```

# Using packed data

The data in the `data/yaml` subdirectory is embedded using go-bindata.  Once you include this library in your project, you won't need to access the data directory.  To add or update the data, make changes to the YAML files then run:

```
go-bindata -pkg gountries data/yaml/*
```

# Testing

Has a pretty solid test coverage but is constantly improving.

# Todo
- [ ] Province/County querying (partially complete)
- [x] Measurement between coordinates
- [ ] GeoJSON information
- [ ] Suggestions?
