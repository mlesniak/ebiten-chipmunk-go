// An example of using Chipmunk with Ebiten.
package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
)

const width = 800
const height = 600
const boxWidth = 40
const boxHeight = 40

type Box struct {
	x, y float64
	w, h float64
	r    float64
}

var space *cp.Space // Simulation space
var boxImage *ebiten.Image
var lastClicked = 0 // Frames to wait until the next mouse click is registered.

func init() {
	boxImage, _, _ = ebitenutil.NewImageFromFile("box.png", ebiten.FilterLinear)
}

func main() {
	initPhysicsEngine()
	addFloor()

	if err := ebiten.Run(update, width, height, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}

func initPhysicsEngine() {
	space = cp.NewSpace()
	space.Iterations = 1000
	space.SetGravity(cp.Vector{0, 500})
}

func update(screen *ebiten.Image) error {
	// Input handling.
	checkExit()
	checkSpawnBox()

	// State progression.
	nextStep()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// View update.
	drawBackground(screen)
	drawBoxes(screen)

	return nil
}

// nextStep computes the updated positions of all boxes.
func nextStep() {
	space.Step(1.0 / float64(ebiten.MaxTPS()))
}

// checkExit checks if Escape was pressed and exits hard.
func checkExit() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
}

func checkSpawnBox() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && lastClicked < 0 {
		lastClicked = 20 // FPS fix
		x, y := ebiten.CursorPosition()
		b := randomBox(float64(x), float64(y))
		addBox(b)
	} else {
		lastClicked--
	}
}

// drawBackground fills the background with a color.
func drawBackground(screen *ebiten.Image) error {
	return screen.Fill(color.Gray{Y: 30})
}

func drawBoxes(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	space.EachBody(func(body *cp.Body) {
		switch body.UserData {
		case "box":
			angle := body.Rotation().ToAngle()
			angle = angle / 180 * math.Pi
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(boxWidth)/2, -float64(boxHeight)/2)
			op.GeoM.Rotate(body.Angle())
			op.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(boxImage, op)
		case "line":
			ebitenutil.DrawLine(screen, 0, body.Position().Y, width, body.Position().Y, colornames.Gray)
		}
	})
}

// addFloor adds the engine object for the floor.
func addFloor() {
	deltaY := 20.
	bf := cp.NewStaticBody()
	bf.SetPosition(cp.Vector{width / 2, height - deltaY})
	bf.UserData = "line"
	sf := cp.NewBox(bf, width, 5, 0.0)
	sf.SetFriction(1.0)
	sf.SetElasticity(0.0)
	space.AddBody(bf)
	space.AddShape(sf)
}

// addBox adds a new Box to the engine system.
func addBox(box Box) *cp.Body {
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

// randomBox creates a new box at the given position with a random rotation.
func randomBox(px, py float64) Box {
	return Box{
		x: px,
		y: py,
		w: boxWidth,
		h: boxHeight,
		r: rand.Float64() * 2 * math.Pi,
	}
}
