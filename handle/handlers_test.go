package handle_test

import (
    "testing"
	"errors"
    // "sync"

    "github.com/SinisterSup/kv-datastore/kvs"
    "github.com/SinisterSup/kv-datastore/handle"
)

func TestSetHandler(t *testing.T) {
    kvs := &kvs.KeyValueStore{Store: make(map[string]*kvs.QueueChannel)}
    cases := []struct {
        name     string
        parts    []string
        expected string
        done     bool
        err      error
    }{
        {
            name:     "Invalid args",
            parts:    []string{"key1"},
            expected: "",
            done:     true,
            err:      errors.New("invalid number of arguments for set"),
        },
        {
            name:     "2 arguments",
            parts:    []string{"key1", "value1"},
            expected: "value set for key: key1",
            done:     true,
            err:      nil,
        },
        {
            name:     "3 arguments",
            parts:    []string{"key2", "value2", "NX"},
            expected: "value set for key: key2",
            done:     true,
            err:      nil,
        },
        {
            name:     "3 arguments",
            parts:    []string{"key2", "someVal", "NX"},
            expected: "Already satisfies condition for NX or XX",
            done:     false,
            err:      nil,
        },
        {
            name:     "3 arguments",
            parts:    []string{"key2", "newVal", "XX"},
            expected: "value set for key: key2",
            done:     true,
            err:      nil,
        },
        {
            name:     "3 arguments",
            parts:    []string{"newKey", "newVal", "XX"},
            expected: "Already satisfies condition for NX or XX",
            done:     false,
            err:      nil,
        },
        {
            name:     "3 arguments",
            parts:    []string{"key3", "value3", ""},
            expected: "value set for key: key3",
            done:     true,
            err:      nil,
        },
        {
            name:     "4 arguments",
            parts:    []string{"key4", "value4", "EX", "10"},
            expected: "value set for key: key4",
            done:     true,
            err:      nil,
        },
        {
            name:     "4 arguments",
            parts:    []string{"key4", "value4", "Hi", "10"},
            expected: "",
            done:     false,
            err:      errors.New("invalid command"),
        },
        {
            name:     "4 arguments",
            parts:    []string{"key4", "value4", "EX", "A1"},
            expected: "",
            done:     false,
            err:      errors.New("invalid time"),
        },
        {
            name:     "5 arguments",
            parts:    []string{"key5", "value5", "EX", "10", "NX"},
            expected: "value set for key: key5",
            done:     true,
            err:      nil,
        },
        {
            name:     "5 arguments",
            parts:    []string{"key5", "value5", "HX", "10", "NX"},
            expected: "",
            done:     false,
            err:      errors.New("invalid command"),
        },
        {
            name:     "5 arguments",
            parts:    []string{"key5", "value5", "EX", "NA", "NX"},
            expected: "",
            done:     false,
            err:      errors.New("invalid time"),
        },
        {
            name:     "5 arguments",
            parts:    []string{"key5", "someVal", "EX", "10", "NX"},
            expected: "Already satisfies condition for NX or XX",
            done:     false,
            err:      nil,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            actual, done, err := handle.SetHandler(c.parts, kvs)
			if actual != c.expected || done != c.done || (err != nil && c.err != nil && err.Error() != c.err.Error()) {
				t.Errorf("Expected: %v, %v, %v, but Got: %v, %v, %v", c.expected, c.done, c.err, actual, done, err)
			}
        })
    }
}

func TestGetHandler(t *testing.T) {
    s := &kvs.KeyValueStore{Store: make(map[string]*kvs.QueueChannel)}
    s.Set("key", "value", 9999, "")
    tests := []struct {
        name     string
        parts    []string
        expected string
        done     bool
        err      error
    }{
        {
            name: "Invalid args",
            parts: []string{"key", "value", "NA"},
            expected: "",
            done: true,
            err: errors.New("invalid number of arguments for get"),
        },
        {
            name: "valid Key",
            parts: []string{"key"},
            expected: "value",
            done: true,
            err: nil,
        },
        {
            name: "Unknown Key",
            parts: []string{"Unknown"},
            expected: "",
            done: false,
            err: errors.New("key not found"),
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            actual, done, err := handle.GetHandler(test.parts, s)
            if actual != test.expected || done != test.done || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
				t.Errorf("Expected: %v, %v, %v, but Got: %v, %v, %v", test.expected, test.done, test.err, actual, done, err)
			}
        })
    }
}

func TestQpushHandler(t *testing.T) {
    s := &kvs.KeyValueStore{Store: make(map[string]*kvs.QueueChannel)}
    tests := []struct {
        name     string
        parts    []string
        expected string
        done     bool
        err      error
    }{
        {
            name: "Invalid args",
            parts: []string{"OnlyKey"},
            expected: "",
            done: true,
            err: errors.New("invalid number of arguments for qpush"),
        },
        {
            name: "Key - Values...",
            parts: []string{"key", "value1", "value2", "value3"},
            expected: "values pushed to queue",
            done: true,
            err: nil,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            actual, done, err := handle.QpushHandler(test.parts, s)
            if actual != test.expected || done != test.done || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
                t.Errorf("Expected: %v, %v, %v, but Got: %v, %v, %v", test.expected, test.done, test.err, actual, done, err)
            }
        })
    }
}

func TestQpopHandler(t *testing.T) {
    s := &kvs.KeyValueStore{Store: make(map[string]*kvs.QueueChannel)}

    someErr := s.Qpush("key", []string{"value1", "value2"})
    if someErr != nil {
        t.Errorf("Error while pushing values to queue")
    }

    tests := []struct {
        name     string
        parts    []string
        expected string
        done     bool
        err      error
    }{
        {
            name: "Invalid args",
            parts: []string{"OnlyKey", "SomeValue"},
            expected: "",
            done: true,
            err: errors.New("invalid number of arguments for qpop"),
        },
        {
            name: "Unknown Key",
            parts: []string{"Unknown"},
            expected: "",
            done: false,
            err: errors.New("key not found"),
        },
        {
            name: "Valid Key",
            parts: []string{"key"},
            expected: "value2",
            done: true,
            err: nil,
        },
        {
            name: "Valid Key",
            parts: []string{"key"},
            expected: "value1",
            done: true,
            err: nil,
        },
        {
            name: "Empty Queue",
            parts: []string{"key"},
            expected: "",
            done: false,
            err: errors.New("queue is empty"),
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            actual, done, err := handle.QpopHandler(test.parts, s)
            if actual != test.expected || done != test.done || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
                t.Errorf("Expected: %v, %v, %v, but Got: %v, %v, %v", test.expected, test.done, test.err, actual, done, err)
            }
        })
    }
}

func TestBqpopHandler(t *testing.T) {
    s := &kvs.KeyValueStore{Store: make(map[string]*kvs.QueueChannel)}
    someErr := s.Qpush("key", []string{"value1", "value2"})
    if someErr != nil {
        t.Errorf("Error while pushing values to queue")
    }

    tests := []struct {
        name     string
        parts    []string
        expected string
        done     bool
        err      error
    }{
        {
            name: "Invalid args",
            parts: []string{"OnlyKey"},
            expected: "",
            done: true,
            err: errors.New("invalid number of arguments for bqpop"),
        },
        {
            name: "Invalid timeout",
            parts: []string{"OnlyKey", "NotANumber"},
            expected: "",
            done: false,
            err: errors.New("invalid timeout request"),
        },
        {
            name: "Unknown Key",
            parts: []string{"Unknown", "0"},
            expected: "",
            done: true,
            err: nil,
        },
        {
            name: "Valid Key",
            parts: []string{"key", "0"},
            expected: "value1",
            done: true,
            err: nil,
        },
        {
            name: "Valid Key",
            parts: []string{"key", "0"},
            expected: "value2",
            done: true,
            err: nil,
        },
        {
            name: "Empty Queue",
            parts: []string{"key", "0"},
            expected: "",
            done: true,
            err: nil,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            actual, done, err := handle.BqpopHandler(test.parts, s)
            if actual != test.expected || done != test.done || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
                t.Errorf("Expected: %v, %v, %v, but Got: %v, %v, %v", test.expected, test.done, test.err, actual, done, err)
            }
        })
    }

    // var wg sync.WaitGroup
    // wg.Add(1)
    // go func() {
    //     t.Log("Testing Concurrent blocking queue pop")
    //     actual, done, err := handle.BqpopHandler([]string{"NewKey", "10"}, s)
    //     someErr := s.Qpush("NewKey", []string{"value1", "value2"})
    //     if someErr != nil {
    //         t.Errorf("Error while pushing values to queue")
    //     }
    //     if actual != "value1" || done != true || err != nil {
    //         t.Errorf("Expected: %v, %v, %v, but Got: %v, %v, %v", "value1", true, nil, actual, done, err)
    //     }
    //     wg.Done()
    // }()
    // wg.Wait()
}