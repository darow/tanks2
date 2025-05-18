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
	SCREEN_SIZE_WIDTH  = 2000
	SCREEN_SIZE_HEIGHT = 1200
	CHARACTER_WIDTH    = 70
)

var (
	CHARACTER_IMAGE *ebiten.Image
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ebiten.SetWindowSize(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	ebiten.SetWindowTitle("tank")

	wall1 := Wall{
		x:          1,
		y:          1,
		horizontal: false,
	}
	wall2 := Wall{
		x:          2,
		y:          1,
		horizontal: true,
	}

	cs := ControlSettings{
		rotateRightButton:  ebiten.KeyRight,
		rotateLeftButton:   ebiten.KeyLeft,
		moveForwardButton:  ebiten.KeyUp,
		moveBackwardButton: ebiten.KeyDown,
		shootButton:        ebiten.KeySpace,
	}

	game := &Game{
		character: Character{
			ControlSettings: cs,
			x:               400,
			y:               400,
		},

		boardImage: ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT),
		tiles: Tiles{
			bullets: make(map[uint16]Bullet),
			walls:   map[Wall]struct{}{wall1: {}, wall2: {}},
		},
	}

	charImg, _, err := image.Decode(bytes.NewReader(images.Tank_png))
	if err != nil {
		log.Fatal(err)
	}
	resizedMouseImage := resize.Resize(CHARACTER_WIDTH, 0, charImg, resize.Lanczos3)
	CHARACTER_IMAGE = ebiten.NewImageFromImage(resizedMouseImage)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
