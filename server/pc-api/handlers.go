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
	"gopkg.in/mgo.v2"
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
		err = artisan.NewError("database error", -1, err)
	}
	return
}

func handleShareToSquare(ctx context.Context, param artisan.JsonMap) (result interface{}, err error) {
	md5 := param.Get("md5").String()
	sliceMd5 := param.Get("sliceMd5").String()
	title := param.Get("title").String()
	duration := param.Get("duration").Int()
	fileSize := param.Get("fileSize").Int64()

	// 查找该文件是否被vip账号存储过
	data, err := dao.VipSaveFileDao.GetByMd5(md5)
	if err != nil && err != mgo.ErrNotFound {
		err = artisan.NewError("database error", -1, err)
		return
	}
	// 没有存储过, 使用秒传进行存储
	if err == mgo.ErrNotFound {
		data = dao.VipSaveFileModel{
			Username:  "",
			Md5:       md5,
			SliceMd5:  sliceMd5,
			FileSize:  0,
			Fid:       "",
			AddAt:     time.Now().Unix(),
			HitAt:     time.Now().Unix(),
			DeletedAt: 0,
		}
		viper := pan_viper.GetVip()
		data.Username = viper.Username()
		data.Fid, data.FileSize, err = viper.SaveFileByMd5(md5, sliceMd5, data.GetSavePath(), fileSize)
		if err != nil {
			err = artisan.NewError("vip账号转存文件错误", -1, err)
			return
		}
		err = dao.VipSaveFileDao.Insert(data)
		if err != nil {
			err = artisan.NewError("database error", -1, err)
			return
		}
	} else { // 存储过, 更新命中时间戳
		err = dao.VipSaveFileDao.Hit(data)
		if err != nil {
			err = artisan.NewError("database error", -1, err)
			return
		}
	}
	// 写入分享表
	share := dao.FileShareModel{
		Uk:       middleware.ContextLoginInfo(ctx).Uk(),
		Title:    title,
		Md5:      data.Md5,
		SliceMd5: data.SliceMd5,
		FileSize: data.FileSize,
		ShareAt:  time.Now().Unix(),
		ExpireAt: time.Now().Add(time.Hour * 24 * time.Duration(duration)).Unix(),
	}
	dao.FileShareDao.Insert(share)
	if err != nil {
		err = artisan.NewError("database error", -1, err)
		return
	}
	result = share
	return
}
