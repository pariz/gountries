package gountries

import "fmt"

func ExampleSubdivisions() {

	se, _ := query.FindCountryByAlpha("SWE")

	subd := se.SubDivisions()

	for _, d := range subd {
		fmt.Println(d.Name)
	}

	// Output:
	//Västerbottens län
	//Uppsala län
	//Södermanlands län
	//Gotlands län
	//Dalarnas län
	//Stockholms län
	//Norrbottens län
	//Skåne län
	//Hallands län
	//Västra Götalands län
	//Örebro län
	//Jämtlands län
	//Östergötlands län
	//Västmanlands län
	//Gävleborgs län
	//Västernorrlands län
	//Jönköpings län
	//Kronobergs län
	//Kalmar län
	//Blekinge län
	//Värmlands län

}
