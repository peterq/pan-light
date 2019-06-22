package pc_api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/conf"
	"github.com/peterq/pan-light/server/dao"
	"github.com/peterq/pan-light/server/pan-viper"
	"github.com/peterq/pan-light/server/pc-api/middleware"
	"strings"
	"time"
)

type gson = map[string]interface{}

func handleLoginToken(ctx context.Context, param artisan.JsonMap) (result interface{}, err error) {
	uk := param.Get("uk").String()
	filename := artisan.Md5bin([]byte(fmt.Sprint(time.Now().UnixNano())))
	filename = filename[:8]
	token, e := middleware.NewJwtToken(time.Minute*5, map[string]interface{}{
		"need-login": uk,
		"filename":   filename,
	})
	result = gson{
		"token":    token,
		"filename": filename,
	}
	if e != nil {
		err = artisan.NewError("generate jwt error", -1, e)
	}
	return
}

func handleLogin(ctx context.Context, param artisan.JsonMap) (result interface{}, err error) {
	link := param.Get("link").String()
	secret := param.Get("secret").String()
	token := param.Get("token").String()

	claim, e := middleware.ParseToken(token)
	if e != nil {
		err = artisan.NewError("token invalid", -1, e)
		return
	}
	uk := claim["need-login"].(string)
	fn := claim["filename"].(string)

	ukShare, filename, share, err := pan_viper.GetVip().LoadShareFilenameAndUk(link, secret)
	if err != nil {
		err = artisan.NewError("读取分享列表错误", -1, err)
		return
	}
	if ukShare != uk {
		err = artisan.NewError("uk not match", -1, err)
		return
	}
	if !strings.Contains(filename, fn) {
		err = artisan.NewError("share file invalid", -1, err)
		return
	}
	dao.UserDao.UpInsert(dao.UserModel{
		Uk:           uk,
		MarkUsername: share["linkusername"].(string),
		Avatar:       share["photo"].(string),
		IsVip:        share["is_master_vip"].(float64) == 1,
		IsSuperVip:   share["is_master_svip"].(float64) == 1,
	})
	result, err = middleware.NewJwtToken(time.Hour*24*30, gson{
		"type": conf.JwtPcLogin,
		"uk":   uk,
	})
	return
}

func handleRefreshToken(ctx context.Context, param artisan.JsonMap) (result interface{}, err error) {
	token := middleware.PcJwtHandler.Get(ctx)
	claims := token.Claims.(jwt.MapClaims)
	if claims.VerifyExpiresAt(time.Now().Add(time.Hour*24*5).Unix(), true) {
		result, err = middleware.NewJwtToken(time.Hour*24*30, gson{
			"type": conf.JwtPcLogin,
			"uk":   middleware.ContextLoginInfo(ctx).Uk(),
		})
		return
	}
	result = token.Raw
	return
}

func handleFeedBack(ctx context.Context, param artisan.JsonMap) (result interface{}, err error) {
	content := param.Get("content").String()
	err = dao.FeedbackDao.Insert(dao.FeedbackModel{
		Uk:      middleware.ContextLoginInfo(ctx).Uk(),
		Content: content,
	})
	if err != nil {
		err = artisan.NewError("database error", -1, nil)
	}
	return
}
