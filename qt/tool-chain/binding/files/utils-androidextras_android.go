package androidextras

import (
	"C"
	"strings"
	"unsafe"

	"github.com/peterq/pan-light/qt"
)

func assertion(key int, input ...interface{}) (unsafe.Pointer, func()) {
	if len(input) > key {
		switch deduced := input[key].(type) {
		case string:
			jObject := QAndroidJniObject_FromString(deduced)

			return jObject.Object(), func() { jObject.DestroyQAndroidJniObject() }

		case []string:
			jObject := QAndroidJniObject_FromString(strings.Join(deduced, ",,,"))
			jObject2 := jObject.CallObjectMethod2("split", "(Ljava/lang/String;)[Ljava/lang/String;", ",,,")
			jObject.DestroyQAndroidJniObject()

			return jObject2.Object(), func() { jObject2.DestroyQAndroidJniObject() }

		case bool:
			return unsafe.Pointer(uintptr(C.char(int8(qt.GoBoolToInt(deduced))))), nil

		case int16:
			return unsafe.Pointer(uintptr(C.short(deduced))), nil

		case uint16:
			return unsafe.Pointer(uintptr(C.ushort(deduced))), nil

		case int:
			return unsafe.Pointer(uintptr(C.int(int32(deduced)))), nil

		case uint:
			return unsafe.Pointer(uintptr(C.uint(uint32(deduced)))), nil

		case int32:
			return unsafe.Pointer(uintptr(C.int(deduced))), nil

		case uint32:
			return unsafe.Pointer(uintptr(C.uint(deduced))), nil

		case int64:
			return unsafe.Pointer(uintptr(C.longlong(deduced))), nil

		case uint64:
			return unsafe.Pointer(uintptr(C.ulonglong(deduced))), nil

		case float32:
			return unsafe.Pointer(uintptr(C.float(deduced))), nil

		case float64:
			return unsafe.Pointer(uintptr(C.double(deduced))), nil

		case uintptr:
			return unsafe.Pointer(deduced), nil

		case unsafe.Pointer:
			return deduced, nil

		case *QAndroidJniObject:
			return deduced.Object(), nil
		}
	}
	return nil, nil
}
