package dao

import (
	"errors"
	"github.com/peterq/pan-light/server/conf"
	"gopkg.in/mgo.v2"
	"time"
)

var FeedbackDao = &feedbackDao{}

type FeedbackModel struct {
	Uk        string `json:"uk" bson:"uk"` // 百度 uk
	Timestamp int64
	Content   string
}

type feedbackDao struct{}

func (d *feedbackDao) Insert(model FeedbackModel) (err error) {
	if model.Uk == "" {
		return errors.New("uk empty")
	}
	if model.Timestamp == 0 {
		model.Timestamp = time.Now().Unix()
	}
	s := conf.MongodbSession.Clone()
	defer s.Refresh()
	collection := d.collection(s)
	err = collection.Insert(&model)
	return
}

func (d *feedbackDao) collection(s *mgo.Session) *mgo.Collection {
	return s.DB(conf.Conf.Database).C(conf.CollectionFeedback)
}
