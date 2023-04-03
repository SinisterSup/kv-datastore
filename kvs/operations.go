package kvs

import (
	"strings"
	"sync"
	"time"
)

type KeyValueStore struct {
	mu    sync.Mutex
	Store map[string]*QueueChannel
}

type KeyValueItem struct {
	value      string
	expiration *time.Time
}

type QueueChannel struct {
	queue []*KeyValueItem
	channel chan string
}

func (s *KeyValueStore) Set(key, value string, expiration int, condition string) bool {
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
	s.Store[key] = &QueueChannel{
		queue:   []*KeyValueItem{{value: value, expiration: &exp}},
		channel: make(chan string, 25),
	}

	return true
}

func (s *KeyValueStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.Store[key]

	if !exists {
		return "", false
	}

	if item.queue[0].expiration != nil && time.Now().After(*item.queue[0].expiration) {
		delete(s.Store, key)
		return "", false
	}

	return item.queue[0].value, true
}

func (s *KeyValueStore) Qpush(key string, values []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range values {
		exp := time.Now().Add(24 * time.Hour)
		item := &KeyValueItem{value: val, expiration: &exp}

		if _, exists := s.Store[key]; exists {
			select {
				case s.Store[key].channel <- val: // send the value to the channel
					s.Store[key].queue = append(s.Store[key].queue, item)
				default:
					s.Store[key].queue = append(s.Store[key].queue, item)
			}
		} else {
			channel := make(chan string, 25)
			s.Store[key] = &QueueChannel{
				queue:   []*KeyValueItem{item},
				channel: channel,
			}
			select {
				case s.Store[key].channel <- val: // send the value to the channel
					// s.Store[key].queue = append(s.Store[key].queue, item)
				default:
					s.Store[key].queue = append(s.Store[key].queue, item)
			}
		}
	}

	return nil
}

func (s *KeyValueStore) Qpop(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.Store[key]

	if !exists {
		return "key not found", false
	}

	n := len(item.queue)
	if n == 0 {
		return "queue is empty", false
	}

	val := item.queue[n-1].value
	item.queue = item.queue[:n-1]

	return val, true
}

func (s *KeyValueStore) Bqpop(key string, timeout time.Duration) (string) {
	s.mu.Lock()
	item, exists := s.Store[key]

	resultChan := make(chan string, 1)
	go func(key string, item *QueueChannel) {
		time.Sleep(timeout)

		if !exists {
			resultChan <- "" // "key not found
			return
		}

		n := len(item.queue)
		if n == 0 {
			resultChan <- "" // "queue is empty"
			return 
		}

		select {
			case val := <-item.channel:
				item.queue = item.queue[1:]   // If you wanna pop from front of the queue, use this line
				// item.queue = item.queue[:n-1] // If you wanna pop from back of the queue, use this line
				resultChan <- val
				return
			default:
				val := item.queue[0].value // If you wanna pop from front of the queue, use this line
				item.queue = item.queue[1:]
			// val := item.queue[n-1].value  // If you wanna pop from back of the queue, use this line
			// item.queue = item.queue[:n-1]
				resultChan <- val
				return
		}
	}(key, item)

	s.mu.Unlock()	
	
	popVal := <- resultChan
	return popVal

}


func (s *KeyValueStore) StartCleanupLoop(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	for {
		<-ticker.C
		s.mu.Lock()
		for key, item := range s.Store {
			if item.queue[0].expiration != nil && time.Now().After(*item.queue[0].expiration) {
				delete(s.Store, key)
			}
		} 
		s.mu.Unlock()
	}
}