package asciivalues

import (
	"image"
	"image/color"
	"os"
	"reflect"
	"testing"
)

func TestBuildCharacterTable(t *testing.T) {
	testCases := []struct {
		max           int32
		expectedTable [][]rune
	}{
		{
			max: 19,
			// The first real character in the utf-8 table is 20, which
			// is a space. If the max is lower than 20, we should
			// receive an empty table
			expectedTable: [][]rune{},
		},
		{
			max:           startRune + 0,
			expectedTable: [][]rune{{32}},
		},
		{
			max: startRune + 3,
			expectedTable: [][]rune{
				{rune(32), rune(33)},
				{rune(34), rune(35)},
			},
		},
		{
			max: startRune + 24,
			expectedTable: [][]rune{
				{rune(32), rune(33), rune(34), rune(35), rune(36)},
				{rune(37), rune(38), rune(39), rune(40), rune(41)},
				{rune(42), rune(43), rune(44), rune(45), rune(46)},
				{rune(47), rune(48), rune(49), rune(50), rune(51)},
				{rune(52), rune(53), rune(54), rune(55), rune(56)},
			},
		},
	}

	for _, test := range testCases {
		table := BuildCharacterTable(test.max)
		if !reflect.DeepEqual(table, test.expectedTable) {
			t.Errorf("Table %+v does not match expected %+v", table, test.expectedTable)
		}
	}
}

func TestImageFromRuneTable(t *testing.T) {
	// Convert a pixel RGBA to grayscale and return 0 or 1 as per our
	// pixel test standard
	var colorToPixelTest = func(c color.Color) int8 {
		gray, _, _, _ := color.GrayModel.Convert(c).RGBA()
		if gray > 0 {
			return 1
		}
		return 0
	}

	testCases := []struct {
		charTable [][]rune

		// A pixel test is a rough approximation, instead of from 0-255
		// we'll approximate from 0-1. 0 being black and 1 being
		// anything not black
		pixelTests [][]int8
	}{
		{
			charTable:  [][]rune{{' '}},
			pixelTests: [][]int8{{0}},
		},
		{
			charTable: [][]rune{
				{'#', ' '},
				{' ', '#'},
			},
			pixelTests: [][]int8{
				{1, 0},
				{0, 1},
			},
		},
	}

	// Load font, there's one in the freetype package we can use. We'll
	// use this font for testing.
	fontFile := "/src/github.com/golang/freetype/testdata/luxisr.ttf"
	fontFile = os.Getenv("GOPATH") + fontFile

	face, err := LoadFont(fontFile)
	if err != nil {
		t.Fatalf("Error reading font file: %s", err.Error())
	}

	for _, test := range testCases {

		// Grab image from the function under test
		img := ImageFromRuneTable(test.charTable, face)

		// Test each pixel
		for i, row := range test.pixelTests {
			for j, expectedPixel := range row {

				rgba := img.At(j*fontSize, i*fontSize)
				pixel := colorToPixelTest(rgba)

				if pixel != expectedPixel {
					t.Errorf("Pixel %d does not match expected %d at [%d, %d]",
						pixel, expectedPixel, i, j)
				}
			}
		}
	}
}

func TestGenerateRuneValueBucket(t *testing.T) {
	{
		charTable := [][]rune{
			{'0', '1'},
			{'2', '3'},
		}

		// If R = G = B = A, then A = the Grayscale value
		rgba := image.NewRGBA(image.Rect(0, 0, len(charTable)*fontSize, len(charTable)*fontSize))
		for i := 0; i < fontSize; i++ {
			for j := 0; j < fontSize; j++ {
				rgba.Set(0*fontSize+j, 0*fontSize+i, color.RGBA{0, 0, 0, 255})
				rgba.Set(0*fontSize+j, 1*fontSize+i, color.RGBA{22, 22, 22, 255})
				rgba.Set(1*fontSize+j, 0*fontSize+i, color.RGBA{55, 55, 55, 255})
				rgba.Set(1*fontSize+j, 1*fontSize+i, color.RGBA{255, 255, 255, 255})
			}
		}

		expected := make([]rune, 256)
		expected[0] = '0'
		expected[55] = '1'
		expected[22] = '2'
		expected[255] = '3'

		bucket := GenerateRuneValueBucket(charTable, rgba)

		if !reflect.DeepEqual(bucket, expected) {
			t.Errorf("Bucket %+v does not match expected %+v", bucket, expected)
		}
	}
}

func TestFillInTheBlanks(t *testing.T) {

	testCases := []struct {
		input    []rune
		expected []rune
	}{
		{
			input: []rune{
				rune(32),
				rune(33),
				rune(0),
				rune(0),
			},
			expected: []rune{
				rune(32),
				rune(32),
				rune(33),
				rune(33),
			},
		},
		{
			input: []rune{
				rune(32),
				rune(33),
				rune(0),
				rune(0),
				rune(34),
				rune(38),
				rune(0),
				rune(0),
				rune(0),
				rune(0),
				rune(0),
				rune(0),
			},
			expected: []rune{
				rune(32),
				rune(32),
				rune(33),
				rune(33),
				rune(33),
				rune(33),
				rune(33),
				rune(33),
				rune(34),
				rune(34),
				rune(38),
				rune(38),
			},
		},
	}

	for _, test := range testCases {
		out := FillInTheBlanks(test.input)

		if !reflect.DeepEqual(out, test.expected) {
			t.Errorf("%+v does not match %+v", out, test.expected)
		}
	}
}
