package cron

import (
	"github.com/coraldane/mymon/g"
	"time"
)

func Heartbeat() {
	t := time.NewTicker(time.Duration(g.Config().Interval) * time.Second).C
	for {
		for _, server := range g.Config().DBServerList {
			go FetchData(server)
		}
		<-t
	}
}
