package gountries

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPopulateCountriesPanicsOnMalformedYAML guards against silent data
// corruption: if the auto-sync ever lands a broken YAML file in data/yaml/,
// we want loud failure at startup — not a Query that's quietly missing
// random countries.
func TestPopulateCountriesPanicsOnMalformedYAML(t *testing.T) {
	fsys := fstest.MapFS{
		"countries/se.yaml": {Data: []byte("name:\n  common: Sweden\n")}, // valid
		"countries/xx.yaml": {Data: []byte("name: [unclosed\n  bracket: true\n: : :")}, // malformed
	}

	defer func() {
		r := recover()
		require.NotNil(t, r, "expected populateCountries to panic on malformed YAML; got nil")
		msg, ok := r.(error)
		if !ok {
			msg = nil
		}
		// Either an error or a string — both should reference the bad file
		// so the operator can find it.
		got := ""
		if msg != nil {
			got = msg.Error()
		} else if s, ok := r.(string); ok {
			got = s
		}
		assert.Contains(t, strings.ToLower(got), "xx.yaml",
			"panic message should name the failing file; got %q", got)
	}()

	populateCountries(fsys, ".")
}

// TestPopulateSubdivisionsPanicsOnMalformedYAML mirrors the above for the
// subdivision loader.
func TestPopulateSubdivisionsPanicsOnMalformedYAML(t *testing.T) {
	fsys := fstest.MapFS{
		"subdivisions/se.yaml": {Data: []byte("- name: Foo\n  code: F\n")},
		"subdivisions/xx.yaml": {Data: []byte(":::garbage:::")},
	}

	defer func() {
		r := recover()
		require.NotNil(t, r, "expected populateSubdivisions to panic on malformed YAML")
		got := ""
		if e, ok := r.(error); ok {
			got = e.Error()
		} else if s, ok := r.(string); ok {
			got = s
		}
		assert.Contains(t, strings.ToLower(got), "xx.yaml")
	}()

	populateSubdivisions(fsys, ".")
}
