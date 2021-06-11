package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/token"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type UserService struct {
	userRepo     user.Repository
	userRoleRepo userrole.Repository
	roleRepo     role.Repository
	redis        *redis.Client
}

func NewUserService(userRepo user.Repository, userRoleRepo userrole.Repository, roleRepo role.Repository, redis *redis.Client) *UserService {
	return &UserService{
		userRepo,
		userRoleRepo,
		roleRepo,
		redis,
	}
}

func (svc *UserService) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	var (
		registerReq RegisterReq
		RES         RegisterRes
		err         error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&registerReq); err != nil {
		logger.ErrorLog(" - RegisterHandler json.NewDecoder err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !registerReq.valid() {
		logger.ErrorLog(" - RegisterHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLog(" - RegisterHandler bcrypt.GenerateFromPassword err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.Id, err = svc.userRepo.Create(user.User{
		Username: registerReq.Username,
		Name:     registerReq.Name,
		Password: string(hashPasswordByte),
	}); err != nil {
		logger.ErrorLog(" - RegisterHandler svc.userRepo.Create err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}

func (svc *UserService) LoginHandler(w http.ResponseWriter, req *http.Request) {
	var (
		loginReq LoginReq
		RES      LoginRes

		user      user.User
		userRoles []userrole.UserRole
		roles     []string
		td        token.TokenDetails
		err       error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		logger.ErrorLog(" - LoginHandler json.NewDecoder err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if user, err = svc.userRepo.GetByUsername(loginReq.Username); err != nil {
		logger.ErrorLog(" - LoginHandler svc.userRepo.GetByUsername err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if userRoles, err = svc.userRoleRepo.GetByUser(user.Id); err != nil {
		logger.ErrorLog(" - LoginHandler svc.userRoleRepo.GetByUser err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	for _, ur := range userRoles {
		var role role.Role

		if role, err = svc.roleRepo.GetByID(ur.RoleId); err != nil {
			logger.ErrorLog(" - LoginHandler svc.userRoleRepo.GetByUser err: " + err.Error())
			RES.Error = err.Error()
			return
		}

		roles = append(roles, role.Name)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		logger.ErrorLog(" - LoginHandler bcrypt.CompareHashAndPassword err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if td, err = token.CreateToken(user.Id, roles); err != nil {
		logger.ErrorLog(" - LoginHandler CreateToken err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = svc.saveAuth(user.Id, td); err != nil {
		logger.ErrorLog(" - LoginHandler svc.saveAuth err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	RES.AccessToken = td.AccessToken
	RES.RefreshToken = td.RefreshToken
}

func (svc *UserService) saveAuth(userId int, td token.TokenDetails) (err error) {
	atExpires := time.Unix(td.AtExpires, 0)
	// rtExpires := time.Unix(td.RtExpires, 0)
	now := time.Now()
	ctx := context.Background()

	if err = svc.redis.Set(ctx, token.AccessTokenCachePrefix+td.AccessUuid, userId, atExpires.Sub(now)).Err(); err != nil {
		return err
	}

	// Refresh Token is not managed yet
	// if err = svc.redis.Set(ctx, token.RefreshTokenCachePrefix+td.RefreshUuid, userId, rtExpires.Sub(now)).Err(); err != nil {
	// 	return err
	// }

	return nil
}

func (svc *UserService) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES StandardRes

		ad     token.AccessDetails
		valid  bool
		claims jwt.MapClaims
		err    error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if claims, valid, err = token.ParseFromRequest(req); err != nil {
		logger.ErrorLog(" - LoginHandler token.ParseFromRequest err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !valid {
		logger.ErrorLog(" - LoginHandler invalid token")
		RES.Error = err.Error()
		return
	}

	if ad, err = token.GetAccessDetails(claims); err != nil {
		logger.ErrorLog(" - LoginHandler token.GetAccessDetails err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = svc.redis.Del(req.Context(), token.AccessTokenCachePrefix+ad.AccessUuid).Err(); err != nil {
		logger.ErrorLog(" - LoginHandler svc.redis.Del err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
