package main

import (
	"fmt"

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

type Scene struct {
	blocks [800]*sdl.Rect
}

var (
	screen         Screen
	gTexture       *sdl.Texture
	selectedRect   *sdl.Rect
	mousePos       sdl.Point
	blocks         map[string]Point2D
	TEXTURE_HEIGHT = 200
	TEXTURE_WIDTH  = 200
	scene          = new(Scene)
)

func main() {
	initializeBlocks()

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
	defer window.Destroy()
	defer screen.window.Destroy()

	// Create renderer
	renderer, err := sdl.CreateRenderer(screen.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	screen.renderer = renderer
	screen.renderer.SetDrawColor(0, 0, 0, 0)

	// how do I just assign the return of CreateRenderer to a global variable
	// without making another variable?
	defer screen.renderer.Destroy()
	defer renderer.Destroy()

	// Create texture
	texture, err := img.LoadTexture(screen.renderer, "./texture_blocks.jpg")
	if err != nil {
		panic("error with texture")
	}

	gTexture = texture
	defer texture.Destroy()
	defer gTexture.Destroy()

	// leftMouseButtonDown := false

	screen.renderer.Clear()
	scene.blocks[0] = drawBlock("dirt", 0, 0)
	scene.blocks[1] = drawBlock("ice", 50, 0)
	screen.renderer.Present()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseMotionEvent:
				mousePos.X = event.X
				mousePos.Y = event.Y
				fmt.Println(mousePos)
				break
			}
		}

		// screen.renderer.Clear()
		// drawBlock("dirt", 0, 0)
		// drawBlock("ice", 50, 0)
		// screen.renderer.Present()

		sdl.Delay(33)
	}
}

func convertEventToMousePosition(event *sdl.MouseMotionEvent) sdl.Point {
	return sdl.Point{
		X: event.X,
		Y: event.Y,
	}
}

func drawBlock(blockName string, x int32, y int32) *sdl.Rect {
	pixel := sdl.Rect{X: x, Y: y, W: 50, H: 50}
	rect := sdl.Rect{
		X: int32(blocks[blockName].x),
		Y: int32(blocks[blockName].y),
		W: int32(TEXTURE_WIDTH),
		H: int32(TEXTURE_HEIGHT),
	}

	block := getAreaTexture(rect, gTexture)
	if block == nil {
		panic("block is nil")
	}

	screen.renderer.Copy(block, nil, &pixel)

	return &pixel
}

func initializeBlocks() {
	// TODO: Initialize the rest of the blocks
	blocks = make(map[string]Point2D)
	blocks["dirt"] = Point2D{x: 0, y: 0}
	blocks["ice"] = Point2D{x: 200, y: 0}
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
