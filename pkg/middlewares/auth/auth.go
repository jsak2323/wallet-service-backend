package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	ctxLib "github.com/btcid/wallet-services-backend-go/pkg/lib/context"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

type authMiddleware struct {
	roleRepo           role.Repository
	permissionRepo     permission.Repository
	rolePermissionRepo rolepermission.Repository
}

func NewAuthMiddleware(
	roleRepo role.Repository,
	permissionRepo permission.Repository,
	rolePermissionRepo rolepermission.Repository,
) *authMiddleware {
	return &authMiddleware{roleRepo, permissionRepo, rolePermissionRepo}
}

var skippedRouteNames = map[string]bool{
	"login":              true,
	"cronhealthcheck":    true,
	"fireblockscallback": true,
}

func skipRoute(name string) bool {
	if _, ok := skippedRouteNames[name]; ok {
		return true
	}

	return false
}

func skipHost(host string) bool {
	if strings.Split(host, ":")[0] == "localhost" {
		return true
	}

	return false
}

func (am *authMiddleware) Authenticate(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// to allow access from localhost
		w.Header().Set("Access-Control-Allow-Origin", "*")

		req = req.WithContext(ctxLib.SetRqId(req.Context()))
		if skipRoute(mux.CurrentRoute(req).GetName()) || skipHost(req.Host) {
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

		accessDetails, valid, err := jwt.ParseFromRequest(req)
		if valid {
			logger.InfoLog(" - AUTH -- User is authenticated req: ", req)
		} else if err != nil {
			logger.InfoLog(" - AUTH -- User is unauthenticated err: "+err.Error()+" req: ", req)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			logger.InfoLog(" - AUTH -- User is unauthenticated req: ", req)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), "access_details", accessDetails))
		hf.ServeHTTP(w, req)
	})
}

func (am *authMiddleware) Authorize(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			authorized bool
			routeRoles []string
			err        error

			routeName = mux.CurrentRoute(req).GetName()
			ad, _     = req.Context().Value("access_details").(jwt.AccessDetails)
		)

		if skipRoute(mux.CurrentRoute(req).GetName()) || skipHost(req.Host) {
			hf.ServeHTTP(w, req)
			return
		}

		handleResponse := func() {
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(StandardRes{
					Error: "Unauthorized user account for resource",
				})
			}
		}
		defer handleResponse()

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
		routeRole, err := am.roleRepo.GetById(rp.RoleId)
		if err != nil {
			return []string{}, err
		}

		routeRoles = append(routeRoles, routeRole.Name)
	}

	return routeRoles, nil
}
