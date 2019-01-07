package imgui

// #cgo linux LDFLAGS: -L./cimgui -l:cimgui.a -lstdc++ -lm
// #include "native.h"
import "C"
import (
	"unsafe"

	"github.com/FooSoft/lazarus/graphics"
	"github.com/FooSoft/lazarus/math"
)

func DialogBegin(label string) bool {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))
	return bool(C.igBegin(labelC, nil, 0))
}

func DialogEnd() {
	C.igEnd()
}

func Button(label string) bool {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))
	return bool(C.igButton(labelC, C.ImVec2{}))
}

func Image(texture graphics.Texture) {
	ImageSized(texture, texture.Size())
}

func ImageSized(texture graphics.Texture, size math.Vec2i) {
	C.igImage(
		C.nativeHandleCast(C.uintptr_t(texture.Id())),
		C.ImVec2{x: C.float(size.X), y: C.float(size.Y)},
		C.ImVec2{0, 0},
		C.ImVec2{1, 1},
		C.ImVec4{1, 1, 1, 1},
		C.ImVec4{0, 0, 0, 0},
	)
}

func SliderInt(label string, value *int, min, max int) bool {
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))
	valueC := C.int(*value)
	result := bool(C.igSliderInt(labelC, &valueC, (C.int)(min), (C.int)(max), nil))
	*value = int(valueC)
	return result
}

func Text(label string) {
	labelStartC := C.CString(label)
	labelEndC := (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(labelStartC)) + uintptr(len(label))))
	defer C.free(unsafe.Pointer(labelStartC))
	C.igTextUnformatted(labelStartC, labelEndC)
}
