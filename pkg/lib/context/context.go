package context

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"
)

type AccessDetails struct {
	AccessUuid  string
	UserId      int
	Roles       []string
	Permissions []string
}

func (ad *AccessDetails) GetAccessUuid() string {
	return ad.AccessUuid
}

func (ad *AccessDetails) GetUserId() int {
	return ad.UserId
}

func (ad *AccessDetails) GetRoles() []string {
	return ad.Roles
}

func (ad *AccessDetails) GetPermissions() []string {
	return ad.Permissions
}

func ValidateAccessDetailsContext(ctx context.Context) (*AccessDetails, bool) {
	if ctx.Value("access_details") != nil {
		ad := &AccessDetails{
			AccessUuid:  ctx.Value("access_details").(jwt.AccessDetails).AccessUuid,
			UserId:      ctx.Value("access_details").(jwt.AccessDetails).UserId,
			Roles:       ctx.Value("access_details").(jwt.AccessDetails).Roles,
			Permissions: ctx.Value("access_details").(jwt.AccessDetails).Permissions,
		}
		return ad, true
	}
	return nil, false
}
