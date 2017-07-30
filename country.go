package gountries

import (
	"strings"
)

// Country contains all countries and their country codes
type Country struct {
	Name struct {
		BaseLang `yaml:",inline"`
		Native   map[string]BaseLang
	} `json:"name"`

	EuMember    bool
	LandLocked  bool
	Nationality string

	TLDs []string `json:"tld"`

	Languages    map[string]string
	Translations map[string]BaseLang
	Currencies   []string `json:"currency"`
	Borders      []string

	// Grouped meta
	Codes
	Geo
	Coordinates

	// Private
	subdivisions []SubDivision
}

// MeasurableCoordinates provides long/lat for country struct
// it does not store a pointer receiver to Country, as it needs to implement the
// Measurer interface
func (c Country) MeasurableCoordinates() (lat, long float64) {

	return c.Coordinates.Latitude, c.Coordinates.Longitude

}

// BorderingCountries returns the bordering countries the given Country
func (c *Country) BorderingCountries() (countries []Country) {

	query := New()

	for _, cca3 := range c.Borders {

		if country, err := query.FindCountryByAlpha(cca3); err == nil {
			countries = append(countries, country)
		}

	}

	return

}

// SubDivisions returns the subdivisions for the given Country
func (c *Country) SubDivisions() (subdivisions []SubDivision) {

	query := New()

	if res := query.Subdivisions[strings.ToLower(c.Alpha2)]; res != nil {

		subdivisions = res

	}

	return

}

// Codes contains various code representations
type Codes struct {
	Alpha2 string `json:"cca2"`
	Alpha3 string `json:"cca3"`
	CIOC   string
	CCN3   string

	//CountryCode         string // Yaml
	CallingCodes []string `json:"callingCode"`

	InternationalPrefix string // Yaml
}
