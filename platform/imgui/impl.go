package imgui

import (
	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func (c *Context) BeginFrame() {
	SetDisplaySize(c.displaySize)

	frequency := sdl.GetPerformanceFrequency()
	currentTime := sdl.GetPerformanceCounter()
	if c.lastTime > 0 {
		SetDeltaTime(float32(currentTime-c.lastTime) / float32(frequency))
	} else {
		SetDeltaTime(1.0 / 60.0)
	}
	c.lastTime = currentTime

	x, y, state := sdl.GetMouseState()
	SetMousePosition(math.Vec2i{X: int(x), Y: int(y)})
	for i, button := range []uint32{sdl.BUTTON_LEFT, sdl.BUTTON_RIGHT, sdl.BUTTON_MIDDLE} {
		SetMouseButtonDown(i, c.buttonsDown[i] || (state&sdl.Button(button)) != 0)
		c.buttonsDown[i] = false
	}

	NewFrame()
}

func (c *Context) ProcessEvent(event sdl.Event) (bool, error) {
	switch event.GetType() {
	case sdl.MOUSEWHEEL:
		wheelEvent := event.(*sdl.MouseWheelEvent)
		var delta math.Vec2i
		if wheelEvent.X > 0 {
			delta.X++
		} else if wheelEvent.X < 0 {
			delta.X--
		}
		if wheelEvent.Y > 0 {
			delta.Y++
		} else if wheelEvent.Y < 0 {
			delta.Y--
		}
		SetMouseDelta(delta)
		return true, nil
	case sdl.MOUSEBUTTONDOWN:
		buttonEvent := event.(*sdl.MouseButtonEvent)
		switch buttonEvent.Button {
		case sdl.BUTTON_LEFT:
			c.buttonsDown[0] = true
			break
		case sdl.BUTTON_RIGHT:
			c.buttonsDown[1] = true
			break
		case sdl.BUTTON_MIDDLE:
			c.buttonsDown[2] = true
			break
		}
		return true, nil
	case sdl.TEXTINPUT:
		inputEvent := event.(*sdl.TextInputEvent)
		AddInputCharacters(string(inputEvent.Text[:]))
		return true, nil
	case sdl.KEYDOWN:
		keyEvent := event.(*sdl.KeyboardEvent)
		SetKeyDown(int(keyEvent.Keysym.Scancode), true)
		modState := sdl.GetModState()
		SetAltDown(modState&sdl.KMOD_ALT != 0)
		SetCtrlDown(modState&sdl.KMOD_CTRL != 0)
		SetShiftDown(modState&sdl.KMOD_SHIFT != 0)
	case sdl.KEYUP:
		keyEvent := event.(*sdl.KeyboardEvent)
		SetKeyDown(int(keyEvent.Keysym.Scancode), false)
		modState := sdl.GetModState()
		SetAltDown(modState&sdl.KMOD_ALT != 0)
		SetCtrlDown(modState&sdl.KMOD_CTRL != 0)
		SetShiftDown(modState&sdl.KMOD_SHIFT != 0)
		return true, nil
	}

	return false, nil
}

// OpenGL2 Render function.
// Note that this implementation is little overcomplicated because we are saving/setting up/restoring every OpenGL singleton explicitly, in order to be able to run within any OpenGL engine that doesn't do so.
func (c *Context) EndFrame() error {
	// We are using the OpenGL fixed pipeline to make the example code simpler to read!
	// Setup render singleton: alpha-blending enabled, no face culling, no depth testing, scissor enabled, vertex/texcoord/color pointers, polygon fill.
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	gl.PushAttrib(gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT | gl.TRANSFORM_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.COLOR_MATERIAL)
	gl.Enable(gl.SCISSOR_TEST)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	gl.Enable(gl.TEXTURE_2D)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// You may want this if using this code in an OpenGL 3+ context where shaders may be bound
	// gl.UseProgram(0)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right). DisplayMin is typically (0,0) for single viewport apps.
	gl.Viewport(0, 0, int32(c.bufferSize.X), int32(c.bufferSize.Y))
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(c.displaySize.X), float64(c.displaySize.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	drawData := Render()
	drawData.ScaleClipRects(math.Vec2f{
		X: float32(c.bufferSize.X) / float32(c.displaySize.X),
		Y: float32(c.bufferSize.Y) / float32(c.displaySize.Y),
	})
	drawData.Draw(c.bufferSize)

	// Restore modified state
	gl.DisableClientState(gl.COLOR_ARRAY)
	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.PopAttrib()
	gl.PolygonMode(gl.FRONT, uint32(lastPolygonMode[0]))
	gl.PolygonMode(gl.BACK, uint32(lastPolygonMode[1]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])

	return nil
}
