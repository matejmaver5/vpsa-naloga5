package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/matejmaver5/vpsa-naloga5/vreme"
	"github.com/urfave/cli/v3"
)

func prikaziVreme(_ context.Context, ukaz *cli.Command) error {
	minTemp := ukaz.Float64("minTemp")
	maxTemp := ukaz.Float64("maxTemp")
	stMeritev := ukaz.Int("stMeritev")
	if !(minTemp <= maxTemp) || math.IsInf(minTemp, 0) || math.IsInf(maxTemp, 0) {
		return fmt.Errorf("meji temperatur morata biti končni števili in minTemp ne sme presegati maxTemp")
	}
	if stMeritev < 1 {
		return fmt.Errorf("stMeritev mora biti vsaj 1")
	}

	postaje := make(map[string]vreme.Postaja)
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Bežigrad", "Ljubljana")
	vreme.DodajPostajo(postaje, "KP01", "Koper Luka", "Koper")
	vreme.DodajPostajo(postaje, "KR01", "Kredarica", "Triglav")
	vreme.DodajPostajo(postaje, "MB01", "Maribor Tabor", "Maribor")

	for _, meritev := range []struct {
		oznaka      string
		temperatura float64
	}{
		{"LJ01", 18.2}, {"LJ01", 21.5}, {"LJ01", 19.8},
		{"KP01", 24.5}, {"KP01", 27}, {"KP01", 26.1},
		{"KR01", 2}, {"KR01", -1.5}, {"MB01", 12.4},
		{"LJ01", 16.5}, {"KP01", 28.4}, {"KR01", 3.5},
		{"LJ01", 60}, {"XX01", 20},
	} {
		vreme.DodajMeritev(postaje, meritev.oznaka, meritev.temperatura, minTemp, maxTemp)
	}

	fmt.Println()
	vreme.IzpisMeritev(postaje)
	fmt.Println("\nPOROČILO:")
	vreme.IzpisPorocila(postaje, stMeritev)
	return nil
}

func main() {
	ukaz := &cli.Command{
		Name:  "vreme",
		Usage: "Izpiše meritve in poročilo vremenskih postaj.",
		Flags: []cli.Flag{
			&cli.Float64Flag{Name: "minTemp", Value: -50, Usage: "Najnižja dovoljena temperatura v °C"},
			&cli.Float64Flag{Name: "maxTemp", Value: 50, Usage: "Najvišja dovoljena temperatura v °C"},
			&cli.IntFlag{Name: "stMeritev", Value: 3, Usage: "Najmanjše število meritev za statistiko"},
		},
		Action: prikaziVreme,
	}
	if err := ukaz.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
