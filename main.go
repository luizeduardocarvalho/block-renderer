package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Point2D struct {
	x int
	y int
}

type Block struct {
	blockType string
	block     *sdl.Rect
}

type Screen struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

type Scene struct {
	blocks [40]Block
}

var (
	screen         Screen
	gTexture       *sdl.Texture
	selectedRect   *sdl.Rect
	mousePos       sdl.Point
	blocks         map[string]Point2D
	TEXTURE_HEIGHT = 200
	TEXTURE_WIDTH  = 200
	scene          = Scene{}
)

func main() {
	for i := 0; i < len(scene.blocks); i++ {
		scene.blocks[i] = Block{
			block:     nil,
			blockType: "",
		}
	}

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

	scene.blocks[0].block = drawBlock("dirt", 0, 0)
	scene.blocks[0].blockType = "dirt"
	scene.blocks[1].block = drawBlock("ice", 50, 0)
	scene.blocks[1].blockType = "ice"

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

				// update the selected block based on the mouse position
				if selectedRect != nil {
					selectedRect.X = mousePos.X - selectedRect.W/2
					selectedRect.Y = mousePos.Y - selectedRect.H/2
				}
				break
			case *sdl.MouseButtonEvent:
				if event.Type == sdl.MOUSEBUTTONDOWN {
					dragObject()
				}
				if event.Type == sdl.MOUSEBUTTONUP {
					selectedRect = nil
				}
				break
			}
		}

		screen.renderer.Clear()
		for _, block := range scene.blocks {
			if block.block != nil {
				drawBlock(block.blockType, block.block.X, block.block.Y)
			}
		}
		screen.renderer.Present()

		sdl.Delay(33)
	}
}

func dragObject() {
	// loop the blocks to find the selected one
	for _, block := range scene.blocks {
		if block.block != nil {
			if (mousePos.X > block.block.X && mousePos.X < (block.block.X+block.block.W)) &&
				mousePos.Y > block.block.Y && mousePos.Y < (block.block.Y+block.block.H) {
				selectedRect = block.block
			}
		}
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
