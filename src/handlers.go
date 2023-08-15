package main

import (
	"bufio"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	resultsTemplate = "results.html"
	resultTemplate  = "result.html"
	reportTemplate  = "report.html"
	indexTemplate   = "index.html"
	tasksTemplate   = "tasks.html"
)

func handleWebSocket(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Processing: %s\n", msg)
		command := string(msg)
		if command == "exit" {
			break
		}
		os.Remove("/tmp/sniper.log")

		command += " > /tmp/sniper.log"
		cmd := exec.Command("bash", "-c", command)

		cmd.Start()

		file, _ := os.Open("/tmp/sniper.log")
		defer file.Close()

		scanner := bufio.NewScanner(file)

		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				conn.WriteMessage(websocket.TextMessage, []byte(html.EscapeString(line)))
			}

			conn.WriteMessage(websocket.TextMessage, []byte("[DONE]"))
			fmt.Println("[DONE]")
		}()

		err = cmd.Wait()
		if err != nil {
			return
		}
	}
}

func handleApiResults(c *gin.Context) {
	query := strings.Trim(c.Query("q"), " \t\n")
	results, err := getResults(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})

}

func handleResults(c *gin.Context) {
	query := c.Query("query")
	results, err := getResults(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, resultsTemplate, gin.H{
		"results": results,
	})
}

func handleApiSingleResult(c *gin.Context) {
	host := c.Param("hostName")
	report, err := createReport(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

func handleSingleResult(c *gin.Context) {
	host := c.Param("hostName")

	report, err := createReport(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, resultTemplate, gin.H{
		"title":  "Report for " + host,
		"report": report,
	})
}

func handleApiTasks(c *gin.Context) {
	query := strings.Trim(c.Query("q"), " \t\n")
	tasks, err := getTasks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, tasksTemplate, gin.H{"tasks": tasks})
}

func handleTasks(c *gin.Context) {
	query := c.Query("query")
	tasks, err := getTasks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, tasksTemplate, gin.H{
		"tasks": tasks,
	})
}

func handleApiSingleTask(c *gin.Context) {
	host := c.Param("hostName")
	task, err := getTask(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func handleSingleTask(c *gin.Context) {
	// TODO
	host := c.Param("id")

	report, err := getTask(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, reportTemplate, gin.H{
		"report": report,
	})
}
