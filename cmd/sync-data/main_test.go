package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// repoPath resolves a path relative to the gountries repo root, regardless of
// where `go test` was invoked. This lets golden-file tests find upstream/
// submodule data even when run from cmd/sync-data/.
func repoPath(rel string) string {
	_, here, _, _ := runtime.Caller(0) // .../cmd/sync-data/main_test.go
	root := filepath.Join(filepath.Dir(here), "..", "..")
	return filepath.Join(root, rel)
}

func TestMergeCountry(t *testing.T) {
	t.Run("primary fields come from gem, supplemental from mledoze", func(t *testing.T) {
		gc := gemCountry{
			Alpha2:              "SE",
			Alpha3:              "SWE",
			Continent:           "Europe",
			CountryCode:         "46",
			CurrencyCode:        "SEK",
			EUMember:            true,
			InternationalPrefix: "00",
			IOC:                 "SWE",
			ISOLongName:         "The Kingdom of Sweden",
			ISOShortName:        "Sweden",
			LanguagesOfficial:   []string{"sv"},
			Nationality:         "Swedish",
			Number:              "752",
			Region:              "Europe",
			Subregion:           "Northern Europe",
			UnofficialNames:     []string{"Sweden", "Sverige"},
			Geo: gemGeo{
				Latitude:     60.1,
				Longitude:    18.6,
				MinLatitude:  55.0,
				MaxLatitude:  69.0,
				MinLongitude: 10.5,
				MaxLongitude: 24.1,
			},
		}
		me := mledozeEntry{
			CCA2:         "SE",
			TLD:          []string{".se"},
			Borders:      []string{"FIN", "NOR"},
			Capital:      []string{"Stockholm"},
			Area:         450295,
			Landlocked:   false,
			Languages:    map[string]string{"swe": "Swedish"},
			AltSpellings: []string{"SE", "Kingdom of Sweden", "Sverige"},
			Translations: map[string]mledozeName{
				"deu": {Common: "Schweden", Official: "Königreich Schweden"},
			},
		}
		me.Name.Native = map[string]mledozeName{
			"swe": {Common: "Sverige", Official: "Konungariket Sverige"},
		}

		oc := mergeCountry(gc, me)

		// Names: gem's iso_short_name is canonical; mledoze provides natives
		assert.Equal(t, "Sweden", oc.Name.Common)
		assert.Equal(t, "The Kingdom of Sweden", oc.Name.Official)
		assert.Equal(t, "Sverige", oc.Name.Native["swe"].Common)

		// Codes
		assert.Equal(t, "SE", oc.Codes.Alpha2)
		assert.Equal(t, "SWE", oc.Codes.Alpha3)
		assert.Equal(t, "SWE", oc.Codes.CIOC)
		assert.Equal(t, "752", oc.Codes.CCN3)
		assert.Equal(t, []string{"46"}, oc.Codes.CallingCodes)
		assert.Equal(t, "00", oc.Codes.InternationalPrefix)

		// Geo from both
		assert.Equal(t, "Europe", oc.Geo.Region)
		assert.Equal(t, "Northern Europe", oc.Geo.SubRegion)
		assert.Equal(t, "Europe", oc.Geo.Continent)
		assert.Equal(t, "Stockholm", oc.Geo.Capital)
		assert.Equal(t, float64(450295), oc.Geo.Area)

		// Mledoze-only fields
		assert.Equal(t, []string{".se"}, oc.TLDs)
		assert.Equal(t, []string{"FIN", "NOR"}, oc.Borders)
		assert.False(t, oc.LandLocked)

		// Currency from gem (single string → slice)
		assert.Equal(t, []string{"SEK"}, oc.Currencies)

		// EU + nationality from gem
		assert.True(t, oc.EuMember)
		assert.Equal(t, "Swedish", oc.Nationality)

		// Coordinates
		assert.Equal(t, 60.1, oc.Coordinates.Latitude)
		assert.Equal(t, 10.5, oc.Coordinates.MinLongitude)
		assert.Equal(t, 24.1, oc.Coordinates.MaxLongitude)
	})

	t.Run("translations are uppercased", func(t *testing.T) {
		me := mledozeEntry{
			CCA2: "SE",
			Translations: map[string]mledozeName{
				"deu": {Common: "Schweden"},
				"fra": {Common: "Suède"},
			},
		}
		oc := mergeCountry(gemCountry{Alpha2: "SE", Alpha3: "SWE", ISOShortName: "Sweden"}, me)
		assert.Contains(t, oc.Translations, "DEU")
		assert.Contains(t, oc.Translations, "FRA")
		assert.NotContains(t, oc.Translations, "deu")
		assert.Equal(t, "Schweden", oc.Translations["DEU"].Common)
	})

	t.Run("AltSpellings dedup is case-insensitive across both sources", func(t *testing.T) {
		gc := gemCountry{
			Alpha2: "US", Alpha3: "USA", ISOShortName: "United States of America",
			UnofficialNames: []string{"United States", "USA", "Estados Unidos"},
		}
		me := mledozeEntry{
			CCA2:         "US",
			AltSpellings: []string{"US", "USA", "United States of America", "America"},
		}
		oc := mergeCountry(gc, me)

		// "USA" appears in both sources but only once in output.
		count := 0
		for _, s := range oc.AltSpellings {
			if s == "USA" {
				count++
			}
		}
		assert.Equal(t, 1, count, "USA should be deduped: %v", oc.AltSpellings)

		// Order preserves gem-first.
		assert.Equal(t, "United States", oc.AltSpellings[0])
		assert.Contains(t, oc.AltSpellings, "America")
		assert.Contains(t, oc.AltSpellings, "Estados Unidos")
	})

	t.Run("AltSpellings dedup ignores empty strings", func(t *testing.T) {
		gc := gemCountry{Alpha2: "X", Alpha3: "XXX", UnofficialNames: []string{"", "Foo", ""}}
		me := mledozeEntry{CCA2: "X", AltSpellings: []string{""}}
		oc := mergeCountry(gc, me)
		assert.Equal(t, []string{"Foo"}, oc.AltSpellings)
	})

	t.Run("languages merge gem codes (no name) with mledoze (with name)", func(t *testing.T) {
		gc := gemCountry{
			Alpha2: "SE", Alpha3: "SWE",
			LanguagesOfficial: []string{"sv"},
		}
		me := mledozeEntry{
			CCA2:      "SE",
			Languages: map[string]string{"swe": "Swedish"},
		}
		oc := mergeCountry(gc, me)
		assert.Equal(t, "Swedish", oc.Languages["swe"])
		assert.Contains(t, oc.Languages, "sv")
		assert.Empty(t, oc.Languages["sv"], "gem-only codes should be present but unnamed")
	})

	t.Run("falls back to mledoze names when gem names are empty", func(t *testing.T) {
		gc := gemCountry{Alpha2: "ZZ", Alpha3: "ZZZ"} // gem present but missing names
		me := mledozeEntry{CCA2: "ZZ"}
		me.Name.Common = "Zedland"
		me.Name.Official = "Republic of Zedland"
		oc := mergeCountry(gc, me)
		assert.Equal(t, "Zedland", oc.Name.Common)
		assert.Equal(t, "Republic of Zedland", oc.Name.Official)
	})

	t.Run("works with missing mledoze entry", func(t *testing.T) {
		gc := gemCountry{
			Alpha2: "SE", Alpha3: "SWE",
			ISOShortName: "Sweden", ISOLongName: "The Kingdom of Sweden",
		}
		// Empty mledoze — common when an alpha2 exists in gem but not mledoze.
		oc := mergeCountry(gc, mledozeEntry{})

		assert.Equal(t, "Sweden", oc.Name.Common)
		assert.Empty(t, oc.Borders, "no borders without mledoze")
		assert.Empty(t, oc.TLDs, "no tlds without mledoze")
		assert.Empty(t, oc.Translations)
	})

	t.Run("capital is first element of mledoze.Capital array", func(t *testing.T) {
		gc := gemCountry{Alpha2: "ZA", Alpha3: "ZAF"}
		me := mledozeEntry{
			CCA2:    "ZA",
			Capital: []string{"Pretoria", "Cape Town", "Bloemfontein"},
		}
		oc := mergeCountry(gc, me)
		assert.Equal(t, "Pretoria", oc.Geo.Capital)
	})

	t.Run("missing currency_code yields nil currencies", func(t *testing.T) {
		gc := gemCountry{Alpha2: "AQ", Alpha3: "ATA"} // Antarctica has no currency
		oc := mergeCountry(gc, mledozeEntry{})
		assert.Nil(t, oc.Currencies)
	})
}

func TestMergeSubdivision(t *testing.T) {
	t.Run("populates all fields including coordinates from geo", func(t *testing.T) {
		sd := gemSubdivision{
			Name:            "Stockholms län",
			Code:            "AB",
			UnofficialNames: []string{"Stockholm"},
			Geo: &gemGeo{
				Latitude:     59.6,
				Longitude:    18.1,
				MinLatitude:  58.7,
				MaxLatitude:  60.2,
				MinLongitude: 17.2,
				MaxLongitude: 19.3,
			},
		}
		out := mergeSubdivision("SE", sd)

		assert.Equal(t, "Stockholms län", out.Name)
		assert.Equal(t, "AB", out.Code)
		assert.Equal(t, "SE", out.CountryAlpha2)
		assert.Equal(t, []string{"Stockholm"}, out.Names)
		assert.Equal(t, 59.6, out.Coordinates.Latitude)
		assert.Equal(t, 17.2, out.Coordinates.MinLongitude)
		assert.Equal(t, 19.3, out.Coordinates.MaxLongitude)
	})

	t.Run("subdivision without geo block has zero coordinates", func(t *testing.T) {
		sd := gemSubdivision{Name: "Foo", Code: "F", Geo: nil}
		out := mergeSubdivision("XX", sd)
		assert.Equal(t, "Foo", out.Name)
		assert.Equal(t, "XX", out.CountryAlpha2)
		assert.Equal(t, float64(0), out.Coordinates.Latitude)
	})
}

// TestLoadGemCountriesGoldenFile verifies the loader parses real upstream YAML.
// Skips if submodules aren't checked out (e.g., on a stale CI cache).
func TestLoadGemCountriesGoldenFile(t *testing.T) {
	gem, err := loadGemCountries(repoPath("upstream/countries/lib/countries/data"))
	if err != nil {
		t.Skipf("upstream submodules not initialized: %v", err)
	}
	require.NotEmpty(t, gem)

	se, ok := gem["SE"]
	require.True(t, ok, "SE should be present in upstream gem data")
	assert.Equal(t, "SE", se.Alpha2)
	assert.Equal(t, "SWE", se.Alpha3)
	assert.Equal(t, "Sweden", se.ISOShortName)
	assert.True(t, se.EUMember)
}

func TestLoadMledozeGoldenFile(t *testing.T) {
	m, err := loadMledoze(repoPath("upstream/mledoze-countries/countries.json"))
	if err != nil {
		t.Skipf("upstream submodules not initialized: %v", err)
	}
	require.NotEmpty(t, m)

	se, ok := m["SE"]
	require.True(t, ok, "SE should be present in mledoze data")
	assert.Contains(t, se.Borders, "FIN")
	assert.Contains(t, se.Borders, "NOR")
	assert.Equal(t, []string{".se"}, se.TLD)
}
