package images

import (
	_ "embed"
)

var (
	//go:embed tank1.png
	Tank_png []byte

	//go:embed TankV2.png
	TankV2png []byte
)
