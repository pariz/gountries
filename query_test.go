package gountries

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCountryByName(t *testing.T) {

	var result Country
	var err error

	// Test for lowercase
	//

	result, err = query.FindCountryByName("sweden")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Lowercase country names should match")

	// Test for common name
	result, err = query.FindCountryByName("United States")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "US", "Lowercase country names should match")

	// Test for official name
	result, err = query.FindCountryByName("United States of America")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "US", "Lowercase country names should match")

	// Test for uppercase
	//

	result, err = query.FindCountryByName("SWEDEN")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Uppercase country names should match")

	// Test for invariants
	//

	invariants := []string{"Sweden", "SwEdEn", "SWEden"}

	for _, invariant := range invariants {

		result, err = query.FindCountryByName(invariant)

		if err != nil {
			t.Fail()
		}

		assert.Equal(t, result.Alpha2, "SE", fmt.Sprintf("Invariants of country names, eg sWeden,SWEDEN,swEdEn should match. %s did not match", invariant))

	}

}

func TestFindCountryByAlpha(t *testing.T) {

	var result Country
	var err error

	// Test for lowercase
	//

	result, err = query.FindCountryByAlpha("se")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Lowercase country names should match")

	// Test for uppercase
	//

	result, err = query.FindCountryByAlpha("SE")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Uppercase country names should match")

	// Test for invariants
	//

	result, err = query.FindCountryByAlpha("Se")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Invariants of country names, eg sWeden,SWEDEN,swEdEn should match")

	// Test for wrong code types (wrong length)
	//

	result, err = query.FindCountryByAlpha("SEE")

	if err != nil {
		assert.EqualError(t, err, "gountries error. Could not find country with code: SEE")

	} else {
		t.Fail()
	}

	// Test for wrong code types: too long
	//

	result, err = query.FindCountryByAlpha("SEEE")

	if err != nil {
		assert.EqualError(t, err, "gountries error. Invalid code format: SEEE")

	} else {
		t.Fail()
	}

	// Test for wrong code types: too short
	//

	result, err = query.FindCountryByAlpha("S")

	if err != nil {
		assert.EqualError(t, err, "gountries error. Invalid code format: S")

	} else {
		t.Fail()
	}

}

func TestFindAllCountries(t *testing.T) {

	assert.Len(t, query.FindAllCountries(), 249)

}

func TestFindCountries(t *testing.T) {

	country := Country{}
	country.Alpha2 = "SE"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 1)

	assert.Equal(t, countries[0].Alpha2, "SE", fmt.Sprintf("Countries did not return expected result %s: %s", "SE", countries[0].Alpha2))

}

func TestFindCountriesByRegion(t *testing.T) {

	country := Country{}
	country.Geo.Region = "Europe"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 52) // 52 is not the exact number of countries in Europe. Fix this later

}

func TestFindCountriesByContinent(t *testing.T) {

	country := Country{}
	country.Geo.Continent = "Europe"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 52) // 52 is not the exact number of countries in Europe. Fix this later

}

func TestFindCountriesBySubRegion(t *testing.T) {

	country := Country{}
	country.Geo.SubRegion = "Eastern Asia"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 8) // 8 is not the exact number of countries in Eastern Asia. Fix this later

}

func TestFindCountryByNativeName(t *testing.T) {

	var result Country
	var err error

	// Test for common name
	//

	result, err = query.FindCountryByNativeName("Sverige")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Common native country names should match")

	// Test for common name
	result, err = query.FindCountryByNativeName("Konungariket Sverige")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Official native country names should match")

	// Test for lowercase
	//

	result, err = query.FindCountryByNativeName("sverige")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Uppercase native country names should match")

	// Test for uppercase
	//

	result, err = query.FindCountryByNativeName("SVERIGE")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Alpha2, "SE", "Uppercase native country names should match")

	// Test for invariants
	//

	invariants := []string{"sVEriGE", "SveRIge", "SVErige"}

	for _, invariant := range invariants {

		result, err = query.FindCountryByNativeName(invariant)

		if err != nil {
			t.Fail()
		}

		assert.Equal(t, result.Alpha2, "SE", fmt.Sprintf("Invariants of native country names, eg sVEriGE,SveRIge,SVErige should match. %s did not match", invariant))

	}

}

func ExampleFindCountriesBorderingCountries() {

	country := Country{}
	//country.Alpha3 = "AUT"
	country.Borders = []string{
		"DEU",
	}

	countries := query.FindCountries(country)
	var c []string
	for _, country := range countries {
		c = append(c, country.Name.Common)
	}
	sort.Strings(c)
	for _, name := range c {
		fmt.Println(name)
	}

	// Output:
	//Austria
	//Belgium
	//Czech Republic
	//Denmark
	//France
	//Luxembourg
	//Netherlands
	//Poland
	//Switzerland

}

func ExampleFindCountriesBorderingCountries2() {

	country := Country{
		Borders: []string{
			"DEU",
			"CHE",
		},
	}

	countries := query.FindCountries(country)
	var c []string
	for _, country := range countries {
		c = append(c, country.Name.Common)
	}
	sort.Strings(c)
	for _, name := range c {
		fmt.Println(name)
	}

	// Output:
	//Austria
	//France

}

var result Country

func BenchmarkCountryLookupByName(b *testing.B) {

	q := New()
	var names []string
	for key := range q.Countries {
		names = append(names, q.Countries[key].Name.Common)
	}
	for n := 0; n <= b.N; n++ {
		randIndex := rand.Intn(len(q.Countries))
		c, err := q.FindCountryByName(names[randIndex])
		if err != nil {
			b.Fail()
		}
		result = c
	}
}

func TestFindSubdivisionCountryByNameOk(t *testing.T) {
	subdivisionName := "Braga"
	subdivisionCountry := "Portugal"

	country, err := query.FindSubdivisionCountryByName(subdivisionName)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, subdivisionCountry, country.Name.Common, fmt.Sprintf("Search %s subdivision country should return %s", subdivisionName, subdivisionCountry))
}

func TestFindSubdivisionCountryByNameError(t *testing.T) {
	subdivisionName := "123-bragaaaaa-123"

	_, err := query.FindSubdivisionCountryByName(subdivisionName)
	if err != nil {
		assert.Equal(t, "gountries error. Invalid subdivision name: "+subdivisionName, err.Error())
	} else {
		t.Fail()
	}
}

func TestFindCountryByCallingCode(t *testing.T) {
	country, _ := query.FindCountryByName("India")
	tests := []struct {
		name            string
		callingCode     string
		expectedError   error
		expectedCountry Country
	}{
		{
			name:            "Find country by valid calling code should return country when calling code is present for countries",
			callingCode:     "91",
			expectedError:   nil,
			expectedCountry: country,
		},
		{
			name:            "Find country by invalid country code should return error when searching for country",
			callingCode:     "abc",
			expectedError:   fmt.Errorf("gountries error. Could not find country with callingCode: abc"),
			expectedCountry: Country{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualCountry, actualError := query.FindCountryByCallingCode(test.callingCode)
			if !(reflect.DeepEqual(test.expectedCountry, actualCountry) &&
				reflect.DeepEqual(test.expectedError, actualError)) {
				t.Fail()
			}
		})
	}
}

func TestFindCountriesByCurrency(t *testing.T) {
	tests := []struct {
		name            string
		currency        string
		expectedError   error
		expectedCountry int
	}{
		{
			name:            "Find countries by valid currency should return countries when currency is present for countries",
			currency:        "EUR",
			expectedError:   nil,
			expectedCountry: 34,
		},
		{
			name:            "Find countries by unknown country code should return error when searching countries",
			currency:        "ERT",
			expectedError:   fmt.Errorf("gountries error. Could not find countries with currency name: ERT"),
			expectedCountry: 0,
		},
		{
			name:            "Find countries by invalid currency length code should return error when searching countries",
			currency:        "CNYE",
			expectedError:   fmt.Errorf("gountries error. Invalid currency format: CNYE"),
			expectedCountry: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualCountry, actualError := query.FindCountriesByCurrency(test.currency)
			if !(reflect.DeepEqual(test.expectedCountry, len(actualCountry)) &&
				reflect.DeepEqual(test.expectedError, actualError)) {
				t.Fail()
			}
		})
	}
}
