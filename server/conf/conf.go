package conf

import (
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2"
	"os"
	"strings"
)

type conf struct {
	Debug      bool
	AppSecret  string
	MongodbUri string
	Database   string
	Redis      redis.RingOptions
}

var Conf *conf
var IrisConf iris.Configuration
var MongodbSession *mgo.Session
var Redis *redis.Ring

func init() {
	confFile, ok := os.LookupEnv("pan_light_server_conf")
	if !ok {
		panic("the conf path must be specified")
	}
	IrisConf = iris.YAML(confFile)
	Conf = &conf{
		Debug:      getConf("debug").(bool),
		AppSecret:  getConf("app-secret").(string),
		MongodbUri: getConf("mongodb-uri").(string),
		Database:   getConf("database").(string),
		Redis: redis.RingOptions{
			Addrs: map[string]string{
				"main": getConf("redis.addr").(string),
			},
			Password: getConf("redis.pwd").(string),
			DB:       getConf("redis.db").(int),
		},
	}
	connectMongo()
	connectRedis()
}

func connectMongo() {
	var err error
	MongodbSession, err = mgo.Dial(Conf.MongodbUri)
	if err != nil {
		panic(err)
	}
	MongodbSession.Refresh()
}

func connectRedis() {
	Redis = redis.NewRing(&Conf.Redis)
}

func getConf(key string) interface{} {
	p := strings.Split(key, ".")
	cnf := IrisConf.Other
	var parent map[interface{}]interface{}
	for idx, name := range p {
		if len(p) == 1 {
			return cnf[name]
		}
		if idx == len(p)-1 {
			return parent[name]
		}
		if idx == 0 {
			parent = cnf[name].(map[interface{}]interface{})
		} else {
			parent = parent[name].(map[interface{}]interface{})
		}
	}
	return nil
}
