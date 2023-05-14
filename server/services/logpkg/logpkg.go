package logpkg

import (
	"log"
	"os"
)

var logFile *os.File

func ErrorLog(message string) {
	log.Println("ERROR: " + message)
}

func InfoLog(message string) {
	log.Println("INFO: " + message)
}

func DebugLog(message string) {
	log.Println("DEBUG: " + message)
}

func InitLogger(filename string) {
	var err error
	logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		ErrorLog("Error opening file for logging: " + err.Error())
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logFile)
	}
}

func CloseLogger() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			ErrorLog("Error closing logpkg file: " + err.Error())
		}
	}
}
