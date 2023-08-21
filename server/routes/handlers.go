package routes

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

var mutex sync.Mutex

var (
	Wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	connections = make(map[*websocket.Conn]bool)
)

func printer(status *Status, start string, target string) {
	for status.Running {
		read := exec.Command("bash", "-c", "tail -n 200 "+sniper.OutPath+"/sniper-"+target+"-"+start+".txt")

		out, err := read.Output()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		for conn := range connections {
			mutex.Lock()
			conn.WriteMessage(websocket.TextMessage, out)
			mutex.Unlock()
		}

		time.Sleep(1 * time.Second)
	}
}

func HandleWebSocket(c *gin.Context) {
	conn, err := Wsupgrader.Upgrade(c.Writer, c.Request, nil)
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
			mutex.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, []byte("[DONE]")); err != nil {
				fmt.Println("Failed to send message to client:", err)
			}
			mutex.Unlock()
		}

		fmt.Println("[DONE]")
	}
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
