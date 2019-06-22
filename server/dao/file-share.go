package dao

import (
	"github.com/peterq/pan-light/server/artisan/cache"
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
	"time"
)

type FileShareModel struct {
	Uk       string
	Title    string
	Md5      string
	SliceMd5 string `bson:"slice_md5"`
	FileSize int64  `bson:"file_size"`
	ShareAt  int64  `bson:"share_at"`
	ExpireAt int64  `bson:"expire_at"`
	Official bool
	HotIndex int64 `bson:"hot_index"` // 热度指数
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

type array = []interface{}

func (d *fileShareDao) List(count, offset int, order string) (data []gson, err error) {
	if order == "hottest" {
		return d.HotList()
	}
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	sort := gson{"share_at": -1}
	condition := gson{
		"expire_at": gson{
			"$gt": time.Now().Unix(),
		},
		"share_at": gson{
			"$gt": offset,
		},
	}
	if order == "official" {
		condition["official"] = true
	}
	collection := d.collection(s)
	err = collection.Pipe([]gson{
		{"$match": condition},
		{"$sort": sort},
		{"$limit": count},
		{
			"$lookup": gson{
				"from": conf.CollectionUser,
				"let":  gson{"uk": "$uk"},
				"pipeline": []gson{
					{
						"$match": gson{
							"$expr": gson{"$eq": []string{"$uk", "$$uk"}},
						},
					},
					{
						"$project": gson{
							"avatar": 1, "mark_username": 1,
						},
					},
					{"$limit": 1},
				},
				"as": "user",
			},
		},
		{
			"$replaceRoot": gson{
				"newRoot": gson{"$mergeObjects": array{
					"$$ROOT",
					gson{"user": nil},
					gson{"user": gson{"$arrayElemAt": array{"$user", 0}}},
				}},
			},
		},
	}).All(&data)
	return
}

func (d *fileShareDao) HotList() (data []gson, err error) {
	cacheKey := "cache-file-share-host-list"
	err = cache.RedisGet(cacheKey, &data)
	if err == nil {
		return
	}
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Pipe([]gson{
		{"$match": gson{
			"expire_at": gson{
				"$gt": time.Now().Unix(),
			}},
		},
		{"$sort": gson{"hot_index": -1, "share_at": -1}},
		{"$limit": 50},
		{
			"$lookup": gson{
				"from": conf.CollectionUser,
				"let":  gson{"uk": "$uk"},
				"pipeline": []gson{
					{
						"$match": gson{
							"$expr": gson{"$eq": []string{"$uk", "$$uk"}},
						},
					},
					{
						"$project": gson{
							"avatar": 1, "mark_username": 1,
						},
					},
					{"$limit": 1},
				},
				"as": "user",
			},
		},
		{
			"$replaceRoot": gson{
				"newRoot": gson{"$mergeObjects": array{
					"$$ROOT",
					gson{"user": nil},
					gson{"user": gson{"$arrayElemAt": array{"$user", 0}}},
				}},
			},
		},
	}).All(&data)
	if err != nil {
		return
	}
	if data == nil {
		data = []gson{}
	} else {
		cache.RedisSet(cacheKey, data, time.Minute)
	}
	return
}

var FileShareDao = &fileShareDao{}
