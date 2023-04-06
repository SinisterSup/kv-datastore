package handle

import (
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/SinisterSup/kv-datastore/kvs"
)


func SetHandler(parts []string, kvs *kvs.KeyValueStore) (string, bool, error) {
	n := len(parts)

	if n < 2 || n > 5 {
		return "", true, errors.New("invalid number of arguments for set")
	}

	key, value := parts[0], parts[1]

	switch {
	case n == 2:
		kvs.Set(key, value, 99999, "")
		returnString := "value set for key: " + key
		return returnString, true, nil
	case n == 3:
		hasSet := kvs.Set(key, value, 99999, parts[2])
		if hasSet {
			returnString := "value set for key: " + key
			return returnString, true, nil
		} else {
			returnString := "Already satisfies condition for NX or XX"
			return returnString, false, nil
		}
	case n == 4:
		if strings.EqualFold(parts[2], "EX") {
			timeInt, _ := strconv.Atoi(parts[3]) // Integer value of time
			kvs.Set(key, value, timeInt, "")
			returnString := "value set for key: " + key
			return returnString, true, nil
		} else {
			return "", false, errors.New("invalid command")
		}
	case n == 5:
		if strings.EqualFold(parts[2], "EX") {
			timeInt, _ := strconv.Atoi(parts[3])
			hasSet := kvs.Set(key, value, timeInt, parts[4])
			if hasSet {
				returnString := "value set for key: " + key
				return returnString, true, nil
			} else {
				returnString := "Already satisfies condition for NX or XX"
				return returnString, false, nil
			}
		} else {
			return "", false, errors.New("invalid command")
		}
	}
	return "", false, errors.New("invalid command")
}

func GetHandler(parts []string, kvs *kvs.KeyValueStore) (string, bool, error) {
	n := len(parts)

	if n != 1 {
		return "", true, errors.New("invalid number of arguments for get")
	}

	key := parts[0]
	val, ok := kvs.Get(key)
	if !ok {
		return "", false, errors.New("key not found")
	} else {
		return val, true, nil
	}
}

func QpushHandler(parts []string, kvs *kvs.KeyValueStore) (string, bool, error) {
	n := len(parts)

	if n < 2 {
		return "", true, errors.New("invalid number of arguments for qpush")
	}

	key := parts[0]
	values := parts[1:]

	if err := kvs.Qpush(key, values); err != nil {
		return "", false, err
	} else {
		return "values pushed to queue", true, nil
	}
}

func QpopHandler(parts []string, kvs *kvs.KeyValueStore) (string, bool, error) {
	n := len(parts)

	if n != 1 {
		return "", true, errors.New("invalid number of arguments for qpop")
	}

	key := parts[0]
	val, ok := kvs.Qpop(key)
	if !ok {
		return "", false, errors.New(val)
	} else {
		return val, true, nil
	}
}

func BqpopHandler(parts []string, kvs *kvs.KeyValueStore) (string, bool, error) {
  n := len(parts)

  if n != 2 {
	return "", true, errors.New("invalid number of arguments for bqpop")
  }

  key := parts[0]
  t, err := strconv.ParseFloat(parts[1], 64)
  if err != nil {
	return "", false, errors.New("invalid timeout request")
  }
  timeout := convFloatToTime(t)
  
  val := kvs.Bqpop(key, timeout)
  return val, true, nil
}


func convFloatToTime(t1 float64) (time.Duration) {
  timeout := time.Duration(t1) * time.Second
  return timeout
  // unixTime := time.Unix(0, future.UnixNano())
  // return unixTime
}
