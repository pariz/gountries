package gountries

import (
	"fmt"
	"strings"
)

// Query holds a reference to the QueryHolder struct
var queryInited = false
var queryInstance *Query

// Query contains queries for countries, cities, etc.
type Query struct {
	Countries    []Country
	Subdivisions map[string][]SubDivision
}

// FindCountryByName fincs a country by given name
func (q *Query) FindCountryByName(name string) (result Country, err error) {

	for _, country := range q.Countries {

		if strings.ToLower(country.Name.Common) == strings.ToLower(name) {
			return country, nil
		}
	}

	return Country{}, fmt.Errorf("Could not find country with name %s", name)

}

// FindCountryByCode fincs a country by given code
func (q *Query) FindCountryByAlpha(code string) (result Country, err error) {

	for _, country := range q.Countries {
		var countryCode = country.Alpha2

		// If code is 2 characters, its CCA2, if 3 its CCA3
		if len(code) == 2 {
			countryCode = country.Alpha2
		} else if len(code) == 3 {
			countryCode = country.Alpha3
		} else {
			return Country{}, fmt.Errorf("%s is an invalid code format", code)
		}

		if strings.ToLower(countryCode) == strings.ToLower(code) {
			return country, nil
		}
	}

	return Country{}, fmt.Errorf("Could not find country with code %s", code)

}

func (q Query) FindCountries(c Country) (countries []Country) {

	for _, country := range q.Countries {

		// Name
		//

		if c.Name.Common != "" && strings.ToLower(c.Name.Common) == strings.ToLower(country.Name.Common) {
			continue
		}

		// Alpha
		//

		if c.Alpha2 != "" && c.Alpha2 != country.Alpha2 {
			continue
		}

		if c.Alpha3 != "" && c.Alpha3 != country.Alpha3 {
			continue
		}

		// Geo
		//

		if c.Geo.Continent != "" && strings.ToLower(c.Geo.Continent) != strings.ToLower(country.Geo.Continent) {
			continue
		}

		if c.Geo.Region != "" && strings.ToLower(c.Geo.Region) != strings.ToLower(country.Geo.Region) {
			continue
		}

		if c.Geo.SubRegion != "" && strings.ToLower(c.Geo.SubRegion) != strings.ToLower(country.Geo.SubRegion) {
			continue
		}

		// Misc
		//

		if c.InternationalPrefix != "" && strings.ToLower(c.InternationalPrefix) != strings.ToLower(country.InternationalPrefix) {
			continue
		}

		// Bordering countries
		//

		allMatch := false

		if len(c.BorderingCountries()) > 0 {

			for _, c1 := range c.BorderingCountries() {

				match := false

				for _, c2 := range country.BorderingCountries() {
					match = c1.Alpha2 == c2.Alpha2

					if match == true {
						break
					}
				}

				if match == true {
					allMatch = true
				} else {
					allMatch = false
					break
				}

			}

			if allMatch == false {
				continue
			}

		}

		// Append if all matches
		//

		countries = append(countries, country)

	}

	return
}
