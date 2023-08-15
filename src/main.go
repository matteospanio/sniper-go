/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	functions := template.FuncMap{
		"upper": func(str string) string {
			return strings.ToUpper(str)
		},
		"lower": func(str string) string {
			return strings.ToLower(str)
		},
		"capitalize": func(str string) string {
			return strings.Title(str)
		},
	}

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}

	components, err := filepath.Glob(templatesDir + "/components/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFilesFuncs(filepath.Base(include), functions, files...)
	}

	for _, component := range components {
		r.AddFromFilesFuncs(filepath.Base(component), functions, component)
	}

	return r
}

var (
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	connections = make(map[*websocket.Conn]bool)
)

const (
	distPath      = "./dist"
	screensPath   = "./screens"
	templatesPath = "./templates"
)

// GIN server
func main() {
	PORT, exists := os.LookupEnv("PORT")
	if !exists {
		PORT = "8080"
	}

	router := gin.Default()
	router.HTMLRender = loadTemplates(templatesPath)
	router.Static("/assets", distPath)
	router.Static("/screens", screensPath)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, indexTemplate, gin.H{
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
		api.GET("/tasks", handleApiTasks)
		api.GET("/tasks/:hostName", handleApiSingleTask)
	}

	router.GET("/results", handleResults)
	router.GET("/results/:hostName", handleSingleResult)
	router.GET("/tasks", handleTasks)
	router.GET("/tasks/:id", handleSingleTask)

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
