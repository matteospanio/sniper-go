package sniper

import (
	"strings"

	"github.com/matteospanio/sniper-go/utils"
)

// Target name and IP address
type Target struct {
	Name []string `json:"name"`
	IP   []string `json:"ip"`
}

/* Read the vulnerability report from sn1per
 * Return the content of the domains/targets-all-sorted.txt file
 * and the content of the ips/ips-all-sorted.txt file
 */
func getTarget(host string) Target {
	result := Target{}

	// find the host name
	targets, _ := readSniperFile(host, "domains/targets-all-sorted.txt")
	for _, target := range strings.Split(targets, "\n") {
		if !utils.IsEmpty(target) {
			result.Name = append(result.Name, target)
		}
	}

	// find the host IP
	ips, _ := readSniperFile(host, "ips/ips-all-sorted.txt")
	for _, ip := range strings.Split(ips, "\n") {
		if !utils.IsEmpty(ip) {
			result.IP = append(result.IP, ip)
		}
	}

	return result
}
