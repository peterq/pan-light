package qt

import (
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

var (
	Logger = log.New(os.Stderr, "", log.Ltime)

	signals      = make(map[unsafe.Pointer]map[string]interface{})
	signalsJNI   = make(map[string]map[string]interface{})
	signalsMutex = new(sync.Mutex)

	objects      = make(map[unsafe.Pointer]interface{})
	objectsMutex = new(sync.Mutex)

	objectsTemp      = make(map[unsafe.Pointer]interface{})
	objectsTempMutex = new(sync.Mutex)
)

func init() { runtime.LockOSThread() }

func ExistsSignal(cPtr unsafe.Pointer, signal string) (exists bool) {
	signalsMutex.Lock()
	_, exists = signals[cPtr][signal]
	signalsMutex.Unlock()
	return
}

func LendSignal(cPtr unsafe.Pointer, signal string) (s interface{}) {
	signalsMutex.Lock()
	s = signals[cPtr][signal]
	signalsMutex.Unlock()
	return
}

func lendSignalJNI(cPtr, signal string) (s interface{}) {
	signalsMutex.Lock()
	s = signalsJNI[cPtr][signal]
	signalsMutex.Unlock()
	return
}

func GetSignal(cPtr interface{}, signal string) interface{} {
	if dcPtr, ok := cPtr.(unsafe.Pointer); ok {
		if signal == "destroyed" || strings.HasPrefix(signal, "~") {
			defer DisconnectAllSignals(dcPtr, signal)
		}
		return LendSignal(dcPtr, signal)
	}
	return lendSignalJNI(cPtr.(string), signal)
}

func ConnectSignal(cPtr interface{}, signal string, function interface{}) {
	if dcPtr, ok := cPtr.(unsafe.Pointer); ok {
		signalsMutex.Lock()
		if s, exists := signals[dcPtr]; !exists {
			signals[dcPtr] = map[string]interface{}{signal: function}
		} else {
			s[signal] = function
		}
		signalsMutex.Unlock()
	} else {
		connectSignalJNI(cPtr.(string), signal, function)
	}
}

func connectSignalJNI(cPtr, signal string, function interface{}) {
	signalsMutex.Lock()
	if s, exists := signalsJNI[cPtr]; !exists {
		signalsJNI[cPtr] = map[string]interface{}{signal: function}
	} else {
		s[signal] = function
	}
	signalsMutex.Unlock()
}

func DisconnectSignal(cPtr interface{}, signal string) {
	if dcPtr, ok := cPtr.(unsafe.Pointer); ok {
		signalsMutex.Lock()
		delete(signals[dcPtr], signal)
		signalsMutex.Unlock()
	} else {
		disconnectSignalJNI(cPtr.(string), signal)
	}
}

func disconnectSignalJNI(cPtr, signal string) {
	signalsMutex.Lock()
	delete(signalsJNI[cPtr], signal)
	signalsMutex.Unlock()
}

func DisconnectAllSignals(cPtr unsafe.Pointer, signal string) {
	signalsMutex.Lock()
	if s, exists := signals[cPtr]["destroyed"]; signal != "destroyed" && exists {
		signals[cPtr] = map[string]interface{}{"destroyed": s}
	} else {
		delete(signals, cPtr)
	}
	signalsMutex.Unlock()
	if signal == "destroyed" {
		Unregister(cPtr)
	}
}

func DumpSignals() {
	Debug("##############################\tSIGNALSTABLE_START\t##############################")
	signalsMutex.Lock()
	for cPtr, entry := range signals {
		Debug(cPtr, entry)
	}
	signalsMutex.Unlock()
	Debug("##############################\tSIGNALSTABLE_END\t##############################")
}

func CountSignals() (c int) {
	signalsMutex.Lock()
	c = len(signals)
	signalsMutex.Unlock()
	return
}

func GoBoolToInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func Recover(fn string) {
	if recover() != nil {
		Debug("RECOVERED:", fn)
	}
}

func Debug(fn ...interface{}) {
	if strings.ToLower(os.Getenv("QT_DEBUG")) == "true" || runtime.GOARCH == "js" || runtime.GOARCH == "wasm" {
		Logger.Println(fn...)
	}
}

func ClearSignals() {
	signalsMutex.Lock()
	signals = make(map[unsafe.Pointer]map[string]interface{})
	signalsMutex.Unlock()
}

func Register(cPtr unsafe.Pointer, gPtr interface{}) {
	objectsMutex.Lock()
	objects[cPtr] = gPtr
	objectsMutex.Unlock()
}

func Receive(cPtr unsafe.Pointer) (o interface{}, ok bool) {
	objectsMutex.Lock()
	o, ok = objects[cPtr]
	objectsMutex.Unlock()
	return
}

func Unregister(cPtr unsafe.Pointer) {
	objectsMutex.Lock()
	delete(objects, cPtr)
	objectsMutex.Unlock()
}

func RegisterTemp(cPtr unsafe.Pointer, gPtr interface{}) {
	objectsTempMutex.Lock()
	objectsTemp[cPtr] = gPtr
	objectsTempMutex.Unlock()
}

func ReceiveTemp(cPtr unsafe.Pointer) (o interface{}, ok bool) {
	objectsTempMutex.Lock()
	o, ok = objectsTemp[cPtr]
	objectsTempMutex.Unlock()
	return
}

func UnregisterTemp(cPtr unsafe.Pointer) {
	objectsTempMutex.Lock()
	delete(objectsTemp, cPtr)
	objectsTempMutex.Unlock()
}
