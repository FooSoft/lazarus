package platform

type Scene interface {
	Name() string
}

type SceneCreator interface {
	Create(window *Window) error
}

type SceneAdvancer interface {
	Advance(window *Window) error
}

type SceneDestroyer interface {
	Destroy(window *Window) error
}

func sceneName(scene Scene) string {
	if scene == nil {
		return "<nil>"
	} else {
		return scene.Name()
	}
}
