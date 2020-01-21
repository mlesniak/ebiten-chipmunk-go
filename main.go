package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"
)

var space *cp.Space // Simulation space

const width = 800
const height = 600
const numBoxes = 10

type Box struct {
	x, y float64
	w, h float64
	r    float64
}

var boxImage *ebiten.Image
var boxes = []Box{}

func update(screen *ebiten.Image) error {
	checkExit()
	space.Step(1.0 / float64(ebiten.MaxTPS()))

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	background(screen)
	drawBoxes(screen)

	return nil
}

func drawBoxes(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	space.EachBody(func(body *cp.Body) {
		switch body.UserData {
		case "box":
			op.GeoM.Reset()
			//log.Printf("%v: %f/%f\n", body, body.Position().X, body.Position().Y)
			op.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(boxImage, op)
		case "line":
			//log.Printf("%v: %f/%f\n", body, body.Position().X, body.Position().Y)
			ebitenutil.DrawLine(screen, 0, body.Position().Y, width, body.Position().Y, colornames.Yellow)
		}
	})
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
	initBoxes()

	space = cp.NewSpace()
	space.Iterations = 1
	space.SetGravity(cp.Vector{0, 800})

	// Add floor
	bf := cp.NewStaticBody()
	h := 100.0
	bf.SetPosition(cp.Vector{width / 2, height - h})
	bf.UserData = "line"
	sf := cp.NewBox(bf, width, h/2, 1.0)
	space.AddBody(bf)
	space.AddShape(sf)

	for _, box := range boxes {
		body := cp.NewBody(1.0, cp.INFINITY)
		body.SetPosition(cp.Vector{X: box.x, Y: box.y})
		body.UserData = "box"

		shape := cp.NewBox(body, box.w, box.h, 0.0)
		shape.SetElasticity(1.0)
		shape.SetFriction(0)

		space.AddBody(body)
		space.AddShape(shape)
		log.Printf("Creating body %v\n", body)
	}

	if err := ebiten.Run(update, width, height, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}

func initBoxes() {
	rand.Seed(time.Now().Unix())
	boxImage, _, _ = ebitenutil.NewImageFromFile("box.png", ebiten.FilterDefault)
	w, h := boxImage.Size()
	boxWidth := float64(w)
	boxHeight := float64(h)
	for i := 0; i < numBoxes; i++ {
		log.Println("Creating box", i)
		boxes = append(boxes, Box{
			x: rand.Float64() * (width - boxWidth),
			y: rand.Float64()*50 + boxHeight/2.0,
			w: boxWidth,
			h: boxHeight,
			r: rand.Float64() * 2 * math.Pi,
		})
	}
}
