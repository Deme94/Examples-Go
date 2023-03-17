package logger

import (
	"API-REST/services/conf"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
)

var Logger *log.Logger

func Setup() error {
	// Read conf
	dir := conf.Conf.GetString("logDir")
	filename := conf.Conf.GetString("logFileName")
	ext := conf.Conf.GetString("logFileExt")

	// Default logger if not specified is stdout
	if dir == "" || filename == "" {
		Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
		return nil
	}
	// Create directory if not exist
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !strings.Contains(fmt.Sprint(err), "file exists") {
		return err
	}
	// Open file
	path := dir + "/" + filename + "." + ext
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	// Create logger
	Logger = log.New(file, "", log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Logger.SetOutput(&lumberjack.Logger{
		Filename: path,
		MaxSize:  1, // megabytes after which new file is created
		//MaxBackups: 3,  // number of backups
		//MaxAge:     28, //days
	})

	return nil
}
