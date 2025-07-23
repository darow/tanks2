package main

import (
	"bytes"
	"flag"
	"image"
	_ "image/png"
	"log"
	"syscall"

	images "myebiten/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	DEBUG_MODE               = false
	FEATURE_DECREASING_TANKS = false
)

var (
	REGULAR_FONT font.Face

	SCREEN_SIZE_WIDTH  = 2560
	SCREEN_SIZE_HEIGHT = 1420

	CHARACTER_IMAGE_TO_RESIZE image.Image

	CONNECTION_MODE             = flag.String("mode", "offline", "offline / server / client")
	SERVER_MODE_PORT            = flag.String("server_mode_port", "8080", "IF TRUE THEN GAME IS IN HOST MODE AND WAITING FOR CONNECTION OF OTHER PLAYER")
	CLIENT_CONNECT_MODE_ADDRESS = flag.String("client_connect_mode_address", "localhost:8080", "IF SETTED THEN GAME TRYING TO CONNECT TO HOST")
	SUCCESS_CONNECTION          bool
)

func main() {
	flag.Parse()
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setScreenSizeParams()

	ebiten.SetWindowSize(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	ebiten.SetWindowTitle("tanks in maze")

	var err error
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	REGULAR_FONT, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	CHARACTER_IMAGE_TO_RESIZE, _, err = image.Decode(bytes.NewReader(images.TankV2png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	charImage := ebiten.NewImageFromImage(resizedCharacterImage)

	cs1 := ControlSettings{
		rotateRightButton:  ebiten.KeyD,
		rotateLeftButton:   ebiten.KeyA,
		moveForwardButton:  ebiten.KeyW,
		moveBackwardButton: ebiten.KeyS,
		shootButton:        ebiten.KeySpace,
	}

	cs2 := ControlSettings{
		rotateRightButton:  ebiten.KeyRight,
		rotateLeftButton:   ebiten.KeyLeft,
		moveForwardButton:  ebiten.KeyUp,
		moveBackwardButton: ebiten.KeyDown,
		shootButton:        ebiten.KeySlash,
	}

	game := &Game{
		boardImage: ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT),

		boardSizeX: 7,
		boardSizeY: 4,

		things: Things{
			bullets: make(map[int]Bullet, 20),
			walls:   map[Wall]struct{}{},
		},

		characters: []*Character{
			{
				id: 0,
				input: Input{
					ControlSettings: cs1,
				},
				x: 400,
				y: 400,

				currentWidth: CHARACTER_WIDTH,
				charImg:      charImage,
			},
			{
				id: 1,
				input: Input{
					ControlSettings: cs2,
				},
				x: 700,
				y: 500,

				currentWidth: CHARACTER_WIDTH,
				charImg:      charImage,
			},
		},
		charactersScores: []uint{0, 0},
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
