package yuma

type Instance interface {
	Create() error
	Delete() error
}