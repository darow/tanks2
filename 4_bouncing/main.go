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
	CHARACTER_IMAGE_TO_RESIZE image.Image
	CHARACTER_IMAGE           *ebiten.Image
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

	var err error
	CHARACTER_IMAGE_TO_RESIZE, _, err = image.Decode(bytes.NewReader(images.Tank_png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	CHARACTER_IMAGE = ebiten.NewImageFromImage(resizedCharacterImage)

	cs2 := ControlSettings{
		rotateRightButton:  ebiten.KeyD,
		rotateLeftButton:   ebiten.KeyA,
		moveForwardButton:  ebiten.KeyW,
		moveBackwardButton: ebiten.KeyS,
		shootButton:        ebiten.KeySpace,
	}

	cs1 := ControlSettings{
		rotateRightButton:  ebiten.KeyRight,
		rotateLeftButton:   ebiten.KeyLeft,
		moveForwardButton:  ebiten.KeyUp,
		moveBackwardButton: ebiten.KeyDown,
		shootButton:        ebiten.KeySlash,
	}

	game := &Game{
		characters: []*Character{
			{
				input: Input{
					ControlSettings: cs1,
				},
				x: 400,
				y: 400,

				currentWidth: CHARACTER_WIDTH,
				charImg:      CHARACTER_IMAGE,
			},
			{
				input: Input{
					ControlSettings: cs2,
				},
				x: 700,
				y: 500,

				currentWidth: CHARACTER_WIDTH,
				charImg:      CHARACTER_IMAGE,
			},
		},

		boardImage: ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT),
		tiles: Tiles{
			bullets: make(map[int]Bullet),
			walls:   map[Wall]struct{}{wall1: {}, wall2: {}},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
