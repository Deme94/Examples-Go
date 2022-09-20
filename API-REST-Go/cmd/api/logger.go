package main

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
)

func logger() *log.Logger {
	os.Mkdir("./log", os.ModePerm)
	filename := "./log/logs.txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(file, "", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	logger.SetOutput(&lumberjack.Logger{
		Filename: filename,
		MaxSize:  1, // megabytes after which new file is created
		//MaxBackups: 3,  // number of backups
		//MaxAge:     28, //days
	})

	return logger
}
