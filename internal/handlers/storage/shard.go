package storage

import (
	"sync"
)

type shard struct {
	mutex sync.RWMutex
	m     map[string]string
}

func NewShard() *shard {
	return &shard{
		m: make(map[string]string),
	}
}
