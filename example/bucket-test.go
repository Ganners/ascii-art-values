package main

import (
	"fmt"
	"log"

	"github.com/Ganners/ascii-art-values"
)

func main() {
	face, err := asciivalues.LoadFont("Hack-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	table := asciivalues.BuildCharacterTable(10000)
	img := asciivalues.ImageFromRuneTable(table, face)
	bucket := asciivalues.GenerateRuneValueBucket(table, img)

	filled := 0
	missing := 0

	zeroRune := rune(0)

	for v, r := range bucket {
		if r == zeroRune {
			missing++
		} else {
			filled++
		}
		fmt.Printf("Value %d: '%s'\n", v, string(r))
	}

	log.Printf("Missing: %d, Filled: %d\n", missing, filled)
}
