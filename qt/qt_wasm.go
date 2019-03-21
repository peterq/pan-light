// +build js,wasm

package qt

import (
	"syscall/js"
	"unsafe"
)

func init() {
	WASM.Set("_callbackReleaseTypedArray", js.NewCallback(func(_ js.Value, args []js.Value) interface{} {
		(*js.TypedArray)(unsafe.Pointer(uintptr(args[0].Int()))).Release()
		return nil
	}))
}

var Global = js.Global()
var Module = Global.Call("eval", "Module")

//TODO: func MakeWrapper(i interface{}) *js.Value

//

var WASM = Module //TODO: remove
