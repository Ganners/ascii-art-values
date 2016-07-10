package main

import (
	"bufio"
	"image/png"
	"log"
	"os"

	"github.com/Ganners/ascii-art-values"
)

func main() {
	face, err := asciivalues.LoadFont("Hack-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	img := asciivalues.ImageFromRuneTable(
		asciivalues.BuildCharacterTable(20000),
		face,
	)

	file, err := os.Create("debug-image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b := bufio.NewWriter(file)
	err = png.Encode(b, img)
	if err != nil {
		log.Fatal(err)
	}

	err = b.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
