package logs

import (
	"log"
	"os"
)

type MyLog struct {
	log.Logger
}

var MyDebugLog *MyLog
var MyInfoLog *MyLog
var MyErrorLog *MyLog

func init() {
	MyDebugLog = new(MyLog)
	MyDebugLog.SetPrefix("【DEBUG】 ")
	MyDebugLog.SetOutput(os.Stdout)
	MyDebugLog.SetFlags(log.LstdFlags)

	MyInfoLog = new(MyLog)
	MyInfoLog.SetPrefix("【INFO】 ")
	MyInfoLog.SetOutput(os.Stdout)
	MyInfoLog.SetFlags(log.LstdFlags)

	MyErrorLog = new(MyLog)
	MyErrorLog.SetPrefix("【ERROR】 ")
	MyErrorLog.SetOutput(os.Stdout)
	MyErrorLog.SetFlags(log.LstdFlags)
}

// check printlnError
func (l MyLog) CheckPrintlnError(msg string, err error) {
	if err != nil {
		l.Println(msg, err)
	}
}

// check fatallnError
func (l MyLog) CheckFatallnError(msg string, err error) {
	if err != nil {
		l.Fatalln(msg, err)
	}
}

// check paniclnError
func (l *MyLog) CheckPaniclnError(msg string, err error) {
	if err != nil {
		l.Logger.Panicln(msg, err)
	}
}
