package util

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func RawToCoin(raw string, maxDecimal int) string {
	if raw == "0" {
		return "0"
	}

	raw = strings.TrimLeft(raw, "0")

	decimalString := ""
	if len(raw) <= maxDecimal { // number is less than one
		decimalString = "0."
		for i := 0; i < (maxDecimal - len(raw)); i++ {
			decimalString = decimalString + "0"
		}
		decimalString = decimalString + raw
		decimalString = strings.TrimRight(decimalString, "0")

	} else { // number is greater than one
		numberPart := raw[0:(len(raw) - maxDecimal)]

		decimalPart := raw[(len(raw) - maxDecimal):]
		decimalPart = strings.TrimRight(decimalPart, "0")

		decimalString = numberPart
		if decimalPart != "" {
			decimalString = decimalString + "." + decimalPart
		}
	}

	_, ok := new(big.Float).SetString(decimalString) // check if number is valid
	if !ok {
		return "0"
	}

	return decimalString
}

func CoinToRaw(decimal string, maxDecimal int) string {
	split := strings.Split(decimal, ".")
	decimal = strings.TrimLeft(decimal, "0")

	if len(split) > 1 {
		decimal = strings.TrimRight(decimal, "0")
		if len(split[1]) > maxDecimal { // reduce decimal count when it is greater than max decimal
			trimmedDecimalPart := split[1][:len(split[1])-(len(split[1])-maxDecimal)]
			decimal = split[0] + "." + trimmedDecimalPart

			return CoinToRaw(decimal, maxDecimal)
		}
	}

	rawString := ""
	if string(decimal[0]) == "." { // number is less than one
		rawString = decimal[1:]
		for i := 0; i < ((maxDecimal + 1) - len(decimal)); i++ {
			rawString = rawString + "0"
		}
	} else { // number is greater than one
		decimalPart := ""
		if len(split) > 1 {
			decimalPart = strings.TrimRight(split[1], "0")
		}
		rawString = strings.ReplaceAll(decimal, ".", "")
		for i := 0; i < (maxDecimal - len(decimalPart)); i++ {
			rawString = rawString + "0"
		}
	}

	rawString = strings.TrimLeft(rawString, "0")
	_, ok := new(big.Float).SetString(rawString) // check if number is valid
	if !ok {
		return "0"
	}

	return rawString
}

func AddCurrency(a, b string) (res string, err error) {
	return calculateCurrency("add", a, b)
}

func SubCurrency(a, b string) (res string, err error) {
	return calculateCurrency("sub", a, b)
}

func CmpCurrency(a, b string) (int, error) {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}

	aBig, ok := big.NewFloat(0).SetString(a)
	if !ok {
		return -2, errs.AddTrace(errors.New("fail big.SetString(" + a + ")"))
	}

	bBig, ok := big.NewFloat(0).SetString(b)
	if !ok {
		return -2, errs.AddTrace(errors.New("fail big.SetString(" + b + ")"))
	}

	return aBig.Cmp(bBig), nil
}

func PercentCurrency(a, b string) (string, error) {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "1"
	}

	aBig, ok := big.NewFloat(0).SetString(a)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + a + ")"))
	}

	bBig, ok := big.NewFloat(0).SetString(b)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + b + ")"))
	}

	if bBig.Cmp(big.NewFloat(0)) <= 0 {
		return "0", errs.AddTrace(errors.New("division by zero"))
	}

	aBig.Mul(aBig.Quo(aBig, bBig), big.NewFloat(100))

	return aBig.Text('f', 2), nil
}

func FormatCurrency(value string) (result string) {
	parts := strings.Split(value, ".")
	intPart := parts[0]
	decPart := ""

	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	for i := len(intPart); i > 0; i-- {
		j := len(intPart) - i
		if (j+1)%3 == 0 && j < len(intPart)-1 {
			result = "," + string(intPart[i-1]) + result
		} else {
			result = string(intPart[i-1]) + result
		}
	}

	return result + decPart
}

func calculateCurrency(method, a, b string) (res string, err error) {
	const maxDecimal = 8

	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}

	valueA := CoinToRaw(a, lenDec(a))
	valueB := CoinToRaw(b, lenDec(b))

	if lenDec(a) != lenDec(b) {
		valueA, valueB, err = matchDecimal(valueA, valueB, lenDec(a), lenDec(b))
		return "0", errs.AddTrace(err)
	}

	aBig, ok := new(big.Int).SetString(valueA, 10)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + a + ")"))
	}

	bBig, ok := new(big.Int).SetString(valueA, 10)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + b + ")"))
	}

	calculate := ""
	switch types := method; types {
	case "sub":
		calculate = aBig.Sub(aBig, bBig).String()
	case "add":
		calculate = aBig.Add(aBig, bBig).String()
		_ = aBig.Int64()
	default:
		return "0", errs.AddTrace(errors.New("wrong method of calculate currency"))
	}

	if lenDec(a) > lenDec(b) {
		res = RawToCoin(calculate, lenDec(a))
	} else {
		res = RawToCoin(calculate, lenDec(b))
	}

	if lenDec(a) > maxDecimal || lenDec(b) > maxDecimal {
		split := strings.Split(res, ".")
		afTrunc, err := truncateText(split[1], maxDecimal)
		if err != nil {
			return "0", errs.AddTrace(err)
		}
		res = split[0] + "." + afTrunc
	}

	return res, nil
}

func lenDec(req string) int {
	var lenDec int

	parts := strings.Split(req, ".")
	if len(parts) > 1 {
		resA := strings.TrimRight(parts[1], "0")
		lenDec = len(resA)
	}
	return lenDec
}

func truncateText(s string, max int) (string, error) {
	if max > len(s) {
		return "", errs.AddTrace(errors.New("slice bounds out of range [:" + fmt.Sprintf("%v", max) + "] with length " + fmt.Sprintf("%v", len(s))))
	}
	return s[:max], nil
}

func matchDecimal(a, b string, lenDecA, lenDecB int) (string, string, error) {

	if lenDecA < 0 || lenDecB < 0 {
		return a, b, errs.AddTrace(errors.New("length decimal can't be minus, len decimal A: " + fmt.Sprintf("%v", lenDecA) + "len decimal B: " + fmt.Sprintf("%v", lenDecB)))
	}

	if lenDecA > lenDecB {
		for i := int(0); i < (lenDecA - lenDecB); i++ {
			b = b + "0"
		}
	}

	if lenDecA < lenDecB {
		for i := int(0); i < (lenDecB - lenDecA); i++ {
			a = a + "0"
		}
	}

	return a, b, nil

}
