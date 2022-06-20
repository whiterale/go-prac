package storage

import "github.com/whiterale/go-prac/internal/agent/buffer"

type Storage struct {
	*buffer.Buffer
}

func Init() *Storage {
	buf := buffer.Init()
	return &Storage{buf}
}
