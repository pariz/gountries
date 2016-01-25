package gountries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

// NewGountries creates an Query object and unmarshals the json file.
func NewGountries() (query *QueryHolder, err error) {

	if query != nil {
		return
	}

	query = &QueryHolder{}

	err = unmarshalCountries(query)

	return
}

func unmarshalCountries(query *QueryHolder) (err error) {

	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	bytes, err := ioutil.ReadFile(path.Dir(f) + "/countries.json")

	if err != nil {
		return
	}

	fmt.Println("Got data", query)

	json.Unmarshal(bytes, &query.Countries)

	return
}
