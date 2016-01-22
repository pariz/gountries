package gountries

// Query holds a reference to the QueryHolder struct
var Query = &query{}

// QueryHolder contains queries for countries, cities, etc.
type query struct {
	Countries []Country
}
