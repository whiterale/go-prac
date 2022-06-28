package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/whiterale/go-prac/internal"
	"github.com/whiterale/go-prac/internal/agent/buffer"
)

type Storage struct {
	*buffer.Buffer
	IsSync bool
}

func Init() *Storage {
	buf := buffer.Init()
	return &Storage{buf, false}
}

func (s *Storage) Update(mtype string, mid string, val interface{}) error {
	err := s.Buffer.Update(mtype, mid, val)
	if err != nil {
		return err
	}
	if s.IsSync {
		s.DumpToFile("/tmp/devops-metrics-db.json")
	}
	return nil
}

func InitFromFile(fname string) (*Storage, error) {
	f, err := os.Open(fname)
	if err != nil {
		// File does not exist, create empty
		return &Storage{buffer.Init(), false}, err
	}
	defer f.Close()

	raw, _ := ioutil.ReadAll(f)
	metrics := make(map[string]map[string]*internal.Metric)

	err = json.Unmarshal(raw, &metrics)
	if err != nil {
		return nil, err
	}
	buf := buffer.InitWithData(metrics)
	return &Storage{buf, false}, nil
}

func (s *Storage) DumpToFile(fname string) error {
	jsonData, err := json.Marshal(s.Buffer.GetRawMetrics())
	if err != nil {
		return err
	}
	ioutil.WriteFile(fname, jsonData, 0644)
	return nil
}

func (s *Storage) StartSync(fname string, stop chan struct{}) {
	storeTick := time.NewTicker(1 * time.Second)
	defer func() {
		storeTick.Stop()
	}()
	for {
		select {
		case <-stop:
			return
		case <-storeTick.C:
			s.DumpToFile(fname)
		}
	}
}
