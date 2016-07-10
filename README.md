ASCII Art - Values
==================

When drawing ASCII art, each character can be treated as a pixel. It
will have a particular value which determines it's brightness and so
on (at least when programatically drawing an ASCII image).

This will determine the values of a large number of pixels, and will
generate a table from value (0 -> 255) to a rune value which can be
printed on-screen.

To try out the examples, simply `go run example/*.go` with whichever
file you please. The image-to-ascii.go demonstrates a practical use.
