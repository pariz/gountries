package gountries

// Query holds a reference to the QueryHolder struct
var queryInited = false
var queryInstance *Query

// Query contains queries for countries, cities, etc.
type Query struct {
	Countries    []Country
	Subdivisions []SubDivision
}
