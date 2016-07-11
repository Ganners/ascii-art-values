package asciivalues

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	startRune = 32

	// For debugging purposes
	fontSize = 2
)

// Given a filepath to a font file (expecting a .ttf) will return the font.Font
// interface, or an error
func LoadFont(file string) (font.Face, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    fontSize, // We only want the font to be 1x1 pixel!
		Hinting: font.HintingNone,
		DPI:     72,
	})
	return face, nil
}

// Generates a 2D square of runes
func BuildCharacterTable(maxRune int32) [][]rune {

	if maxRune < startRune {
		return [][]rune{}
	}

	size := 1 + maxRune - startRune // Rune starts from 32 (space)
	square := int(math.Ceil(math.Sqrt(float64(size))))

	charTable := make([][]rune, square)
	for i := 0; i < square; i++ {
		charTable[i] = make([]rune, square)
	}

	for i := 0; i < square; i++ {
		for j := 0; j < square; j++ {
			index := (i * square) + j
			charTable[i][j] = rune(index + startRune)
		}
	}

	return charTable
}

// Using a table of runes, builds an image matching that definition
func ImageFromRuneTable(charTable [][]rune, face font.Face) *image.RGBA {

	// Assume a perfect square, so the number of rows should match the
	// number of columns
	sizeX := len(charTable) * fontSize
	sizeY := (len(charTable) * fontSize)

	rgba := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))

	// Fill the canvas with black
	draw.Draw(rgba, rgba.Bounds(), image.Black, image.ZP, draw.Src)

	// Define a font drawer onto this rgba image
	drawer := &font.Drawer{
		Dst:  rgba,
		Src:  image.White,
		Face: face,
	}

	// For each character
	for i, row := range charTable {
		for j, char := range row {
			_, _, _, _, ok := face.Glyph(drawer.Dot, char)
			if !ok {
				continue
			}

			// Skip if it's less than the start rune
			if char < startRune {
				continue
			}
			drawer.Dot.X = fixed.I(((0 + j) * fontSize))
			drawer.Dot.Y = fixed.I(((1 + i) * fontSize))
			drawer.DrawString(string(char))
		}
	}

	return rgba
}

// Look at all of the pixels in the image and place their respective rune which
// created that value into the return at that value
func GenerateRuneValueBucket(charTable [][]rune, rgba *image.RGBA) []rune {

	size := len(charTable)
	emptyRune := rune(0)
	bucket := make([]rune, 256)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {

			startY := (i * fontSize)
			endY := (startY + fontSize)

			startX := j * fontSize
			endX := startX + fontSize

			numPixels := float64(fontSize * fontSize)
			pixelAverage := 0.0

			for i2 := startY; i2 < endY; i2++ {
				for j2 := startX; j2 < endX; j2++ {
					gray, _, _, _ := color.GrayModel.Convert(rgba.At(j2, i2)).RGBA()
					pixelAverage += float64(gray / 256)
				}
			}

			gray := int(pixelAverage / numPixels)

			if bucket[gray] == emptyRune {
				bucket[gray] = charTable[i][j]
			}
		}
	}

	return bucket
}

// Fills in the missing values as well as scaling to the full range of
// value
func FillInTheBlanks(in []rune) []rune {

	zeroRune := rune(0)

	out := make([]rune, len(in))
	missingFromEnd := 0

	// Find the last filled in value
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] != zeroRune {
			break
		}
		missingFromEnd++
	}

	realLength := len(in) - missingFromEnd
	ratio := len(in) / realLength

	// Fill in
	for i := 0; i < realLength; i++ {
		if in[i] != zeroRune {
			out[i*ratio] = in[i]
		}
	}

	previousValue := rune(0)
	for i := 0; i < len(in); i++ {
		if out[i] == zeroRune {
			out[i] = previousValue
		} else {
			previousValue = out[i]
		}
	}

	return out
}
