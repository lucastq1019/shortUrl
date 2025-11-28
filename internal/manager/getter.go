package manager

type ClientGetter interface {
	GetClient(serviceName string) (interface{}, error)
}
