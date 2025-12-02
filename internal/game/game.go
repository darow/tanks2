package game

import (
	"image"
	"image/color"

	"myebiten/internal/models"
	"myebiten/internal/websocket/client"
	"myebiten/internal/websocket/server"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	STATE_GAME_ENDING_TIMER_SECONDS = 4
	ITEM_SPAWN_INTERVAL             = 5
)

var (
	REGULAR_FONT              font.Face
	CHARACTER_IMAGE_TO_RESIZE image.Image

	SCREEN_SIZE_WIDTH  = 2560
	SCREEN_SIZE_HEIGHT = 1420
)

const (
	STATE_MAZE_CREATING = iota
	STATE_GAME_RUNNING
	STATE_GAME_ENDING
)

const (
	MENU_SCENE_ID  = 1
	LOBBY_SCENE_ID = 2
	MAIN_SCENE_ID  = 3
)

const (
	ROOT_AREA_ID         = "root_area"
	MAIN_PLAYING_AREA_ID = "main_playing_area"
	MAZE_AREA_ID         = "maze_area"
	UI_AREA1_ID          = "ui_area_1"
	SCORE_AREA_ID        = "score_area"
	SCORE_AREA_1_ID      = "score_area_1"
	SCORE_AREA_2_ID      = "score_area_2"
	SCORE_AREA_3_ID      = "score_area_3"
	SCORE_AREA_4_ID      = "score_area_4"
)

var noChars = true // zaglushka

var wallsToCheck []*models.Wall = make([]*models.Wall, 12)
var (
	COLOR_BLACK = color.RGBA{0x0f, 0x0f, 0x0f, 0xff}
)

type Game struct {
	server   *server.Server
	client   *client.Client
	connMode string

	scenes      map[int]models.Scene `json:"-"`
	activeScene models.Scene         `json:"-"`
}

func CreateGame(connectionMode, serverPort, address string) *Game {
	game := Game{}

	menuScene := &LobbyScene{}
	lobbyScene := &LobbyScene{}
	mainScene := CreateMainScene()

	mainScene.getConnectionMode = game.getConnectionMode
	mainScene.getGameClient = game.getClient
	mainScene.getGameServer = game.getServer

	game.scenes = make(map[int]models.Scene, 3)

	game.scenes[MENU_SCENE_ID] = menuScene
	game.scenes[LOBBY_SCENE_ID] = lobbyScene
	game.scenes[MAIN_SCENE_ID] = mainScene

	switch connectionMode {
	case CONNECTION_MODE_SERVER:
		game.server = server.New(serverPort)
	case CONNECTION_MODE_CLIENT:
		game.client = client.New(address)
		go mainScene.ReceiveMazeUpdates()
	default:
	}
	game.connMode = connectionMode

	game.SetActiveScene(MAIN_SCENE_ID)

	return &game
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	// Zaglushka
	if noChars {
		g.CreateCharacter(0)
		g.CreateCharacter(1)
		noChars = false
	}
	return g.activeScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	image := g.activeScene.Draw()

	screen.Clear()
	screen.DrawImage(image, &ebiten.DrawImageOptions{})
}

func (g *Game) getConnectionMode() string {
	return g.connMode
}

func (g *Game) getServer() *server.Server {
	return g.server
}

func (g *Game) getClient() *client.Client {
	return g.client
}

func (g *Game) CreateCharacter(id int) {
	g.scenes[MAIN_SCENE_ID].(*MainScene).CreateCharacter(id)
}

func (g *Game) SetActiveScene(sceneID int) {
	g.activeScene = g.scenes[sceneID]
}
