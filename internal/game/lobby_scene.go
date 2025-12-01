package game

import (
	"myebiten/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
)

type LobbyScene struct {
	models.SceneUI
}

func (menuScene *LobbyScene) Update() error {
	return nil
}

func (menuScene *LobbyScene) Draw() *ebiten.Image {
	return nil
}
