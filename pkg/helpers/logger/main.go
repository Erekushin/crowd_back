package logger

import (
	"os"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func createFile(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func changer(fileName string) error {
	createFile(fileName)
	if file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		log.SetOutput(file)
	} else {
		return err
	}
	return nil
}

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	os.MkdirAll("logs", os.ModePerm)

	filePre := "logs/log_"
	fileName := filePre + time.Now().Format("2006_01_02") + ".log"
	changer(fileName)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		fileName = filePre + time.Now().Format("2006_01_02") + ".log"
		changer(fileName)
	})
	c.Start()
}
