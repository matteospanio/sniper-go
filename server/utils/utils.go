/*
 * sniper-go
 *
 * A simple web interface for sniper
 * author: matteospanio
 */

package utils

import (
	"regexp"
	"strings"
)

func IsIP(s string) bool {
	// Regex for IPv4
	ip4Regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	// Regex for IPv6
	ip6Regex := regexp.MustCompile(`^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}$`)

	if ip4Regex.MatchString(s) || ip6Regex.MatchString(s) {
		return true
	}

	return false
}

func IsEmpty(s string) bool {
	trimmed := strings.TrimSpace(s)
	return len(trimmed) == 0
}
