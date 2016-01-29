package gountries

// SubDivision contains Country subdivison information
type SubDivision struct {
	Name  string
	Names []string
	Code  string

	CountryAlpha2 string

	Coordinates
}
