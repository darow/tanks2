package main

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	RotateRight  bool
	RotateLeft   bool
	MoveForward  bool
	MoveBackward bool
	Shoot        bool

	ControlSettings
}

func (in *Input) Update() {
	if inpututil.IsKeyJustPressed(in.rotateRightButton) {
		in.RotateRight = true
	}
	if in.RotateRight && inpututil.IsKeyJustReleased(in.rotateRightButton) {
		in.RotateRight = false
	}

	if inpututil.IsKeyJustPressed(in.rotateLeftButton) {
		in.RotateLeft = true
	}
	if in.RotateLeft && inpututil.IsKeyJustReleased(in.rotateLeftButton) {
		in.RotateLeft = false
	}

	if inpututil.IsKeyJustPressed(in.moveBackwardButton) {
		in.MoveBackward = true
	}
	if in.MoveBackward && inpututil.IsKeyJustReleased(in.moveBackwardButton) {
		in.MoveBackward = false
	}

	if inpututil.IsKeyJustPressed(in.moveForwardButton) {
		in.MoveForward = true
	}
	if in.MoveForward && inpututil.IsKeyJustReleased(in.moveForwardButton) {
		in.MoveForward = false
	}

	if inpututil.IsKeyJustPressed(in.shootButton) {
		in.Shoot = true
	}
}
