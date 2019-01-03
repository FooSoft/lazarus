package imgui

import (
	"github.com/FooSoft/lazarus/graphics"
	"github.com/FooSoft/lazarus/math"
)

func (*Context) Button(label string) bool {
	// return imgui.Button(label)
	return false
}

func (c *Context) Image(texture graphics.Texture) {
	// c.ImageSized(texture, texture.Size())
}

func (*Context) ImageSized(texture graphics.Texture, size math.Vec2i) {
	// imgui.Image(imgui.TextureID(texture.Id()), imgui.Vec2{X: float32(size.X), Y: float32(size.Y)})
}

func (*Context) SliderInt(label string, value *int, min, max int) bool {
	// temp := int32(*value)
	// result := imgui.SliderInt(label, &temp, int32(min), int32(max))
	// *value = int(temp)
	// return result
	return false
}

func (*Context) Text(label string) {
	// imgui.Text(label)
}
