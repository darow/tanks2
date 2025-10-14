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

	DRAWING_OFFSET_X = 300
	DRAWING_OFFSET_Y = 50
	DRAWING_SCALE    = 0.8

	CHARACTER_IMAGE_TO_RESIZE image.Image

	CONNECTION_MODE  = flag.String("mode", "offline", "offline / server / client")
	SERVER_MODE_PORT = flag.String("server_mode_port", "8080", "IF TRUE THEN GAME IS IN HOST MODE AND WAITING FOR CONNECTION OF OTHER PLAYER")

	ADDRESS            = flag.String("address", "localhost:8080", "IF SETTED THEN GAME TRYING TO CONNECT TO HOST")
	SUCCESS_CONNECTION bool
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
	if err != nil {
		log.Fatal(err)
	}

	CHARACTER_IMAGE_TO_RESIZE, _, err = image.Decode(bytes.NewReader(images.TankV2png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	charImage := ebiten.NewImageFromImage(resizedCharacterImage)

	cs1 := ControlSettings{
		rotateRightButton:  ebiten.KeyF,
		rotateLeftButton:   ebiten.KeyS,
		moveForwardButton:  ebiten.KeyE,
		moveBackwardButton: ebiten.KeyD,
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

		boardSizeX: 6,
		boardSizeY: 4,

		leftAlive: 2,

		Bullets: make([]*Bullet, 10),

		Characters: []*Character{
			{
				GameObject: GameObject{
					id:       0,
					active:   true,
					position: Vector2D{400, 400},
					rotation: 0.0,
				},

				hitbox: RectangleHitbox{},
				sprite: ImageSprite{charImage},
				input: Input{
					ControlSettings: cs1,
				},
			},
			{
				GameObject: GameObject{
					id:       1,
					active:   true,
					position: Vector2D{700, 500},
					rotation: 0.0,
				},

				hitbox: RectangleHitbox{},
				sprite: ImageSprite{charImage},
				input: Input{
					ControlSettings: cs2,
				},
			},
		},
		CharactersScores: []uint{0, 0},
	}

	for i := range game.Bullets {
		game.Bullets[i] = &Bullet{
			hitbox: CircleHitbox{BULLET_RADIUS},
			sprite: BallSprite{BULLET_RADIUS},
		}
	}

	for _, char := range game.Characters {
		char.weapon = DefaultWeapon{game.Bullets, 5}
	}

	if *CONNECTION_MODE != CONNECTION_MODE_OFFLINE && !SUCCESS_CONNECTION {
		game.makeSuccessConnection()
		//if *CONNECTION_MODE == CONNECTION_MODE_SERVER {
		//	game.Characters[0].input.ControlSettings = ControlSettings{}
		//}
		//if *CONNECTION_MODE == CONNECTION_MODE_CLIENT {
		//	game.Characters[1].input.ControlSettings = ControlSettings{}
		//}
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
	SCREEN_SIZE_HEIGHT = GetSystemMetrics(SM_CYSCREEN)
}
