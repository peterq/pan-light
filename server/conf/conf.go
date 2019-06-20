package conf

import (
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"os"
	"strings"
)

type conf struct {
	AppSecret  string
	MongodbUri string
	Database   string
}

var Conf *conf
var IrisConf iris.Configuration
var MongodbSession *mgo.Session

func init() {
	confFile, ok := os.LookupEnv("pan_light_server_conf")
	if !ok {
		panic("the conf path must be specified")
	}
	IrisConf = iris.YAML(confFile)
	Conf = &conf{
		AppSecret:  getConf("app-secret").(string),
		MongodbUri: getConf("mongodb-uri").(string),
		Database:   getConf("database").(string),
	}
	var err error
	MongodbSession, err = mgo.Dial(Conf.MongodbUri)
	if err != nil {
		panic(err)
	}
	MongodbSession.Refresh()
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
		parent = parent[name].(map[interface{}]interface{})
	}
	return nil
}
