// +build !plugin

package gui

import (
	"github.com/peterq/pan-light/pc/dep"
	_ "github.com/peterq/pan-light/pc/functions"
	_ "github.com/peterq/pan-light/pc/gui/bridge"
	_ "github.com/peterq/pan-light/pc/gui/comp"
	_ "github.com/peterq/pan-light/pc/gui/qml"
	"github.com/peterq/pan-light/qt/bindings/core"
	"github.com/peterq/pan-light/qt/bindings/gui"
	"github.com/peterq/pan-light/qt/bindings/qml"
	"github.com/peterq/pan-light/qt/bindings/quick"
	"log"
	"os"
)

func StartGui() {
	// 开启高清
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	quick.QQuickWindow_SetDefaultAlphaBuffer(true) // 悬浮窗需要此设置

	// 下面2句话居然能解决windows 异常退出的bug
	core.QCoreApplication_SetOrganizationName("PeterQ") //needed to fix an QML Settings issue on windows
	if os.Getenv("pan_light_render_exception_fix") == "true" {
		quick.QQuickWindow_SetSceneGraphBackend(quick.QSGRendererInterface__Software)
	}

	//rccFile := "E:\\pan-light\\qml.rcc"
	//bin, _ := ioutil.ReadFile(rccFile)
	//go func() {
	//	for range time.Tick(2 * time.Second) {
	//		n, _ := ioutil.ReadFile(rccFile)
	//		if !bytes.Equal(bin, n) {
	//			os.Exit(2)
	//		}
	//	}
	//}()
	//core.QResource_RegisterResource(rccFile, "/")

	app := gui.NewQGuiApplication(len(os.Args), os.Args)

	engine := qml.NewQQmlApplicationEngine(nil)
	engine.Load(core.NewQUrl3("qrc:/main.qml", 0))
	dep.OnClose(func() {
		log.Println("will exit ui")
		app.Exit(dep.ExitCode)
	})
	app.Exec()
}
