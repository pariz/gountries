package gountries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubdivisions(t *testing.T) {

	se, _ := query.FindCountryByAlpha("SWE")

	subd := se.SubDivisions()

	assert.Len(t, subd, 21)

	found, err := se.FindSubdivisionByName(subd[0].Name)
	assert.NoError(t, err)
	assert.Equal(t, subd[0], found)

	for _, n := range subd[0].Names {
		found, err = se.FindSubdivisionByName(n)
		assert.NoError(t, err)
		assert.Equal(t, subd[0], found)
	}

	found, err = se.FindSubdivisionByCode(subd[0].Code)
	assert.NoError(t, err)
	assert.Equal(t, subd[0], found)
}
