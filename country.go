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

type Country2 struct {
	Name    CountryName
	Demonym string

	TLDs []string

	Languages    map[string]string
	Translations map[string]BaseLang

	Alpha2 string `json:"cca2"`
	Alpha3 string `json:"cca3"`
	CIOC   string

	Capital string

	Borders []string

	Geo CountryGeo

	LandLocked bool

	Area int64

	// Yaml
	//

	Currency            string
	CountryCode         int
	InternationalPrefix string
	Continent           string
	EuMember            bool
	Region              string
	SubRegion           string
}

type CountryGeo struct {
	Longitude string
	Latitude  string

	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
	LatitudeF    float64
	LongitudeF   float64
}

// Country contains all countries and their country codes
type Country struct {
	Name    CountryName
	Demonym string

	Code        string
	TLD         []string
	Alpha2      string `json:"cca2"`
	Alpha3      string `json:"cca3"`
	CIOC        string
	Currency    []string
	CallingCode []string

	Capital   string
	Region    string
	SubRegion string

	Languages    map[string]string
	Translations map[string]BaseLang
	AltSpellings []string

	LatLng []float64

	LandLocked bool

	Borders []string // Bordering countries as slice of CCA3 strings

	Area int64

	// Yaml
	//

	Geo CountryGeo

	//Currency            string
	CountryCode         int
	InternationalPrefix string
	Continent           string
	EuMember            bool
	//Region              string
	//SubRegion           string
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
