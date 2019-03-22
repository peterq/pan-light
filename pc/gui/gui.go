// +build !plugin

package gui

import (
	_ "github.com/peterq/pan-light/pc/functions"
	_ "github.com/peterq/pan-light/pc/gui/bridge"
	_ "github.com/peterq/pan-light/pc/gui/comp"
	_ "github.com/peterq/pan-light/pc/gui/qml"
	"github.com/peterq/pan-light/qt/bindings/core"
	"github.com/peterq/pan-light/qt/bindings/gui"
	"github.com/peterq/pan-light/qt/bindings/qml"
	"github.com/peterq/pan-light/qt/bindings/quick"
	"os"
)

func StartGui() {
	// 开启高清
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	quick.QQuickWindow_SetDefaultAlphaBuffer(true) // 悬浮窗需要此设置

	app := gui.NewQGuiApplication(len(os.Args), os.Args)

	engine := qml.NewQQmlApplicationEngine(nil)
	engine.Load(core.NewQUrl3("qrc:/main.qml", 0))
	app.Exec()
}
