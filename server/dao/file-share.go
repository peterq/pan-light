package dao

import (
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
)

type FileShareModel struct {
	Uk       string
	Title    string
	Md5      string
	SliceMd5 string `bson:"slice_md5"`
	FileSize int64  `bson:"file_size"`
	ShareAt  int64  `bson:"share_at"`
	ExpireAt int64  `bson:"expire_at"`
}

type fileShareDao struct{}

func (*fileShareDao) collection(s *mgo.Session) *mgo.Collection {
	return s.DB(conf.Conf.Database).C(conf.CollectionFileShare)
}

func (d *fileShareDao) Insert(data FileShareModel) (err error) {
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Insert(data)
	return
}

var FileShareDao = &fileShareDao{}
