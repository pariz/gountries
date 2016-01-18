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

// Country contains all countries and their country codes
type Country struct {
	Name    CountryName
	Demonym string

	Code string
	TLDs []string

	Languages    map[string]string
	Translations map[string]BaseLang

	Alpha2 string `json:"cca2"`
	Alpha3 string `json:"cca3"`
	CIOC   string

	Capital string

	Borders []string

	LandLocked bool

	Area int64

	// Yaml
	//

	Currency            string `yaml:"currency"`
	CountryCode         int    `yaml:"country_code"`
	InternationalPrefix string `yaml:"international_prefix"`
	Continent           string `yaml:"continent"`
	EuMember            bool   `yaml:"eu_member"`
	Region              string `yaml:"region"`
	SubRegion           string `yaml:"subregion"`

	Longitude string `yaml:"longitude"`
	Latitude  string `yaml:"latitude"`

	MinLongitude float64 `yaml:"min_longitude"`
	MinLatitude  float64 `yaml:"min_latitude"`
	MaxLongitude float64 `yaml:"max_longitude"`
	MaxLatitude  float64 `yaml:"max_latitude"`
	LatitudeF    float64 `yaml:"latitude_dec"`
	LongitudeF   float64 `yaml:"longitude_dec"`
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
