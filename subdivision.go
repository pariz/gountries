package gountries

// SubDivision contains Country subdivison information
type SubDivision struct {
	Name  string
	Names []string
	Code  string

	CountryAlpha2 string

	Coordinates
}

// MeasurableCoordinates provides long/lat for country struct
func (sd *SubDivision) MeasurableCoordinates() (lat, long float64) {

	return sd.Coordinates.Latitude, sd.Coordinates.Longitude

}
