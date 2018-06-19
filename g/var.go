package g

import (
	log "github.com/cihub/seelog"
	"os"
)

func Hostname(hostname string) string {
	if "" != hostname {
		return hostname
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error("ERROR: os.Hostname() fail", err)
	}
	return hostname
}
