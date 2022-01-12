package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func MapJsonToStruct(input interface{}, output interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{TagName: "json", Result: output}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return errs.AddTrace(errors.New("mapstructure.NewDecoder(decoderConfig) err :" + err.Error()))
	}

	derr := decoder.Decode(input)
	if derr != nil {
		return errs.AddTrace(errors.New("decoder.Decode(input) err: " + derr.Error()))
	}

	return nil
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func UniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errs.AddTrace(errors.New("ciphertext too short"))
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Microtime() float64 {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return float64(now.Unix()) + micSeconds
}

func GetMinuteDiffFromNow(datetime string) (float64, error) {
	minuteDiff := 0.0

	timeNow := time.Now()

	layout := "2006-01-02 15:04:05"
	time, err := time.Parse(layout, datetime)
	if err != nil {
		return minuteDiff, err
	}

	minuteDiff = timeNow.Sub(time).Minutes()
	return minuteDiff, nil
}
