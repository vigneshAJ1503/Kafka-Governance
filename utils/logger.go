package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New(os.Stdout, "[kafka-gov] ", log.LstdFlags|log.Lshortfile)
	Logger.Println("Logger initialized")
}
