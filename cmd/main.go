package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"syscall"

	"myebiten/internal/game"
	"myebiten/internal/models"
	"myebiten/internal/weapons"
	images "myebiten/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	NUMBER_OF_CHARACTERS = 4

	REGULAR_FONT              font.Face
	CHARACTER_IMAGE_TO_RESIZE image.Image

	SCREEN_SIZE_WIDTH  = 2560
	SCREEN_SIZE_HEIGHT = 1420

	DEBUG_MODE = flag.Bool("debug", true, "true / false")

	CONNECTION_MODE  = flag.String("mode", "offline", "offline / server / client")
	SERVER_MODE_PORT = flag.String("server_mode_port", "8080", "IF TRUE THEN GAME IS IN HOST MODE AND WAITING FOR CONNECTION OF OTHER PLAYER")

	ADDRESS = flag.String("address", "localhost:8080", "IF SET THEN GAME TRYING TO CONNECT TO HOST")
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	outputFileName := fmt.Sprintf("%s.txt", *CONNECTION_MODE)
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Recovered from panic: ", r)
		}
	}()

	setScreenSizeParams()

	ebiten.SetWindowSize(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)
	ebiten.SetWindowTitle("tanks in maze")

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

	bullets := make([]*models.Bullet, weapons.BULLETS_COUNT*NUMBER_OF_CHARACTERS)
	for i := range bullets {
		bullets[i] = models.CreateBullet(models.BULLET_RADIUS)
	}

	characters := createCharacters(bullets)

	UIScore1 := models.UIText{}
	UIScore2 := models.UIText{}

	mainScene := buildMainScene(UIScore1, UIScore2)
	scenes := map[int]*models.Scene{game.MAIN_SCENE_ID: mainScene}

	tanksGame := game.CreateGame(bullets, characters, scenes)

	tanksGame.SetActiveScene(game.MAIN_SCENE_ID)
	tanksGame.MakeSuccessConnection(*CONNECTION_MODE, *SERVER_MODE_PORT, *ADDRESS)

	ebiten.SetFullscreen(false)

	if err := ebiten.RunGame(tanksGame); err != nil {
		log.Fatal(err)
	}
}

func buildMainScene(score1, score2 models.Drawable) *models.Scene {
	ebitenImage := ebiten.NewImage(SCREEN_SIZE_WIDTH, SCREEN_SIZE_HEIGHT)

	scene := models.CreateScene(ebitenImage, float64(SCREEN_SIZE_HEIGHT), float64(SCREEN_SIZE_WIDTH))

	rootArea := scene.GetRootArea()

	mainArea := rootArea.NewArea(
		rootArea.Height*0.8,
		rootArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: rootArea.Height / 10},
			Scale:  1.0,
		})
	scene.AddDrawingArea(game.MAIN_PLAYING_AREA_ID, mainArea)

	UIArea1 := rootArea.NewArea(
		rootArea.Height*0.1,
		rootArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: 0.0},
			Scale:  1.0,
		})
	scene.AddDrawingArea(game.UI_AREA1_ID, UIArea1)

	ScoreArea := rootArea.NewArea(
		rootArea.Height*0.1,
		rootArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.0, Y: rootArea.Height * 0.9},
			Scale:  1.0,
		})
	scene.AddDrawingArea(game.SCORE_AREA_ID, ScoreArea)

	scoreArea1 := ScoreArea.NewArea(
		0.99*ScoreArea.Height,
		0.2*ScoreArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.2 * ScoreArea.Width, Y: 0.5 * ScoreArea.Height},
			Scale:  1.0,
		})
	scene.AddDrawingArea(game.SCORE_AREA_1_ID, scoreArea1)

	scoreArea2 := ScoreArea.NewArea(
		0.99*ScoreArea.Height,
		0.2*ScoreArea.Width,
		models.DrawingSettings{
			Offset: models.Vector2D{X: 0.6 * ScoreArea.Width, Y: 0.5 * ScoreArea.Height},
			Scale:  1.0,
		})
	scene.AddDrawingArea(game.SCORE_AREA_2_ID, scoreArea2)

	scene.AddObject(score1, game.SCORE_AREA_1_ID)
	scene.AddObject(score2, game.SCORE_AREA_2_ID)

	return scene
}

func createCharacters(bullets []*models.Bullet) []*models.Character {
	chars := make([]*models.Character, 0, 2)

	CHARACTER_IMAGE_TO_RESIZE, _, err := image.Decode(bytes.NewReader(images.TankV2png))
	if err != nil {
		log.Fatal(err)
	}
	resizedCharacterImage := resize.Resize(models.CHARACTER_WIDTH, 0, CHARACTER_IMAGE_TO_RESIZE, resize.Lanczos3)
	charImage := ebiten.NewImageFromImage(resizedCharacterImage)

	cs1 := models.ControlSettings{
		RotateRightButton:  ebiten.KeyF,
		RotateLeftButton:   ebiten.KeyS,
		MoveForwardButton:  ebiten.KeyE,
		MoveBackwardButton: ebiten.KeyD,
		ShootButton:        ebiten.KeySpace,
	}

	defaultWeapon := weapons.DefaultWeapon{
		Clip:     bullets[:weapons.BULLETS_COUNT],
		Cooldown: 5,
	}

	char1 := models.CreateCharacter(0, charImage, &defaultWeapon, cs1)
	chars = append(chars, &char1)

	cs2 := models.ControlSettings{
		RotateRightButton:  ebiten.KeyRight,
		RotateLeftButton:   ebiten.KeyLeft,
		MoveForwardButton:  ebiten.KeyUp,
		MoveBackwardButton: ebiten.KeyDown,
		ShootButton:        ebiten.KeySlash,
	}

	defaultWeapon = weapons.DefaultWeapon{
		Clip:     bullets[weapons.BULLETS_COUNT:],
		Cooldown: 5,
	}

	char2 := models.CreateCharacter(0, charImage, &defaultWeapon, cs2)
	chars = append(chars, &char2)

	return chars
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
