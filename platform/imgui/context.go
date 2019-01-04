package imgui

// #cgo linux CFLAGS: -I./cimgui
// #cgo linux LDFLAGS: -L./cimgui -l:cimgui.a -lstdc++ -lm
// #include "native.h"
import "C"
import (
	"log"
	"unsafe"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	pointerSize     = unsafe.Sizeof(C.uintptr_t(0))
	drawCommandSize = unsafe.Sizeof(C.ImDrawCmd{})
	vertexSize      = unsafe.Sizeof(C.ImDrawVert{})
	vertexOffsetPos = unsafe.Offsetof(C.ImDrawVert{}.pos)
	vertexOffsetUv  = unsafe.Offsetof(C.ImDrawVert{}.uv)
	vertexOffsetCol = unsafe.Offsetof(C.ImDrawVert{}.col)
)

var keyMapping = map[int]C.int{
	C.ImGuiKey_Tab:        sdl.SCANCODE_TAB,
	C.ImGuiKey_LeftArrow:  sdl.SCANCODE_LEFT,
	C.ImGuiKey_RightArrow: sdl.SCANCODE_RIGHT,
	C.ImGuiKey_UpArrow:    sdl.SCANCODE_UP,
	C.ImGuiKey_DownArrow:  sdl.SCANCODE_DOWN,
	C.ImGuiKey_PageUp:     sdl.SCANCODE_PAGEUP,
	C.ImGuiKey_PageDown:   sdl.SCANCODE_PAGEDOWN,
	C.ImGuiKey_Home:       sdl.SCANCODE_HOME,
	C.ImGuiKey_End:        sdl.SCANCODE_END,
	C.ImGuiKey_Insert:     sdl.SCANCODE_INSERT,
	C.ImGuiKey_Delete:     sdl.SCANCODE_DELETE,
	C.ImGuiKey_Backspace:  sdl.SCANCODE_BACKSPACE,
	C.ImGuiKey_Space:      sdl.SCANCODE_BACKSPACE,
	C.ImGuiKey_Enter:      sdl.SCANCODE_RETURN,
	C.ImGuiKey_Escape:     sdl.SCANCODE_ESCAPE,
	C.ImGuiKey_A:          sdl.SCANCODE_A,
	C.ImGuiKey_C:          sdl.SCANCODE_C,
	C.ImGuiKey_V:          sdl.SCANCODE_V,
	C.ImGuiKey_X:          sdl.SCANCODE_X,
	C.ImGuiKey_Y:          sdl.SCANCODE_Y,
	C.ImGuiKey_Z:          sdl.SCANCODE_Z,
}

var singleton struct {
	nativeContext *C.ImGuiContext
	nativeIo      *C.ImGuiIO
	refCount      int
}

type Context struct {
	buttonsDown [3]bool
	lastTime    uint64
	fontTexture uint32
	displaySize math.Vec2i
	bufferSize  math.Vec2i
}

func New(displaySize, bufferSize math.Vec2i) (*Context, error) {
	singleton.refCount++
	if singleton.refCount == 1 {
		log.Println("imgui global create")
		singleton.nativeContext = C.igCreateContext(nil)
		singleton.nativeIo = C.igGetIO()

		for imguiKey, nativeKey := range keyMapping {
			singleton.nativeIo.KeyMap[imguiKey] = nativeKey
		}
	}

	log.Println("imgui context create")
	c := &Context{displaySize: displaySize, bufferSize: bufferSize}

	var imageData *C.uchar
	var imageWidth, imageHeight C.int
	C.ImFontAtlas_GetTexDataAsRGBA32(singleton.nativeIo.Fonts, &imageData, &imageWidth, &imageHeight, nil)

	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &c.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, c.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(imageWidth), int32(imageHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(imageData))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))

	return c, nil
}

func (c *Context) Destroy() error {
	if c == nil || c.fontTexture == 0 {
		return nil
	}

	log.Println("imgui context destroy")
	gl.DeleteTextures(1, &c.fontTexture)
	singleton.nativeIo.Fonts.TexID = C.ImTextureID(uintptr(c.fontTexture))
	c.fontTexture = 0

	singleton.refCount--
	if singleton.refCount == 0 {
		log.Println("imgui global destroy")
		C.igDestroyContext(singleton.nativeContext)
		singleton.nativeContext = nil
		singleton.nativeIo = nil
	}

	return nil
}

func (c *Context) SetDisplaySize(displaySize math.Vec2i) {
	c.displaySize = displaySize
}

func (c *Context) SetBufferSize(bufferSize math.Vec2i) {
	c.bufferSize = bufferSize
}

func (c *Context) BeginFrame() {
	singleton.nativeIo.Fonts.TexID = C.ImTextureID(uintptr(c.fontTexture))
	singleton.nativeIo.DisplaySize.x = C.float(c.displaySize.X)
	singleton.nativeIo.DisplaySize.y = C.float(c.displaySize.Y)

	currentTime := sdl.GetPerformanceCounter()
	if c.lastTime > 0 {
		singleton.nativeIo.DeltaTime = C.float(float32(currentTime-c.lastTime) / float32(sdl.GetPerformanceFrequency()))
	} else {
		singleton.nativeIo.DeltaTime = C.float(1.0 / 60.0)
	}
	c.lastTime = currentTime

	x, y, state := sdl.GetMouseState()
	singleton.nativeIo.MousePos.x = C.float(x)
	singleton.nativeIo.MousePos.y = C.float(y)
	for i, button := range []uint32{sdl.BUTTON_LEFT, sdl.BUTTON_RIGHT, sdl.BUTTON_MIDDLE} {
		singleton.nativeIo.MouseDown[i] = C.bool(c.buttonsDown[i] || (state&sdl.Button(button)) != 0)
		c.buttonsDown[i] = false
	}

	C.igNewFrame()
}

func (c *Context) ProcessEvent(event sdl.Event) (bool, error) {
	switch event.GetType() {
	case sdl.MOUSEWHEEL:
		wheelEvent := event.(*sdl.MouseWheelEvent)
		if wheelEvent.X > 0 {
			singleton.nativeIo.MouseDelta.x++
		} else if wheelEvent.X < 0 {
			singleton.nativeIo.MouseDelta.x--
		}
		if wheelEvent.Y > 0 {
			singleton.nativeIo.MouseDelta.y++
		} else if wheelEvent.Y < 0 {
			singleton.nativeIo.MouseDelta.y--
		}
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
		C.ImGuiIO_AddInputCharactersUTF8(singleton.nativeIo, (*C.char)(unsafe.Pointer(&inputEvent.Text[0])))
		return true, nil
	case sdl.KEYDOWN:
		keyEvent := event.(*sdl.KeyboardEvent)
		singleton.nativeIo.KeysDown[keyEvent.Keysym.Scancode] = true
		modState := sdl.GetModState()
		singleton.nativeIo.KeyCtrl = C.bool(modState&sdl.KMOD_CTRL != 0)
		singleton.nativeIo.KeyAlt = C.bool(modState&sdl.KMOD_ALT != 0)
		singleton.nativeIo.KeyShift = C.bool(modState&sdl.KMOD_SHIFT != 0)
	case sdl.KEYUP:
		keyEvent := event.(*sdl.KeyboardEvent)
		singleton.nativeIo.KeysDown[keyEvent.Keysym.Scancode] = false
		modState := sdl.GetModState()
		singleton.nativeIo.KeyCtrl = C.bool(modState&sdl.KMOD_CTRL != 0)
		singleton.nativeIo.KeyAlt = C.bool(modState&sdl.KMOD_ALT != 0)
		singleton.nativeIo.KeyShift = C.bool(modState&sdl.KMOD_SHIFT != 0)
		return true, nil
	}

	return false, nil
}

func (c *Context) EndFrame() error {
	C.igRender()

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

	gl.Viewport(0, 0, int32(c.bufferSize.X), int32(c.bufferSize.Y))
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(c.displaySize.X), float64(c.displaySize.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	drawData := C.igGetDrawData()
	C.ImDrawData_ScaleClipRects(
		drawData,
		C.ImVec2{
			x: C.float(float32(c.bufferSize.X) / float32(c.displaySize.X)),
			y: C.float(float32(c.bufferSize.Y) / float32(c.displaySize.Y)),
		},
	)

	for i := C.int(0); i < drawData.CmdListsCount; i++ {
		commandList := *(**C.ImDrawList)(unsafe.Pointer(uintptr(unsafe.Pointer(drawData.CmdLists)) + pointerSize*uintptr(i)))
		vertexBuffer := unsafe.Pointer(commandList.VtxBuffer.Data)
		indexBuffer := unsafe.Pointer(commandList.IdxBuffer.Data)
		indexBufferOffset := uintptr(indexBuffer)

		gl.VertexPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetPos)))
		gl.TexCoordPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetUv)))
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetCol)))

		for j := C.int(0); j < commandList.CmdBuffer.Size; j++ {
			command := (*C.ImDrawCmd)(unsafe.Pointer(uintptr(unsafe.Pointer(commandList.CmdBuffer.Data)) + drawCommandSize*uintptr(j)))
			gl.Scissor(
				int32(command.ClipRect.x),
				int32(c.bufferSize.Y)-int32(command.ClipRect.w),
				int32(command.ClipRect.z-command.ClipRect.x),
				int32(command.ClipRect.w-command.ClipRect.y),
			)
			gl.BindTexture(gl.TEXTURE_2D, uint32(uintptr(command.TextureId)))
			gl.DrawElements(gl.TRIANGLES, int32(command.ElemCount), gl.UNSIGNED_SHORT, unsafe.Pointer(indexBufferOffset))
			indexBufferOffset += uintptr(command.ElemCount * 2)
		}
	}

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
