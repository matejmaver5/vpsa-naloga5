package vreme

import (
	"math"
	"testing"
)

func ExampleDodajPostajo() {
	postaje := make(map[string]Postaja)
	DodajPostajo(postaje, "KR01", "Kredarica", "Triglav")
	DodajMeritev(postaje, "KR01", -2, -2, 2)
	DodajMeritev(postaje, "KR01", 3, -2, 2)
	IzpisPorocila(postaje, 2)
	DodajMeritev(postaje, "KR01", 2, -2, 2)
	IzpisPorocila(postaje, 2)
	// Output:
	// Neveljavna temperatura 3.0 °C. Dovoljeno območje je od -2 do 2 °C.
	// Kredarica: premalo meritev
	// Kredarica: min -2.0 °C, max 2.0 °C, povprečje 0.0 °C -> Mrzlo
}

func TestDodajMeritev(t *testing.T) {
	primeri := []struct {
		ime         string
		oznaka      string
		temperatura float64
		veljavna    bool
	}{
		{"obicajna", "LJ01", 16.5, true},
		{"spodnja meja", "LJ01", -50, true},
		{"zgornja meja", "LJ01", 50, true},
		{"prenizka", "LJ01", -50.1, false},
		{"previsoka", "LJ01", 50.1, false},
		{"neobstojeca postaja", "XX01", 20, false},
		{"ni stevilo", "LJ01", math.NaN(), false},
		{"neskoncnost", "LJ01", math.Inf(1), false},
	}
	for _, primer := range primeri {
		t.Run(primer.ime, func(t *testing.T) {
			postaje := map[string]Postaja{
				"LJ01": {ime: "Ljubljana", kraj: "Ljubljana", meritve: []float64{10, 20}},
			}
			DodajMeritev(postaje, primer.oznaka, primer.temperatura, -50, 50)
			postaja := postaje["LJ01"]
			zelenaDolzina := 2
			if primer.veljavna {
				zelenaDolzina = 3
			}
			if len(postaje) != 1 || len(postaja.meritve) != zelenaDolzina {
				t.Fatalf("nepricakovano stanje postaj: %v", postaje)
			}
			if postaja.meritve[0] != 10 || postaja.meritve[1] != 20 || postaja.ime != "Ljubljana" || postaja.kraj != "Ljubljana" {
				t.Fatalf("obstojeci podatki so se spremenili: %v", postaja)
			}
			if primer.veljavna && postaja.meritve[2] != primer.temperatura {
				t.Fatalf("meritev ni shranjena: %v", postaja.meritve)
			}
		})
	}
}

func TestStatistika(t *testing.T) {
	primeri := []struct {
		ime       string
		meritve   []float64
		min       float64
		max       float64
		povprecje float64
		ok        bool
	}{
		{"brez meritev", nil, 0, 0, 0, false},
		{"ena meritev", []float64{12.4}, 0, 0, 0, false},
		{"dve meritvi", []float64{10, 20}, 0, 0, 0, false},
		{"tri meritve", []float64{2, -1.5, 3.5}, -1.5, 3.5, 4.0 / 3, true},
		{"pozitivne", []float64{18.2, 21.5, 19.8, 16.5}, 16.5, 21.5, 19, true},
		{"negativne", []float64{-10, -2, -6}, -10, -2, -6, true},
		{"enake", []float64{5, 5, 5}, 5, 5, 5, true},
	}
	for _, primer := range primeri {
		t.Run(primer.ime, func(t *testing.T) {
			postaje := map[string]Postaja{"P": {meritve: primer.meritve}}
			min, max, povprecje, ok := statistika(postaje, "P", 3)
			if ok != primer.ok || min != primer.min || max != primer.max || math.Abs(povprecje-primer.povprecje) > 1e-9 {
				t.Fatalf("dobljeno (%v, %v, %v, %v), pricakovano (%v, %v, %v, %v)",
					min, max, povprecje, ok, primer.min, primer.max, primer.povprecje, primer.ok)
			}
		})
	}
	if _, _, _, ok := statistika(nil, "XX01", 3); ok {
		t.Fatal("neobstojeca postaja ne sme imeti statistike")
	}
}

func ExampleIzpisMeritev() {
	IzpisMeritev(map[string]Postaja{
		"LJ01": {ime: "Ljubljana Bežigrad", kraj: "Ljubljana", meritve: []float64{18.2, 21.5, 19.8}},
		"MB01": {ime: "Maribor Tabor", kraj: "Maribor", meritve: []float64{12.4}},
	})
	// Unordered output:
	// MERITVE:
	// LJ01 - Ljubljana Bežigrad (Ljubljana): [18.2 21.5 19.8]
	// MB01 - Maribor Tabor (Maribor): [12.4]
}

func ExampleIzpisPorocila() {
	IzpisPorocila(map[string]Postaja{
		"A": {ime: "Pod 5", meritve: []float64{4.9, 4.9, 4.9}},
		"B": {ime: "Meja 5", meritve: []float64{5, 5, 5}},
		"C": {ime: "Pod 15", meritve: []float64{14.9, 14.9, 14.9}},
		"D": {ime: "Meja 15", meritve: []float64{15, 15, 15}},
		"E": {ime: "Pod 25", meritve: []float64{24.9, 24.9, 24.9}},
		"F": {ime: "Meja 25", meritve: []float64{25, 25, 25}},
		"G": {ime: "Premalo", meritve: []float64{10, 20}},
	}, 3)
	// Unordered output:
	// Pod 5: min 4.9 °C, max 4.9 °C, povprečje 4.9 °C -> Mrzlo
	// Meja 5: min 5.0 °C, max 5.0 °C, povprečje 5.0 °C -> Hladno
	// Pod 15: min 14.9 °C, max 14.9 °C, povprečje 14.9 °C -> Hladno
	// Meja 15: min 15.0 °C, max 15.0 °C, povprečje 15.0 °C -> Toplo
	// Pod 25: min 24.9 °C, max 24.9 °C, povprečje 24.9 °C -> Toplo
	// Meja 25: min 25.0 °C, max 25.0 °C, povprečje 25.0 °C -> Vroče
	// Premalo: premalo meritev
}
