package overlay

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	rbaseoverlay "github.com/mistweaverco/caline/resources/baseoverlay"
	pins "github.com/mistweaverco/caline/resources/pins"
)

const (
	totalPixels      = 752
	width            = 800
	minutesInDay     = 1440
	widthPerMinute   = float64(totalPixels) / float64(minutesInDay)
	height           = 24
	indicatorWidthPx = 3
)

var (
	baseoverlay     *ebiten.Image
	currentlocation *ebiten.Image
	startTime       time.Time
	endTime         time.Time
)

func getCurrentIndicatorPosition(startTime, endTime time.Time) (int, int, error) {
	// Get the current time
	now := time.Now()

	// Check if the current time is within the range
	if now.Before(startTime) || now.After(endTime) {
		return 0, 0, fmt.Errorf("current time is outside the specified range")
	}

	// Calculate the total duration in minutes
	totalDuration := endTime.Sub(startTime).Minutes()

	// Calculate the width per minute
	widthPerMinute := float64(totalPixels) / totalDuration

	// Calculate minutes since the start time
	minutesSinceStart := now.Sub(startTime).Minutes()

	// Calculate the x-coordinate
	xPos := int(minutesSinceStart * widthPerMinute)

	// Since we want a 3px wide indicator, adjust the position
	startPos := xPos - indicatorWidthPx/2
	endPos := startPos + indicatorWidthPx

	// Ensure the indicator is within bounds
	if startPos < 0 {
		startPos = 0
	}
	if endPos > totalPixels {
		endPos = totalPixels
	}

	return startPos, endPos, nil
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(rbaseoverlay.Baseoverlay_png))
	if err != nil {
		log.Fatal(err)
	}
	baseoverlay = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(pins.CurrentLocation))
	if err != nil {
		log.Fatal(err)
	}
	currentlocation = ebiten.NewImageFromImage(img)
}

func init() {
	// Nothing here
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

	// Center the overlay on the screen
	// Also on the top of the screen
	ebiten.SetWindowPosition((sw-width)/2, 0)
	// Check if mouse is pressed
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Println("Mouse pressed")
		fmt.Println(ebiten.CursorPosition())
		fmt.Println(currentlocation.Bounds())
	}
	return nil
}

func (m *overlay) Draw(screen *ebiten.Image) {
	baseoverlayOptions := &ebiten.DrawImageOptions{}
	pos, _, err := getCurrentIndicatorPosition(startTime, endTime)
	if err == nil {
		currentlocationOptions := &ebiten.DrawImageOptions{}
		currentlocationOptions.GeoM.Translate(float64(pos), 0)
		baseoverlay.DrawImage(currentlocation, currentlocationOptions)
	}
	screen.DrawImage(baseoverlay, baseoverlayOptions)
}

func Start() {
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(width, height)
	// Set to true to pass through mouse events to the window below
	// but then we can't capture mouse events
	ebiten.SetWindowMousePassthrough(false)

	currentYear, currentMonth, currentDay := time.Now().Local().Date()
	// TODO: Make this configurable via a config file (maybe caline.yml)
	startTime = time.Date(currentYear, currentMonth, currentDay, 7, 0, 0, 0, time.Local)
	endTime = time.Date(currentYear, currentMonth, currentDay, 19, 0, 0, 0, time.Local)

	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true
	if err := ebiten.RunGameWithOptions(&overlay{}, op); err != nil {
		log.Fatal(err)
	}
}
