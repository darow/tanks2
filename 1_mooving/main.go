package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/png"
	"log"
	images "myebiten/resources"
)

const screenSizeHeight = 1200
const screenSizeWidth = 2000

var (
	characterImage *ebiten.Image
)

type Game struct {
	boardImage *ebiten.Image
	x          float64
	y          float64

	xIncrease bool
	xDecrease bool
	yIncrease bool
	yDecrease bool
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.xIncrease = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		g.xIncrease = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.xDecrease = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		g.xDecrease = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.yDecrease = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		g.yDecrease = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.yIncrease = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		g.yIncrease = false
	}

	if g.xIncrease {
		g.x++
	}
	if g.xDecrease {
		g.x--
	}

	if g.yIncrease {
		g.y++
	}
	if g.yDecrease {
		g.y--
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(screenSizeWidth, screenSizeHeight)
	}

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(1)
	op.GeoM.Reset()

	x := g.x
	y := g.y
	op.GeoM.Translate(x, y)

	g.boardImage.Clear()
	g.boardImage.Fill(color.RGBA{0xff, 0xff, 0xca, 0xff})
	g.boardImage.DrawImage(characterImage, op)

	screen.Clear()
	screen.DrawImage(g.boardImage, &ebiten.DrawImageOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeWidth, screenSizeHeight
}

func main() {
	ebiten.SetWindowSize(screenSizeWidth, screenSizeHeight)
	ebiten.SetWindowTitle("title")
	game := &Game{}

	//path, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(path)
	//f, err := os.Open("./resources/m2.png")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//img, _, err := image.Decode(f)
	//if err != nil {
	//	log.Fatal(err)
	//}

	img, _, err := image.Decode(bytes.NewReader(images.M2_png))
	if err != nil {
		log.Fatal(err)
	}
	resizedImage := resize.Resize(100, 0, img, resize.Lanczos3)

	characterImage = ebiten.NewImageFromImage(resizedImage)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
