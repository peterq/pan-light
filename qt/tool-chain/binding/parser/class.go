package parser

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

type Class struct {
	Name       string      `xml:"name,attr"`
	Status     string      `xml:"status,attr"`
	Access     string      `xml:"access,attr"`
	Abstract   bool        `xml:"abstract,attr"`
	Bases      string      `xml:"bases,attr"`
	Module     string      `xml:"module,attr"`
	Brief      string      `xml:"brief,attr"`
	Functions  []*Function `xml:"function"`
	Enums      []*Enum     `xml:"enum"`
	Variables  []*Variable `xml:"variable"`
	Properties []*Variable `xml:"property"`
	Classes    []*Class    `xml:"class"`
	Since      string      `xml:"since,attr"`

	DocModule    string
	Stub         bool
	WeakLink     map[string]struct{}
	Export       bool
	Fullname     string
	Pkg          string
	Path         string
	HasFinalizer bool

	Constructors []string
	Derivations  []string

	sync.Mutex
}

func (c *Class) register(m *Module) {
	c.DocModule = c.Module
	c.Module = m.Project
	c.Pkg = m.Pkg
	State.ClassMap[c.Name] = c

	for _, sc := range c.Classes {
		if sc.Name != "PaintContext" { //TODO: remove and support all sub classes
			continue
		}

		sc.Fullname = fmt.Sprintf("%v::%v", c.Name, sc.Name)
		sc.register(m)
	}
}

func (c *Class) derivation() {
	for _, b := range c.GetBases() {
		if bc, e := State.ClassMap[b]; e {
			bc.Derivations = append(bc.Derivations, c.Name)
		}
	}
}

func (c *Class) GetBases() []string {
	if c == nil || c.Bases == "" {
		return make([]string, 0)
	}
	if strings.Contains(c.Bases, ",") {
		return strings.Split(c.Bases, ",")
	}
	return []string{c.Bases}
}

func (c *Class) GetAllBases() []string {
	var out, _ = c.GetAllBasesRecursiveCheckFailed(0)
	return out
}

func (c *Class) GetAllBasesRecursiveCheckFailed(i int) ([]string, bool) {
	var input = make([]string, 0)

	i++
	if i > 100 {
		return input, true
	}

	for _, b := range c.GetBases() {
		var bc, ok = State.ClassMap[b]
		if !ok {
			continue
		}

		input = append(input, b)
		var bs, isRecursive = bc.GetAllBasesRecursiveCheckFailed(i)
		if isRecursive {
			return input, true
		}
		for _, sbc := range bs {
			input = append(input, sbc)
		}
	}

	return input, false
}

func (c *Class) IsSubClassOfQObject() bool {
	return c.IsSubClassOf("QObject")
}

func (c *Class) IsSubClassOf(class string) bool {
	if c == nil {
		return false
	}

	for _, bcn := range append([]string{c.Name}, c.GetAllBases()...) {
		if bcn == class {
			return true
		}
	}

	return false
}

func (c *Class) isSubClass() bool { return c.Fullname != "" }

func (c *Class) GetAllDerivations() []string {

	var input = make([]string, 0)

	for _, b := range c.Derivations {
		var bc, exists = State.ClassMap[b]
		if !exists {
			continue
		}

		input = append(input, b)
		for _, sbc := range bc.GetAllDerivations() {
			input = append(input, sbc)
		}
	}

	return input
}

func (c *Class) GetAllDerivationsInSameModule() []string {

	var input = make([]string, 0)

	for _, i := range c.GetAllDerivations() {
		if State.ClassMap[i].Module == c.Module && i != "QWinEventNotifier" && i != "QFutureWatcher" {
			input = append(input, i)
		}
	}

	return input
}

func (c *Class) HasFunction(f *Function) bool {
	for _, cf := range c.Functions {
		if cf.Name == f.Name && cf.Virtual == f.Virtual &&
			cf.Meta == f.Meta &&
			cf.Output == f.Output && len(cf.Parameters) == len(f.Parameters) {

			var similar = true
			for i, cfp := range cf.Parameters {
				if cfp.Value != f.Parameters[i].Value {
					similar = false
				}
			}
			if similar {
				return true
			}
		}
	}

	return false
}

func (c *Class) HasFunctionWithName(n string) bool {
	return c.HasFunctionWithNameAndOverloadNumber(n, "")
}

func (c *Class) HasFunctionWithNameAndOverloadNumber(n string, num string) bool {
	for _, f := range c.Functions {
		if strings.ToLower(f.Name) == strings.ToLower(n) && f.OverloadNumber == num {
			return true
		}
	}
	return false
}

func (c *Class) IsPolymorphic() bool { return len(c.GetBases()) > 1 }

func (c *Class) HasConstructor() bool {
	for _, f := range c.Functions {
		if f.Meta == CONSTRUCTOR || f.Meta == COPY_CONSTRUCTOR || f.Meta == MOVE_CONSTRUCTOR {
			return true
		}
	}
	return false
}

func (c *Class) HasDestructor() bool {
	for _, f := range c.Functions {
		if f.Meta == DESTRUCTOR {
			return true
		}
	}
	return false
}

func (c *Class) HasCallbackFunctions() bool {
	for _, bcn := range append([]string{c.Name}, c.GetAllBases()...) {
		var bc, ok = State.ClassMap[bcn]
		if !ok {
			continue
		}
		for _, f := range bc.Functions {
			if f.Virtual == IMPURE || f.Virtual == PURE || f.Meta == SIGNAL || f.Meta == SLOT {
				return true
			}
		}
	}
	return false
}

func (c *Class) IsSupported() bool {
	if c == nil {
		return false
	}

	switch c.Name {
	case "QCborStreamReader", "QCborStreamWriter", "QCborValue", "QScopeGuard", "QTest",
		"QImageReaderWriterHelpers", "QPasswordDigestor", "QDtls", "QDtlsClientVerifier":
		c.Access = "unsupported_isBlockedClass"
		return false
	}

	if utils.QT_VERSION_NUM() >= 5080 {
		switch c.Name {
		case "QSctpServer", "QSctpSocket", "Http2", "QAbstractExtensionFactory":
			{
				c.Access = "unsupported_isBlockedClass"
				return false
			}
		}
	}

	if UseJs() {
		if strings.HasPrefix(c.Name, "QOpenGLFunctions_") || strings.HasPrefix(c.Name, "QSsl") || strings.HasPrefix(c.Name, "QNetwork") {
			c.Access = "unsupported_isBlockedClass"
			return false
		}
		switch c.Name {
		case "QThreadPool", "QSharedMemory", "QPluginLoader", "QSemaphore", "QSemaphoreReleaser",
			"QSystemSemaphore", "QThread", "QWaitCondition", "QUnhandledException", "QFileSystemModel",
			"QLibrary", "QUdpSocket", "QHttpMultiPart", "QHttpPart", "QOpenGLTimeMonitor", "QOpenGLTimerQuery",
			"QRemoteObjectHost": //TODO: only block in 5.11 ?
			{
				c.Access = "unsupported_isBlockedClass"
				return false
			}
		}
	}

	switch c.Name {
	case
		"QString", "QStringList", //mapped to primitive

		"QExplicitlySharedDataPointer", "QFuture", "QDBusPendingReply", "QDBusReply", "QFutureSynchronizer", //needs template
		"QGlobalStatic", "QMultiHash", "QQueue", "QMultiMap", "QScopedPointer", "QSharedDataPointer",
		"QScopedArrayPointer", "QSharedPointer", "QThreadStorage", "QScopedValueRollback", "QVarLengthArray",
		"QWeakPointer", "QWinEventNotifier",

		"QFlags", "QException", "QStandardItemEditorCreator", "QSGSimpleMaterialShader", "QGeoCodeReply", "QFutureWatcher", //other
		"QItemEditorCreator", "QGeoCodingManager", "QGeoCodingManagerEngine", "QQmlListProperty",

		"QPlatformGraphicsBuffer", "QPlatformSystemTrayIcon", "QRasterPaintEngine", "QSupportedWritingSystems", "QGeoLocation", //file not found or QPA API
		"QAbstractOpenGLFunctions",

		"QProcess", "QProcessEnvironment", //TODO: iOS

		"QRemoteObjectPackets",

		"QStaticByteArrayMatcher", "QtDummyFutex", "QtLinuxFutex",
		"QShaderLanguage",

		"AndroidNfc", "OSXBluetooth",

		"QtROClientFactory", "QtROServerFactory",

		"QWebViewFactory", "QGeoServiceProviderFactoryV2",

		"QtDwmApiDll", "QWinMime",

		"QAbstract3DGraph": //TODO: only for arch with pkg_config
		{
			c.Access = "unsupported_isBlockedClass"
			return false
		}
	}

	switch {
	case
		c.Name == "QOpenGLFunctions_ES2", strings.HasPrefix(c.Name, "QPlace"), //file not found or QPA API

		strings.HasPrefix(c.Name, "QAtomic"), //other

		strings.HasSuffix(c.Name, "terator"), strings.Contains(c.Brief, "emplate"), //needs template

		strings.HasPrefix(c.Name, "QVulkan"):

		{
			c.Access = "unsupported_isBlockedClass"
			return false
		}
	}

	if strings.HasPrefix(c.Name, "QOpenGL") && (os.Getenv("DEB_TARGET_ARCH_CPU") == "arm" || (strings.HasPrefix(State.Target, "android") && utils.QT_FAT())) { //TODO: block indiv classes for fat android build instead
		c.Access = "unsupported_isBlockedClass"
		return false
	}

	if utils.QT_VERSION_NUM() <= 5042 {
		if c.Name == "QQmlAbstractProfilerAdapter" ||
			c.Name == "QGraphicsLayout" ||
			c.Name == "QQmlAbstractProfilerAdapter" ||
			c.Name == "QDeclarativeMultimedia" {
			c.Access = "unsupported_isBlockedClass"
			return false
		}
	}

	if strings.HasPrefix(c.Name, "QOpenGLFunctions_") && !utils.QT_GEN_OPENGL() {
		c.Access = "unsupported_isBlockedClass"
		return false
	}

	if State.Minimal {
		return c.Export
	}

	return true
}

func (c *Class) GetFunction(fname string) *Function {
	for _, f := range c.Functions {
		if f.Name == fname {
			return f
		}
	}
	return nil
}

//

func (c *Class) Hash() string {
	h := sha1.New()
	h.Write([]byte(c.Path))
	return hex.EncodeToString(h.Sum(nil)[:3])
}
