package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/whiterale/go-prac/internal"
	"github.com/whiterale/go-prac/internal/agent/buffer"
)

type Storage struct {
	*buffer.Buffer
}

func Init() *Storage {
	buf := buffer.Init()
	return &Storage{buf}
}

func InitFromFile(fname string) (*Storage, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	raw, _ := ioutil.ReadAll(f)
	metrics := make(map[string]map[string]*internal.Metric)

	err = json.Unmarshal(raw, &metrics)
	if err != nil {
		return nil, err
	}
	buf := buffer.InitWithData(metrics)
	return &Storage{buf}, nil
}

func (s *Storage) DumpToFile(fname string) error {

	jsonData, err := json.Marshal(s.Buffer.GetRawMetrics())
	if err != nil {
		return err
	}
	ioutil.WriteFile(fname, jsonData, 0644)
	return nil
}
