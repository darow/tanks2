//go:build !windows

package main

import "myebiten/internal/game"

func setScreenSizeParams() {
	// Default screen size for macOS (and Linux): 1920x1080
	game.SCREEN_SIZE_WIDTH = 1920
	game.SCREEN_SIZE_HEIGHT = 1080
}
