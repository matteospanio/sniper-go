package routes

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matteospanio/sniper-go/sniper"
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

func HandleScreenshots(c *gin.Context) {
	host := c.Param("host")
	filename := c.Param("filename")

	c.File(filepath.Join(sniper.ReportPath, host, "screenshots", filename))
}

func HandleApiResults(c *gin.Context) {
	query := strings.Trim(c.Query("q"), " \t\n")
	results, err := sniper.GetResults(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})

}

func HandleResults(c *gin.Context) {
	query := c.Query("query")
	results, err := sniper.GetResults(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, resultsTemplate, gin.H{
		"results": results,
		"section": "Results",
	})
}

func HandleApiSingleResult(c *gin.Context) {
	host := c.Param("hostName")
	report, err := sniper.CreateReport(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

func HandleSingleResult(c *gin.Context) {
	host := c.Param("hostName")

	report, err := sniper.CreateReport(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, resultTemplate, gin.H{
		"title":  "Report for " + host,
		"host":   host,
		"report": report,
	})
}

func HandleApiTasks(c *gin.Context) {
	query := strings.Trim(c.Query("q"), " \t\n")
	tasks, err := sniper.GetTasks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func HandleTasks(c *gin.Context) {
	query := c.Query("query")
	tasks, err := sniper.GetTasks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, tasksTemplate, gin.H{
		"section": "Tasks",
		"tasks":   tasks,
	})
}

func HandleApiSingleTask(c *gin.Context) {
	host := c.Param("hostName")
	task, err := sniper.GetTask(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func HandleSingleTask(c *gin.Context) {
	// TODO
	host := c.Param("id")

	report, err := sniper.GetTask(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, reportTemplate, gin.H{
		"report": report,
	})
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, indexTemplate, gin.H{
		"title": "Sn1per web interface",
	})
}
