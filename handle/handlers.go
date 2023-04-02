package handle

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	// "time"

	"github.com/SinisterSup/kv-datastore/kvs"
	"github.com/gin-gonic/gin"
)

func SetHandler(parts []string, kvs *kvs.KeyValueStore, c *gin.Context) {
	n := len(parts)

	if n < 2 || n > 5 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid number of arguments for set"})
	}

	key, value := parts[0], parts[1]

	switch {
	case n == 2:
		kvs.Set(key, value, 99999, "")
		c.IndentedJSON(http.StatusOK, gin.H{"message": "value set for key " + key})
	case n == 3:
		hasSet := kvs.Set(key, value, 99999, parts[2])
		if hasSet {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "value set for key " + key})
		} else {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "satisfies condition"})
		}
	case n == 4:
		if strings.EqualFold(parts[2], "EX") {
			timeInt, _ := strconv.Atoi(parts[3])
			kvs.Set(key, value, timeInt, "")
			c.IndentedJSON(http.StatusOK, gin.H{"message": "value set for key " + key})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid command"})
		}
	case n == 5:
		if strings.EqualFold(parts[2], "EX") {
			timeInt, _ := strconv.Atoi(parts[3])
			hasSet := kvs.Set(key, value, timeInt, parts[4])
			if hasSet {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "value set for key " + key})
			} else {
				c.IndentedJSON(http.StatusConflict, gin.H{"message": "satisfies condition"})
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid command"})
		}
	}
}

func GetHandler(parts []string, kvs *kvs.KeyValueStore, c *gin.Context) {
	n := len(parts)

	if n != 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid number of arguments for get"})
	}

	key := parts[0]
	val, ok := kvs.Get(key)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "key not found"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"value": val})
	}
}

func QpushHandler(parts []string, kvs *kvs.KeyValueStore, c *gin.Context) {
	n := len(parts)

	if n < 2 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid number of arguments for qpush"})
	}

	key := parts[0]
	values := parts[1:]

	if err := kvs.Qpush(key, values); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "values pushed to queue"})
	}
}

func QpopHandler(parts []string, kvs *kvs.KeyValueStore, c *gin.Context) {
	n := len(parts)

	if n != 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid number of arguments for qpop"})
	}

	key := parts[0]
	val, ok := kvs.Qpop(key)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": val})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"value": val})
	}
}

func BqpopHandler(parts []string, kvs *kvs.KeyValueStore, c *gin.Context) {
  n := len(parts)

  if n != 2 {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid number of arguments for bqpop"})
  }

  key := parts[0]
  t, err := strconv.Atoi(parts[1])
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid timeout request"})
  }
  timeout := convIntToTime(t)
  
  val, ok := kvs.Bqpop(key, timeout)
  if !ok {
    c.IndentedJSON(http.StatusNotFound, gin.H{"error": val})
  } else {
    c.IndentedJSON(http.StatusOK, gin.H{"value": val})
  }
}

func convIntToTime(t1 int) (time.Duration) {
  timeout := time.Duration(t1) * time.Second
  return timeout
  // unixTime := time.Unix(0, future.UnixNano())
  // return unixTime
}
