// +build !js

package qt

type jsValue interface {
	Call(...string) jsValue
	Int() int
}

var Global jsValue
var Module jsValue
