package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Point2D struct {
	x int
	y int
}

type Screen struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

var (
	screen         Screen
	gTexture       *sdl.Texture
	blocks         map[string]Point2D
	TEXTURE_HEIGHT = 200
	TEXTURE_WIDTH  = 200
)

func main() {
	blocks = make(map[string]Point2D)
	blocks["dirt"] = Point2D{x: 0, y: 0}
	blocks["ice"] = Point2D{x: 200, y: 0}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Create window
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	screen.window = window
	defer screen.window.Destroy()

	// Create renderer
	renderer, err := sdl.CreateRenderer(screen.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	screen.renderer = renderer
	screen.renderer.SetDrawColor(0, 0, 0, 0)
	defer screen.renderer.Destroy()

	// Create texture
	texture, err := img.LoadTexture(screen.renderer, "./texture_blocks.jpg")
	if err != nil {
		panic("error with texture")
	}

	gTexture = texture
	defer gTexture.Destroy()

	var dirtRect sdl.Rect
	dirtRect.X = int32(blocks["dirt"].x)
	dirtRect.Y = int32(blocks["dirt"].y)
	dirtRect.W = int32(TEXTURE_WIDTH)
	dirtRect.H = int32(TEXTURE_HEIGHT)

	block := getAreaTexture(dirtRect, gTexture)
	if block == nil {
		panic("block is nil")
	}

	var iceRect sdl.Rect
	iceRect.X = int32(blocks["ice"].x)
	iceRect.Y = int32(blocks["ice"].y)
	iceRect.W = int32(TEXTURE_WIDTH)
	iceRect.H = int32(TEXTURE_HEIGHT)

	iceBlock := getAreaTexture(iceRect, gTexture)
	if iceBlock == nil {
		panic("block is nil")
	}

	var texturePixel sdl.Rect
	texturePixel.X = 0
	texturePixel.Y = 0
	texturePixel.W = 50
	texturePixel.H = 50

	var secondPixel sdl.Rect
	secondPixel.X = 50
	secondPixel.Y = 0
	secondPixel.W = 50
	secondPixel.H = 50

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		screen.renderer.Clear()
		screen.renderer.Copy(block, nil, &texturePixel)
		screen.renderer.Copy(iceBlock, nil, &secondPixel)
		screen.renderer.Present()

		sdl.Delay(33)
	}
}

func getAreaTexture(rect sdl.Rect, source *sdl.Texture) *sdl.Texture {
	result, err := screen.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic("area texture error")
	}

	screen.renderer.SetRenderTarget(result)
	screen.renderer.Copy(source, &rect, nil)

	// Reset the target to default(the screen)
	screen.renderer.SetRenderTarget(nil)

	return result
}
