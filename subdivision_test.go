package gountries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubdivisions(t *testing.T) {

	se, _ := query.FindCountryByAlpha("SWE")

	subd := se.SubDivisions()

	assert.Len(t, subd, 21)

}
