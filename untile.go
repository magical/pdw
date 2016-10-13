// usage: go run untile.go -o untiled/ images/*.png

package main

import (
	"fmt"
	"image"
	"image/color"
	//"image/draw"
	"flag"
	"image/png"
	"os"
	"path/filepath"
)

func readPng(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

type Untiled struct {
	image.Image
	Offset int
}

func (m *Untiled) At(x, y int) color.Color {
	p := image.Pt(x+m.Offset, y+m.Offset)
	p = p.Mod(m.Image.Bounds())
	return m.Image.At(p.X, p.Y)
}

func main() {
	outdir := flag.String("o", ".", "output directory")
	flag.Parse()

	for _, filename := range flag.Args() {
		m, err := readPng(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		outfilename := filepath.Join(*outdir, filepath.Base(filename))
		if filepath.Clean(filename) == filepath.Clean(outfilename) {
			fmt.Println("not overwriting", filename)
			continue
		}

		f, err := os.Create(outfilename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = png.Encode(f, &Untiled{m, 150})
		if err != nil {
			fmt.Println(err)
		}

		f.Close()
	}
}
