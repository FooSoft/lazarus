package platform

type Scene interface {
	Name() string
}

type SceneCreator interface {
	Create() error
}

type SceneAdvancer interface {
	Advance() error
}

type SceneDestroyer interface {
	Destroy() error
}
