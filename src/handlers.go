package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

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

type Status struct {
	Running bool
}

func printer(status *Status, start string, target string) {
	for status.Running {
		read := exec.Command("bash", "-c", "tail -n 200 "+sniperOutPath+"/sniper-"+target+"-"+start+".txt")
		out, err := read.Output()
		if err != nil {
			fmt.Println(err)
			continue
		}

		for conn := range connections {
			conn.WriteMessage(websocket.TextMessage, out)
		}

		time.Sleep(1 * time.Second)
	}
}

func handleWebSocket(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	defer func() {
		conn.Close()
		delete(connections, conn)
	}()

	connections[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		target := string(msg)

		cmd := exec.Command("bash", "-c", "sniper -t "+target)
		fmt.Println("Running: ", cmd.String())

		status := Status{Running: true}

		date := time.Now().Format("2006-01-02-15-04")
		date = strings.Replace(date, "-", "", -1)

		cmd.Start()
		go printer(&status, date, target)

		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
			continue
		}

		status.Running = false

		for conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("[DONE]")); err != nil {
				fmt.Println("Failed to send message to client:", err)
			}
		}

		fmt.Println("[DONE]")
	}
}

func handleScreenshots(c *gin.Context) {
	host := c.Param("host")
	filename := c.Param("filename")

	c.File(sniperReportPath + "/" + host + "/" + filename)
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
		"section": "Results",
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

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func handleTasks(c *gin.Context) {
	query := c.Query("query")
	tasks, err := getTasks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, tasksTemplate, gin.H{
		"section": "Tasks",
		"tasks":   tasks,
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
