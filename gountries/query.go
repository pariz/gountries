package gountries

// Query holds a reference to the QueryHolder struct
var query *QueryHolder

// QueryHolder contains queries for countries, cities, etc.
type QueryHolder struct {
	Countries []Country
}
