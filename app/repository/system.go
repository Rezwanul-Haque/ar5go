package repository

type ISystem interface {
	DBCheck() (bool, error)
	CacheCheck() bool
}
