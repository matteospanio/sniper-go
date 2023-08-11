package main

import (
	"bufio"
	"fmt"
	"html"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		cmd := exec.Command("bash", "-c", command)

		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		scanner := bufio.NewScanner(stdout)

		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				conn.WriteMessage(websocket.TextMessage, []byte(html.EscapeString(line)))
			}

			conn.WriteMessage(websocket.TextMessage, []byte("[DONE]"))
			fmt.Println("[DONE]")
		}()

		cmd.Wait()
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

	c.HTML(http.StatusOK, results_template, gin.H{
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

	c.HTML(http.StatusOK, report_template, gin.H{
		"report": report,
	})
}
