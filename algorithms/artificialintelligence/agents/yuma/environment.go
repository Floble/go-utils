package yuma

type Environment interface {
	GetYuma() *Yuma
	GetExecutor() Executor
	GetOmega() int
	GetSigma() int
	GetInstances(int) map[int][]Instance
	Initialize() error
	CleanUp() error
	CleanResults() error
	CreateInstance(int, int, int) error
	DeleteInstance(int, int) error
	DeleteAllInstances(int) error
	TakeAction(int, int, int, []string, bool) (error, bool, float64, int)
}