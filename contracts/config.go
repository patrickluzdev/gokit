package contracts

type Config interface {
	Get(key string) string
	GetWithDefault(key, defaultValue string) string
	Set(key, value string)
	All() map[string]string
	GetInt(key string) int
	GetBool(key string) bool
}
