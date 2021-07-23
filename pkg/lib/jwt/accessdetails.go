package jwt

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type AccessDetails struct {
	AccessUuid  string
	UserId      int
	Roles       []string
	Permissions []string
}

func getAccessDetails(claims jwt.MapClaims) (ad AccessDetails, err error) {
	for key, value := range claims {
		switch key {
		case "roles":
			rolesValue, ok := value.(string)
			if !ok {
				err = errors.New("malformed roles token claims")
				break
			}

			ad.Roles = strings.Split(rolesValue, ",")
		case "user_id":
			userIdValue, ok := value.(float64)
			if !ok {
				err = errors.New("malformed user_id token claims")
				break
			}

			ad.UserId = int(userIdValue)
		case "access_uuid":
			accessUuid, ok := value.(string)
			if !ok {
				err = errors.New("malformed uuid token claims")
				break
			}

			ad.AccessUuid = accessUuid
		default:
			continue
		}
	}

	return ad, err
}
