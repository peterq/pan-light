package parser

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func (c *Class) fix() {
	c.fixFunctions((*Function).fix)

	c.fixEnums()

	//c.fixGeneral()
	c.fixGeneral_Version()

	c.fixLinkage()

	c.fixBases()

	c.FixGenericHelper()

	c.fixFunctions((*Function).fixGeneral_AfterClasses)
}

func (c *Class) fixFunctions(fix func(*Function)) {
	for _, f := range c.Functions {
		fix(f)
	}
}

//TODO: (*Enum) IsSupported
//TODO: merge into (*Class).removeEnums_Version
func (c *Class) fixEnums() {
	for _, e := range c.Enums {
		if e.Fullname == "QVariant::Type" {
			e.Status = "active"
			for i := len(e.Values) - 1; i >= 0; i-- {
				if v := e.Values[i]; v.Name == "LastCoreType" || v.Name == "LastGuiType" {
					e.Values = append(e.Values[:i], e.Values[i+1:]...)
				}
			}
		}
		if e.Fullname == "QFileSystemModel::Roles" && utils.QT_VERSION_NUM() <= 5042 {
			for i := len(e.Values) - 1; i >= 0; i-- {
				if v := e.Values[i]; v.Name == "FileIconRole" {
					e.Values = append(e.Values[:i], e.Values[i+1:]...)
				}
			}
		}
		if utils.QT_MACPORTS() {
			if e.Fullname == "QWebSettings::WebAttribute" {
				for i := len(e.Values) - 1; i >= 0; i-- {
					if v := e.Values[i]; v.Name == "MediaSourceEnabled" || v.Name == "MediaEnabled" || v.Name == "WebSecurityEnabled" || v.Name == "FullScreenSupportEnabled" {
						e.Values = append(e.Values[:i], e.Values[i+1:]...)
					}
				}
			}
			if e.Fullname == "QWebPage::MessageSource" {
				for i := len(e.Values) - 1; i >= 0; i-- {
					if v := e.Values[i]; v.Name == "MessageSource" || v.Name == "MessageLevel" {
						e.Values = append(e.Values[:i], e.Values[i+1:]...)
					}
				}
			}
		}
	}
}

func (c *Class) fixGeneral_Version() {
	switch c.Name {
	case "QStyle":
		{
			for _, f := range c.Functions {
				if f.Name != "standardIcon" {
					continue
				}

				var tmpF = *f
				tmpF.Name = "standardPixmap"
				tmpF.Output = "QPixmap"
				tmpF.Fullname = fmt.Sprintf("%v::%v", c.Name, tmpF.Name)

				c.Functions = append(c.Functions, &tmpF)
			}
		}

	case "QScxmlCppDataModel":
		{
			for _, s := range []struct{ Name, Output string }{
				{"evaluateToString", "QString"},
				{"evaluateToBool", "bool"},
				{"evaluateToVariant", "QVariant"},
				{"evaluateToVoid", "void"},
				{"evaluateAssignment", "void"},
				{"evaluateInitialization", "void"},
			} {
				c.Functions = append(c.Functions, &Function{
					Name:     s.Name,
					Fullname: fmt.Sprintf("%v::%v", c.Name, s.Name),
					Access:   "public",
					Virtual:  PURE,
					Meta:     PLAIN,
					Output:   s.Output,
					Parameters: []*Parameter{
						{
							Name:  "id",
							Value: "QScxmlExecutableContent::EvaluatorId",
						},
						{
							Name:  "ok",
							Value: "bool*",
						},
					},
					Signature: "()",
				})
			}

			c.Functions = append(c.Functions, &Function{
				Name:     "evaluateForeach",
				Fullname: fmt.Sprintf("%v::evaluateForeach", c.Name),
				Access:   "public",
				Virtual:  PURE,
				Meta:     PLAIN,
				Output:   "void",
				Parameters: []*Parameter{
					{
						Name:  "id",
						Value: "QScxmlExecutableContent::EvaluatorId",
					},
					{
						Name:  "ok",
						Value: "bool*",
					},
					{
						Name:  "body",
						Value: "ForeachLoopBody*",
					},
				},
				Signature: "()",
			})
		}

	case "QOperatingSystemVersion":
		{
			if utils.QT_VERSION_NUM() >= 5093 {
				c.Functions = append(c.Functions, &Function{
					Name:     "QOperatingSystemVersion",
					Fullname: fmt.Sprintf("%v::QOperatingSystemVersion", c.Name),
					Access:   "public",
					Virtual:  "non",
					Meta:     CONSTRUCTOR,
					Status:   "active",
					Parameters: []*Parameter{
						{
							Name:  "other",
							Value: "const QOperatingSystemVersion &amp;",
						},
					},
					Signature: "QOperatingSystemVersion(const QOperatingSystemVersion &amp;other) = default",
				})
			}
		}
	case "QDesktopWidget":
		{
			for _, f := range c.Functions {
				f.Status = "active"
			}
		}
	}
}

func (c *Class) fixLinkage() {
	switch c.Module {
	case "QtCore":
		{
			c.WeakLink = map[string]struct{}{
				"QtGui":     {},
				"QtWidgets": {},
			}
		}

	case "QtGui":
		{
			c.WeakLink = map[string]struct{}{
				"QtWidgets":    {},
				"QtMultimedia": {},
			}
		}
	}
}

var (
	pkgConfigIncludeDir      string
	pkgConfigIncludeDirMutex = new(sync.Mutex)
)

func (c *Class) fixBases() {
	if c.Module == MOC || c.Pkg != "" {
		return
	}

	if c.Bases == "QList" {
		c.Bases = ""
		return
	}

	switch c.Name {
	case "QChart", "QLegend":
		c.Bases = "QGraphicsWidget"

	case "QChartView":
		c.Bases = "QGraphicsView"

	case "QCandlestickModelMapper", "QHBarModelMapper", "QHBoxPlotModelMapper",
		"QVBarModelMapper", "QVBoxPlotModelMapper", "QTextToSpeech":
		c.Bases = "QObject"
	}

	if c.Module == "QtCharts" {
		return
	}

	//if utils.QT_VERSION() == "5.8.0" {
	if c.Name == "QDesignerCustomWidgetInterface" ||
		c.Name == "QDesignerCustomWidgetCollectionInterface" {
		return
	}
	//}

	var (
		prefixPath = utils.QT_INSTALL_PREFIX(State.Target)
		infixPath  = "include"
		suffixPath = string(filepath.Separator)
	)

	switch State.Target {
	case "darwin":
		{
			if utils.QT_NIX() {
				for _, qmakepath := range strings.Split(os.Getenv("QMAKEPATH"), string(filepath.ListSeparator)) {
					if utils.ExistsFile(filepath.Join(qmakepath, infixPath, c.DocModule+suffixPath+c.Name)) {
						prefixPath = qmakepath
						break
					}
				}
			} else if !utils.QT_STATIC() {
				infixPath = "lib"
				suffixPath = ".framework/Headers/"
			}
		}

	default:
		{
			if utils.QT_SAILFISH() {
				prefixPath = "/srv/mer/targets/SailfishOS-" + utils.QT_SAILFISH_VERSION() + "-i486/usr/include/qt5"
				infixPath = ""
			} else if utils.QT_PKG_CONFIG() {
				pkgConfigIncludeDirMutex.Lock()
				if pkgConfigIncludeDir == "" {
					pkgConfigIncludeDir = strings.TrimSpace(utils.RunCmd(exec.Command("pkg-config", "--variable=includedir", "Qt5Core"), "parser.class_includedir"))
				}
				prefixPath = pkgConfigIncludeDir
				pkgConfigIncludeDirMutex.Unlock()
			} else if strings.HasPrefix(State.Target, "sailfish") && runtime.GOOS == "darwin" {
				infixPath = "lib"
				suffixPath = ".framework/Headers/"
			}
		}
	}

	//TODO: remove
	switch c.Name {
	case "Qt", "QtGlobalStatic", "QUnicodeTools", "QHooks", "QModulesPrivate", "QtMetaTypePrivate", "QUnicodeTables", "QAndroidJniEnvironment", "QAndroidJniObject", "QAndroidActivityResultReceiver", "QSupportedWritingSystems", "QAbstractOpenGLFunctions":
		{
			c.Bases = ""
			return
		}

	case "QFutureWatcher", "QDBusAbstractInterface":
		{
			c.Bases = "QObject"
			return
		}

	case "QDBusPendingReply":
		{
			c.Bases = "QDBusPendingCall"
			return
		}

	case "QRasterPaintEngine":
		{
			c.Bases = "QPaintEngine"
			return
		}

	case "QUiLoader", "QEGLNativeContext", "QWGLNativeContext", "QGLXNativeContext", "QEglFSFunctions", "QWindowsWindowFunctions", "QCocoaNativeContext", "QXcbWindowFunctions", "QCocoaWindowFunctions":
		{
			if utils.QT_PKG_CONFIG() && runtime.GOOS != "darwin" {
				c.Bases = getBasesFromHeader(utils.LoadOptional(filepath.Join(prefixPath, c.Module, strings.ToLower(c.Name)+".h")), c.Name, c.Module)
			} else {
				c.Bases = getBasesFromHeader(utils.Load(filepath.Join(prefixPath, "include", c.Module, strings.ToLower(c.Name)+".h")), c.Name, c.Module)
			}
			return
		}

	case "QPlatformSystemTrayIcon", "QPlatformGraphicsBuffer":
		{
			if utils.QT_PKG_CONFIG() && runtime.GOOS != "darwin" {
				c.Bases = getBasesFromHeader(utils.LoadOptional(filepath.Join(prefixPath, c.Module, utils.QT_VERSION(), c.Module, "qpa", strings.ToLower(c.Name)+".h")), c.Name, c.Module)
			} else {
				c.Bases = getBasesFromHeader(utils.Load(filepath.Join(prefixPath, infixPath, c.Module+suffixPath+utils.QT_VERSION(), "QtGui", "qpa", strings.ToLower(c.Name)+".h")), c.Name, c.Module)
			}
			return
		}

	case "QColumnView", "QLCDNumber", "QWebEngineUrlSchemeHandler", "QWebEngineUrlRequestInterceptor", "QWebEngineCookieStore", "QWebEngineUrlRequestInfo", "QWebEngineUrlRequestJob":
		{
			for _, m := range append(LibDeps[strings.TrimPrefix(c.Module, "Qt")], strings.TrimPrefix(c.Module, "Qt")) {
				m = fmt.Sprintf("Qt%v", m)
				if utils.QT_PKG_CONFIG() && runtime.GOOS != "darwin" {
					if utils.ExistsFile(filepath.Join(prefixPath, m, strings.ToLower(c.Name)+".h")) {
						c.Bases = getBasesFromHeader(utils.LoadOptional(filepath.Join(prefixPath, m, strings.ToLower(c.Name)+".h")), c.Name, c.Module)
						return
					}
				} else {
					if utils.ExistsFile(filepath.Join(prefixPath, infixPath, m+suffixPath+c.Name)) {
						c.Bases = getBasesFromHeader(utils.Load(filepath.Join(prefixPath, infixPath, m+suffixPath+strings.ToLower(c.Name)+".h")), c.Name, c.Module)
						return
					}
				}
			}
		}
	}

	//TODO:
	var libs = append(LibDeps[strings.TrimPrefix(c.Module, "Qt")], strings.TrimPrefix(c.Module, "Qt"))
	for i, v := range libs {
		if v == "TestLib" {
			libs[i] = "Test"
		}
	}

	var found bool
	for _, m := range libs {
		m = fmt.Sprintf("Qt%v", m)
		if utils.QT_PKG_CONFIG() && runtime.GOOS != "darwin" {
			if utils.ExistsFile(filepath.Join(prefixPath, m, c.Name)) {

				var f = utils.LoadOptional(filepath.Join(prefixPath, m, c.Name))
				if f != "" {
					found = true
					c.Bases = getBasesFromHeader(utils.LoadOptional(filepath.Join(prefixPath, m, strings.Split(f, "\"")[1])), c.Name, m)
				}
				break
			}
		} else {
			if utils.ExistsFile(filepath.Join(prefixPath, infixPath, m+suffixPath+c.Name)) {
				var f = utils.Load(filepath.Join(prefixPath, infixPath, m+suffixPath+c.Name))
				if f != "" {
					found = true
					c.Bases = getBasesFromHeader(utils.Load(filepath.Join(filepath.Join(prefixPath, infixPath, m+suffixPath), strings.Split(f, "\"")[1])), c.Name, m)
				}
				break
			}
		}
	}

	if !found && c.Name != "SailfishApp" && c.Fullname == "" {
		utils.Log.WithField("module", strings.TrimPrefix(c.Module, "Qt")).WithField("class", c.Name).Debugln("failed to find header file")
	}
}

func getBasesFromHeader(f string, n string, m string) string {

	f = strings.Replace(f, "\r", "", -1)

	if strings.HasSuffix(n, "Iterator") {
		return ""
	}

	for i, l := range strings.Split(f, "\n") {

		//TODO: reduce
		if strings.Contains(l, "class "+n) ||
			strings.Contains(l, "class Q_"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"_EXPORT "+n) || strings.Contains(l, "class Q"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"_EXPORT "+n) ||
			strings.Contains(l, "class Q_"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"CORE_EXPORT "+n) || strings.Contains(l, "class Q"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"CORE_EXPORT "+n) ||
			strings.Contains(l, "class QDESIGNER_SDK_EXPORT "+n) || strings.Contains(l, "class QDESIGNER_EXTENSION_EXPORT "+n) || strings.Contains(l, "class QDESIGNER_UILIB_EXPORT "+n) || strings.Contains(l, "class  "+n) || strings.Contains(l, "class Q_"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"_EXPORT  "+n) || strings.Contains(l, "class Q"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"_EXPORT  "+n) || strings.Contains(l, "class QDESIGNER_SDK_EXPORT  "+n) || strings.Contains(l, "class QDESIGNER_EXTENSION_EXPORT  "+n) || strings.Contains(l, "class QDESIGNER_UILIB_EXPORT  "+n) {

			if strings.Contains(l, n+" ") || strings.Contains(l, n+":") || strings.HasSuffix(l, n) {

				l = normalizedClassDeclaration(f, i)

				if !strings.Contains(l, ":") {
					return ""
				}

				if strings.Contains(l, "<") {
					l = strings.Split(l, "<")[0]
				}

				if strings.Contains(l, "/") {
					l = strings.Split(l, "/")[0]
				}

				var tmp = strings.Split(l, ":")[1]

				for _, s := range []string{"{", "}", "#ifndef", "QT_NO_QOBJECT", "#else", "#endif", "class", "Q_" + strings.ToUpper(strings.TrimPrefix(m, "Qt")) + "_EXPORT " + n, "public", "protected", "private", "  ", " "} {
					tmp = strings.Replace(tmp, s, "", -1)
				}

				return strings.TrimSpace(tmp)
			}
		}
	}

	for _, l := range strings.Split(f, "\n") {
		if strings.Contains(l, "struct "+n) || strings.Contains(l, "struct Q_"+strings.ToUpper(strings.TrimPrefix(m, "Qt"))+"_EXPORT "+n) {
			return ""
		}
	}

	for _, l := range strings.Split(f, "\n") {
		if strings.Contains(l, "namespace "+n) {
			return ""
		}
	}

	utils.Log.WithField("module", strings.TrimPrefix(m, "Qt")).WithField("class", n).Debugln("failed to parse header")
	return ""
}

func normalizedClassDeclaration(f string, is int) string {
	var bb = new(bytes.Buffer)
	defer bb.Reset()

	for i, l := range strings.Split(f, "\n") {
		if i >= is {
			fmt.Fprint(bb, l)
			if strings.Contains(l, "{") {
				break
			}
		}
	}
	return bb.String()
}
