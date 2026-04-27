package gountries

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed all:data/yaml
var dataFS embed.FS

// embeddedRoot is the path inside dataFS that contains "countries" and
// "subdivisions" directories.
const embeddedRoot = "data/yaml"

// New creates a Query backed by the YAML data embedded in the package binary.
func New() *Query {
	return newFromFS(dataFS, embeddedRoot)
}

// NewFromPath creates a Query that reads YAML data from the given filesystem
// path. The directory must contain "countries" and "subdivisions"
// subdirectories matching the embedded layout.
//
// Like New, this is memoized for the lifetime of the process: the first
// successful call wins, and subsequent calls (with any path) return the same
// Query.
func NewFromPath(dataPath string) *Query {
	return newFromFS(os.DirFS(dataPath), ".")
}

func newFromFS(fsys fs.FS, root string) *Query {
	queryInitOnce.Do(func() {
		queryInstance = &Query{
			Countries: populateCountries(fsys, root),
		}

		queryInstance.NameToAlpha2 = populateNameIndex(queryInstance.Countries)
		queryInstance.Alpha3ToAlpha2 = populateAlphaIndex(queryInstance.Countries)
		queryInstance.NativeNameToAlpha2 = populateNativeNameIndex(queryInstance.Countries)
		queryInstance.CallingCodeToAlpha2 = populateCallingCodeIndex(queryInstance.Countries)
		queryInstance.CurrencyToAlpha2 = populateCallingCurrencyIndex(queryInstance.Countries)

		subdivisions := populateSubdivisions(fsys, root)
		for k := range queryInstance.Countries {
			c := queryInstance.Countries[k]
			c.subdivisions = subdivisions[strings.ToLower(c.Alpha2)]
			c.nameToSubdivision = map[string]SubDivision{}
			c.codeToSubdivision = map[string]SubDivision{}
			for _, s := range c.subdivisions {
				for _, n := range s.Names {
					c.nameToSubdivision[strings.ToLower(n)] = s
				}
				c.nameToSubdivision[strings.ToLower(s.Name)] = s
				c.codeToSubdivision[strings.ToLower(s.Code)] = s
			}
			queryInstance.Countries[k] = c
		}
	})

	return queryInstance
}

func populateNameIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)
	for alpha2, country := range countries {
		index[strings.ToLower(country.Name.Common)] = alpha2
		index[strings.ToLower(country.Name.Official)] = alpha2
		// Also index alternate spellings so common short forms (e.g.
		// "United States" for an iso_short_name of "United States of America")
		// resolve correctly. First-write wins so canonical names always
		// outrank alt-spellings on collision.
		for _, alt := range country.AltSpellings {
			key := strings.ToLower(alt)
			if _, ok := index[key]; !ok {
				index[key] = alpha2
			}
		}
	}
	return index
}

func populateAlphaIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)
	for alpha2, country := range countries {
		index[country.Codes.Alpha3] = alpha2
	}
	return index
}

func populateCountries(fsys fs.FS, root string) map[string]Country {
	countries := make(map[string]Country)
	dir := path.Join(root, "countries")

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		panic(fmt.Errorf("gountries: read countries dir %q: %w", dir, err))
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var c Country
		if err := yaml.Unmarshal(data, &c); err != nil {
			continue
		}
		countries[c.Codes.Alpha2] = c
	}

	return countries
}

func populateSubdivisions(fsys fs.FS, root string) map[string][]SubDivision {
	out := make(map[string][]SubDivision)
	dir := path.Join(root, "subdivisions")

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		panic(fmt.Errorf("gountries: read subdivisions dir %q: %w", dir, err))
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var subdivisions []SubDivision
		if err := yaml.Unmarshal(data, &subdivisions); err != nil {
			continue
		}
		alpha2 := strings.TrimSuffix(e.Name(), ".yaml")
		out[alpha2] = subdivisions
	}

	return out
}

func populateNativeNameIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)
	for alpha2, country := range countries {
		for _, nativeNames := range country.Name.Native {
			nativeName := strings.ToLower(nativeNames.Common)
			officialNativeName := strings.ToLower(nativeNames.Official)
			index[nativeName] = alpha2
			index[officialNativeName] = alpha2
		}
	}
	return index
}

func populateCallingCodeIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)
	for alpha2, country := range countries {
		for _, callingCode := range country.CallingCodes {
			index[callingCode] = alpha2
		}
	}
	return index
}

func populateCallingCurrencyIndex(countries map[string]Country) map[string][]Country {
	index := make(map[string][]Country)
	for _, country := range countries {
		for _, currency := range country.Currencies {
			index[currency] = append(index[currency], country)
		}
	}
	return index
}
