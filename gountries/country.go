package gountries

import (
	"fmt"
	"strings"
)

// FindCountryByName fincs a country by given name
func (q *QueryHolder) FindCountryByName(name string) (result Country, err error) {

	for _, country := range q.Countries {

		if strings.ToLower(country.Name.Common) == strings.ToLower(name) {
			return country, nil
		}
	}

	return Country{}, fmt.Errorf("Could not find country with name %s", name)

}

// FindCountryByCode fincs a country by given code
func (q *QueryHolder) FindCountryByCode(code string) (result Country, err error) {

	for _, country := range q.Countries {
		var countryCode = country.Code

		// If code is 2 characters, its CCA2, if 3 its CCA3
		if len(code) == 2 {
			countryCode = country.Code
		} else if len(code) == 3 {
			countryCode = country.CCA3
		} else {
			return Country{}, fmt.Errorf("%s is an invalid code format", code)
		}

		if strings.ToLower(countryCode) == strings.ToLower(code) {
			return country, nil
		}
	}

	return Country{}, fmt.Errorf("Could not find country with code %s", code)

}

// Country contains all countries and their country codes
type Country struct {
	Name        CountryName
	Code        string `json:"cca2"`
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

	Borders []string // Bordering countries as slice of CCA3 strings

	Area int64
}

// BorderingCountries gets the bordering countries for this country
func (c *Country) BorderingCountries() (countries []Country) {

	for _, cca3 := range c.Borders {

		if country, err := query.FindCountryByCode(cca3); err == nil {
			countries = append(countries, country)
		}

	}

	return

}

// CountryName contains the common, official and native intepretations of the country name
type CountryName struct {
	BaseLang
	Native map[string]BaseLang
}

// BaseLang is a basic structure for common language formatting in the JSON file
type BaseLang struct {
	Common   string
	Official string
}
