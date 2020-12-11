package auth

import (
    "strings"
    "strconv"
    "net/http"

    "github.com/btcid/wallet-services-backend-go/cmd/config"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

const emailSubjectPrefix string = "[WALLETSERVICE]"

func handleUnauthorizedIp(req *http.Request) {
    ip := strings.Split(req.RemoteAddr, ":")[0]

    logger.Log(" - AUTH -- Sending notification email ...")
    
    subject := emailSubjectPrefix+" Request from suspicious IP address: "+ip
    message := "A request from suspicious IP address was recorded with following detail: "+
    "\n IP Address: "+ip+
    "\n URL: "+req.URL.String()

    recipients := config.CONF.NotificationEmails

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - AUTH -- Is unauthorized ip notification email sent: "+strconv.FormatBool(isEmailSent))
}


