package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	rotateRight  bool
	rotateLeft   bool
	moveForward  bool
	moveBackward bool

	ControlSettings
}

func (inp *Input) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		inp.rotateRight = true
	}
	if inp.rotateRight && inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		inp.rotateRight = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		inp.rotateLeft = true
	}
	if inp.rotateLeft && inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		inp.rotateLeft = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		inp.moveBackward = true
	}
	if inp.moveBackward && inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		inp.moveBackward = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		inp.moveForward = true
	}
	if inp.moveForward && inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		inp.moveForward = false
	}
}
