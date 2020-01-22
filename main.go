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

var space *cp.Space // Simulation space
var boxImage *ebiten.Image
var boxes = []Box{}
var lastClicked = 0 // Frames to wait until the next mouse click is registered.

func main() {
	boxImage, _, _ = ebitenutil.NewImageFromFile("box.png", ebiten.FilterDefault)
	initBoxes()

	space = cp.NewSpace()
	space.Iterations = 10
	space.SetGravity(cp.Vector{0, 800})

	addFloor()

	for _, box := range boxes {
		body := addBoxToPhysics(box)
		log.Printf("Creating body %v\n", body)
	}

	if err := ebiten.Run(update, width, height, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	// Input handling.
	checkExit()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && lastClicked < 0 {
		lastClicked = 10 // FPS fix
		x, y := ebiten.CursorPosition()
		b := addBox(float64(x), float64(y))
		addBoxToPhysics(b)
	} else {
		lastClicked--
	}

	// Next step in engine.
	space.Step(1.0 / float64(ebiten.MaxTPS()))

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Show results.
	background(screen)
	drawBoxes(screen)

	return nil
}

func drawBoxes(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(200.0/255.0, 200.0/255.0, 200.0/255.0, 1)

	space.EachBody(func(body *cp.Body) {
		switch body.UserData {
		case "box":
			//fmt.Printf("%.2f\n", body.Rotation().ToAngle())
			op.GeoM.Reset()
			op.GeoM.Rotate(body.Angle())
			op.GeoM.Translate(body.Position().X-boxWidth/2, body.Position().Y-boxHeight/2)
			screen.DrawImage(boxImage, op)
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

func addFloor() {
	floorHeight := 100.0
	bf := cp.NewStaticBody()
	bf.SetPosition(cp.Vector{width / 2, height + floorHeight/2})
	bf.UserData = "line"
	sf := cp.NewBox(bf, width, floorHeight, 0.0)
	space.AddBody(bf)
	space.AddShape(sf)
}

func addBoxToPhysics(box Box) *cp.Body {
	body := cp.NewBody(1000.0, cp.INFINITY)
	body.SetPosition(cp.Vector{X: box.x, Y: box.y})
	//body.SetAngle(box.r)
	body.UserData = "box"

	shape := cp.NewBox(body, box.w, box.h, 0.0)
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
		r: rand.Float64() * 2 * math.Pi,
	}
	boxes = append(boxes, box)
	return box
}
