package main

import (
	"fmt"
	"net/http"
	"strings"
	// "strconv"
	// "sync"
	// "time"

	"github.com/SinisterSup/kv-datastore/handle"
	"github.com/SinisterSup/kv-datastore/kvs"
	"github.com/gin-gonic/gin"
)

type Command struct {
	Cmnd string `json:"command"`
}

func ParseCommand(cmd string) (string, []string) {
	theCommand := strings.Trim(cmd, " ")
	cmdParts := strings.Split(theCommand, " ")
	operation := cmdParts[0]
	contents := cmdParts[1:]
	return operation, contents
}

func main() {
	myStore := &kvs.KeyValueStore{
		Store: make(map[string]*kvs.QueueChannel),
	}

	fmt.Println("Starting server...")
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		var cmd Command
		if err := c.ShouldBindJSON(&cmd); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		operation, contents := ParseCommand(cmd.Cmnd)

		switch operation {
		case "SET": 
			message, done, err := handle.SetHandler(contents, myStore)
			if err != nil {
				if done {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return 
			}
			if done {
				c.IndentedJSON(http.StatusOK, gin.H{"message": message})
			} else {
				c.IndentedJSON(http.StatusNotModified, gin.H{"message": message})
			}

		case "QPUSH": 
			message, done, err := handle.QpushHandler(contents, myStore)
			if err != nil {
				if done {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return 
			}
			c.IndentedJSON(http.StatusOK, gin.H{"message": message})

		default:
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid command"})
		}
	})


	router.GET("/", func(c *gin.Context) { 
		var getcmd Command
		if err := c.ShouldBindJSON(&getcmd); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		operation, contents := ParseCommand(getcmd.Cmnd)

		switch operation {
		case "GET":
			val, done, err := handle.GetHandler(contents, myStore)
			if err != nil {
				if done {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
				}
				return 
			}
			c.IndentedJSON(http.StatusOK, gin.H{"value": val})

		case "QPOP":
			val, done, err := handle.QpopHandler(contents, myStore)
			if err != nil {
				if done {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
				}
				return
			}
			c.IndentedJSON(http.StatusOK, gin.H{"value": val})

		case "BQPOP":
			val, done, err := handle.BqpopHandler(contents, myStore)
			if err != nil {
				if done {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return
			}
			c.IndentedJSON(http.StatusOK, gin.H{"value": val})

		default:	
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid command"})
		}
	})

	router.Run(":8080")
	// myStore.StartCleanupLoop(10) // Cleans up the expired keys every 10 seconds
}
