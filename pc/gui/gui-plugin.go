// +build plugin

package main

import (
	"github.com/peterq/pan-light/pc/gui/bridge"
	_ "github.com/peterq/pan-light/pc/gui/bridge"
	_ "github.com/peterq/pan-light/pc/gui/comp"
	"github.com/peterq/pan-light/qt/bindings/core"
	"github.com/peterq/pan-light/qt/bindings/gui"
	"github.com/peterq/pan-light/qt/bindings/qml"
	"github.com/peterq/pan-light/qt/bindings/quick"
	"os"
)

func StartGui(rccFile, mainQml string) {
	// 开启高清
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	quick.QQuickWindow_SetDefaultAlphaBuffer(true) // 悬浮窗需要此设置

	// 加载qml
	core.QResource_RegisterResource(rccFile, "/")

	app := gui.NewQGuiApplication(len(os.Args), os.Args)

	engine := qml.NewQQmlApplicationEngine(nil)
	engine.Load(core.NewQUrl3(mainQml, 0))
	app.Exec()
}

func SyncRouteRegister(routes map[string]func(map[string]interface{}) interface{}) {
	bridge.SyncRouteRegister(routes)
}

func AsyncRouteRegister(routes map[string]func(map[string]interface{},
	func(interface{}), func(interface{}), func(interface{}), chan interface{})) {
	bridge.AsyncRouteRegister(routes)
}

func NotifyQml(event string, data map[string]interface{}) {
	bridge.NotifyQml(event, data)
}
