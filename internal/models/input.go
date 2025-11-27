package models

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ControlSettings struct {
	RotateRightButton  ebiten.Key
	RotateLeftButton   ebiten.Key
	MoveForwardButton  ebiten.Key
	MoveBackwardButton ebiten.Key
	ShootButton        ebiten.Key
}

type Input struct {
	RotateRight  bool
	RotateLeft   bool
	MoveForward  bool
	MoveBackward bool
	Shoot        bool

	ControlSettings
}

func (in *Input) Update() {
	if inpututil.IsKeyJustPressed(in.RotateRightButton) {
		in.RotateRight = true
	}
	if in.RotateRight && inpututil.IsKeyJustReleased(in.RotateRightButton) {
		in.RotateRight = false
	}

	if inpututil.IsKeyJustPressed(in.RotateLeftButton) {
		in.RotateLeft = true
	}
	if in.RotateLeft && inpututil.IsKeyJustReleased(in.RotateLeftButton) {
		in.RotateLeft = false
	}

	if inpututil.IsKeyJustPressed(in.MoveBackwardButton) {
		in.MoveBackward = true
	}
	if in.MoveBackward && inpututil.IsKeyJustReleased(in.MoveBackwardButton) {
		in.MoveBackward = false
	}

	if inpututil.IsKeyJustPressed(in.MoveForwardButton) {
		in.MoveForward = true
	}
	if in.MoveForward && inpututil.IsKeyJustReleased(in.MoveForwardButton) {
		in.MoveForward = false
	}

	if inpututil.IsKeyJustPressed(in.ShootButton) {
		in.Shoot = true
	}
}

func (in *Input) Reset() {
	in.RotateLeft = false
	in.RotateRight = false
	in.MoveBackward = false
	in.MoveForward = false
	in.Shoot = false
}
