/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package main

import (
	"io"
	"os"
	"regexp"
)

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

func isEmpty(s string) bool {
	re, err := regexp.Compile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*$`)
	if err != nil {
		panic(err)
	}

	return !re.MatchString(s)
}

func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
