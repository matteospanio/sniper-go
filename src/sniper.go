/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// The vulnerability report from sn1per
type ReportSummary struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Info     int `json:"info"`
	Score    int `json:"score"`
}

type Vulnerability struct {
	Name        string   `json:"name"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	Remediation string   `json:"remediation"`
	References  []string `json:"references"`
}

type Report struct {
	Host            Target          `json:"host"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         ReportSummary   `json:"summary"`
	Date            time.Time       `json:"date"`
	Screens         []string        `json:"screens"`
}

// Target's name and IP address
type Target struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

const sniper_report_path = "/usr/share/sniper/loot/workspace"

/*
 * Create a report from the scan
 * Return a Report struct
 */
func createReport(host string) (Report, error) {
	var result Report

	date, err := getDate(host)
	if err != nil {
		return Report{}, err
	}
	result.Date = date

	target, err := getTarget(host)
	if err != nil {
		return Report{}, err
	}
	result.Host = target

	summary, err := readSummaryFile(host)
	if err != nil {
		return Report{}, err
	}
	result.Summary = summary

	result.Screens = getScreenshots(host)

	return result, nil
}

func readSummaryFile(host string) (ReportSummary, error) {
	// Parse file
	content, err := readSniperFile(host, "vulnerabilities/vulnerability-report.txt")
	if err != nil {
		return ReportSummary{}, err
	}

	re := regexp.MustCompile(`Critical: (\d+)\nHigh: (\d+)\nMedium: (\d+)\nLow: (\d+)\nInfo: (\d+)\nScore: (\d+)`)
	matches := re.FindStringSubmatch(content)

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
		Critical: crt,
		High:     hig,
		Medium:   med,
		Low:      low,
		Info:     inf,
		Score:    scr,
	}, nil

}

/* Read the vulnerability report from sn1per
 * Return the content of the domains/targets-all-sorted.txt file
 * and the content of the ips/ips-all-sorted.txt file
 */
func getTarget(host string) (Target, error) {
	result := Target{}

	// find the host name
	targets, err := readSniperFile(host, "domains/targets-all-sorted.txt")
	if err != nil {
		return Target{}, err
	}

	// find the host IP
	ips, err := readSniperFile(host, "ips/ips-all-sorted.txt")
	if err != nil {
		return Target{}, err
	}

	result.Name = strings.Split(targets, "\n")[0]
	result.IP = strings.Split(ips, "\n")[0]

	return result, nil
}

/* Get the start moment of the scan
 * Return the date from the scans/tasks.txt file
 */
func getDate(host string) (time.Time, error) {

	content, err := readSniperFile(host, "scans/tasks.txt")
	if err != nil {
		return time.Now(), err
	}

	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2})`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 3 || matches == nil {
		return time.Now(), fmt.Errorf("regxep: no matches for %s", host)
	}

	date := matches[1]
	hour := matches[2]

	result, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, hour))
	if err != nil {
		return time.Now(), err
	}

	return result, nil
}

func getScreenshots(host string) []string {
	var result []string

	files, err := os.ReadDir(
		fmt.Sprintf(
			"%s/%s/screenshots",
			sniper_report_path,
			host,
		))
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		result = append(result, file.Name())
	}

	return result
}

func getDetails(host string) ([]Vulnerability, error) {
	var result []Vulnerability

	content, err := readSniperFile(host, "vulnerabilities/vulnerability-report.txt")
	if err != nil {
		return []Vulnerability{}, err
	}

	// skip first 10 lines of content
	lines := strings.Split(content, "\n")[10:]
	for _, line := range lines {
		if line == "" {
			continue
		}

		// parse line
		re := regexp.MustCompile(`(\w+) (.*)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 || matches == nil {
			return []Vulnerability{}, fmt.Errorf("regxep: no matches for %s", host)
		}

		// parse references
		references := strings.Split(matches[2], ", ")
		for i := range references {
			references[i] = strings.TrimSpace(references[i])
		}

		// append to result
		result = append(result, Vulnerability{
			Name:        matches[1],
			Severity:    matches[1],
			Description: matches[2],
			Remediation: "",
			References:  references,
		})
	}

	return result, nil
}

func readSniperFile(host string, filePath string) (string, error) {
	content, err := os.ReadFile(
		fmt.Sprintf(
			"%s/%s/%s",
			sniper_report_path,
			host,
			filePath,
		))
	if err != nil {
		return "", err
	}

	return string(content), nil
}
