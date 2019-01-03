package imgui

// #cgo linux CFLAGS: -I./cimgui
// #cgo linux LDFLAGS: -L./cimgui -l:cimgui.a -lstdc++ -lm
// #include "native.h"
import "C"
import (
	"unsafe"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type imGuiContext = C.ImGuiContext
type imDrawData = C.ImDrawData
type imGuiIO = C.ImGuiIO
type imTextureId = C.ImTextureID
type imVec2 = C.ImVec2

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

func (d *imDrawData) ScaleClipRects(scale math.Vec2f) {
	C.ImDrawData_ScaleClipRects(d, imVec2{x: C.float(scale.X), y: C.float(scale.Y)})
}

func (d *imDrawData) Draw(bufferSize math.Vec2i) {
	vert := C.ImDrawVert{}
	vertexSize := unsafe.Sizeof(vert)
	vertexOffsetPos := unsafe.Offsetof(vert.pos)
	vertexOffsetUv := unsafe.Offsetof(vert.uv)
	vertexOffsetCol := unsafe.Offsetof(vert.col)

	for i := C.int(0); i < d.CmdListsCount; i++ {
		commandList := C.getDrawList(d.CmdLists, i)
		vertexBuffer := unsafe.Pointer(commandList.VtxBuffer.Data)
		indexBuffer := unsafe.Pointer(commandList.IdxBuffer.Data)
		indexBufferOffset := uintptr(indexBuffer)

		gl.VertexPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetPos)))
		gl.TexCoordPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetUv)))
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetCol)))

		for j := C.int(0); j < commandList.CmdBuffer.Size; j++ {
			command := C.getDrawCmd(commandList.CmdBuffer.Data, j)
			gl.Scissor(
				int32(command.ClipRect.x),
				int32(bufferSize.Y)-int32(command.ClipRect.w),
				int32(command.ClipRect.z-command.ClipRect.x),
				int32(command.ClipRect.w-command.ClipRect.y),
			)

			gl.BindTexture(gl.TEXTURE_2D, uint32(uintptr(command.TextureId)))
			gl.DrawElements(gl.TRIANGLES, int32(command.ElemCount), gl.UNSIGNED_SHORT, unsafe.Pointer(indexBufferOffset))

			indexBufferOffset += uintptr(command.ElemCount * 2)
		}
	}

	// for _, commandList := range drawData.CommandLists() {
	// 	vertexBuffer, _ := commandList.VertexBuffer()
	// 	indexBuffer, _ := commandList.IndexBuffer()
	// 	indexBufferOffset := uintptr(indexBuffer)

	// 	gl.VertexPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetPos)))
	// 	gl.TexCoordPointer(2, gl.FLOAT, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetUv)))
	// 	gl.ColorPointer(4, gl.UNSIGNED_BYTE, int32(vertexSize), unsafe.Pointer(uintptr(vertexBuffer)+uintptr(vertexOffsetCol)))

	// 	for _, command := range commandList.Commands() {
	// 		if command.HasUserCallback() {
	// 			command.CallUserCallback(commandList)
	// 		} else {
	// 			clipRect := command.ClipRect()
	// 			gl.Scissor(int32(clipRect.X), int32(c.bufferSize.Y)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
	// 			gl.BindTexture(gl.TEXTURE_2D, uint32(command.TextureID()))
	// 			gl.DrawElements(gl.TRIANGLES, int32(command.ElementCount()), uint32(drawType), unsafe.Pointer(indexBufferOffset))
	// 		}

	// 		indexBufferOffset += uintptr(command.ElementCount() * indexSize)
	// 	}
	// }
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

func SetAltDown(down bool) {
	io := IO()
	io.KeyAlt = C.bool(down)
}

func SetShiftDown(down bool) {
	io := IO()
	io.KeyShift = C.bool(down)
}

func SetCtrlDown(down bool) {
	io := IO()
	io.KeyCtrl = C.bool(down)
}

func SetKeyDown(key int, down bool) {
	io := IO()
	io.KeysDown[key] = C.bool(down)
}

func SetMouseDelta(delta math.Vec2i) {
	io := IO()
	io.MouseDelta.x = C.float(delta.X)
	io.MouseDelta.y = C.float(delta.Y)
}

func SetFontTexture(textureId uintptr) {
	io := IO()
	io.Fonts.TexID = imTextureId(textureId)
}

func AddInputCharacters(characters string) {
	io := IO()
	cs := C.CString(characters)
	defer C.free(unsafe.Pointer(cs))
	C.ImGuiIO_AddInputCharactersUTF8(io, cs)
}
