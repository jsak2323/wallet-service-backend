package logging

import (
    "fmt"
    "os"
    "time"
    "github.com/sirupsen/logrus"
    "net/http"
)

var logsDir = "/logs/"
var log = logrus.New()

var currentTime time.Time
var currentDate string

func updateTime() {
    currentTime = time.Now()
    currentDate = currentTime.Format("01-02-2006")
}

func setupLogger() {
    pwd, _ := os.Getwd()
    file, err := os.OpenFile(pwd+logsDir+"app-"+currentDate+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
        log.Out = file
    } else {
        log.Info("Failed to log to file, using default stderr")
    }
    log.SetFormatter(&logrus.TextFormatter{
        DisableColors       : true,
        TimestampFormat     : "2 Jan 2006 15:04:05",
    })
}

func Log(msg string){
    updateTime()
    setupLogger()

    fmt.Println(msg)
    log.Info(msg)
}

func InfoLog(msg string, req *http.Request) {
    updateTime()
    setupLogger()

    fmt.Println(msg)

    log.WithFields(logrus.Fields{
        "Method"        : req.Method,
        "RemoteAddr"    : req.RemoteAddr,
    }).Info(msg)
}

func ErrorLog(msg string) {
    updateTime()
    setupLogger()
    
    fmt.Println(msg)
    log.Error(msg)

    // todo: add send email for error notification
}
