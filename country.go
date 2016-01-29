package gountries

import "strings"

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

// BorderingCountries gets the bordering countries for this country
func (c *Country) BorderingCountries() (countries []Country) {

	query := New()

	for _, cca3 := range c.Borders {

		if country, err := query.FindCountryByAlpha(cca3); err == nil {
			countries = append(countries, country)
		}

	}

	return

}

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
