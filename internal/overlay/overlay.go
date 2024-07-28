package overlay

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	rbaseoverlay "github.com/mistweaverco/caline/resources/baseoverlay"
)

const (
	width  = 800
	height = 24
)

var (
	baseoverlay *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img1, _, err := image.Decode(bytes.NewReader(rbaseoverlay.Baseoverlay_png))
	if err != nil {
		log.Fatal(err)
	}
	baseoverlay = ebiten.NewImageFromImage(img1)
}

func init() {
	// nothing here
}

type overlay struct {
	x16  int
	y16  int
	vx16 int
	vy16 int
}

func (m *overlay) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (m *overlay) Update() error {
	sw, _ := ebiten.Monitor().Size()

	// center the overlay on the screen
	// it also should be top most
	ebiten.SetWindowPosition((sw-width)/2, 0)
	return nil
}

func (m *overlay) Draw(screen *ebiten.Image) {
	img := baseoverlay
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}

func Start() {
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowMousePassthrough(true)

	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true
	if err := ebiten.RunGameWithOptions(&overlay{}, op); err != nil {
		log.Fatal(err)
	}
}
