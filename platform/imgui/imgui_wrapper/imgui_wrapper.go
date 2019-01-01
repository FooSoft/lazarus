package imgui

// #define CIMGUI_DEFINE_ENUMS_AND_STRUCTS
// #include "cimgui/cimgui.h"
// #cgo linux LDFLAGS: -L./cimgui -l:cimgui.a -lstdc++ -lm
import "C"
import "unsafe"

type Context = C.ImGuiContext
type DrawData = C.ImDrawData
type FontAtlas = C.ImFontAtlas
type Io = C.ImGuiIO

func CreateContext(fontAtlas *FontAtlas) *Context {
	c := C.igCreateContext(fontAtlas)
	return c
}

func (c *Context) Destroy() {
	C.igDestroyContext(c)
}

func NewFrame() {
	C.igNewFrame()
}

func Render() {
	C.igRender()
}

func GetDrawData() *DrawData {
	return C.igGetDrawData()
}

func GetIo() *Io {
	return C.igGetIO()
}

func (fa *FontAtlas) GetTexDataAsRGBA32() (pixels unsafe.Pointer, width, height int32) {
	var data *C.uint8_t
	var w, h C.int
	C.ImFontAtlas_GetTexDataAsRGBA32(fa, &data, &w, &h, nil)
	return unsafe.Pointer(data), int32(w), int32(h)
}
