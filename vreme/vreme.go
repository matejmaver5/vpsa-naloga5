// Paket vreme hrani meritve vremenskih postaj in izpisuje temperaturna poročila.
package vreme

import "fmt"

// Postaja hrani ime, kraj in izmerjene temperature.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo doda postajo brez meritev, če je oznaka še prosta.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	if _, obstaja := postaje[oznaka]; !obstaja {
		postaje[oznaka] = Postaja{ime: ime, kraj: kraj}
	}
}

// DodajMeritev shrani temperaturo v območju [minTemp, maxTemp] ali izpiše razlog za zavrnitev.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) {
	postaja, ok := postaje[oznaka]
	if !ok {
		fmt.Printf("Postaja %s ne obstaja.\n", oznaka)
		return
	}

	if !(temperatura >= minTemp && temperatura <= maxTemp) {
		fmt.Printf("Neveljavna temperatura %.1f °C. Dovoljeno območje je od %g do %g °C.\n", temperatura, minTemp, maxTemp)
		return
	}

	postaja.meritve = append(postaja.meritve, temperatura)

	postaje[oznaka] = postaja
}

func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	postaja, ok := postaje[oznaka]
	if !ok || stMeritev < 1 || len(postaja.meritve) < stMeritev {
		return 0, 0, 0, false
	}

	najnizja = postaja.meritve[0]
	najvisja = postaja.meritve[0]
	vsota := 0.0

	for _, temperatura := range postaja.meritve {
		if temperatura < najnizja {
			najnizja = temperatura
		}
		if temperatura > najvisja {
			najvisja = temperatura
		}
		vsota += temperatura
	}

	povprecna = vsota / float64(len(postaja.meritve))
	return najnizja, najvisja, povprecna, true
}

// IzpisMeritev izpiše oznako, ime, kraj in vse meritve vsake postaje.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for oznaka, postaja := range postaje {
		fmt.Printf("%s - %s (%s): %v\n", oznaka, postaja.ime, postaja.kraj, postaja.meritve)
	}
}

// IzpisPorocila izpiše statistiko in opis vremena za vsako postajo z vsaj stMeritev meritvami.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	for oznaka, postaja := range postaje {
		najnizja, najvisja, povprecna, ok := statistika(postaje, oznaka, stMeritev)
		if !ok {
			fmt.Printf("%s: premalo meritev\n", postaja.ime)
			continue
		}

		opis := ""
		switch {
		case povprecna >= 25:
			opis = "Vroče"
		case povprecna >= 15:
			opis = "Toplo"
		case povprecna >= 5:
			opis = "Hladno"
		default:
			opis = "Mrzlo"
		}

		fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n",
			postaja.ime, najnizja, najvisja, povprecna, opis)
	}
}
