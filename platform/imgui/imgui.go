package imgui

import (
	imgui "github.com/FooSoft/imgui-go"
	"github.com/FooSoft/lazarus/graphics"
	"github.com/FooSoft/lazarus/math"
)

func (*Context) DialogBegin(label string) bool {
	return imgui.Begin(label)
}

func (*Context) DialogEnd() {
	imgui.End()
}

func (*Context) Button(label string) bool {
	return imgui.Button(label)
}

func (*Context) Image(texture graphics.Handle, size math.Vec2i) {
	imgui.Image(imgui.TextureID(texture), imgui.Vec2{X: float32(size.X), Y: float32(size.Y)})
}

func (*Context) SliderInt(label string, value *int, min, max int) bool {
	temp := int32(*value)
	result := imgui.SliderInt(label, &temp, int32(min), int32(max))
	*value = int(temp)
	return result
}

func (*Context) Text(label string) {
	imgui.Text(label)
}
