package platform

type Scene interface{}

type SceneCreator interface {
	Create(window *Window) error
}

type SceneAdvancer interface {
	Advance(window *Window) error
}

type SceneDestroyer interface {
	Destroy(window *Window) error
}
