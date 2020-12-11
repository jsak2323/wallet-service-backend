package util

import (
    "io"
    "time"
    "errors"
    "reflect"
    "math/big"
    "crypto/aes"
    "crypto/rand"
    "crypto/cipher"

    "github.com/mitchellh/mapstructure"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

func MapJsonToStruct(input interface{}, output interface{}) error {
    decoderConfig := &mapstructure.DecoderConfig{ TagName: "json", Result: output }
    decoder, err := mapstructure.NewDecoder(decoderConfig)
    if err != nil { return errors.New("mapstructure.NewDecoder(decoderConfig) err :"+err.Error()) }

    derr := decoder.Decode(input)
    if derr != nil { return errors.New("decoder.Decode(input) err: "+derr.Error()) }

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

func WeiToEther(wei string) string {
    bigFloat, ok := new(big.Float).SetString(wei)
    if !ok { return "0" }

    divider, ok := new(big.Float).SetString("1000000000000000000")
    if !ok { return "0" }

    eth := new(big.Float).Quo(bigFloat, divider)
    return eth.String()
}

func EtherToWei(ether string) uint64 {
    bigFloat, ok := new(big.Float).SetString(ether)
    if !ok { return uint64(0) }

    multiplier, ok := new(big.Float).SetString("1000000000000000000")
    if !ok { return uint64(0) }

    wei := new(big.Float).Mul(bigFloat, multiplier)
    weiuint64, _ := wei.Uint64()
    return weiuint64
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
        return nil, errors.New("ciphertext too short")
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

func GetRpcConfigByType(SYMBOL string, rpcConfigType string) (rc.RpcConfig, error) {
    for _, rpcConfig := range config.CURR[SYMBOL].RpcConfigs {
        if rpcConfig.Type == rpcConfigType || rpcConfig.Type == "master" {
            return rpcConfig, nil
        }
    }
    return rc.RpcConfig{}, errors.New("RpcConfig not found.")
}


