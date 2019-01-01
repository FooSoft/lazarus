package graphics

import "github.com/FooSoft/lazarus/math"

type TextureId uintptr

type Texture interface {
	Id() TextureId
	Size() math.Vec2i
	Destroy() error
}
