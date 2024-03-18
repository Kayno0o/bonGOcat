package main

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"os"
	"strconv"
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

	position = 0
)

type Game struct {
}

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
		registered = false
	})

	s := hook.Start()
	<-hook.Process(s)
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float64(bw) * ratio), int(float64(bh) * ratio)
}

func readImg(path string) (*ebiten.Image, error) {
	imgBytes, err := embeddedImages.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}

	return img, nil
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

	img, err := readImg("assets/left.png")
	if err != nil {
		panic(err)
	}
	left = img

	img, err = readImg("assets/middle.png")
	if err != nil {
		panic(err)
	}
	middle = img

	img, err = readImg("assets/right.png")
	if err != nil {
		panic(err)
	}
	right = img

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
