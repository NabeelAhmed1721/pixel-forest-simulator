package main

import (
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Forest Simulator",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	treeSpritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	var treeFrames []pixel.Rect

	for x := treeSpritesheet.Bounds().Min.X; x < treeSpritesheet.Bounds().Max.X; x += 32 {
		for y := treeSpritesheet.Bounds().Min.Y; y < treeSpritesheet.Bounds().Max.Y; y += 32 {
			treeFrames = append(treeFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 5000.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		trees        []*pixel.Sprite
		matrices     []pixel.Matrix
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyA) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyD) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyS) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyW) {
			camPos.Y += camSpeed * dt
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(colornames.Forestgreen)

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			// Random tree
			tree := pixel.NewSprite(treeSpritesheet, treeFrames[rand.Intn(len(treeFrames))])
			trees = append(trees, tree) // append tree sprite to array
			mouse := cam.Unproject(win.MousePosition())

			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 16).Moved(mouse))
		}

		if win.JustPressed(pixelgl.MouseButtonRight) {
			trees = nil
			matrices = nil
		}

		for i, tree := range trees {
			tree.Draw(win, matrices[i])
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
