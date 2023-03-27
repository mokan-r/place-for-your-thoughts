package main

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
)

func main() {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 100})

	dc := gg.NewContext(300, 300)

	grad := gg.NewLinearGradient(20, 320, 400, 20)
	grad.AddColorStop(0, color.RGBA{R: 36, B: 20, A: 255})
	grad.AddColorStop(0.24, color.RGBA{R: 49, G: 6, B: 93, A: 255})
	grad.AddColorStop(0.55, color.RGBA{R: 9, G: 9, B: 121, A: 255})
	grad.AddColorStop(1, color.RGBA{G: 212, B: 255, A: 255})
	dc.SetFillStyle(grad)

	dc.DrawCircle(150, 146, 120)

	dc.Fill()

	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("MR", 150, 130, 0.5, 0.5)
	dc.SavePNG("amazing_logo.png")
}
