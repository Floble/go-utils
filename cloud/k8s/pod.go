package k8s

type Pod struct {
	id string
	port int
}

func NewPod() *Pod {
	return new(Pod)
}

func (pod *Pod) GetID() string {
	return pod.id
}

func (pod *Pod) SetID(id string) {
	pod.id = id
}

func (pod *Pod) Create() error {
	// NOT RELEVANT
	return nil
}

func (pod *Pod) Delete() error {
	// NOT RELEVANT
	return nil
}

func (pod *Pod) GetPublicIP() string {
	// NOT RELEVANT
	return ""
}

func (pod *Pod) SetPublicIP(publicIP string) {
	// NOT RELEVANT
}

func (pod *Pod) GetPrivateIP() string {
	// NOT RELEVANT
	return ""
}

func (pod *Pod) SetPrivateIP(privateIP string) {
	// NOT RELEVANT
}

func (pod *Pod) GetPort() int {
	return pod.port
}

func (pod *Pod) SetPort(port int) {
	pod.port = port
}