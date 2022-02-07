package yuma

type Executor interface {
	GetRepository() string
	SetPlaybook(string)
	CreateExecutionOrder(string, []string, string) error
	DeleteExecutionOrder() error
	CreateEnvironmentDescription(description interface{}) error
	RemoveEnvironmentDescription() error
	Execute(string, []string, string) bool
}