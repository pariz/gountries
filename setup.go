package gountries

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// New creates an Query object and unmarshals the json file.
func New() *Query {

	if queryInited == false {
		queryInstance = &Query{
			Countries:    populateCountries(),
			Subdivisions: populateSubdivisions(),
		}
		queryInited = true
	}

	return queryInstance
}

func populateCountries() (countries []Country) {

	// Load countries into memory
	//
	//
	countries = []Country{}

	countriesPath := filepath.Join("data", "yaml", "countries")

	if info, err := ioutil.ReadDir(countriesPath); err == nil {

		for _, v := range info {
			//fmt.Println(v.IsDir())

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
		panic("Error loading Countries")
	}
	return
}

func populateSubdivisions() (subdivisions []SubDivision) {

	// Load countries into memory
	//

	subdivisions = []SubDivision{}

	countriesPath := filepath.Join("data", "yaml", "subdivisions")

	if info, err := ioutil.ReadDir(countriesPath); err == nil {

		for _, v := range info {
			//fmt.Println(v.IsDir())

			if !v.IsDir() {

				if file, err := ioutil.ReadFile(filepath.Join(countriesPath, v.Name())); err == nil {

					subdivision := SubDivision{}
					if err := yaml.Unmarshal(file, &subdivision); err == nil {

						// Save
						subdivisions = append(subdivisions, subdivision)

					}

				}

			}

		}

	} else {
		panic("Error loading Subdivisions")
	}

	return
}
