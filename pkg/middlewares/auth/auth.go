package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/token"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type authMiddleware struct {
	roleRepo           role.Repository
	permissionRepo     permission.Repository
	rolePermissionRepo rolepermission.Repository
	redis              *redis.Client
}

func NewAuthMiddleware(
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	rolePermissionRepo rolepermission.Repository,
	redis *redis.Client,
) *authMiddleware {
	return &authMiddleware{roleRepo, permissionRepo, rolePermissionRepo, redis}
}

var skippedRouteNames = map[string]bool{
	"login":    true,
	"register": true,
}

func skipRoute(name string) bool {
	if _, ok := skippedRouteNames[name]; ok {
		return true
	}

	return false
}

func (am *authMiddleware) Authenticate(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// to allow access from localhost
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if skipRoute(mux.CurrentRoute(req).GetName()) {
			hf.ServeHTTP(w, req)
			return
		}

		// AUTHORIZE IP
		if isIpAuthorized := isIpAuthorized(req); isIpAuthorized { // if ip is authorized, continue
			logger.InfoLog(" - AUTH -- IP is authorized.", req)

		} else { // if not authorized, send notification email and stop request
			logger.InfoLog(" - AUTH -- IP is unauthorized.", req)
			handleUnauthorizedIp(req)
			return
		}

		claims, valid, err := token.ParseFromRequest(req)
		if valid {
			logger.InfoLog(" - AUTH -- User is authenticated req: ", req)
		} else if err != nil {
			logger.InfoLog(" - AUTH -- User is unauthenticated err: "+err.Error()+" req: ", req)
			return
		} else {
			logger.InfoLog(" - AUTH -- User is unauthenticated req: ", req)
			return
		}

		hf.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), "claims", claims)))
	})
}

func (am *authMiddleware) Authorize(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			authorized bool
			routeRoles []string
			err        error

			routeName    = mux.CurrentRoute(req).GetName()
			mapClaims, _ = req.Context().Value("claims").(jwt.MapClaims)
			ad           token.AccessDetails
		)

		if skipRoute(mux.CurrentRoute(req).GetName()) {
			hf.ServeHTTP(w, req)
			return
		}

		handleResponse := func() {
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
		defer handleResponse()

		if ad, err = token.GetAccessDetails(mapClaims); err != nil {
			logger.ErrorLog("- AUTH -- User -- extractToken err: " + err.Error())
			return
		}

		fmt.Println("TOKEN", token.AccessTokenCachePrefix+ad.AccessUuid)

		if _, err = am.redis.Get(req.Context(), token.AccessTokenCachePrefix+ad.AccessUuid).Result(); err != nil {
			logger.ErrorLog("- AUTH -- User -- am.redis.Get err: " + err.Error())
			return
		}

		if routeRoles, err = am.getRouteRoles(routeName); err != nil {
			logger.ErrorLog("- AUTH -- User -- am.getRouteRoles err: " + err.Error())
			return
		}

		for _, ur := range ad.Roles {
			if authorized, _ = util.InArray(ur, routeRoles); !authorized {
				continue
			}
		}

		if !authorized {
			logger.InfoLog(" - AUTH -- Authorize is unauthorized.", req)
			return
		}

		logger.InfoLog(" - AUTH -- Authorize is authorized.", req)

		hf.ServeHTTP(w, req)
	})
}

func (am *authMiddleware) getRouteRoles(routeName string) (routeRoles []string, err error) {

	permission, err := am.permissionRepo.GetByName(routeName)
	if err != nil {
		return []string{}, err
	}

	rolePermissions, err := am.rolePermissionRepo.GetByPermission(permission.Id)
	if err != nil {
		return []string{}, err
	}

	for _, rp := range rolePermissions {
		routeRole, err := am.roleRepo.GetByID(rp.RoleId)
		if err != nil {
			return []string{}, err
		}

		routeRoles = append(routeRoles, routeRole.Name)
	}

	return routeRoles, nil
}
