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
	subdivisions      []SubDivision
	nameToSubdivision map[string]SubDivision
	codeToSubdivision map[string]SubDivision
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
	return c.subdivisions
}

// FindSubdivisionByName fincs a country by given name
func (c *Country) FindSubdivisionByName(name string) (result SubDivision, err error) {
	s, exists := c.nameToSubdivision[strings.ToLower(name)]
	if !exists {
		return SubDivision{}, makeError("Could not find subdivision with name", name)
	}
	return s, nil
}

// FindSubdivisionByCode fincs a country by given code
func (c *Country) FindSubdivisionByCode(code string) (result SubDivision, err error) {
	s, exists := c.codeToSubdivision[strings.ToLower(code)]
	if !exists {
		return SubDivision{}, makeError("Could not find subdivision with code", code)
	}
	return s, nil
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
