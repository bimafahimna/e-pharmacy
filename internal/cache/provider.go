package cache

type Provider interface {
	Ping() error
	SetFromStruct(key string, obj interface{}, expiration int32) error
	Set(key string, value []byte, expiration int32) error
	Get(key string) ([]byte, error)
	FlushAll() error
	DeleteAll() error
}
