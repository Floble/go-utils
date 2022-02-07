package yuma

type Environment interface {
	GetYuma() *Yuma
	GetExecutor() Executor
	GetOmega() int
	GetSigma() int
	GetInstance() Instance
	CreateInstance(int) error
	DeleteInstance() error
	TakeAction(int, int, []string, bool) (error, bool, float64, int)
}