package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"os"

	"myebiten/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	DEBUG_MODE = flag.Bool("debug", true, "true / false")

	CONNECTION_MODE  = flag.String("mode", "offline", "offline / server / client")
	SERVER_MODE_PORT = flag.String("server_mode_port", "8080", "IF TRUE THEN GAME IS IN HOST MODE AND WAITING FOR CONNECTION OF OTHER PLAYER")

	ADDRESS       = flag.String("address", "localhost:8080", "IF SET THEN GAME TRYING TO CONNECT TO HOST")
	PLAYERS_COUNT = flag.Int("players_count", game.DEFAULT_PLAYERS_COUNT, "PLAYERS COUNT FROM 2 TO 10")
	PLAYER_ID     = flag.Int("player_id", 1, "CLIENT PLAYER ID FROM 1 TO players_count-1")
)

func main() {
	flag.Parse()

	ebiten.SetTPS(300)
	if *CONNECTION_MODE == game.CONNECTION_MODE_CLIENT {
		fmt.Println("Running in client mode")
		ebiten.SetTPS(400)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	outputFileName := fmt.Sprintf("%s%d.txt", *CONNECTION_MODE, *PLAYER_ID)
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
	ebiten.SetFullscreen(false)

	tanksGame := game.CreateGame(*CONNECTION_MODE, *SERVER_MODE_PORT, *ADDRESS, *PLAYERS_COUNT, *PLAYER_ID)
	if err := ebiten.RunGame(tanksGame); err != nil {
		log.Fatal(err)
	}
}
