/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
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
		wshandler(c.Writer, c.Request)
	})

	api := router.Group("/api")
	{
		api.GET("/results", func(c *gin.Context) {
			var results []ReportSummary

			fileList, err := os.ReadDir("./data")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			for _, file := range fileList {
				report, err := readSummaryFile(file.Name())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			}
			c.JSON(http.StatusOK, gin.H{"data": report})
		})
	}

	router.GET("/results", func(c *gin.Context) {
		results, err := http.Get("http://0.0.0.0:8080/api/results")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer results.Body.Close()

		if results.StatusCode != http.StatusOK {
			c.HTML(http.StatusOK, "<div><p>No results found</p></div>", gin.H{})
		}

		var dataResponse map[string][]ReportSummary
		err = json.NewDecoder(results.Body).Decode(&dataResponse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
