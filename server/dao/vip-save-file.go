package dao

import (
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type VipSaveFileModel struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Username  string        `bson:"username,omitempty"` // vip 用户名
	Md5       string        `bson:"md5,omitempty"`
	SliceMd5  string        `bson:"slice_md5,omitempty"`
	FileSize  int64         `bson:"file_size,omitempty"`
	Fid       string        `bson:"fid,omitempty"`
	AddAt     int64         `bson:"add_at,omitempty"` // 转存时间
	HitAt     int64         `bson:"hit_at,omitempty"` // 重复转存命中时间
	DeletedAt int64         `bson:"deleted_at"`       // 删除时间
}

func (f *VipSaveFileModel) GetSavePath() string {
	log.Printf("%#v", f)
	return "/pan-light-save/" + f.Md5[:2] + "/" + f.Md5[2:4] + "/" + f.Md5
}

type vipSaveFileDao struct{}

func (*vipSaveFileDao) collection(s *mgo.Session) *mgo.Collection {
	return s.DB(conf.Conf.Database).C(conf.CollectionVipSaveFile)
}

func (d *vipSaveFileDao) GetByMd5(md5 string) (data VipSaveFileModel, err error) {
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Pipe([]bson.M{
		{
			"$match": bson.M{
				"deleted_at": 0,
				"md5":        md5,
			},
		},
		{
			"$lookup": bson.M{
				"from":         conf.CollectionVip,
				"localField":   "username",
				"foreignField": "username",
				"as":           "viper",
			},
		},
		{
			"$match": bson.M{
				"viper.enabled": true,
			},
		},
	}).One(&data)
	return
}

func (d *vipSaveFileDao) Insert(data VipSaveFileModel) (err error) {
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Insert(data)
	return
}

func (d *vipSaveFileDao) Hit(data VipSaveFileModel) (err error) {
	data.HitAt = time.Now().Unix()
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Update(bson.M{
		"_id": data.Id,
	}, bson.M{
		"$set": bson.M{
			"hit_at": data.HitAt,
		},
	})
	return
}

var VipSaveFileDao = &vipSaveFileDao{}
