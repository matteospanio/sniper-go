package sniper

import "time"

type Report struct {
	Host            Target          `json:"host"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         ReportSummary   `json:"summary"`
	Date            time.Time       `json:"date"`
	Screens         []string        `json:"screens"`
	History         []Scan          `json:"history"`
}

func (r *Report) GetScreens() []string {
	return r.Screens
}

/*
 * Create a report from the scan
 * Return a Report struct
 */
func CreateReport(host string) (Report, error) {
	var result Report

	date, err := getLastScan(host)
	if err != nil {
		return Report{}, err
	}
	result.Date = date

	target := getTarget(host)
	if err != nil {
		return Report{}, err
	}
	result.Host = target

	summary, err := readSummaryFile(host)
	if err != nil {
		return Report{}, err
	}
	result.Summary = summary

	history, err := getHistory(host)
	if err != nil {
		return Report{}, err
	}
	result.History = history

	vuln, err := GetDetails(host)
	if err != nil {
		return result, err
	}
	result.Vulnerabilities = vuln

	result.Screens = getScreenshots(host)

	return result, nil
}
