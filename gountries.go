// Package gountries provides offline lookups for ISO-3166-1 countries,
// ISO-3166-2 subdivisions, ISO-4217 currencies, and ISO-6709 coordinates,
// plus translations and bordering-country relationships. Data is embedded
// in the package binary; no network is required at runtime.
package gountries

// BaseLang is a basic structure for common language formatting in the JSON file
type BaseLang struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

// Geo contains geographical information
type Geo struct {
	Region    string `json:"region"`
	SubRegion string `json:"subregion"`
	Continent string // Yaml
	Capital   string `json:"capital"`

	Area float64
}

// Coordinates contains the coordinates for both Country and SubDivision
type Coordinates struct {
	LongitudeString string `json:"longitude"`
	LatitudeString  string `json:"latitude"`

	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
	Latitude     float64
	Longitude    float64
}

// MeasureDistancePythagoras measures distances betweeen two countries using
// Pythagoras equirect angular equation
func MeasureDistancePythagoras(m1 Measurer, m2 Measurer) (distance float64) {

	m1Long, m1Lat := m1.MeasurableCoordinates()
	m2Long, m2Lat := m2.MeasurableCoordinates()

	return CalculatePythagorasEquirectangular(m1Lat, m1Long, m2Lat, m2Long)

}

// MeasureDistanceHaversine measures distances betweeen two countries using
// Havesine equation
func MeasureDistanceHaversine(m1 Measurer, m2 Measurer) (distance float64) {

	m1Long, m1Lat := m1.MeasurableCoordinates()
	m2Long, m2Lat := m2.MeasurableCoordinates()

	return CalculateHaversine(m1Lat, m1Long, m2Lat, m2Long)

}

// Measurer provides coordinates for measurements
type Measurer interface {
	MeasurableCoordinates() (lat, long float64)
}
