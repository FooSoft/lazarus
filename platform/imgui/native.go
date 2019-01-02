package imgui

// #define CIMGUI_DEFINE_ENUMS_AND_STRUCTS
// #include "cimgui/cimgui.h"
// #cgo linux LDFLAGS: -L./cimgui -l:cimgui.a -lstdc++ -lm
import "C"
import (
	"unsafe"

	"github.com/FooSoft/lazarus/math"
	"github.com/veandco/go-sdl2/sdl"
)

type imGuiContext = C.ImGuiContext
type imDrawData = C.ImDrawData
type imGuiIO = C.ImGuiIO
type imTextureId = C.ImTextureID

func CreateContext() *imGuiContext {
	c := C.igCreateContext(nil)

	keys := map[int]C.int{
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

	io := IO()
	for imguiKey, nativeKey := range keys {
		io.KeyMap[imguiKey] = nativeKey
	}

	return c
}

func (c *imGuiContext) Destroy() {
	C.igDestroyContext(c)
}

func FontImage() (unsafe.Pointer, int, int) {
	io := IO()

	var pixels *C.uchar
	var width, height C.int
	C.ImFontAtlas_GetTexDataAsRGBA32(io.Fonts, &pixels, &width, &height, nil)

	return unsafe.Pointer(pixels), int(width), int(height)
}

func NewFrame() {
	C.igNewFrame()
}

func Render() *imDrawData {
	C.igRender()
	return C.igGetDrawData()
}

func IO() *imGuiIO {
	return C.igGetIO()
}

func SetDeltaTime(time float32) {
	io := IO()
	io.DeltaTime = C.float(time)
}

func SetDisplaySize(size math.Vec2i) {
	io := IO()
	io.DisplaySize.x = C.float(size.X)
	io.DisplaySize.y = C.float(size.Y)
}

func SetMousePosition(position math.Vec2i) {
	io := IO()
	io.MousePos.x = C.float(position.X)
	io.MousePos.y = C.float(position.Y)
}

func SetMouseButtonDown(index int, down bool) {
	io := IO()
	io.MouseDown[index] = C.bool(down)
}

func SetKeyState(ctrl, shift, alt bool) {
	io := IO()
	io.KeyCtrl = C.bool(ctrl)
	io.KeyShift = C.bool(shift)
	io.KeyAlt = C.bool(alt)
}

func SetKeyDown(key int, down bool) {
	io := IO()
	io.KeysDown[key] = C.bool(down)
}

func SetFontTexture(textureId uintptr) {
	io := IO()
	io.Fonts.TexID = imTextureId(textureId)
}
