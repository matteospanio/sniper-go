package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ReportSummary struct {
	Host     string `json:"host"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
	Info     int    `json:"info"`
	Score    int    `json:"score"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func main() {
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
		wshandler(c.Writer, c.Request)
	})
	router.GET("/results", func(c *gin.Context) {
		results := []ReportSummary{}

		fileList, err := os.ReadDir("./data")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}

		for _, file := range fileList {
			log.Println(file.Name())
			report, err := readSummaryFile(file.Name())
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			results = append(results, report)
		}
		c.JSON(200, gin.H{"results": results})
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

func readSummaryFile(host string) (ReportSummary, error) {
	// Read file
	os.ReadDir(fmt.Sprintf("/data/%s", host))
	report := ReportSummary{}

	// Parse file
	content, err := ioutil.ReadFile(
		fmt.Sprintf(
			"./data/%s/vulnerabilities/vulnerability-report-%s.txt",
			host,
			host,
		))
	if err != nil {
		return report, err
	}

	content_str := string(content)

	re := regexp.MustCompile(`Critical: (\d+)\nHigh: (\d+)\nMedium: (\d+)\nLow: (\d+)\nInfo: (\d+)\nScore: (\d+)`)
	matches := re.FindStringSubmatch(content_str)

	if len(matches) < 7 || matches == nil {
		return report, fmt.Errorf("regxep: no matches for %s", host)
	}

	crt, _ := strconv.Atoi(matches[1])
	hig, _ := strconv.Atoi(matches[2])
	med, _ := strconv.Atoi(matches[3])

	low, _ := strconv.Atoi(matches[4])
	inf, _ := strconv.Atoi(matches[5])
	scr, _ := strconv.Atoi(matches[6])

	return ReportSummary{
		Host:     host,
		Critical: crt,
		High:     hig,
		Medium:   med,
		Low:      low,
		Info:     inf,
		Score:    scr,
	}, nil

}
