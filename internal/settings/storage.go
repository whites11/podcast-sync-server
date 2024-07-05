package settings

type Storage interface {
	Get(string) (string, error)
	Set(string, string) error
}
