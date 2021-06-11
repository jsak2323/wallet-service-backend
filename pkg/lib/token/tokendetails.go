package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

const AccessTokenCachePrefix = "aToken-"
const RefreshTokenCachePrefix = "rToken-"

func CreateToken(userId int, roles []string) (td TokenDetails, err error) {
	td = TokenDetails{
		AccessUuid:  uuid.NewString(),
		AtExpires:   time.Now().Add(time.Minute * 15).Unix(),
		RefreshUuid: uuid.NewString(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 2).Unix(),
	}

	fmt.Println("UUID", td.AccessUuid)

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user_id":     userId,
		"roles":       strings.Join(roles, ","),
		"access_uuid": td.AccessUuid,
		"exp":         td.AtExpires,
	})

	td.AccessToken, err = at.SignedString([]byte(config.CONF.JWTSecret))
	if err != nil {
		return TokenDetails{}, err
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user_id":      userId,
		"roles":        strings.Join(roles, ","),
		"refresh_uuid": td.RefreshUuid,
		"exp":          td.RtExpires,
	})

	td.RefreshToken, err = rt.SignedString([]byte(config.CONF.JWTRefreshSecret))
	if err != nil {
		return TokenDetails{}, err
	}
	return td, nil
}

func getTokenFromHeader(req *http.Request) (token string) {
	var (
		bearerStr    string
		bearerStrArr []string
	)

	bearerStr = req.Header.Get("Authorization")
	bearerStrArr = strings.Split(bearerStr, " ")

	if len(bearerStrArr) != 2 {
		return ""
	}

	return bearerStrArr[1]
}

func ParseFromRequest(req *http.Request) (claims jwt.MapClaims, valid bool, err error) {
	var (
		tokenStr string
		token    *jwt.Token
	)

	if tokenStr = getTokenFromHeader(req); tokenStr == "" {
		return nil, false, errors.New("malformed authorization header")
	}

	if token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.CONF.JWTSecret), nil
	}); err != nil {
		return nil, false, errors.New("jwt.Parse err: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false, errors.New("malformed token claims")
	}

	if token == nil {
		return nil, false, errors.New("nil token")
	}

	return claims, true, nil
}
