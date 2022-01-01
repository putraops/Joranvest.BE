package commons

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Logger() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	var filename string = "logs/logger - " + (time.Now()).Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	log.SetFormatter(Formatter)
	if err != nil {
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}
}
