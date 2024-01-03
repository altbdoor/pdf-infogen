# pdf-infogen

_Too generic of a name for a generic action._

In short:

1. Takes in a bunch of data from `coords.json`
1. Each item in `coords.json` is:
   - `field` just a descriptive name of what is being filled in
   - `value` the text you want to fill in
   - `position` is the `[X, Y]` coordinate on where to start
   - `block` on whether the input is a square block
   - `cellSize` is the width and height of the block
   - `fontSize` is the font size
1. Based on the data, render them on the passed in PDF file
1. Output it as `output.pdf`

## Usage

First, edit the `coords.json` as you see fit. Then, run the binary.

```sh
# for normal usage
./infogen form.pdf

# to see the box boundaries
DEBUG=1 ./infogen form.pdf
```

## Building

```sh
go build -o infogen
```
