package main

import (
	"bytes"
	"flag"
	"image"
	_ "image/png"
	"log"
	"syscall"

	"myebiten/internal/models"
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

	MAX_BOARD_HEIGHT = 7
	MAX_BOARD_WIDTH  = 12

	MIN_BOARD_HEIGHT = 3
	MIN_BOARD_WIDTH  = 3

	DRAWING_OFFSET = models.Vector2D{X: 300, Y: 50}
	DRAWING_SCALE  = 1.0

	CHARACTER_IMAGE_TO_RESIZE image.Image

	CONNECTION_MODE  = flag.String("mode", "offline", "offline / server / client")
	SERVER_MODE_PORT = flag.String("server_mode_port", "8080", "IF TRUE THEN GAME IS IN HOST MODE AND WAITING FOR CONNECTION OF OTHER PLAYER")

	ADDRESS            = flag.String("address", "localhost:8080", "IF SET THEN GAME TRYING TO CONNECT TO HOST")
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

	bullets := make([]*models.Bullet, 20)

	for i := range bullets {
		bullets[i] = &models.Bullet{
			R: float64(BULLET_RADIUS),
		}
	}

	characters := []*Character{
		{
			GameObject: models.GameObject{
				ID:       0,
				Active:   true,
				Position: models.Vector2D{X: 400, Y: 400},
				Rotation: 0.0,
			},

			Hitbox: RectangleHitbox{CHARACTER_WIDTH, CHARACTER_WIDTH},
			Sprite: ImageSprite{charImage},
			weapon: &DefaultWeapon{bullets, 5},
			input: Input{
				ControlSettings: cs1,
			},
		},
		{
			GameObject: models.GameObject{
				ID:       1,
				Active:   true,
				Position: models.Vector2D{X: 700, Y: 500},
				Rotation: 0.0,
			},

			Hitbox: RectangleHitbox{CHARACTER_WIDTH, CHARACTER_WIDTH},
			Sprite: ImageSprite{charImage},
			weapon: &DefaultWeapon{bullets, 5},
			input: Input{
				ControlSettings: cs2,
			},
		},
	}

	image := ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)

	mainArea := &models.DrawingArea{
		BoardImage: image,
		DrawingSettings: models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: float64(SCREEN_SIZE_HEIGHT) / 10},
			Scale:  1.0,
		},
		Height: float64(SCREEN_SIZE_HEIGHT) * 0.8,
		Width:  float64(SCREEN_SIZE_WIDTH),
	}

	UIArea1 := &models.DrawingArea{
		BoardImage: image,
		DrawingSettings: models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: 0.0},
			Scale:  1.0,
		},
		Height: float64(SCREEN_SIZE_HEIGHT) * 0.1,
		Width:  float64(SCREEN_SIZE_WIDTH),
	}

	UIArea2 := &models.DrawingArea{
		BoardImage: image,
		DrawingSettings: models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: float64(SCREEN_SIZE_HEIGHT) * 0.9},
			Scale:  1.0,
		},
		Height: float64(SCREEN_SIZE_HEIGHT) * 0.1,
		Width:  float64(SCREEN_SIZE_WIDTH),
	}

	game := &Game{
		boardImage:       image,
		leftAlive:        2,
		Bullets:          bullets,
		Characters:       characters,
		CharactersScores: []uint{0, 0},
		mainArea:         mainArea,
		UIArea1:          UIArea1,
		UIArea2:          UIArea2,
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
	// ebiten.SetFullscreen(true)

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
		SM_CXMAXIMIZED = 0 // Width of maximized window
		SM_CYMAXIMIZED = 1 // Height of maximized window
	)

	SCREEN_SIZE_WIDTH = GetSystemMetrics(SM_CXMAXIMIZED)
	SCREEN_SIZE_HEIGHT = GetSystemMetrics(SM_CYMAXIMIZED)
}
