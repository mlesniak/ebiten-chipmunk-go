package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	_ "image/png"
)

var space *cp.Space // Simulation space

const width = 800
const height = 600
const numBoxes = 1
const boxWidth = 40
const boxHeight = 40

type Box struct {
	x, y float64
	w, h float64
	r    float64
}

var boxImage *ebiten.Image
var boxes = []Box{}

var down = 0

func update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && down < 0 {
		down = 10 // FPS fix
		x, y := ebiten.CursorPosition()
		b := addBox(float64(x), float64(y))
		addBoxToPhysics(b)
	} else {
		down--
	}

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
	space.EachBody(func(body *cp.Body) {
		switch body.UserData {
		case "box":
			ebitenutil.DrawRect(screen, body.Position().X-boxHeight/2, body.Position().Y-boxWidth/2, boxWidth, boxHeight, color.RGBA{80, 80, 80, 255})
		case "line":
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

	bf := cp.NewStaticBody()
	bf.SetPosition(cp.Vector{width / 2, height})
	bf.UserData = "line"
	sf := cp.NewBox(bf, width, 0, 0.0)
	space.AddBody(bf)
	space.AddShape(sf)

	for _, box := range boxes {
		body := addBoxToPhysics(box)
		log.Printf("Creating body %v\n", body)
	}

	if err := ebiten.Run(update, width, height, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}

func addBoxToPhysics(box Box) *cp.Body {
	body := cp.NewBody(10.0, cp.INFINITY)
	body.SetPosition(cp.Vector{X: box.x, Y: box.y})
	body.UserData = "box"

	shape := cp.NewBox(body, box.w, box.h, box.r)
	shape.SetElasticity(1.0)
	shape.SetFriction(0)

	space.AddBody(body)
	space.AddShape(shape)
	return body
}

func initBoxes() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < numBoxes; i++ {
		addBox(rand.Float64()*(width-boxWidth)+boxWidth, rand.Float64()*(height-boxHeight)+boxHeight)
	}
}

func addBox(px, py float64) Box {
	box := Box{
		x: px,
		y: py,
		w: boxWidth,
		h: boxHeight,
	}
	boxes = append(boxes, box)
	return box
}
