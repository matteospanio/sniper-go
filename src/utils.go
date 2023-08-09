/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
)

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

func GenerateReport() {
	// TODO
}

func isIP(s string) bool {
	// Regex for IPv4
	ip4Regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	// Regex for IPv6
	ip6Regex := regexp.MustCompile(`^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}$`)

	if ip4Regex.MatchString(s) || ip6Regex.MatchString(s) {
		return true
	}

	return false
}
