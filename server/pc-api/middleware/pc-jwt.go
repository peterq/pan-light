package middleware

import (
	"errors"
	"runtime/debug"
	"time"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/peterq/pan-light/server/conf"
	"github.com/peterq/pan-light/server/dao"
)

type PcLoginInfo struct {
	uk   string
	user *dao.UserModel
}

func (p *PcLoginInfo) User() *dao.UserModel {
	if p.user == nil {
		p.user, _ = dao.UserDao.GetByUk(p.uk)
	}
	return p.user
}

func (p *PcLoginInfo) Uk() string {
	return p.uk
}

var PcJwtHandler = jwtmiddleware.New(jwtmiddleware.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.AppSecret), nil
	},

	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(ctx iris.Context, err error) {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]interface{}{
			"success": false,
			"message": "jwt token check error: " + err.Error(),
			"code":    iris.StatusUnauthorized,
		})
		debug.PrintStack()
	},
})

// 校验pc端jwt
func PcJwtAuth(ctx iris.Context) {
	if err := PcJwtHandler.CheckJWT(ctx); err != nil {
		ctx.StopExecution()
		return
	}

	token := PcJwtHandler.Get(ctx)
	claim := token.Claims.(jwt.MapClaims)

	if claim["type"].(string) != conf.JwtPcLogin {
		ctx.StopExecution()
		return
	}

	loginInfo := &PcLoginInfo{
		uk: claim["uk"].(string),
	}
	ctx.Values().Set(conf.CtxPcLogin, loginInfo)

	ctx.Next()
}

func NewJwtToken(expireIn time.Duration, info map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(expireIn).Unix()
	claims["iat"] = time.Now().Unix()
	if info != nil {
		for key, value := range info {
			claims[key] = value
		}
	}
	token.Claims = claims

	return token.SignedString([]byte(conf.Conf.AppSecret))
}

func ParseToken(token string) (claim jwt.MapClaims, err error) {
	parsedToken, err := jwt.Parse(token, PcJwtHandler.Config.ValidationKeyGetter)
	if err != nil {
		return
	}
	if !parsedToken.Valid {
		err = errors.New("token invalid")
		return
	}
	claim = parsedToken.Claims.(jwt.MapClaims)
	if expired := claim.VerifyExpiresAt(time.Now().Unix(), true); !expired {
		err = errors.New("token is expired")
		return
	}
	return
}

func ContextLoginInfo(ctx iris.Context) *PcLoginInfo {
	return ctx.Values().Get(conf.CtxPcLogin).(*PcLoginInfo)
}
