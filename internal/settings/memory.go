package settings

type Memory struct {
	data map[string]string
}

func NewMemoryStorage() *Memory {
	return &Memory{
		data: make(map[string]string),
	}
}

func (m *Memory) Get(name string) (string, error) {
	return m.data[name], nil
}

func (m *Memory) Set(name string, value string) error {
	m.data[name] = value

	return nil
}
