package molecule

type Container struct {
	id string
}

func NewContainer(id string) *Container {
	container := new(Container)
	container.id = id

	return container
}

func (container *Container) GetID() string {
	return container.id
}

func (container *Container) SetID(id string) {
	container.id = id
}

func (container *Container) Create() error {
	// NOT RELEVANT
	return nil
}

func (container *Container) Delete() error {
	// NOT RELEVANT
	return nil
}