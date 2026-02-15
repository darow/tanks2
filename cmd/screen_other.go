//go:build !windows

package main

import "myebiten/internal/game"

func setScreenSizeParams() {
	// Default screen size for macOS 
	game.SCREEN_SIZE_WIDTH = 1400
	game.SCREEN_SIZE_HEIGHT = 890
}
