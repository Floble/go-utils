package yuma

type Environment interface {
	GetYuma() *Yuma
	GetExecutor() Executor
	GetOmega() int
	GetSigma() int
	GetInstance(int) Instance
	Initialize() error
	CleanUp() error
	CleanResults() error
	CreateInstance(int, int) error
	DeleteInstance(int) error
	TakeAction(int, int, int, []string, bool) (error, bool, float64, int)
}