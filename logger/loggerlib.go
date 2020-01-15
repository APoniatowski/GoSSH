package loggerlib

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// OutputLogger All INFO/WARNING messages will be written in this function
func OutputLogger(servername string, output []byte) {
	currentDate := time.Now()
	dateFormatted := currentDate.Format("2006-01-02")
	path, _ := filepath.Abs("./logs/output/")
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil || os.IsExist(err) {
		okFile, err := os.OpenFile(path+"/"+dateFormatted+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer okFile.Close()
		logger := log.New(okFile, "[INFO: Succeeded] ", log.LstdFlags)
		logger.Print(servername + ": " + string(output))
	} else {
		log.Println(err)
	}
}

// ErrorLogger All ERROR/FATAL messages will be written in this function
func ErrorLogger(servername string, output []byte) {
	currentDate := time.Now()
	dateFormatted := currentDate.Format("2006-01-02")
	path, _ := filepath.Abs("./logs/errors/")
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil || os.IsExist(err) {
		errFile, err := os.OpenFile(path+"/"+dateFormatted+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer errFile.Close()
		logger := log.New(errFile, "[INFO: Failed] ", log.LstdFlags)
		logger.Print(servername + ": " + string(output)
	} else {
		log.Println(err)
	}
}
