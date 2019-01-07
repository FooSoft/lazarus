package imgui

// #include "native.h"
import "C"
import (
	"errors"
	"log"
	"unsafe"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ErrNotInit = errors.New("imgui context was not created")
)

const (
	pointerSize     = unsafe.Sizeof(C.uintptr_t(0))
	drawCommandSize = unsafe.Sizeof(C.ImDrawCmd{})
	indexSize       = unsafe.Sizeof(C.ImDrawIdx(0))
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

var imguiState struct {
	imguiContext *C.ImGuiContext
	imguiIo      *C.ImGuiIO

	fontTexture uint32
	displaySize math.Vec2i
	bufferSize  math.Vec2i
	buttonsDown [3]bool
	lastTime    uint64
}

func Create() error {
	if IsCreated() {
		return nil
	}

	log.Println("imgui create")
	imguiState.imguiContext = C.igCreateContext(nil)
	imguiState.imguiIo = C.igGetIO()

	for imguiKey, nativeKey := range keyMapping {
		imguiState.imguiIo.KeyMap[imguiKey] = nativeKey
	}

	var imageData *C.uchar
	var imageWidth, imageHeight C.int
	C.ImFontAtlas_GetTexDataAsRGBA32(imguiState.imguiIo.Fonts, &imageData, &imageWidth, &imageHeight, nil)

	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &imguiState.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, imguiState.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(imageWidth), int32(imageHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(imageData))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))

	return nil
}

func IsCreated() bool {
	return imguiState.imguiContext != nil
}

func Destroy() error {
	if !IsCreated() {
		return nil
	}

	gl.DeleteTextures(1, &imguiState.fontTexture)
	imguiState.imguiIo.Fonts.TexID = C.nativeHandleCast(C.uintptr_t(imguiState.fontTexture))
	imguiState.fontTexture = 0

	log.Println("imgui destroy")
	C.igDestroyContext(imguiState.imguiContext)
	imguiState.imguiContext = nil
	imguiState.imguiIo = nil

	return nil
}

func BeginFrame(displaySize, bufferSize math.Vec2i) error {
	if !IsCreated() {
		return ErrNotInit
	}

	imguiState.displaySize = displaySize
	imguiState.bufferSize = bufferSize

	imguiState.imguiIo.Fonts.TexID = C.nativeHandleCast(C.uintptr_t(imguiState.fontTexture))
	imguiState.imguiIo.DisplaySize.x = C.float(displaySize.X)
	imguiState.imguiIo.DisplaySize.y = C.float(displaySize.Y)

	currentTime := sdl.GetPerformanceCounter()
	if imguiState.lastTime > 0 {
		imguiState.imguiIo.DeltaTime = C.float(float32(currentTime-imguiState.lastTime) / float32(sdl.GetPerformanceFrequency()))
	} else {
		imguiState.imguiIo.DeltaTime = C.float(1.0 / 60.0)
	}
	imguiState.lastTime = currentTime

	x, y, state := sdl.GetMouseState()
	imguiState.imguiIo.MousePos.x = C.float(x)
	imguiState.imguiIo.MousePos.y = C.float(y)
	for i, button := range []uint32{sdl.BUTTON_LEFT, sdl.BUTTON_RIGHT, sdl.BUTTON_MIDDLE} {
		imguiState.imguiIo.MouseDown[i] = C.bool(imguiState.buttonsDown[i] || (state&sdl.Button(button)) != 0)
		imguiState.buttonsDown[i] = false
	}

	C.igNewFrame()
	return nil
}

func ProcessEvent(event sdl.Event) (bool, error) {
	if !IsCreated() {
		return false, ErrNotInit
	}

	switch event.GetType() {
	case sdl.MOUSEWHEEL:
		wheelEvent := event.(*sdl.MouseWheelEvent)
		imguiState.imguiIo.MouseWheelH += C.float(wheelEvent.X)
		imguiState.imguiIo.MouseWheel += C.float(wheelEvent.Y)
		return true, nil
	case sdl.MOUSEBUTTONDOWN:
		buttonEvent := event.(*sdl.MouseButtonEvent)
		switch buttonEvent.Button {
		case sdl.BUTTON_LEFT:
			imguiState.buttonsDown[0] = true
			break
		case sdl.BUTTON_RIGHT:
			imguiState.buttonsDown[1] = true
			break
		case sdl.BUTTON_MIDDLE:
			imguiState.buttonsDown[2] = true
			break
		}
		return true, nil
	case sdl.TEXTINPUT:
		inputEvent := event.(*sdl.TextInputEvent)
		C.ImGuiIO_AddInputCharactersUTF8(imguiState.imguiIo, (*C.char)(unsafe.Pointer(&inputEvent.Text[0])))
		return true, nil
	case sdl.KEYDOWN:
		keyEvent := event.(*sdl.KeyboardEvent)
		imguiState.imguiIo.KeysDown[keyEvent.Keysym.Scancode] = true
		modState := sdl.GetModState()
		imguiState.imguiIo.KeyCtrl = C.bool(modState&sdl.KMOD_CTRL != 0)
		imguiState.imguiIo.KeyAlt = C.bool(modState&sdl.KMOD_ALT != 0)
		imguiState.imguiIo.KeyShift = C.bool(modState&sdl.KMOD_SHIFT != 0)
	case sdl.KEYUP:
		keyEvent := event.(*sdl.KeyboardEvent)
		imguiState.imguiIo.KeysDown[keyEvent.Keysym.Scancode] = false
		modState := sdl.GetModState()
		imguiState.imguiIo.KeyCtrl = C.bool(modState&sdl.KMOD_CTRL != 0)
		imguiState.imguiIo.KeyAlt = C.bool(modState&sdl.KMOD_ALT != 0)
		imguiState.imguiIo.KeyShift = C.bool(modState&sdl.KMOD_SHIFT != 0)
		return true, nil
	}

	return false, nil
}

func EndFrame() error {
	if !IsCreated() {
		return ErrNotInit
	}

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

	gl.Viewport(0, 0, int32(imguiState.bufferSize.X), int32(imguiState.bufferSize.Y))
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(imguiState.displaySize.X), float64(imguiState.displaySize.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	drawData := C.igGetDrawData()
	C.ImDrawData_ScaleClipRects(
		drawData,
		C.ImVec2{
			x: C.float(float32(imguiState.bufferSize.X) / float32(imguiState.displaySize.X)),
			y: C.float(float32(imguiState.bufferSize.Y) / float32(imguiState.displaySize.Y)),
		},
	)

	for i := C.int(0); i < drawData.CmdListsCount; i++ {
		var (
			commandList  = *(**C.ImDrawList)(unsafe.Pointer(uintptr(unsafe.Pointer(drawData.CmdLists)) + pointerSize*uintptr(i)))
			vertexBuffer = unsafe.Pointer(commandList.VtxBuffer.Data)
			indexBuffer  = unsafe.Pointer(commandList.IdxBuffer.Data)
			elementCount = C.unsigned(0)
		)

		gl.VertexPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetPos)))
		gl.TexCoordPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetUv)))
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetCol)))

		for j := C.int(0); j < commandList.CmdBuffer.Size; j++ {
			command := (*C.ImDrawCmd)(unsafe.Pointer(uintptr(unsafe.Pointer(commandList.CmdBuffer.Data)) + drawCommandSize*uintptr(j)))
			gl.Scissor(
				int32(command.ClipRect.x),
				int32(imguiState.bufferSize.Y)-int32(command.ClipRect.w),
				int32(command.ClipRect.z-command.ClipRect.x),
				int32(command.ClipRect.w-command.ClipRect.y),
			)
			gl.BindTexture(gl.TEXTURE_2D, uint32(uintptr(command.TextureId)))
			gl.DrawElements(
				gl.TRIANGLES,
				int32(command.ElemCount),
				gl.UNSIGNED_SHORT,
				unsafe.Pointer(uintptr(unsafe.Pointer(indexBuffer))+uintptr(elementCount)*uintptr(indexSize)),
			)

			elementCount += command.ElemCount
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
