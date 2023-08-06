/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

/*
 * GIN server
 */
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
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		})
	})
	router.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})
	router.GET("/results", func(c *gin.Context) {
		var results []ReportSummary

		fileList, err := os.ReadDir("./data")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		for _, file := range fileList {
			log.Println(file.Name())
			report, err := readSummaryFile(file.Name())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			results = append(results, report)
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	})
	router.GET("/results/:hostName", func(c *gin.Context) {
		host := c.Param("hostName")
		report, err := readSummaryFile(host)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"data": report})
	})

	err = router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
