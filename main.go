package main

import (
	"image/color"
	"log"
	"os"
)
import "github.com/hajimehoshi/ebiten"

func update(screen *ebiten.Image) error {
	checkExit()
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	background(screen)
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
	if err := ebiten.Run(update, 800, 600, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}
