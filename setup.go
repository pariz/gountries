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

	if queryInited == false {
		queryInstance = &Query{
			Countries:    populateCountries(dataPath),
			Subdivisions: populateSubdivisions(dataPath),
		}
		queryInited = true
	}

	return queryInstance
}

func populateCountries(dataPath string) (countries []Country) {

	// Load countries into memory
	//
	//
	countries = []Country{}

	countriesPath := path.Join(dataPath, "countries")

	if info, err := ioutil.ReadDir(countriesPath); err == nil {

		for _, v := range info {

			if !v.IsDir() {

				if file, err := ioutil.ReadFile(filepath.Join(countriesPath, v.Name())); err == nil {

					country := Country{}
					if err := yaml.Unmarshal(file, &country); err == nil {

						// Save
						countries = append(countries, country)

					}

				}

			}

		}

	} else {
		panic(fmt.Errorf("Error loading Countries: %s", err))
	}
	return
}

func populateSubdivisions(dataPath string) (list map[string][]SubDivision) {

	// Load countries into memory
	//

	list = map[string][]SubDivision{}

	subdivisionsPath := path.Join(dataPath, "subdivisions")

	if info, err := ioutil.ReadDir(subdivisionsPath); err == nil {

		for _, v := range info {
			//fmt.Println(v.IsDir())

			if !v.IsDir() {

				if file, err := ioutil.ReadFile(filepath.Join(subdivisionsPath, v.Name())); err == nil {

					subdivisions := []SubDivision{}

					if err := yaml.Unmarshal(file, &subdivisions); err == nil {

						// Save
						//subdivisions = append(subdivisions, subdivision...)

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
