package snapshot

import (
	"encoding/json"
	"os"
)

type Snapshot struct{
	path string
}

func New(path string) *Snapshot{
	return &Snapshot{
		path: path,
	}
}

func (s *Snapshot) Save(data map[string]string) error{

	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	return encoder.Encode(data)

}

func (s *Snapshot) Load() (map[string]string, error) {

	file, err := os.Open(s.path)

	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}

		return nil, err
	}
	defer file.Close()

	data := make(map[string]string)

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}