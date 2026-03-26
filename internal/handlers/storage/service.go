package storage

import (
	"errors"

	"github.com/cespare/xxhash/v2"
)

type StorageService interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type storage struct {
	shards []*shard
}

const shardsCount = 256

var ErrItemNotFound = errors.New("item not found")

func NewStorageService() StorageService {
	s := make([]*shard, shardsCount)
	for i := range s {
		s[i] = NewShard()
	}

	return &storage{
		shards: s,
	}
}

func (s *storage) Get(key string) (string, error) {
	shard := s.shards[getShardIndex(key)]
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()

	item, ok := shard.m[key]
	if !ok {
		return "", ErrItemNotFound
	}

	return item, nil
}

func (s *storage) Set(key, value string) error {
	shard := s.shards[getShardIndex(key)]
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	shard.m[key] = value

	return nil
}

func (s *storage) Delete(key string) error {
	shard := s.shards[getShardIndex(key)]
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	delete(shard.m, key)

	return nil
}

func getShardIndex(value string) int {
	return int(xxhash.Sum64String(value) % uint64(shardsCount))
}
