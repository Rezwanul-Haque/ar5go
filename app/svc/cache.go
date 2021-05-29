package svc

type ICache interface {
	Set(key string, value interface{}, ttl int) error
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetStruct(key string, outputStruct interface{}) error
	Del(keys ...string) error
	DelPattern(pattern string) error
}
