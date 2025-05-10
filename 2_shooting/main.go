package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	images "myebiten/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

const (
	SCREEN_SIZE_WIDTH     = 2000
	SCREEN_SIZE_HEIGHT    = 1200
	MOUSE_CHARACTER_WIDTH = 70
)

var (
	MOUSE_CHARACTER_IMAGE *ebiten.Image
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ebiten.SetWindowSize(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	ebiten.SetWindowTitle("shooting_try")
	game := &Game{
		x:          400,
		y:          400,
		boardImage: ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT),
		tiles:      Tiles{bullets: make(map[uint16]Bullet)},
	}

	mouseImg, _, err := image.Decode(bytes.NewReader(images.M2_png))
	if err != nil {
		log.Fatal(err)
	}
	resizedMouseImage := resize.Resize(MOUSE_CHARACTER_WIDTH, 0, mouseImg, resize.Lanczos3)
	MOUSE_CHARACTER_IMAGE = ebiten.NewImageFromImage(resizedMouseImage)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
