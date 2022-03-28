package yuma

type Executor interface {
	GetRepository() string
	CreateExecutionOrder(string, string, []string, string) error
	DeleteExecutionOrder(string) error
	CreateEnvironmentDescription(int, interface{}) error
	RemoveEnvironmentDescription(int) error
	Execute(int, string, []string, string) bool
}