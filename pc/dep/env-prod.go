// +build production

package dep

var Env = tEnv{
	Dev:               false,
	Platform:          "",
	ApiBase:           "https://pan-light.peterq.cn",
	InternalServerUrl: "",
	VersionString:     "v0.0.1-preview",
	DataPath:          "",
	Version:           20190626001,
	ClientUA:          "pan-light/v0.0.1-preview;build 20190626;",
}
