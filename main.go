package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"math/rand"
	"os"

	_ "image/png"
)

const width = 800
const height = 600
const numBoxes = 10

type Box struct {
	x, y float64
	w, h float64
}

var boxImage *ebiten.Image
var boxes = make([]Box, numBoxes)

func update(screen *ebiten.Image) error {
	checkExit()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	background(screen)
	for _, box := range boxes {
		op := &ebiten.DrawImageOptions{}
		w, h := boxImage.Size()
		scale := box.w / float64(w)
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(box.x, box.y)
		screen.DrawImage(boxImage, op)
	}

	return nil
}

func background(screen *ebiten.Image) error {
	return screen.Fill(color.Gray{Y: 30})
}

func checkExit() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
}

func main() {
	boxImage, _, _ = ebitenutil.NewImageFromFile("box.png", ebiten.FilterDefault)
	boxWidth := 40.0
	boxHeight := 40.0
	for i := 0; i < numBoxes; i++ {
		boxes = append(boxes, Box{
			x: rand.Float64() * (width - boxWidth),
			y: rand.Float64()*50 + boxHeight/2.0,
			w: boxWidth,
			h: boxHeight,
		})
	}

	if err := ebiten.Run(update, width, height, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}
