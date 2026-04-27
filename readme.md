# gountries

[![Go Report Card](https://goreportcard.com/badge/pariz/gountries)](https://goreportcard.com/report/pariz/gountries)

> **Heads up — v1.0.0 is breaking.** The data pipeline has been rebuilt and `bindata.go` was removed. See [CHANGELOG.md](CHANGELOG.md) for the migration. To stay on the old behavior, pin: `go get github.com/pariz/gountries@v0.1.6`.

Countries (ISO-3166-1), Country Subdivisions (ISO-3166-2), Currencies (ISO 4217), Geo Coordinates (ISO-6709), translations, country borders, and related metadata, exposed as Go struct data.

Data is sourced from two actively maintained upstream projects and refreshed weekly:
- [countries/countries](https://github.com/countries/countries) (MIT) — primary schema
- [mledoze/countries](https://github.com/mledoze/countries) (ODbL) — borders, TLDs, translations

Inspired by the [countries](https://github.com/countries/countries) Ruby gem.

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

# How the data is built

Country and subdivision YAML in `data/yaml/` is generated from two git submodules under `third_party/`:

```
third_party/countries           → countries/countries (MIT)
third_party/mledoze-countries   → mledoze/countries (ODbL)
```

To regenerate locally:

```
git submodule update --init --recursive
go run ./cmd/sync-data
```

The data is embedded into the package binary via `//go:embed data/yaml`, so consumers do not need to ship the YAML alongside their app.

A weekly GitHub Action ([`sync-upstream.yml`](.github/workflows/sync-upstream.yml)) bumps both submodules to their latest tags, regenerates the data, runs the test suite, and opens a PR if anything changed.

# Testing

```
go test ./...
```

`TestDataInvariants` runs after every regen and asserts that the country count, key fields, and canary lookups (Sweden, US, Germany) are all sane — guarding against silent upstream regressions.

# License

Go code: MIT (see [LICENSE](LICENSE)).

Embedded data: derivative of the two upstream sources above. Because mledoze data is ODbL-licensed, the `data/yaml/` directory should be treated as a Derivative Database under ODbL. Querying the data through this Go API is not a Derivative Database — it produces a Produced Work — so your application's license is unaffected.
