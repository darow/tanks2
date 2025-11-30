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

var noChars = false

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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	// Zaglushka
	if noChars {
		g.CreateCharacter(0)
		g.CreateCharacter(1)
		noChars = true
	}
	return g.activeScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	image := g.activeScene.Draw()

	screen.Clear()
	screen.DrawImage(image, &ebiten.DrawImageOptions{})
}

func CreateGame(scenes map[int]models.Scene) *Game {
	return &Game{
		scenes: scenes,
	}
}
