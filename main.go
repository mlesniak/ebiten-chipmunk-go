package main

import (
	"log"
	"os"
)
import "github.com/hajimehoshi/ebiten"

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	return nil
}

func main() {
	if err := ebiten.Run(update, 800, 600, 1, "Physics Demo"); err != nil {
		log.Fatal(err)
	}
}
