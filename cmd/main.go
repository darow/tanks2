package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"syscall"

	"myebiten/internal/game"
	"myebiten/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
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

	ebiten.SetWindowSize(game.SCREEN_SIZE_WIDTH, game.SCREEN_SIZE_HEIGHT)
	ebiten.SetWindowTitle("tanks in maze")

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	game.REGULAR_FONT, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	menuScene := &game.LobbyScene{}
	lobbyScene := &game.LobbyScene{}
	mainScene := game.CreateMainScene()

	scenes := map[int]models.Scene{
		game.MENU_SCENE_ID:  menuScene,
		game.LOBBY_SCENE_ID: lobbyScene,
		game.MAIN_SCENE_ID:  mainScene,
	}

	tanksGame := game.CreateGame(scenes)

	tanksGame.SetActiveScene(game.MAIN_SCENE_ID)
	tanksGame.MakeSuccessConnection(*CONNECTION_MODE, *SERVER_MODE_PORT, *ADDRESS)

	ebiten.SetFullscreen(false)

	if err := ebiten.RunGame(tanksGame); err != nil {
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

	game.SCREEN_SIZE_WIDTH = GetSystemMetrics(SM_CXMAXIMIZED)
	game.SCREEN_SIZE_HEIGHT = GetSystemMetrics(SM_CYMAXIMIZED)
}
