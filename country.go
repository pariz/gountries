package gountries

import (
	"fmt"
	"strings"
)

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
func (q *Query) FindCountryByCode(code string) (result Country, err error) {

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
	Name struct {
		BaseLang `yaml:",inline"`
		Native   map[string]BaseLang
	} `json:"name"`

	EuMember    bool
	LandLocked  bool
	Nationality string

	//Code         string

	TLDs []string `json:"tld"`

	Languages    map[string]string
	Translations map[string]BaseLang

	Currencies []string `json:"currency"`

	Borders []string

	// Grouped meta
	Codes
	Geo
	Coordinates
}

// BorderingCountries gets the bordering countries for this country
func (c *Country) BorderingCountries() (countries []Country) {

	query := New()

	for _, cca3 := range c.Borders {

		if country, err := query.FindCountryByCode(cca3); err == nil {
			countries = append(countries, country)
		}

	}

	return

}

// BaseLang is a basic structure for common language formatting in the JSON file
type BaseLang struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

// SubDivision contains Country subdivison information
type SubDivision struct {
	Name  string
	Names []string
	Code  string

	CountryAlpha2 string

	Coordinates
}

// Geo contains geographical information
type Geo struct {
	Region    string `json:"region"`
	SubRegion string `json:"subregion"`
	Continent string // Yaml
	Capital   string `json:"capital"`

	Area float64
}

// Codes contains various code representations
type Codes struct {
	Alpha2 string `json:"cca2"`
	Alpha3 string `json:"cca3"`
	CIOC   string
	CCN3   string

	//CountryCode         string // Taml
	CallingCodes []string `json:"callingCode"`

	InternationalPrefix string // Yaml
}

// Measurer provides coordinates for measurements
type Measurer interface {
	Coordinates() (minLong, minLat, maxLong, maxLat float64)
}

// Coordinates contains the coordinates for both Country and SubDivision
type Coordinates struct {
	LongitudeString string `json:"longitude"`
	LatitudeString  string `json:"latitude"`

	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
	Latitude     float64
	Longitude    float64
}
