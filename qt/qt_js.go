// +build js,!wasm

package qt

import "github.com/gopherjs/gopherjs/js"

var Global = js.Global
var Module = Global.Call("eval", "Module")

func MakeWrapper(i interface{}) *js.Object {
	o := js.InternalObject(i)
	methods := o.Get("constructor").Get("methods")
	for i := 0; i < methods.Length(); i++ {
		m := methods.Index(i)
		if m.Get("pkg").String() != "" { // not exported
			continue
		}
		if o.Get(m.Get("name").String()) == js.Undefined {
			o.Set(m.Get("name").String(), func(args ...*js.Object) *js.Object {
				return js.Global.Call("$externalizeFunction", o.Get(m.Get("prop").String()), m.Get("typ"), true).Call("apply", o, args)
			})
		}
	}
	return o
}

//

var WASM = Module //TODO: remove
