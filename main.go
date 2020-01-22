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

func init() {
	boxImage, _, _ = ebitenutil.NewImageFromFile("box.png", ebiten.FilterDefault)
}

func main() {

	space = cp.NewSpace()
	space.Iterations = 1000
	space.SetGravity(cp.Vector{0, 500})

	addFloor()

	if err := ebiten.Run(update, width, height, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	// Input handling.
	checkExit()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && lastClicked < 0 {
		lastClicked = 20 // FPS fix
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
			angle := body.Rotation().ToAngle()
			angle = angle / 180 * math.Pi
			//fmt.Printf("BOX (%.2f/%.2f);%.2f\n", body.Position().X, body.Position().Y, angle)
			op.GeoM.Reset()
			// Center of image
			op.GeoM.Translate(-float64(boxWidth)/2, -float64(boxHeight)/2)
			op.GeoM.Rotate(body.Angle())
			op.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(boxImage, op)
		case "line":
			ebitenutil.DrawLine(screen, 0, body.Position().Y, width, body.Position().Y, colornames.Gray)
			//ebitenutil.DrawLine(screen, 0, height-20, width, height-20, colornames.Yellow)
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
	deltaY := 20.
	bf := cp.NewStaticBody()
	bf.SetPosition(cp.Vector{width / 2, height - deltaY})
	bf.UserData = "line"
	sf := cp.NewBox(bf, width, 2, 0.0)
	sf.SetFriction(1.0)
	sf.SetElasticity(0.0)
	space.AddBody(bf)
	space.AddShape(sf)
}

func addBoxToPhysics(box Box) *cp.Body {
	body := cp.NewBody(1.0, 1)
	body.SetPosition(cp.Vector{X: box.x, Y: box.y})
	rad := (rand.Float64() * 360) * 180 / math.Pi
	body.SetAngle(rad)
	body.UserData = "box"

	shape := cp.NewBox(body, box.w, box.h, 0.0)
	shape.SetFriction(1.0)
	shape.SetElasticity(0.0)
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
