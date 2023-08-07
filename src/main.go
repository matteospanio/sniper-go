/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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
		command := strings.Split(string(msg), " ")
		cmd := exec.Command(command[0], command[1:]...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
			return
		}

		cmd.Start()

		scanner := bufio.NewScanner(stdout)

		go func() {
			f, err := os.Create(fmt.Sprintf("./history/%s.txt", strings.Join(command, "_")))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()

			for scanner.Scan() {
				line := scanner.Text()
				f.WriteString(line + "\n")
				conn.WriteMessage(websocket.TextMessage, []byte(html.EscapeString(line)))
			}
		}()

		cmd.Wait()
	}
}

// GIN server
func main() {
	PORT, exists := os.LookupEnv("PORT")
	if !exists {
		PORT = "8080"
	}

	router := gin.Default()

	templates, err := ParseTemplateDir("templates")
	if err != nil {
		panic("Failed to load templates: " + err.Error())
	}
	router.SetHTMLTemplate(templates)
	router.Static("/assets", "./dist")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Sn1per web interface",
		})
	})
	router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c)
	})

	api := router.Group("/api")
	{
		api.GET("/results", func(c *gin.Context) {
			var results []ReportSummary

			fileList, err := os.ReadDir("./data")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, file := range fileList {
				report, err := readSummaryFile(file.Name())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				results = append(results, report)
			}
			c.JSON(http.StatusOK, gin.H{"results": results})
		})
		api.GET("/results/:hostName", func(c *gin.Context) {
			host := c.Param("hostName")
			report, err := readSummaryFile(host)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": report})
		})
	}

	router.GET("/results", func(c *gin.Context) {
		results, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s/api/results", PORT))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer results.Body.Close()

		if results.StatusCode != http.StatusOK {
			c.HTML(http.StatusNotFound, "no_results.html", gin.H{})
			return
		}
		var dataResponse map[string][]ReportSummary
		err = json.NewDecoder(results.Body).Decode(&dataResponse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "results.html", gin.H{
			"results": dataResponse["results"],
		})
	})

	err = router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
