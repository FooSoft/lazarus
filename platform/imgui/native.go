package imgui

// #cgo linux CFLAGS: -I./cimgui
// #cgo linux LDFLAGS: -L./cimgui -l:cimgui.a -lstdc++ -lm
// #include "native.h"
import "C"
import (
	"github.com/FooSoft/lazarus/math"
)

type imGuiContext = C.ImGuiContext
type imDrawData = C.ImDrawData
type imGuiIO = C.ImGuiIO

func (d *imDrawData) Draw(bufferSize math.Vec2i) {
}
