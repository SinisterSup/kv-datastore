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

func main() {
	myStore := &kvs.KeyValueStore{
		Store: make(map[string][]*kvs.KeyValueItem),
	}

	fmt.Println("Starting server...")
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		var cmd Command
		if err := c.ShouldBindJSON(&cmd); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		theCommand := strings.Trim(cmd.Cmnd, " ")
		cmdParts := strings.Split(theCommand, " ")
		operation := cmdParts[0]

		if operation == "SET" {
			handle.SetHandler(cmdParts[1:], myStore, c)
		} else if operation == "QPUSH" {
			handle.QpushHandler(cmdParts[1:], myStore, c)
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid command"})
		}
	})

	router.GET("/", func(c *gin.Context) { 
		var getcmd Command
		if err := c.ShouldBindJSON(&getcmd); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		theCommand := strings.Trim(getcmd.Cmnd, " ")
		cmdParts := strings.Split(theCommand, " ")
		operation := cmdParts[0]

		switch operation {
		case "GET":
			handle.GetHandler(cmdParts[1:], myStore, c)
		case "QPOP":
			handle.QpopHandler(cmdParts[1:], myStore, c)
		case "BQPOP":
			handle.BqpopHandler(cmdParts[1:], myStore, c)
		default:	
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid command"})
		}
	})

	router.Run(":8080")
	myStore.StartCleanupLoop(10) // Cleans up the expired keys every 10 seconds
}
