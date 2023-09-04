package yuma

type Executor interface {
	GetRepository() string
	CreateExecutionOrder(int, string, string, []string, string, string) error
	DeleteExecutionOrder(string) error
	CreateEnvironmentDescription(int, []string) error
	RemoveEnvironmentDescription(int) error
	Execute(int, string, []string, string, string) bool
	DetermineInputs(string) (error, []string)
	DetermineValues(string) (error, []string)
	DetermineOutputs(string) (error, []string)
}