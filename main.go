// MySQL Performance Monitor(For open-falcon)
// Write by coraldane<coraldane@163.com>
package main

import (
	"flag"
	"fmt"
	"github.com/coraldane/mymon/cron"
	"github.com/coraldane/mymon/db"
	"github.com/coraldane/mymon/g"
	log "github.com/cihub/seelog"
	"os"
)

func main() {
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("log.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if err := g.ParseConfig(*cfg); err != nil {
		log.Error(err)
		log.Flush()
		os.Exit(0)
	}

	go db.InitDatabase()
        go cron.Heartbeat()
	select {}
}
