package gountries

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDataInvariants is a fast sanity check that the embedded data is sane
// after every regeneration. It guards against silent upstream regressions —
// a country disappearing, a key field going empty, an alpha2 code mutating
// — without locking us into specific values that may legitimately drift.
//
// Bounds are deliberately loose; the goal is to catch *catastrophic* changes
// (mass deletions, schema breakage), not nuance.
func TestDataInvariants(t *testing.T) {
	q := New()

	t.Run("country count within sensible bounds", func(t *testing.T) {
		// 249 today; ISO-3166-1 ratifications happen but rarely.
		assert.GreaterOrEqual(t, len(q.Countries), 245, "country count dropped below floor")
		assert.LessOrEqual(t, len(q.Countries), 260, "country count above ceiling")
	})

	t.Run("every country has the basics", func(t *testing.T) {
		for alpha2, c := range q.Countries {
			require.Lenf(t, alpha2, 2, "country key %q is not 2 chars", alpha2)
			require.Equalf(t, alpha2, c.Alpha2, "key/Alpha2 mismatch for %q", alpha2)
			require.Lenf(t, c.Alpha3, 3, "%s missing alpha3 (got %q)", alpha2, c.Alpha3)
			require.NotEmptyf(t, c.Name.Common, "%s missing common name", alpha2)
			require.NotEmptyf(t, c.Name.Official, "%s missing official name", alpha2)
		}
	})

	t.Run("indexes are populated", func(t *testing.T) {
		assert.NotEmpty(t, q.NameToAlpha2)
		assert.NotEmpty(t, q.Alpha3ToAlpha2)
		assert.NotEmpty(t, q.NativeNameToAlpha2)
		assert.NotEmpty(t, q.CallingCodeToAlpha2)
		assert.NotEmpty(t, q.CurrencyToAlpha2)

		// alpha3 → alpha2 should cover every country
		assert.Lenf(t, q.Alpha3ToAlpha2, len(q.Countries), "alpha3 index incomplete")
	})

	t.Run("canary lookups", func(t *testing.T) {
		cases := []struct {
			lookup, alpha2, common string
			method                 func(string) (Country, error)
		}{
			{"SE", "SE", "Sweden", q.FindCountryByAlpha},
			{"swe", "SE", "Sweden", q.FindCountryByAlpha},
			{"DEU", "DE", "Germany", q.FindCountryByAlpha},
			{"USA", "US", "United States of America", q.FindCountryByAlpha},
			{"Sweden", "SE", "Sweden", q.FindCountryByName},
			{"United States", "US", "United States of America", q.FindCountryByName}, // alt-spelling
		}
		for _, c := range cases {
			c := c
			t.Run(c.lookup, func(t *testing.T) {
				got, err := c.method(c.lookup)
				require.NoError(t, err)
				assert.Equal(t, c.alpha2, got.Alpha2)
				assert.Equal(t, c.common, got.Name.Common)
			})
		}
	})

	t.Run("Sweden borders are intact", func(t *testing.T) {
		se, err := q.FindCountryByAlpha("SE")
		require.NoError(t, err)
		// Defensive: borders set should at minimum include FIN and NOR.
		assert.Contains(t, se.Borders, "FIN")
		assert.Contains(t, se.Borders, "NOR")
	})

	t.Run("subdivisions present for major countries", func(t *testing.T) {
		// Pick a few we expect to always have subdivisions; if these go to zero
		// the upstream pipeline almost certainly broke.
		mustHaveSubdivisions := []string{"SE", "US", "DE", "GB", "FR", "BR", "IN"}
		for _, a2 := range mustHaveSubdivisions {
			c, err := q.FindCountryByAlpha(a2)
			require.NoError(t, err)
			assert.NotEmptyf(t, c.SubDivisions(), "%s should have subdivisions", a2)
		}
	})
}
