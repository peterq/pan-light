package parser

import (
	"fmt"
	"strings"
)

func (c *Class) add() {
	c.addGeneralFuncs()

	c.addVarAndPropFuncs()

	c.addMocFuncs()
}

func (c *Class) addGeneralFuncs() {
	switch c.Name {
	case "QColor", "QFont", "QImage", "QObject", "QIcon", "QBrush":
		{
			c.Functions = append(c.Functions, &Function{
				Name:       "toVariant",
				Fullname:   fmt.Sprintf("%v::toVariant", c.Name),
				Access:     "public",
				Virtual:    "non",
				Meta:       PLAIN,
				Output:     "QVariant",
				Parameters: []*Parameter{},
				Signature:  "()",
			})
		}

	case "QVariant":
		{
			for _, name := range []string{"toColor", "toFont", "toImage", "toObject", "toIcon", "toBrush"} {
				c.Functions = append(c.Functions, &Function{
					Name:       name,
					Fullname:   fmt.Sprintf("%v::%v", c.Name, name),
					Access:     "public",
					Virtual:    "non",
					Meta:       PLAIN,
					Output:     strings.Replace(name, "to", "Q", -1),
					Parameters: []*Parameter{},
					Signature:  "()",
				})
			}
		}

	case "QQmlEngine":
		{
			//http://doc.qt.io/qt-5/qqmlengine.html#qmlRegisterSingletonType-2
			//int qmlRegisterSingletonType(const QUrl &url, const char *uri, int versionMajor, int versionMinor, const char *qmlName)
			c.Functions = append(c.Functions, &Function{
				Name:      "qmlRegisterSingletonType",
				Fullname:  fmt.Sprintf("%v::qmlRegisterSingletonType", c.Name),
				Access:    "public",
				Virtual:   "non",
				Meta:      PLAIN,
				NonMember: true,
				Static:    true,
				Output:    fmt.Sprintf("int"),
				Parameters: []*Parameter{
					{Name: "url", Value: "const QUrl &"},
					{Name: "uri", Value: "const char *"},
					{Name: "versionMajor", Value: "int"},
					{Name: "versionMinor", Value: "int"},
					{Name: "qmlName", Value: "const char *"},
				},
				Signature: "(const QUrl &url, const char *uri, int versionMajor, int versionMinor, const char *qmlName)",
			})

			//http://doc.qt.io/qt-5/qqmlengine.html#qmlRegisterType-2
			//int qmlRegisterType(const QUrl &url, const char *uri, int versionMajor, int versionMinor, const char *qmlName)
			c.Functions = append(c.Functions, &Function{
				Name:      "qmlRegisterType",
				Fullname:  fmt.Sprintf("%v::qmlRegisterType", c.Name),
				Access:    "public",
				Virtual:   "non",
				Meta:      PLAIN,
				NonMember: true,
				Static:    true,
				Output:    fmt.Sprintf("int"),
				Parameters: []*Parameter{
					{Name: "url", Value: "const QUrl &"},
					{Name: "uri", Value: "const char *"},
					{Name: "versionMajor", Value: "int"},
					{Name: "versionMinor", Value: "int"},
					{Name: "qmlName", Value: "const char *"},
				},
				Signature: "(const QUrl &url, const char *uri, int versionMajor, int versionMinor, const char *qmlName)",
			})
		}

	case "QAndroidJniEnvironment":
		{
			c.Functions = append(c.Functions, &Function{
				Name:       "ExceptionCheck",
				Fullname:   fmt.Sprintf("%v::ExceptionCheck", c.Name),
				Access:     "public",
				Virtual:    "non",
				Meta:       PLAIN,
				Static:     true,
				Output:     "bool",
				Parameters: []*Parameter{},
				Signature:  "()",
			})

			c.Functions = append(c.Functions, &Function{
				Name:       "ExceptionDescribe",
				Fullname:   fmt.Sprintf("%v::ExceptionDescribe", c.Name),
				Access:     "public",
				Virtual:    "non",
				Meta:       PLAIN,
				Static:     true,
				Output:     "void",
				Parameters: []*Parameter{},
				Signature:  "()",
			})

			c.Functions = append(c.Functions, &Function{
				Name:       "ExceptionClear",
				Fullname:   fmt.Sprintf("%v::ExceptionClear", c.Name),
				Access:     "public",
				Virtual:    "non",
				Meta:       PLAIN,
				Static:     true,
				Output:     "void",
				Parameters: []*Parameter{},
				Signature:  "()",
			})

			c.Functions = append(c.Functions, &Function{
				Name:       "ExceptionOccurred",
				Fullname:   fmt.Sprintf("%v::ExceptionOccurred", c.Name),
				Access:     "public",
				Virtual:    "non",
				Meta:       PLAIN,
				Static:     true,
				Output:     "void*",
				Parameters: []*Parameter{},
				Signature:  "()",
			})
		}

	case "QVideoFrame":
		{
			//QImage qt_imageFromVideoFrame(const QVideoFrame &frame)
			/* requires multimedia-private
			c.Functions = append(c.Functions, &Function{
				Name:      "qt_imageFromVideoFrame",
				Fullname:  fmt.Sprintf("%v::qt_imageFromVideoFrame", c.Name),
				Access:    "public",
				Virtual:   "non",
				Meta:      PLAIN,
				NonMember: true,
				Static:    true,
				Output:    fmt.Sprintf("QImage"),
				Parameters: []*Parameter{
					{Name: "frame", Value: "const QVideoFrame &"},
				},
				Signature: "(const QVideoFrame &frame)",
			})
			*/
		}
	}

	//TODO: make general
	if c.Name == "QQmlNetworkAccessManagerFactory" && !c.HasConstructor() {
		c.Functions = append(c.Functions, &Function{
			Name:       c.Name,
			Fullname:   fmt.Sprintf("%v::%v", c.Name, c.Name),
			Access:     "public",
			Virtual:    "non",
			Meta:       CONSTRUCTOR,
			Parameters: []*Parameter{},
			Signature:  "()",
		})
	}
}

func (c *Class) addVarAndPropFuncs() {
	for _, v := range c.Variables {
		c.Functions = append(c.Functions, v.varToFunc()...)
	}
	for _, p := range c.Properties {
		c.Functions = append(c.Functions, p.propToFunc(c)...)
	}
}

func (c *Class) addMocFuncs() {
	if c.Module != MOC {
		return
	}

	if c.HasFunctionWithName("qRegisterMetaType") {
		return
	}

	//http://doc.qt.io/qt-5/qmetatype.html#qRegisterMetaType-1
	//int qRegisterMetaType()
	qRF := &Function{
		Name:           "qRegisterMetaType",
		Fullname:       fmt.Sprintf("%v::qRegisterMetaType", c.Name),
		Access:         "public",
		Virtual:        "non",
		Meta:           PLAIN,
		NonMember:      true,
		NoMocDeduce:    true,
		Static:         true,
		Output:         fmt.Sprintf("int"),
		Parameters:     []*Parameter{},
		Signature:      "()",
		TemplateModeGo: fmt.Sprintf("%v*", c.Name),
	}
	c.Functions = append(c.Functions, qRF)

	//http://doc.qt.io/qt-5/qmetatype.html#qRegisterMetaType
	//int qRegisterMetaType(const char *typeName)
	qRF2 := *qRF
	qRF2.Overload = true
	qRF2.OverloadNumber = "2"
	qRF2.Parameters = []*Parameter{{Name: "typeName", Value: "const char *"}}
	qRF2.Signature = "(const char *typeName)"
	c.Functions = append(c.Functions, &qRF2)

	if c.IsSubClassOf("QCoreApplication") {
		return
	}

	//http://doc.qt.io/qt-5/qqmlengine.html#qmlRegisterType
	//int qmlRegisterType()
	qmlF := &Function{
		Name:           "qmlRegisterType",
		Fullname:       fmt.Sprintf("%v::qmlRegisterType", c.Name),
		Access:         "public",
		Virtual:        "non",
		Meta:           PLAIN,
		NonMember:      true,
		NoMocDeduce:    true,
		Static:         true,
		Output:         fmt.Sprintf("int"),
		Parameters:     []*Parameter{},
		Signature:      "()",
		TemplateModeGo: fmt.Sprintf("%v", c.Name),
	}
	c.Functions = append(c.Functions, qmlF)

	//int qmlRegisterType(const char *uri, int versionMajor, int versionMinor, const char *qmlName)
	qmlF2 := *qmlF
	qmlF2.Overload = true
	qmlF2.OverloadNumber = "2"
	qmlF2.Parameters = []*Parameter{
		{Name: "uri", Value: "const char *"},
		{Name: "versionMajor", Value: "int"},
		{Name: "versionMinor", Value: "int"},
		{Name: "qmlName", Value: "const char *"},
	}
	qmlF2.Signature = "(const char *uri, int versionMajor, int versionMinor, const char *qmlName)"
	c.Functions = append(c.Functions, &qmlF2)
}
