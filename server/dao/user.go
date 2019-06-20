package dao

import (
	"errors"
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
)

type gson = map[string]interface{}

var UserDao = &userDao{}

type UserModel struct {
	Uk           string // 百度 uk
	MarkUsername string // 脱敏 用户名
	Username     string `json:"username,omitempty"` // 用户名
	Avatar       string
	IsVip        bool
	IsSuperVip   bool
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
