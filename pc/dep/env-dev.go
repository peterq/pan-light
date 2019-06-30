// +build !production

package dep

var Env = tEnv{
	Dev:               true,
	Platform:          "",
	ApiBase:           "http://127.0.0.1:8081",
	InternalServerUrl: "",
	VersionString:     "v0.0.0-development",
	DataPath:          "",
	Version:           20190620001,
	ClientUA:          "pan-light/v0.0.0-development;",
}
