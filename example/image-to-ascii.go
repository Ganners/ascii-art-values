package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"os"

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
	bucket = asciivalues.FillInTheBlanks(bucket)

	reader, err := os.Open("example/test-ascii-art.jpg")
	if err != nil {
		log.Fatal("Failed to open image:", err)
	}

	decoded, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal("Failed to open image:", err)
	}

	for i := 0; i < decoded.Bounds().Max.Y; i += 1 {
		for j := 0.0; j < float64(decoded.Bounds().Max.X); j += 0.4 { // Note the 0.4
			val, _, _, _ := color.GrayModel.Convert(decoded.At(int(j), i)).RGBA()
			fmt.Print(string(bucket[val/256]))
		}
		fmt.Printf("\n")
	}
}
