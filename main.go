package main

import (
	"image"
	_ "image/png"
	"math/rand"
	"os"

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
		trees    []*pixel.Sprite
		matrices []pixel.Matrix
	)

	for !win.Closed() {
		win.Clear(colornames.Forestgreen)
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			// Random tree
			tree := pixel.NewSprite(treeSpritesheet, treeFrames[rand.Intn(len(treeFrames))])
			trees = append(trees, tree) // append tree sprite to array
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.MousePosition()))
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
