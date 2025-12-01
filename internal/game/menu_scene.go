package game

import (
	"myebiten/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScene struct {
	models.SceneUI
}

func (menuScene *MenuScene) Update() error {
	return nil
}

func (menuScene *MenuScene) Draw() *ebiten.Image {
	return nil
}
