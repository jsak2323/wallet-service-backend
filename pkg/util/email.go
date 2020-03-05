package util

import (
    // "fmt"
    "encoding/hex"
    "net/smtp"
    "strings"

    "github.com/btcid/wallet-services-backend/cmd/config"
)

var (
    from = config.CONF.MailUser
    host = config.CONF.MailHost
    port = config.CONF.MailPort
)

func SendEmail(subject string, message string, recipients []string) (bool, error) {
    encryptedPassBytes, _ := hex.DecodeString(config.CONF.MailEncryptedPass)

    decryptedPass, err := Decrypt(encryptedPassBytes, []byte(config.CONF.MailEncryptionKey))
    if err != nil { return false, err }

    auth := smtp.PlainAuth("", from, string(decryptedPass), host)

    contents := []byte("To: "+strings.Join(recipients, ";")+"\r\n" +
        "Subject: "+subject+"\r\n" +
        "\r\n" +
        message+"\r\n")

    err = smtp.SendMail(host+":"+port, auth, from, recipients, contents)
    if err != nil { return false, err  }

    return true, nil
}