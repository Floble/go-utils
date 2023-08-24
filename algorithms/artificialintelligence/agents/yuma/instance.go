package yuma

type Instance interface {
	Create() error
	Delete() error
	GetPublicIP() string
	SetPublicIP(string)
	GetPrivateIP() string
	SetPrivateIP(string)
}