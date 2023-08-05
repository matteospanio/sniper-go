package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	templates, err := ParseTemplateDir("templates")
	if err != nil {
		panic("Failed to load templates: " + err.Error())
	}

	router.SetHTMLTemplate(templates)

	router.Static("/assets", "./static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		})
	})

	router.Run("0.0.0.0:8080")
}

func snipe(c *gin.Context) {

	cCp := c.Copy()

	go func() {

		var command = fmt.Sprintf("sniper -t %s", cCp.Query("target"))
		cmd, err := exec.Command("bash", "-c", command).Output()

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"output": string(cmd)})
		}
	}()

	c.JSON(200, gin.H{"message": "Snipe started"})
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}
