package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	xIncrease bool
	xDecrease bool
	yIncrease bool
	yDecrease bool
}

func (inp *Input) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		inp.xIncrease = true
	}
	if inp.xIncrease && inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		inp.xIncrease = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		inp.xDecrease = true
	}
	if inp.xDecrease && inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		inp.xDecrease = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		inp.yDecrease = true
	}
	if inp.yDecrease && inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		inp.yDecrease = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		inp.yIncrease = true
	}
	if inp.yIncrease && inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		inp.yIncrease = false
	}
}
