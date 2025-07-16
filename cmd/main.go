package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"syscall"

	images "myebiten/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

const (
	DEBUG_MODE               = false
	FEATURE_DECREASING_TANKS = false
)

var (
	SCREEN_SIZE_WIDTH  = 2560
	SCREEN_SIZE_HEIGHT = 1420

	CHARACTER_IMAGE_TO_RESIZE image.Image
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setScreenSizeParams()

	ebiten.SetWindowSize(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	ebiten.SetWindowTitle("tanks in maze")

	var err error
	CHARACTER_IMAGE_TO_RESIZE, _, err = image.Decode(bytes.NewReader(images.TankV2png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	charImage := ebiten.NewImageFromImage(resizedCharacterImage)

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
		boardImage: ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT),

		boardSizeX: 10,
		boardSizeY: 7,

		things: Things{
			bullets: make(map[int]Bullet, 20),
			walls: map[Wall]struct{}{
				Wall{
					x:          1,
					y:          1,
					horizontal: false,
				}: {},
				Wall{
					x:          2,
					y:          1,
					horizontal: true,
				}: {},
			},
		},

		characters: []*Character{
			{
				input: Input{
					ControlSettings: cs1,
				},
				x: 400,
				y: 400,

				currentWidth: CHARACTER_WIDTH,
				charImg:      charImage,
			},
			{
				input: Input{
					ControlSettings: cs2,
				},
				x: 700,
				y: 500,

				currentWidth: CHARACTER_WIDTH,
				charImg:      charImage,
			},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func setScreenSizeParams() {
	var (
		user32           = syscall.NewLazyDLL("User32.dll")
		getSystemMetrics = user32.NewProc("GetSystemMetrics")
	)

	GetSystemMetrics := func(nIndex int) int {
		index := uintptr(nIndex)
		ret, _, _ := getSystemMetrics.Call(index)
		return int(ret)
	}

	const (
		SM_CXSCREEN = 0
		SM_CYSCREEN = 1
	)

	SCREEN_SIZE_WIDTH = GetSystemMetrics(SM_CXSCREEN)
	SCREEN_SIZE_HEIGHT = GetSystemMetrics(SM_CYSCREEN) - 20
}
