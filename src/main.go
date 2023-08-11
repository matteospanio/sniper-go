/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	results_template = "results.html"
	report_template  = "report.html"
	index_template   = "index.html"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
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
		c.HTML(http.StatusOK, index_template, gin.H{
			"title": "Sn1per web interface",
		})
	})

	router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c)
	})

	api := router.Group("/api")
	{
		api.GET("/results", handleApiResults)
		api.GET("/results/:hostName", handleApiSingleResult)
	}

	router.GET("/results", handleResults)
	router.GET("/results/:hostName", handleSingleResult)

	err = router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
