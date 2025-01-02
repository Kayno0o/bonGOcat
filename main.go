package main

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	hook "github.com/robotn/gohook"
)

//go:embed assets/*.png
var embeddedImages embed.FS

var (
	bw         = 457
	bh         = 289
	ratio      = 0.2
	registered = false
	left       *ebiten.Image
	middle     *ebiten.Image
	right      *ebiten.Image

	resetChan = make(chan bool)

	position        = 0
	keyPresses      int
	mutex           sync.Mutex
	calculateTicker = time.NewTicker(time.Second)

	// Sliding window to store keypress counts
	keypressCounts = make([]int, 5)
	currentIndex   = 0
)

type Game struct{}

func timer(duration time.Duration, resetChan chan bool) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			position = 0
		case <-resetChan:
			ticker.Reset(duration)
		}
	}
}

func keypress(resetChan chan bool) {
	hook.Register(hook.MouseDown, nil, func(e hook.Event) {
		if !registered {
			return
		}

		// Handle mouse clicks
		if e.Button == 1 {
			position = 1 // Left click sets the position to left
		} else if e.Button == 3 {
			position = -1 // Right click sets the position to right
		}
		resetChan <- true
		registered = false
	})

	hook.Register(hook.KeyDown, nil, func(e hook.Event) {
		if !registered {
			return
		}

		if position == 1 {
			position = -1
		} else {
			position = 1
		}
		resetChan <- true

		mutex.Lock()
		keyPresses++
		mutex.Unlock()

		registered = false
	})

	s := hook.Start()
	<-hook.Process(s)
}

func calculateAverageKPS() {
	for range calculateTicker.C {
		mutex.Lock()
		// Update sliding window with the latest keypress count
		keypressCounts[currentIndex] = keyPresses
		currentIndex = (currentIndex + 1) % len(keypressCounts)

		// Reset the current second's keypress count
		keyPresses = 0
		mutex.Unlock()
	}
}

func getAverageKPS() float64 {
	mutex.Lock()
	defer mutex.Unlock()
	total := 0
	for _, count := range keypressCounts {
		total += count
	}
	return float64(total) / float64(len(keypressCounts))
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0})
	img := middle
	if position == -1 {
		img = left
	}
	if position == 1 {
		img = right
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ratio, ratio)
	screen.DrawImage(img, op)
	registered = true

	// Draw the average on the top-right corner as plain text
	average := strconv.FormatFloat(getAverageKPS(), 'f', 2, 64)
	ebitenutil.DebugPrintAt(screen, average, int(float64(bw)*ratio)-50, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float64(bw) * ratio), int(float64(bh) * ratio)
}

func readImg(path string) *ebiten.Image {
	imgBytes, err := embeddedImages.ReadFile(path)
	if err != nil {
		panic(err)
	}
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(imgBytes))
	if err != nil {
		panic(err)
	}

	return img
}

func main() {
	if len(os.Args) >= 2 {
		ratioStr := os.Args[1]
		if ratioFloat, err := strconv.ParseFloat(ratioStr, 64); err == nil {
			ratio = ratioFloat
		}
	}

	go timer(time.Second, resetChan)
	go keypress(resetChan)
	go calculateAverageKPS()

	left = readImg("assets/left.png")
	middle = readImg("assets/middle.png")
	right = readImg("assets/right.png")

	// Create the window
	ebiten.SetWindowSize(int(float64(bw)*ratio), int(float64(bh)*ratio))
	ebiten.SetWindowTitle("BONGO CAT")
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowIcon([]image.Image{middle})

	game := Game{}
	options := ebiten.RunGameOptions{
		ScreenTransparent: true,
	}

	// Start the game loop
	if err := ebiten.RunGameWithOptions(&game, &options); err != nil {
		panic(err)
	}
}
