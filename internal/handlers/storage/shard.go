package storage

import (
	"sync"
	"time"
)

type item struct {
	value          string
	expirationDate time.Time
}

type shard struct {
	mutex sync.Mutex
	m     map[string]*item
}

func NewShard() *shard {
	return &shard{
		m: make(map[string]*item),
	}
}
