package gountries

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	var err error

	if err != nil {
		panic(err)
	}

	flag.Parse()
	os.Exit(m.Run())
}

func TestFindCountryByName(t *testing.T) {

	var result Country
	var err error

	// Test for lowercase
	//

	result, err = Query.FindCountryByName("sweden")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Code, "SE", "Lowercase country names should match")

	// Test for uppercase
	//

	result, err = Query.FindCountryByName("SWEDEN")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Code, "SE", "Uppercase country names should match")

	// Test for invariants
	//

	invariants := []string{"Sweden", "SwEdEn", "SWEden"}

	for _, invariant := range invariants {

		result, err = Query.FindCountryByName(invariant)

		if err != nil {
			t.Fail()
		}

		assert.Equal(t, result.Code, "SE", fmt.Sprintf("Invariants of country names, eg sWeden,SWEDEN,swEdEn should match. %s did not match", invariant))

	}

}

func TestFindCountryByCode(t *testing.T) {

	var result Country
	var err error

	// Test for lowercase
	//

	result, err = Query.FindCountryByCode("se")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Code, "SE", "Lowercase country names should match")

	// Test for uppercase
	//

	result, err = Query.FindCountryByCode("SE")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Code, "SE", "Uppercase country names should match")

	// Test for invariants
	//

	result, err = Query.FindCountryByCode("Se")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, result.Code, "SE", "Invariants of country names, eg sWeden,SWEDEN,swEdEn should match")

	// Test for wrong code types (wrong length)
	//

	result, err = Query.FindCountryByCode("SEE")

	if err != nil {
		assert.EqualError(t, err, "Could not find country with code SEE")

	} else {
		t.Fail()
	}

	// Test for wrong code types: too long
	//

	result, err = Query.FindCountryByCode("SEEE")

	if err != nil {
		assert.EqualError(t, err, "SEEE is an invalid code format")

	} else {
		t.Fail()
	}

	// Test for wrong code types: too short
	//

	result, err = Query.FindCountryByCode("S")

	if err != nil {
		assert.EqualError(t, err, "S is an invalid code format")

	} else {
		t.Fail()
	}

}

func ExampleBorderCountries() {

	se, _ := Query.FindCountryByCode("SWE")
	for _, country := range se.BorderingCountries() {
		fmt.Println(country.Name.Common)
	}

	// Output:
	// Finland
	// Norway

}
