package dao

import (
	"errors"
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
)

type gson = map[string]interface{}

var UserDao = &userDao{}

type UserModel struct {
	Uk           string `json:"uk" bson:"uk"`                       // 百度 uk
	MarkUsername string `json:"mark_username" bson:"mark_username"` // 脱敏 用户名
	Username     string `json:"username" bson:"username,omitempty"` // 用户名
	Avatar       string `json:"avatar" bson:"avatar"`
	IsVip        bool   `json:"is_vip" bson:"is_vip"`
	IsSuperVip   bool   `json:"is_super_vip" bson:"is_super_vip"`
}

type userDao struct{}

func (d *userDao) UpInsert(model UserModel) (err error) {
	if model.Uk == "" {
		return errors.New("uk empty")
	}
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	_, err = collection.Upsert(gson{
		"uk": model.Uk,
	}, model)
	return
}

func (d *userDao) collection(s *mgo.Session) *mgo.Collection {
	return s.DB(conf.Conf.Database).C(conf.CollectionUser)
}

func (d *userDao) GetByUk(uk string) (*UserModel, error) {
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)

	var model UserModel
	err := collection.Find(gson{
		"uk": uk,
	}).One(&model)

	return &model, err
}
