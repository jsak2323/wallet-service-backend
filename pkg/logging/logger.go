package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	ctxLib "github.com/btcid/wallet-services-backend-go/pkg/lib/context"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
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
		DisableColors:   true,
		TimestampFormat: "2 Jan 2006 15:04:05",
	})
}

// context
func Log(msg string, ctx context.Context) {
	updateTime()
	setupLogger()

	log.Info(msg)
	// logField := logrus.Fields{}
	// if reqId, ok := ctx.Value(ctxLib.RequestIdKey).(string); ok {
	// 	logField["RequestId"] = reqId
	// }

	// log.WithFields(logField).Error(msg)
}

func InfoLog(msg string, req *http.Request) {
	updateTime()
	setupLogger()

	body, err := preprocessingGetRequestBody(req)
	if err != nil {
		log.Println("error getRequestBody:", err)
		return
	}

	logField := logrus.Fields{
		"Method":     req.Method,
		"RemoteAddr": req.RemoteAddr,
		"body":       body,
	}

	if ad, valid := ctxLib.ValidateAccessDetailsContext(req.Context()); valid {
		logField["UserId"] = ad.GetUserId()
	}

	if reqId, ok := req.Context().Value(ctxLib.RequestIdKey).(string); ok {
		logField["RequestId"] = reqId
	}

	log.WithFields(logField).Info(msg)
}

func preprocessingGetRequestBody(req *http.Request) (resp io.ReadCloser, err error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return resp, errs.AddTrace(errors.New("failed read all body"))
	}

	// make real copy
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// make copy response for log
	respPreprocess := make(map[string]interface{})
	respLog := make(map[string]interface{})

	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&respPreprocess)
	if err != nil && err.Error() != "EOF" {
		return resp, errs.AddTrace(err)
	}

	// censored key "password" for log
	for key, value := range respPreprocess {
		if key == "password" {
			value = "-"
		}
		respLog[key] = value
	}

	// reconstruction request body for log
	jsonString, err := json.Marshal(respLog)
	if err != nil {
		return resp, errs.AddTrace(errors.New("failed marshal request body"))
	}
	resp = ioutil.NopCloser(bytes.NewBuffer(jsonString))

	return resp, nil
}

// context
func ErrorLog(msg string, ctx context.Context) {
	updateTime()
	setupLogger()

	logField := logrus.Fields{}
	if reqId, ok := ctx.Value(ctxLib.RequestIdKey).(string); ok {
		logField["RequestId"] = reqId
	}

	log.WithFields(logField).Error(msg)

	go sendErrorNotificationEmail(msg)
}

func sendErrorNotificationEmail(msg string) {
	config.ErrorMailCount += 1

	if config.ErrorMailCount > config.CONF.SessionErrorMailNotifLimit {
		fmt.Println("Error Notification Mail Limit is hit for this session. skipping ...")
		return

	} else {
		const emailSubjectPrefix string = "[WALLETSERVICE]"
		subject := emailSubjectPrefix + " Application Error"
		message := "An error was encountered with following detail: " +
			"\n Error: " + msg

		recipients := config.CONF.NotificationEmails
		util.SendEmail(subject, message, recipients)
	}
}
