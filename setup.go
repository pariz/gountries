package gountries

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// New creates an Query object and unmarshals the json file.
func New() *Query {

	dataPath := filepath.Join("data", "yaml")

	return NewFromPath(dataPath)
}

// NewFromPath creates a Query object from data folder in provided path
func NewFromPath(dataPath string) *Query {

	queryInitOnce.Do(func() {
		queryInstance = &Query{
			Countries: populateCountries(dataPath),
		}

		queryInstance.NameToAlpha2 = populateNameIndex(queryInstance.Countries)
		queryInstance.Alpha3ToAlpha2 = populateAlphaIndex(queryInstance.Countries)
		queryInstance.NativeNameToAlpha2 = populateNativeNameIndex(queryInstance.Countries)

		subdivisions := populateSubdivisions(dataPath)
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
		name := strings.ToLower(country.Name.Common)
		officialName := strings.ToLower(country.Name.Official)
		index[name] = alpha2
		index[officialName] = alpha2
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

func populateCountries(dataPath string) map[string]Country {

	// Load countries into memory
	//
	//
	var countries = make(map[string]Country)

	countriesPath := path.Join(dataPath, "countries")

	// Try packed data first before custom data directory
	if yamlFileList, err := AssetDir("data/yaml/countries"); err == nil {
		return populateCountriesFromPackedData(yamlFileList, "data/yaml/countries")
	}

	if info, err := ioutil.ReadDir(countriesPath); err == nil {

		var file []byte

		for _, v := range info {

			if !v.IsDir() {

				if file, err = ioutil.ReadFile(filepath.Join(countriesPath, v.Name())); err == nil {

					country := Country{}
					if err = yaml.Unmarshal(file, &country); err == nil {

						// Save
						countries[country.Codes.Alpha2] = country

					}

				}

			}

		}

	} else {
		panic(fmt.Errorf("Error loading Countries: %s", err))
	}
	return countries
}

func populateCountriesFromPackedData(fileList []string, path string) map[string]Country {
	var data []byte
	var err error
	var countries = make(map[string]Country)

	for _, yamlFile := range fileList {
		var country Country
		data, err = Asset(filepath.Join(path, yamlFile))
		if err != nil {
			continue
		}
		if err = yaml.Unmarshal(data, &country); err != nil {
			continue
		}
		countries[country.Codes.Alpha2] = country
	}
	return countries
}

func populateSubdivisions(dataPath string) (list map[string][]SubDivision) {

	// Load countries into memory
	//

	list = map[string][]SubDivision{}

	subdivisionsPath := path.Join(dataPath, "subdivisions")

	// Try packed data first before custom data directory
	if yamlFileList, err := AssetDir("data/yaml/subdivisions"); err == nil {
		return populateSubdivisionsFromPackedData(yamlFileList, "data/yaml/subdivisions")
	}

	if info, err := ioutil.ReadDir(subdivisionsPath); err == nil {

		for _, v := range info {

			if !v.IsDir() {

				if file, err := ioutil.ReadFile(filepath.Join(subdivisionsPath, v.Name())); err == nil {

					subdivisions := []SubDivision{}

					if err := yaml.Unmarshal(file, &subdivisions); err == nil {

						// Save
						// subdivisions = append(subdivisions, subdivision...)
						list[strings.Split(v.Name(), ".")[0]] = subdivisions
					}

				}

			}

		}

	} else {
		panic("Error loading Subdivisions")
	}

	return
}

func populateSubdivisionsFromPackedData(fileList []string, path string) map[string][]SubDivision {
	var data []byte
	var err error
	sd := make(map[string][]SubDivision)

	for _, yamlFile := range fileList {
		data, err = Asset(filepath.Join(path, yamlFile))
		if err != nil {
			continue
		}
		var subdivisions []SubDivision
		if err = yaml.Unmarshal(data, &subdivisions); err != nil {
			continue
		}
		alpha2 := strings.Split(yamlFile, ".")[0]

		sd[alpha2] = subdivisions
	}
	return sd
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
