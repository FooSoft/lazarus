package platform

import (
	"errors"
	"log"

	"github.com/FooSoft/lazarus/graphics"
	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/platform/imgui"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ErrWindowExists    = errors.New("only one window can exist at a time")
	ErrWindowNotExists = errors.New("no window has been created")
)

type Scene interface {
	Name() string
}

type SceneCreator interface {
	Create() error
}

type SceneAdvancer interface {
	Advance() error
}

type SceneDestroyer interface {
	Destroy() error
}

var windowState struct {
	sdlWindow    *sdl.Window
	sdlGlContext sdl.GLContext
	scene        Scene
}

func WindowCreate(title string, size math.Vec2i, scene Scene) error {
	if WindowIsCreated() {
		return ErrWindowExists
	}

	var err error
	log.Println("window create")
	if windowState.sdlWindow, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(size.X), int32(size.Y), sdl.WINDOW_OPENGL); err != nil {
		return err
	}

	log.Println("window gl context create")
	if windowState.sdlGlContext, err = windowState.sdlWindow.GLCreateContext(); err != nil {
		WindowDestroy()
		return err
	}

	log.Println("window gl context make current")
	windowState.sdlWindow.GLMakeCurrent(windowState.sdlGlContext)

	log.Println("window gl init")
	if err := gl.Init(); err != nil {
		WindowDestroy()
		return err
	}

	if err := imgui.Create(); err != nil {
		WindowDestroy()
		return err
	}

	if err := WindowSetScene(scene); err != nil {
		WindowDestroy()
		return err
	}

	return nil
}

func WindowSetScene(scene Scene) error {
	if !WindowIsCreated() {
		return ErrWindowNotExists
	}

	if windowState.scene == scene {
		return nil
	}

	sceneName := func(s Scene) string {
		if s == nil {
			return "<nil>"
		} else {
			return s.Name()
		}
	}

	if sceneDestroyer, ok := windowState.scene.(SceneDestroyer); ok {
		log.Printf("window scene notify destroy \"%s\"\n", sceneName(windowState.scene))
		if err := sceneDestroyer.Destroy(); err != nil {
			return err
		}
	}

	log.Printf("window scene transition \"%v\" => \"%v\"\n", sceneName(windowState.scene), sceneName(scene))
	windowState.scene = scene

	if sceneCreator, ok := scene.(SceneCreator); ok {
		log.Printf("window scene notify create \"%s\"\n", sceneName(windowState.scene))
		if err := sceneCreator.Create(); err != nil {
			return err
		}
	}

	return nil
}

func WindowDestroy() error {
	if !WindowIsCreated() {
		return nil
	}

	if err := WindowSetScene(nil); err != nil {
		return err
	}

	if err := imgui.Destroy(); err != nil {
		return err
	}

	log.Println("window gl context destroy")
	sdl.GLDeleteContext(windowState.sdlGlContext)
	windowState.sdlGlContext = nil

	log.Println("window destroy")
	if err := windowState.sdlWindow.Destroy(); err != nil {
		return err
	}
	windowState.sdlWindow = nil

	return nil
}

func WindowRenderTexture(texture graphics.Texture, position math.Vec2i) error {
	if !WindowIsCreated() {
		return ErrWindowNotExists
	}

	size := texture.Size()

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, uint32(texture.Id()))
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(0, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(0, float32(size.Y))
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(float32(size.X), float32(size.Y))
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(float32(size.X), 0)
	gl.End()

	return nil
}

func WindowDisplaySize() (math.Vec2i, error) {
	if !WindowIsCreated() {
		return math.Vec2i{}, ErrWindowNotExists
	}

	width, height := windowState.sdlWindow.GetSize()
	return math.Vec2i{X: int(width), Y: int(height)}, nil
}

func WindowIsCreated() bool {
	return windowState.sdlWindow != nil
}

func windowBufferSize() (math.Vec2i, error) {
	if !WindowIsCreated() {
		return math.Vec2i{}, ErrWindowNotExists
	}

	width, height := windowState.sdlWindow.GLGetDrawableSize()
	return math.Vec2i{X: int(width), Y: int(height)}, nil
}

func windowAdvance() (bool, error) {
	if !WindowIsCreated() {
		return false, ErrWindowNotExists
	}

	displaySize, err := WindowDisplaySize()
	if err != nil {
		return false, err
	}

	bufferSize, err := windowBufferSize()
	if err != nil {
		return false, err
	}

	imgui.BeginFrame(displaySize, bufferSize)

	gl.Viewport(0, 0, int32(displaySize.X), int32(displaySize.Y))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(displaySize.X), float64(displaySize.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	if sceneAdvancer, ok := windowState.scene.(SceneAdvancer); ok {
		if err := sceneAdvancer.Advance(); err != nil {
			return false, err
		}
	}

	imgui.EndFrame()
	windowState.sdlWindow.GLSwap()

	return windowState.scene != nil, nil
}

func windowProcessEvent(event sdl.Event) (bool, error) {
	if !WindowIsCreated() {
		return false, ErrWindowNotExists
	}

	return imgui.ProcessEvent(event)
}
