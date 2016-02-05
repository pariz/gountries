package gountries

import "fmt"

func ExampleSubdivisions() {

	se, _ := query.FindCountryByAlpha("SWE")

	subd := se.SubDivisions()

	for _, d := range subd {
		fmt.Println(d.Name)
	}

	// Output:
	//Östergötlands län
	//Jönköpings län
	//Kalmar län
	//Gotlands län
	//Skåne län
	//Hallands län
	//Stockholms län
	//Uppsala län
	//Västmanlands län
	//Örebro län
	//Kronobergs län
	//Blekinge län
	//Värmlands län
	//Jämtlands län
	//Norrbottens län
	//Västra Götalands län
	//Dalarnas län
	//Gävleborgs län
	//Västernorrlands län
	//Västerbottens län
	//Södermanlands län

}
