// Command sync-data merges country and subdivision data from two upstream
// sources into the canonical YAML files embedded by gountries.
//
// Sources (vendored as git submodules under vendor/):
//   - countries/countries (MIT) — primary schema for ISO-3166-1, calling
//     codes, EU membership, currency, intl prefix, geo bounds.
//   - mledoze/countries (ODbL) — supplemental fields not in the gem:
//     borders, TLDs, language-name lookup, multi-lang translations.
//
// Run from the gountries repo root: `go run ./cmd/sync-data`.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	gemDataPath   = "upstream/countries/lib/countries/data"
	mledozeJSON   = "upstream/mledoze-countries/countries.json"
	outCountriesD = "data/yaml/countries"
	outSubdivD    = "data/yaml/subdivisions"
)

func main() {
	log.SetFlags(0)

	gem, err := loadGemCountries()
	if err != nil {
		log.Fatalf("load gem countries: %v", err)
	}
	log.Printf("loaded %d countries from countries/countries", len(gem))

	mledoze, err := loadMledoze()
	if err != nil {
		log.Fatalf("load mledoze: %v", err)
	}
	log.Printf("loaded %d countries from mledoze/countries", len(mledoze))

	if err := os.MkdirAll(outCountriesD, 0o755); err != nil {
		log.Fatalf("mkdir countries: %v", err)
	}
	if err := os.MkdirAll(outSubdivD, 0o755); err != nil {
		log.Fatalf("mkdir subdivisions: %v", err)
	}

	// Country output
	written := 0
	for alpha2, gc := range gem {
		me := mledoze[alpha2] // may be zero value
		oc := mergeCountry(gc, me)
		path := filepath.Join(outCountriesD, strings.ToLower(alpha2)+".yaml")
		if err := writeYAML(path, oc); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}
		written++
	}
	log.Printf("wrote %d country files to %s", written, outCountriesD)

	// Subdivision output
	subs, err := loadGemSubdivisions()
	if err != nil {
		log.Fatalf("load subdivisions: %v", err)
	}
	subWritten := 0
	totalSubs := 0
	for alpha2, list := range subs {
		out := make([]outSubdivision, 0, len(list))
		for _, s := range list {
			out = append(out, mergeSubdivision(alpha2, s))
		}
		sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
		path := filepath.Join(outSubdivD, strings.ToLower(alpha2)+".yaml")
		if err := writeYAML(path, out); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}
		subWritten++
		totalSubs += len(out)
	}
	log.Printf("wrote %d subdivision files (%d entries) to %s", subWritten, totalSubs, outSubdivD)
}

// ---------- countries/countries (gem) loaders ----------

type gemGeo struct {
	Latitude     float64 `yaml:"latitude"`
	Longitude    float64 `yaml:"longitude"`
	MaxLatitude  float64 `yaml:"max_latitude"`
	MaxLongitude float64 `yaml:"max_longitude"`
	MinLatitude  float64 `yaml:"min_latitude"`
	MinLongitude float64 `yaml:"min_longitude"`
}

type gemCountry struct {
	Alpha2              string   `yaml:"alpha2"`
	Alpha3              string   `yaml:"alpha3"`
	Continent           string   `yaml:"continent"`
	CountryCode         string   `yaml:"country_code"`
	CurrencyCode        string   `yaml:"currency_code"`
	EUMember            bool     `yaml:"eu_member"`
	Geo                 gemGeo   `yaml:"geo"`
	InternationalPrefix string   `yaml:"international_prefix"`
	IOC                 string   `yaml:"ioc"`
	ISOLongName         string   `yaml:"iso_long_name"`
	ISOShortName        string   `yaml:"iso_short_name"`
	LanguagesOfficial   []string `yaml:"languages_official"`
	LanguagesSpoken     []string `yaml:"languages_spoken"`
	Nationality         string   `yaml:"nationality"`
	Number              string   `yaml:"number"`
	Region              string   `yaml:"region"`
	Subregion           string   `yaml:"subregion"`
	UnMember            bool     `yaml:"un_member"`
	UnofficialNames     []string `yaml:"unofficial_names"`
	CapitalRaw          string   `yaml:"-"` // gem stores capital in nameservers/cache, not here
}

type gemSubdivision struct {
	Name            string            `yaml:"name"`
	Code            string            `yaml:"code"`
	UnofficialNames []string          `yaml:"unofficial_names"`
	Geo             *gemGeo           `yaml:"geo"`
	Translations    map[string]string `yaml:"translations"`
	Type            string            `yaml:"type"`
}

func loadGemCountries() (map[string]gemCountry, error) {
	dir := filepath.Join(gemDataPath, "countries")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(map[string]gemCountry, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var wrap map[string]gemCountry
		if err := yaml.Unmarshal(data, &wrap); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		for k, v := range wrap {
			result[strings.ToUpper(k)] = v
		}
	}
	return result, nil
}

func loadGemSubdivisions() (map[string][]gemSubdivision, error) {
	dir := filepath.Join(gemDataPath, "subdivisions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]gemSubdivision, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var wrap map[string]gemSubdivision
		if err := yaml.Unmarshal(data, &wrap); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		alpha2 := strings.ToUpper(strings.TrimSuffix(e.Name(), ".yaml"))
		list := make([]gemSubdivision, 0, len(wrap))
		for code, sd := range wrap {
			if sd.Code == "" {
				sd.Code = code
			}
			list = append(list, sd)
		}
		result[alpha2] = list
	}
	return result, nil
}

// ---------- mledoze loader ----------

type mledozeEntry struct {
	Name struct {
		Common   string                        `json:"common"`
		Official string                        `json:"official"`
		Native   map[string]mledozeName        `json:"native"`
	} `json:"name"`
	TLD          []string                      `json:"tld"`
	CCA2         string                        `json:"cca2"`
	CCN3         string                        `json:"ccn3"`
	CCA3         string                        `json:"cca3"`
	CIOC         string                        `json:"cioc"`
	Currencies   map[string]map[string]string  `json:"currencies"`
	IDD          struct {
		Root     string   `json:"root"`
		Suffixes []string `json:"suffixes"`
	} `json:"idd"`
	Capital      []string                       `json:"capital"`
	AltSpellings []string                       `json:"altSpellings"`
	Region       string                         `json:"region"`
	Subregion    string                         `json:"subregion"`
	Languages    map[string]string              `json:"languages"`
	Translations map[string]mledozeName         `json:"translations"`
	Latlng       []float64                      `json:"latlng"`
	Landlocked   bool                           `json:"landlocked"`
	Borders      []string                       `json:"borders"`
	Area         float64                        `json:"area"`
}

type mledozeName struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

func loadMledoze() (map[string]mledozeEntry, error) {
	data, err := os.ReadFile(mledozeJSON)
	if err != nil {
		return nil, err
	}
	var list []mledozeEntry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	result := make(map[string]mledozeEntry, len(list))
	for _, e := range list {
		result[strings.ToUpper(e.CCA2)] = e
	}
	return result, nil
}

// ---------- output schema (matches gountries' YAML) ----------

type outBaseLang struct {
	Common   string `yaml:"common"`
	Official string `yaml:"official"`
}

type outCodes struct {
	Alpha2              string   `yaml:"alpha2"`
	Alpha3              string   `yaml:"alpha3"`
	CIOC                string   `yaml:"cioc,omitempty"`
	CCN3                string   `yaml:"ccn3,omitempty"`
	CallingCodes        []string `yaml:"callingcodes,omitempty"`
	InternationalPrefix string   `yaml:"internationalprefix,omitempty"`
}

type outGeo struct {
	Region    string  `yaml:"region"`
	SubRegion string  `yaml:"subregion"`
	Continent string  `yaml:"continent"`
	Capital   string  `yaml:"capital"`
	Area      float64 `yaml:"area,omitempty"`
}

type outCoordinates struct {
	LongitudeString string  `yaml:"longitudestring"`
	LatitudeString  string  `yaml:"latitudestring"`
	MinLongitude    float64 `yaml:"minlongitude"`
	MinLatitude     float64 `yaml:"minlatitude"`
	MaxLongitude    float64 `yaml:"maxlongitude"`
	MaxLatitude     float64 `yaml:"maxlatitude"`
	Latitude        float64 `yaml:"latitude"`
	Longitude       float64 `yaml:"longitude"`
}

type outName struct {
	Common   string                 `yaml:"common"`
	Official string                 `yaml:"official"`
	Native   map[string]outBaseLang `yaml:"native,omitempty"`
}

type outCountry struct {
	Name         outName                `yaml:"name"`
	EuMember     bool                   `yaml:"eumember"`
	LandLocked   bool                   `yaml:"landlocked"`
	Nationality  string                 `yaml:"nationality"`
	AltSpellings []string               `yaml:"altspellings,omitempty"`
	TLDs         []string               `yaml:"tlds,omitempty"`
	Languages    map[string]string      `yaml:"languages,omitempty"`
	Translations map[string]outBaseLang `yaml:"translations,omitempty"`
	Currencies   []string               `yaml:"currencies,omitempty"`
	Borders      []string               `yaml:"borders,omitempty"`
	Codes        outCodes               `yaml:"codes"`
	Geo          outGeo                 `yaml:"geo"`
	Coordinates  outCoordinates         `yaml:"coordinates"`
}

type outSubdivision struct {
	Name          string         `yaml:"name"`
	Names         []string       `yaml:"names,omitempty"`
	Code          string         `yaml:"code"`
	CountryAlpha2 string         `yaml:"countryalpha2"`
	Coordinates   outCoordinates `yaml:"coordinates"`
}

// ---------- merge logic ----------

func mergeCountry(gc gemCountry, me mledozeEntry) outCountry {
	oc := outCountry{}

	oc.Name.Common = gc.ISOShortName
	oc.Name.Official = gc.ISOLongName
	if me.CCA2 != "" {
		// fall back to mledoze names if gem fields are empty
		if oc.Name.Common == "" {
			oc.Name.Common = me.Name.Common
		}
		if oc.Name.Official == "" {
			oc.Name.Official = me.Name.Official
		}
		// native names only available from mledoze
		if len(me.Name.Native) > 0 {
			oc.Name.Native = make(map[string]outBaseLang, len(me.Name.Native))
			for k, v := range me.Name.Native {
				oc.Name.Native[k] = outBaseLang{Common: v.Common, Official: v.Official}
			}
		}
	}

	oc.EuMember = gc.EUMember
	oc.LandLocked = me.Landlocked
	oc.Nationality = gc.Nationality

	if len(me.TLD) > 0 {
		oc.TLDs = append([]string(nil), me.TLD...)
	}

	// AltSpellings: dedup union of gem.unofficial_names + mledoze.altSpellings.
	seen := make(map[string]struct{})
	for _, n := range gc.UnofficialNames {
		if n == "" {
			continue
		}
		key := strings.ToLower(n)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		oc.AltSpellings = append(oc.AltSpellings, n)
	}
	for _, n := range me.AltSpellings {
		if n == "" {
			continue
		}
		key := strings.ToLower(n)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		oc.AltSpellings = append(oc.AltSpellings, n)
	}

	// languages: gem has codes only; mledoze has code → name. Combine.
	langs := map[string]string{}
	for code, name := range me.Languages {
		langs[code] = name
	}
	for _, code := range gc.LanguagesOfficial {
		if _, ok := langs[code]; !ok {
			langs[code] = "" // present but unnamed
		}
	}
	if len(langs) > 0 {
		oc.Languages = langs
	}

	if len(me.Translations) > 0 {
		oc.Translations = make(map[string]outBaseLang, len(me.Translations))
		for k, v := range me.Translations {
			// upstream uses 3-letter lowercase ISO 639-3; existing yaml uses uppercase
			oc.Translations[strings.ToUpper(k)] = outBaseLang{Common: v.Common, Official: v.Official}
		}
	}

	if gc.CurrencyCode != "" {
		oc.Currencies = []string{gc.CurrencyCode}
	}

	if len(me.Borders) > 0 {
		oc.Borders = append([]string(nil), me.Borders...)
	}

	oc.Codes = outCodes{
		Alpha2:              gc.Alpha2,
		Alpha3:              gc.Alpha3,
		CIOC:                gc.IOC,
		CCN3:                gc.Number,
		InternationalPrefix: gc.InternationalPrefix,
	}
	if gc.CountryCode != "" {
		oc.Codes.CallingCodes = []string{gc.CountryCode}
	}

	oc.Geo = outGeo{
		Region:    gc.Region,
		SubRegion: gc.Subregion,
		Continent: gc.Continent,
		Area:      me.Area,
	}
	if len(me.Capital) > 0 {
		oc.Geo.Capital = me.Capital[0]
	}

	oc.Coordinates = outCoordinates{
		MinLongitude: gc.Geo.MinLongitude,
		MinLatitude:  gc.Geo.MinLatitude,
		MaxLongitude: gc.Geo.MaxLongitude,
		MaxLatitude:  gc.Geo.MaxLatitude,
		Latitude:     gc.Geo.Latitude,
		Longitude:    gc.Geo.Longitude,
	}

	return oc
}

func mergeSubdivision(alpha2 string, sd gemSubdivision) outSubdivision {
	out := outSubdivision{
		Name:          sd.Name,
		Names:         sd.UnofficialNames,
		Code:          sd.Code,
		CountryAlpha2: alpha2,
	}
	if sd.Geo != nil {
		out.Coordinates = outCoordinates{
			MinLongitude: sd.Geo.MinLongitude,
			MinLatitude:  sd.Geo.MinLatitude,
			MaxLongitude: sd.Geo.MaxLongitude,
			MaxLatitude:  sd.Geo.MaxLatitude,
			Latitude:     sd.Geo.Latitude,
			Longitude:    sd.Geo.Longitude,
		}
	}
	return out
}

// ---------- yaml output helper ----------

func writeYAML(path string, v any) error {
	buf, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, buf, 0o644)
}
