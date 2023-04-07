package kvs

import (
    "testing"
	"time"
)
func TestSetAndGet(t *testing.T) {
	kvs := KeyValueStore{Store: make(map[string]*QueueChannel)}

	key := "test_key"
	value := "test_value"
	expiration := 10

	// Test Set
	if !kvs.Set(key, value, expiration, "NX") {
		t.Errorf("Set() FAILED: must set new key-value pair for NX")
	}
	if kvs.Set(key, "new_test_value", expiration, "NX") {
		t.Errorf("Set() FAILED: must return false when the key already exists for NX")
	}
	if !kvs.Set(key, "new_value", expiration, "XX") {
		t.Errorf("Set() FAILED: must update the key-value pair that already exists for XX")
	}
	if kvs.Set("new_key", "new_value", expiration, "XX") {
		t.Errorf("Set() FAILED: must return false when the key does not exist for XX")
	}

	key1 := "test_key1"
	value1 := "test_value1"
	expiration1 := 0
	if !kvs.Set(key1, value1, expiration1, "") {
		t.Errorf("Set() FAILED: must set new key-value pair with no expiration")
	}

	// Test Get
	val, ok := kvs.Get(key)
	if !ok || val != "new_value" {
		t.Errorf("Get() FAILED: must retrieve the value for the key")
	}
	val1, ok1 := kvs.Get(key1)
	if !ok1 || val1 != "test_value1" {
		t.Errorf("Get() FAILED: must retrieve the value for the key")
	}

	// Test Expiration
	time.Sleep(time.Duration(expiration+1) * time.Second)
	if _, ok := kvs.Get(key); ok {
		t.Errorf("Get() FAILED: must return false for expired key")
	}
}

func BenchmarkSetAndGet(b *testing.B) {
	kvs := KeyValueStore{Store: make(map[string]*QueueChannel)}
	key := "test_key"
	value := "test_value"
	expiration := 10

	for i := 0; i < b.N; i++ {
		kvs.Set(key, value, expiration, "")
		kvs.Get(key)
	}
}

func TestQueueOperations(t *testing.T) {
	kvs := KeyValueStore{Store: make(map[string]*QueueChannel)}

	key := "test_queue"

	// Test Qpush
	values := []string{"value1", "value2", "value3"}
	if err := kvs.Qpush(key, values); err != nil {
		t.Errorf("Qpush() FAILED: to push values to the queue")
	}

	// Test Qpop
	val, ok := kvs.Qpop(key)
	if !ok || val != "value3" {
		t.Errorf("Qpop() FAILED: couldn't pop the value from the queue")
	}
	if _, ok := kvs.Qpop(key); !ok {
		t.Errorf("Qpop() FAILED: couldn't pop the value from the queue")
	}
	if _, ok := kvs.Qpop(key); !ok {
		t.Errorf("Qpop() FAILED: couldn't pop the value from the queue")
	}
	if _, ok := kvs.Qpop(key); ok {
		t.Errorf("Qpop() FAILED: must return false when the queue is empty, but got %v", ok)
	}
	if _, ok := kvs.Qpop("UnknownKey"); ok {
		t.Errorf("Qpop() FAILED: must return false when the key does not exist, but got %v", ok)
	}

	// Test Bqpop
	if val := kvs.Bqpop(key, time.Second); val != "" {
		t.Errorf("Bqpop() FAILED: expected to return null when the queue is empty, but got %v", val)
	}
	if val := kvs.Bqpop("UnknownKey", time.Second); val != "" {
		t.Errorf("Bqpop() FAILED: expected to return null when the key does not exist, but got %v", val)
	}

	kvs.Qpush(key, []string{"value1", "value2", "value3"})

	doneChan := make(chan bool)

	go func() {
		time.Sleep(time.Second)

		val1 := kvs.Bqpop(key, time.Second)
		if val1 != "value1" {
			t.Errorf("Bqpop() FAILED: expected to pop %v, but got %v", "value1", val1)
		}
		val2 := kvs.Bqpop(key, time.Second)
		if val2 != "value2" {
			t.Errorf("Bqpop() FAILED: expected to pop %v, but got %v", "value2", val2)
		}
		val3 := kvs.Bqpop(key, time.Second)
		if val3 != "value3" {
			t.Errorf("Bqpop() FAILED: expected to pop %v, but got %v", "value3", val3)
		}
		emptyVal1 := kvs.Bqpop(key, time.Second)
		if emptyVal1 != "" {
			t.Errorf("Bqpop() FAILED: expected to return null when the queue is empty, but got %v", emptyVal1)
		}
		emptyVal2 := kvs.Bqpop("NewKey", 0*time.Second)
		if emptyVal2 != "" {
			t.Errorf("Bqpop() FAILED: expected to return null when the key does not exist, but got %v", emptyVal2)
		}

		doneChan <- true
	}()

	<-doneChan
}

func BenchmarkQpush(b *testing.B) {
    kvs := KeyValueStore{Store: make(map[string]*QueueChannel)}

    key := "test_queue"
    // Push values to the queue
    values := []string{"value1", "value2", "value3"}
    if err := kvs.Qpush(key, values); err != nil {
        b.Errorf("Qpush() FAILED: to push values to the queue")
    }
}

func BenchmarkQpop(b *testing.B) {
    kvs := KeyValueStore{Store: make(map[string]*QueueChannel)}
    key := "test_queue"
    // Benchmark Qpop
	for i := 0; i < b.N; i++ {	
    	values := []string{"value1", "value2", "value3"}
		if err := kvs.Qpush(key, values); err != nil {
    	    b.Errorf("Qpush() FAILED: to push values to the queue")
    	}

		val, ok := kvs.Qpop(key)
		if val != "value3" {
			b.Errorf("Qpop() FAILED: couldn't pop the value from the queue")
		}
		if !ok {
			b.Errorf("Qpop() FAILED: either the queue is empty or the key does not exist")
		}
	}
}

func BenchmarkBqpop(b *testing.B) {
	kvs := KeyValueStore{Store: make(map[string]*QueueChannel)}

	key := "test_queue"
	values := []string{"value1", "value2", "value3"}
	kvs.Qpush(key, values)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kvs.Bqpop(key, time.Second)
	}
}