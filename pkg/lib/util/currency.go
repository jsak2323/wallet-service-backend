package util

import (
	"errors"
	"math"
	"math/big"
	"strings"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func CoinToRaw(value string, decimal int) (string, error) {
	if value == "" {
		value = "0"
	}

	valueFloat, ok := big.NewFloat(0).SetString(value)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + value + ")"))
	}

	decimalFloat := big.NewFloat(math.Pow(10, float64(decimal)))
	valueFloat.Mul(valueFloat, decimalFloat)

	return valueFloat.Text('f', 0), nil
}

func RawToCoin(value string, decimal int) (string, error) {
	if value == "" {
		value = "0"
	}

	valueFloat, ok := big.NewFloat(0).SetString(value)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + value + ")"))
	}

	decimalFloat := big.NewFloat(math.Pow(10, float64(decimal)))
	valueFloat.Quo(valueFloat, decimalFloat)

	return valueFloat.Text('f', 8), nil
}

func AddCoin(a, b string) (string, error) {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}

	aBig, ok := big.NewFloat(0).SetString(a)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + a + ")"))
	}

	bBig, ok := big.NewFloat(0).SetString(b)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + b + ")"))
	}

	return bBig.Add(aBig, bBig).Text('f', 8), nil
}

func AddIdr(a, b string) (string, error) {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}

	aBig, ok := big.NewFloat(0).SetString(a)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + a + ")"))
	}

	bBig, ok := big.NewFloat(0).SetString(b)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + b + ")"))
	}

	return bBig.Add(aBig, bBig).Text('f', 0), nil
}

func SubIdr(a, b string) (string, error) {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}

	aBig, ok := big.NewFloat(0).SetString(a)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + a + ")"))
	}

	bBig, ok := big.NewFloat(0).SetString(b)
	if !ok {
		return "0", errs.AddTrace(errors.New("fail big.SetString(" + b + ")"))
	}

	return aBig.Sub(aBig, bBig).Text('f', 0), nil
}

func CmpBig(a, b string) (int, error) {
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

func PercentBig(a, b string) (string, error) {
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
