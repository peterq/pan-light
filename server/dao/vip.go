package dao

import (
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
)

type VipModel struct {
	Username string
	Bduss    string
}

type vipDao struct{}

func (*vipDao) collection(s *mgo.Session) *mgo.Collection {
	return s.DB(conf.Conf.Database).C(conf.CollectionVip)
}

func (d *vipDao) GetAll() (data []VipModel, err error) {
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Find(nil).All(&data)
	return
}

var Vip = &vipDao{}
