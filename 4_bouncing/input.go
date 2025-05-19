package main

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	rotateRight  bool
	rotateLeft   bool
	moveForward  bool
	moveBackward bool

	ControlSettings
}

func (in *Input) Update() {
	if inpututil.IsKeyJustPressed(in.rotateRightButton) {
		in.rotateRight = true
	}
	if in.rotateRight && inpututil.IsKeyJustReleased(in.rotateRightButton) {
		in.rotateRight = false
	}

	if inpututil.IsKeyJustPressed(in.rotateLeftButton) {
		in.rotateLeft = true
	}
	if in.rotateLeft && inpututil.IsKeyJustReleased(in.rotateLeftButton) {
		in.rotateLeft = false
	}

	if inpututil.IsKeyJustPressed(in.moveForwardButton) {
		in.moveBackward = true
	}
	if in.moveBackward && inpututil.IsKeyJustReleased(in.moveForwardButton) {
		in.moveBackward = false
	}

	if inpututil.IsKeyJustPressed(in.moveBackwardButton) {
		in.moveForward = true
	}
	if in.moveForward && inpututil.IsKeyJustReleased(in.moveBackwardButton) {
		in.moveForward = false
	}
}
