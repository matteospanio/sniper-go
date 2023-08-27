/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package sniper

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ReportSummary The vulnerability report from sn1per
type ReportSummary struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Info     int `json:"info"`
	Score    int `json:"score"`
}

type Vulnerability struct {
	Severity    string `json:"severity"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Task struct {
	IsRunning bool      `json:"isRunning"`
	Target    Target    `json:"target"`
	Date      time.Time `json:"date"`
}

type Scan struct {
	Date   time.Time `json:"date"`
	Job    string    `json:"scan"`
	Target string    `json:"target"`
}

const (
	ReportPath = "/usr/share/sniper/loot/workspace"
	OutPath    = "/usr/share/sniper/loot/output"
)

func readSummaryFile(host string) (ReportSummary, error) {
	// Parse file
	content, err := readSniperFile(host, fmt.Sprintf("vulnerabilities/vulnerability-report-%s.txt", host))
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

func getHistory(host string) ([]Scan, error) {
	var result []Scan

	scans, err := readSniperFile(host, "scans/tasks.txt")
	if err != nil {
		return result, err
	}

	scanList := strings.Split(scans, "\n")
	for _, scan := range scanList {
		if scan == "" {
			continue
		}

		splitted := strings.SplitN(scan, " ", 3)

		if len(splitted) < 3 {
			return result, errors.New("error parsing scan history for host " + host)
		}

		date, err := time.Parse("2006-01-02 15:04", splitted[2])
		if err != nil {
			return result, err
		}

		result = append(result, Scan{
			Date:   date,
			Job:    splitted[1],
			Target: splitted[0],
		})
	}

	return result, nil
}

/* Get the start moment of the scan
 * Return the date from the scans/tasks.txt file
 */
func getLastScan(host string) (time.Time, error) {

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
			ReportPath,
			host,
		))
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		result = append(result, "/screens/"+host+"/"+file.Name())
	}

	return result
}

func GetDetails(host string) ([]Vulnerability, error) {
	var result []Vulnerability

	content, err := readSniperFile(host, "vulnerabilities/sc0pe-all-vulnerabilities-sorted.txt")
	if err != nil {
		return result, err
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		splittedLine := strings.SplitN(line, ",", 3)

		severity := strings.Split(splittedLine[0], "-")[1]
		name := splittedLine[1]
		description := splittedLine[2]

		// append to result
		result = append(result, Vulnerability{
			Severity:    severity,
			Name:        name,
			Description: description,
		})
	}

	return result, nil
}

func readSniperFile(host string, filePath string) (string, error) {
	content, err := os.ReadFile(
		fmt.Sprintf(
			"%s/%s/%s",
			ReportPath,
			host,
			filePath,
		))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func GetResults(query string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query = strings.Trim(query, " \t\n")

	fileList, err := os.ReadDir(ReportPath)
	if err != nil {
		return results, err
	}

	for _, file := range fileList {
		report, err := readSummaryFile(file.Name())
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return results, err
		}

		date, err := getLastScan(file.Name())
		if err != nil {
			return results, err
		}

		tmp := map[string]interface{}{
			"date":    date,
			"summary": report,
			"host":    file.Name(),
		}

		results = append(results, tmp)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i]["date"].(time.Time).After(results[j]["date"].(time.Time))
	})

	if query != "" {
		var filteredResults []map[string]interface{}
		for _, result := range results {
			if strings.Contains(result["host"].(string), query) {
				filteredResults = append(filteredResults, result)
			}
		}
		results = filteredResults
	}

	return results, nil
}

func isRunning(host string) (bool, error) {
	content, err := readSniperFile(host, "scans/tasks-running.txt")
	if err != nil {
		return false, err
	}

	isRunning, err := strconv.Atoi(strings.Split(content, "\n")[0])
	if err != nil {
		return false, err
	}

	return isRunning > 0, nil
}

func GetTasks(query string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query = strings.Trim(query, " \t\n")

	fileList, err := os.ReadDir(ReportPath)
	if err != nil {
		return results, err
	}

	for _, file := range fileList {
		date, err := getLastScan(file.Name())
		if err != nil {
			return results, err
		}

		running, err := isRunning(file.Name())
		if err != nil {
			return results, err
		}

		tmp := map[string]interface{}{
			"date":      date,
			"host":      file.Name(),
			"isRunning": running,
		}

		results = append(results, tmp)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i]["date"].(time.Time).After(results[j]["date"].(time.Time))
	})

	if query != "" {
		var filteredResults []map[string]interface{}
		for _, result := range results {
			if strings.Contains(result["host"].(string), query) {
				filteredResults = append(filteredResults, result)
			}
		}
		results = filteredResults
	}

	return results, nil
}

func GetTask(id string) (Task, error) {
	var result Task
	result.Target = getTarget(id)

	running, err := isRunning(id)
	if err != nil {
		return result, err
	}

	date, err := getLastScan(id)
	if err != nil {
		return result, err
	}

	result.IsRunning = running
	result.Date = date

	return result, nil
}
