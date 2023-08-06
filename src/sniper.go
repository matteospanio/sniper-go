/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

// The vulnerability report from sn1per
type ReportSummary struct {
	Host     string `json:"host"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
	Info     int    `json:"info"`
	Score    int    `json:"score"`
}

// Target's name and IP address
type Target struct {
	Name string     `json:"name"`
	IP   net.IPAddr `json:"ip"`
}

func readSummaryFile(host string) (ReportSummary, error) {
	// Parse file
	content, err := os.ReadFile(
		fmt.Sprintf(
			"./data/%s/vulnerabilities/vulnerability-report-%s.txt",
			host,
			host,
		))
	if err != nil {
		return ReportSummary{}, err
	}

	contentStr := string(content)

	re := regexp.MustCompile(`Critical: (\d+)\nHigh: (\d+)\nMedium: (\d+)\nLow: (\d+)\nInfo: (\d+)\nScore: (\d+)`)
	matches := re.FindStringSubmatch(contentStr)

	if len(matches) < 7 || matches == nil {
		return ReportSummary{}, fmt.Errorf("regxep: no matches for %s", host)
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
