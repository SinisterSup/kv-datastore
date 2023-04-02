package kvs

import (
	"strings"
	"sync"
	"time"
)

type KeyValueStore struct {
	mu    sync.Mutex
	Store map[string][]*KeyValueItem
}

type KeyValueItem struct {
	value      string
	expiration *time.Time
}

func (s *KeyValueStore) Set(key, value string, expiration int, condition string) (bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.Store[key]

	if strings.EqualFold(condition, "NX") {
		if exists {
			return false
		}
	} else if strings.EqualFold(condition, "XX") {
		if !exists {
			return false
		}
	}

	exp := time.Now().Add(time.Duration(expiration) * time.Second)
	s.Store[key] = []*KeyValueItem{{value: value, expiration: &exp}}

	return true
}

func (s *KeyValueStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.Store[key]

	if !exists {
		return "", false
	}

	if item[0].expiration != nil && time.Now().After(*item[0].expiration) {
		delete(s.Store, key)
		return "", false
	}

	return item[0].value, true
}

func (s *KeyValueStore) Qpush(key string, values []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var valItems []*KeyValueItem

	for _, val := range values {
		valItems = append(valItems, &KeyValueItem{value: val, expiration: nil})
	}

	if _, exists := s.Store[key]; exists {
		s.Store[key] = append(s.Store[key], valItems...)
	} else {
		s.Store[key] = valItems
	}
	
	// To notify waiting threads that new items have been added to the queue.
	itemCond := sync.NewCond(&s.mu)
	itemCond.Broadcast()

	return nil
}

func (s *KeyValueStore) Qpop(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.Store[key]
	n := len(item)

	if !exists {
		return "key not found", false
	}
	if n == 0 {
		return "queue is empty", false
	}

	val := item[n-1].value
	s.Store[key] = item[:n-1]

	return val, true
}

func (s *KeyValueStore) Bqpop(key string, timeout time.Duration) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.Store[key]

	if !exists {
		return "key not found", false
	}

	startTime := time.Now()
	for {
		if len(item) > 0 {
			return s.Qpop(key)
		}
		if timeout == 0 {
			return "", true
		}

		// To wait for new items to be added to the queue.
		itemCond := sync.NewCond(&s.mu)
		itemCond.Wait()

		if time.Since(startTime) >= timeout {
			return "", true
		}
	}	
}

func (s *KeyValueStore) StartCleanupLoop(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	for {
		<-ticker.C
		s.mu.Lock()
		for key, item := range s.Store {
			if item[0].expiration != nil && time.Now().After(*item[0].expiration) {
				delete(s.Store, key)
			}
		} 
		s.mu.Unlock()
	}
}
