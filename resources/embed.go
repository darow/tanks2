package images

import (
	_ "embed"
)

var (
	//go:embed ebiten.png
	Ebiten_png []byte

	//go:embed m2.png
	M2_png []byte

	//go:embed tank1.png
	Tank_png []byte

	//go:embed img.png
	Img_png []byte
)
