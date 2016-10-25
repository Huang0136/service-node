package logs

import (
	"log"
	"os"
)

type MyLogger int

var MyDebugLogger *log.Logger = new(log.Logger)
var MyInfoLogger *log.Logger = new(log.Logger)
var MyErrorLogger *log.Logger = new(log.Logger)

func init() {
	MyDebugLogger.SetPrefix("【DEBUG】 ")
	MyDebugLogger.SetOutput(os.Stdout)
	MyDebugLogger.SetFlags(log.LstdFlags)

	MyInfoLogger.SetPrefix("【INFO】 ")
	MyInfoLogger.SetOutput(os.Stdout)
	MyInfoLogger.SetFlags(log.LstdFlags)

	MyErrorLogger.SetPrefix("【ERROR】 ")
	MyErrorLogger.SetOutput(os.Stdout)
	MyErrorLogger.SetFlags(log.LstdFlags)
}

// check err
func (l *MyLogger) CheckError(msg string, err error) {
	if err != nil {
		l.Fatalln(msg, err)
	}
}
