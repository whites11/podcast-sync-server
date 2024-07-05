package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type JSON struct {
	filename string

	mutex sync.Mutex
}

func NewJSONStorage(filename string) (*JSON, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename can't be empty")
	}
	return &JSON{
		filename: filename,
	}, nil
}

func (j *JSON) Get(name string) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	data, err := j.readFile()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return data[name], nil
}

func (j *JSON) Set(name string, value string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	data, err := j.readFile()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	data[name] = value

	err = j.writeFile(data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (j *JSON) ensureFile() error {
	var _, err = os.Stat(j.filename)
	if os.IsNotExist(err) {
		file, err := os.Create(j.filename)
		if err != nil {
			fmt.Println(err)
		}

		_, err = fmt.Fprintf(file, "{}")
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		defer file.Close()
	} else if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (j *JSON) readFile() (map[string]string, error) {
	err := j.ensureFile()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	file, err := os.Open(j.filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	d := make(map[string]string)

	err = json.Unmarshal(data, &d)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return d, nil
}

func (j *JSON) writeFile(data map[string]string) error {
	err := j.ensureFile()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	jsonFile, err := os.OpenFile(j.filename, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer jsonFile.Close()

	bytes, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = fmt.Fprintf(jsonFile, "%s", bytes)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
