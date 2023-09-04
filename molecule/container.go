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

func (container *Container) GetPublicIP() string {
	// NOT RELEVANT
	return ""
}

func (container *Container) SetPublicIP(publicIP string) {
	// NOT RELEVANT
}

func (container *Container) GetPrivateIP() string {
	// NOT RELEVANT
	return ""
}

func (container *Container) SetPrivateIP(privateIP string) {
	// NOT RELEVANT
}

func (container *Container) GetPort() int {
	// NOT RELEVANT
	return -1
}

func (container *Container) SetPort(port int) {
	// NOT RELEVANT
}