package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func (f *Function) fix() {
	f.fixGeneral()
	f.fixGeneral_Version()

	f.fixOverload()
	f.fixOverload_Version()

	//f.fixGeneric()
}

func (f *Function) fixGeneral() {

	if utils.QT_VERSION_NUM() >= 5110 {
		for _, p := range f.Parameters {
			if p.Value == "" {
				p.Value = p.ValueNew
			}
			if p.Default == "nullptr" {
				p.Default = "Q_NULLPTR"
			}
		}
	}

	//linux fixes

	if f.Fullname == "QThread::start" {
		f.Parameters = make([]*Parameter, 0)
	}

	//virtual fixes

	if f.Virtual == "virtual" {
		f.Virtual = IMPURE
	}

	if f.Virtual == IMPURE || f.Virtual == PURE ||
		f.Meta == SIGNAL || f.Meta == SLOT {

		f.Static = false
	}

	//constructor fixes

	if f.Meta == COPY_CONSTRUCTOR || f.Meta == MOVE_CONSTRUCTOR {
		f.Meta = CONSTRUCTOR
	}

	var class, ok = f.Class()
	if !ok || !class.isSubClass() {
		return
	}

	if f.Meta == CONSTRUCTOR {
		f.Status = "active"
		f.Access = "public"
	}
}

func (f *Function) fixGeneral_AfterClasses() {
	if f.Name != "open" && f.Name != "setGeometry" && f.Name != "setScxmlEvent" && //TODO:
		!f.Static && f.Virtual == "non" && f.Meta == PLAIN && f.IsDerivedFromVirtual() {
		f.Virtual = IMPURE
	}
}

func (f *Function) fixGeneral_Version() {
	switch f.Fullname {
	case "QScxmlCppDataModel::setScxmlEvent":
		{
			f.Virtual = "non"
		}

	case "QGraphicsObject::z", "QGraphicsObject::setZ":
		{
			f.Name = func() string {
				if f.Name == "setZ" {
					return "setZValue"
				}
				return "zValue"
			}()
			f.Fullname = fmt.Sprintf("%v::%v", f.ClassName(), f.Name)
		}

	case "QGraphicsObject::effect", "QGraphicsObject::setEffect":
		{
			f.Name = func() string {
				if f.Name == "setEffect" {
					return "setGraphicsEffect"
				}
				return "graphicsEffect"
			}()
			f.Fullname = fmt.Sprintf("%v::%v", f.ClassName(), f.Name)
		}

	case "QOperatingSystemVersion::QOperatingSystemVersion":
		{
			if len(f.Parameters) > 1 {
				f.OverloadNumber = "2"
				f.Overload = true
			}
		}
	case "QAndroidJniEnvironment::javaVM":
		{
			f.Output = strings.Replace(f.Output, "int *", "JavaVM *", -1)
			f.Signature = strings.Replace(f.Signature, "int *", "JavaVM *", -1)
		}
	case "QAndroidJniObject::fromLocalRef":
		{
			f.Parameters[0].Value = strings.Replace(f.Parameters[0].Value, "int", "jobject", -1)
			f.Signature = strings.Replace(f.Signature, "int", "jobject", -1)
		}
	case "QAndroidJniObject::QAndroidJniObject":
		{
			if strings.HasSuffix(f.Href, "5") {
				f.Parameters[0].Value = strings.Replace(f.Parameters[0].Value, "int", "jobject", -1)
				f.Signature = strings.Replace(f.Signature, "int", "jobject", -1)
			}
		}
	case "QImage::QImage":
		{
			for i := len(f.Parameters) - 1; i >= 0; i-- {
				if strings.HasPrefix(f.Parameters[i].Name, "cleanup") {
					f.Parameters = append(f.Parameters[:i], f.Parameters[i+1:]...)
				}
			}
		}
	case "QWebEnginePage::certificateError":
		{
			f.Signature = strings.Replace(f.Signature, "const int", "const QWebEngineCertificateError", -1)
			f.Parameters[0].Value = strings.Replace(f.Parameters[0].Value, "const int", "const QWebEngineCertificateError", -1)
		}
	case "QTextLayout::formats":
		{
			f.Output = "QVector<QTextLayout::FormatRange>"
		}
	case "QWebEnginePage::print":
		{
			f.Parameters = append(f.Parameters[:len(f.Parameters)-1], f.Parameters[len(f.Parameters):]...)
		}
	}
}

func (f *Function) fixOverload() {

	if strings.Contains(f.Href, "-") {
		tmp, err := strconv.Atoi(strings.Split(f.Href, "-")[1])
		if err == nil && tmp > 0 {
			f.Overload = true
			f.OverloadNumber = strconv.Itoa(tmp + 1)
		}
	}

	if f.OverloadNumber == "1" {
		f.OverloadNumber = "2"
	}

	if f.OverloadNumber != "0" {
		return
	}

	f.Overload = false
	f.OverloadNumber = ""
}

func (f *Function) fixOverload_Version() {
	switch f.Fullname {
	case "QGraphicsDropShadowEffect::setOffset", "QGraphicsScene::setSceneRect",
		"QGraphicsView::setSceneRect", "QQuickItem::setFocus",
		"QAccessibleWidget::setText", "QSvgGenerator::setViewBox",
		"QSvgRenderer::setViewBox":
		{
			var class, ok = f.Class()
			if !ok {
				return
			}

			var count int
			for _, sf := range class.Functions {
				if sf.Fullname != f.Fullname {
					continue
				}

				if sf.Signature != f.Signature {
					count++
					continue
				}

				break
			}
			if count == 0 {
				return
			}

			f.Overload = true
			f.OverloadNumber = strconv.Itoa(count + 1)
		}
	}
}

func (f *Function) fixGeneric() {
	f.fixGenericOutput()
	f.fixGenericInput()
}

func (f *Function) fixGenericOutput() {

	switch CleanValue(f.Output) {
	case "QVariantHash":
		{
			f.Output = "QHash<QString, QVariant>"
		}

	case "QVariantMap":
		{
			f.Output = "QMap<QString, QVariant>"
		}

	case "QJSValueList":
		{
			f.Output = "QList<QJSValue>"
		}

	case "QModelIndexList":
		{
			f.Output = "QList<QModelIndex>"
		}

	case "QVariantList":
		{
			f.Output = "QList<QVariant>"
		}

	case "QObjectList":
		{
			f.Output = "QList<QObject *>"
		}

	case "QMediaResourceList":
		{
			f.Output = "QList<QMediaResource>"
		}

	case "QFileInfoList":
		{
			f.Output = "QList<QFileInfo>"
		}

	case "QWidgetList":
		{
			f.Output = "QList<QWidget *>"
		}

	case "QCameraFocusZoneList":
		{
			//f.Output = "QList<QCameraFocusZone *>" //TODO: uncomment
		}

	case "QList<T>":
		{
			f.TemplateModeGo = "QObject*"
			f.Output = "QList<QObject*>"
		}

	case "QVector<T>":
		{
			f.Output = "QList<QObject*>"
		}

	case "T":
		{
			switch className := f.ClassName(); className {
			case "QObject", "QMediaService":
				{
					f.TemplateModeGo = fmt.Sprintf("%v*", className)
					f.Output = fmt.Sprintf("%v*", className)
				}
			}
		}
	}
}

func (f *Function) fixGenericInput() {
	var skipOG bool
	for _, p := range f.Parameters {
		if strings.HasPrefix(p.Value, "[]") || strings.HasPrefix(p.Value, "map[") {
			skipOG = true
			break
		}
	}

	if len(f.OgParameters) == 0 && !skipOG {
		for _, p := range f.Parameters {
			if p.Default == "..." {
				switch f.Name {
				case "QPaintEngine":
					p.Default = "PaintEngineFeatures()"
				case "QLayoutItem":
					p.Default = "Qt::Alignment()"
				case "QBluetoothLocalDevice", "QWebEngineUrlRequestInterceptor":
					p.Default = "Q_NULLPTR"
				case "QMediaPlayer", "QQuickImageProvider":
					p.Default = "Flags()"
				case "QModbusRequest", "QModbusResponse":
					p.Default = "QByteArray()"
				case "QGeoServiceProvider":
					p.Default = "QVariantMap()"
				default:
					p.Default = "Qt::WindowFlags()"
				}
			}
			if strings.HasPrefix(p.Default, "DECLARE_READING") {
				p.Default = ""
			}
			f.OgParameters = append(f.OgParameters, *p)
		}
	}

	for _, p := range f.Parameters {
		switch CleanValue(p.Value) {
		case "QVariantHash":
			{
				p.Value = "QHash<QString, QVariant>"
			}

		case "QVariantMap":
			{
				p.Value = "QMap<QString, QVariant>"
			}

		case "QJSValueList":
			{
				p.Value = "QList<QJSValue>"
			}

		case "QModelIndexList":
			{
				p.Value = "QList<QModelIndex>"
			}

		case "QVariantList":
			{
				p.Value = "QList<QVariant>"
			}

		case "QObjectList":
			{
				p.Value = "QList<QObject *>"
			}

		case "QMediaResourceList":
			{
				p.Value = "QList<QMediaResource>"
			}

		case "QFileInfoList":
			{
				p.Value = "QList<QFileInfo>"
			}

		case "QWidgetList":
			{
				p.Value = "QList<QWidget *>"
			}

		case "QCameraFocusZoneList":
			{
				//p.Value = "QList<QCameraFocusZone *>" //TODO: uncomment
			}

		case "QList<T>":
			{
				p.Value = "QList<QObject*>"
			}

		case "QVector<T>":
			{
				p.Value = "QList<QObject*>"
			}

		case "T":
			{
				switch className := f.ClassName(); className {
				case "QObject", "QMediaService":
					{
						p.Value = fmt.Sprintf("%v*", className)
					}
				}
			}
		}
	}
}

func (c *Class) FixGenericHelper() {
	for _, cn := range append([]string{c.Name}, c.GetAllBases()...) {
		var rec bool

		var class, e = State.ClassMap[cn]
		if !e {
			continue
		}
		for _, f := range class.Functions {
			//TODO: needed because there could be unfixed subclasses; delay this to later (also check for GetAllBases or GetBases in parser)
			f.fixGeneral()
			f.fixGeneric()

			if IsPackedList(CleanValue(f.Output)) || IsPackedMap(CleanValue(f.Output)) {
				if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_atList", f.Name), f.OverloadNumber) {
					var key, output, isMap = func() (string, string, bool) {
						if IsPackedList(CleanValue(f.Output)) {
							return "int", strings.Split(strings.Split(f.Output, "<")[1], ">")[0], false
						}
						var key, value = UnpackedMapDirty(CleanValue(f.Output))
						return key, value, true
					}()

					params := []*Parameter{{Name: "i", Value: key}}
					if isMap {
						params[0].Name = "v"
						params = append(params, &Parameter{Name: "i", Value: "int"})
						params[0].PureGoType = func() string {
							if strings.Contains(f.PureGoOutput, "[error]") {
								return "error"
							}
							return ""
						}()
					} else {
						params[0].PureGoType = func() string {
							if strings.Contains(f.PureGoOutput, "]error") {
								return "error"
							}
							return ""
						}()
					}

					c.Functions = append(c.Functions, &Function{
						Name:           fmt.Sprintf("__%v_atList", f.Name),
						Fullname:       fmt.Sprintf("%v::__%v_atList", c.Name, f.Name),
						Access:         "public",
						Virtual:        "non",
						Meta:           PLAIN,
						Output:         output,
						Parameters:     params,
						Signature:      "()",
						Container:      strings.Split(f.Output, "<")[0],
						OverloadNumber: f.OverloadNumber,
						Overload:       f.Overload,
						NoMocDeduce:    true,
						PureGoOutput: func() string {
							if strings.Contains(f.PureGoOutput, "]error") {
								return "error"
							}
							return ""
						}(),
						IsMap: isMap,
					})
				}

				if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_setList", f.Name), f.OverloadNumber) {
					var params = func() []*Parameter {
						if IsPackedList(CleanValue(f.Output)) {
							return []*Parameter{{Name: "i", Value: strings.Split(strings.Split(f.Output, "<")[1], ">")[0],
								PureGoType: func() string {
									if strings.Contains(f.PureGoOutput, "]error") {
										return "error"
									}
									return ""
								}()}}
						}
						var key, value = UnpackedMapDirty(CleanValue(f.Output))
						return []*Parameter{{Name: "key", Value: key,
							PureGoType: func() string {
								if strings.Contains(f.PureGoOutput, "[error]") {
									return "error"
								}
								return ""
							}()}, {Name: "i", Value: value,
							PureGoType: func() string {
								if strings.Contains(f.PureGoOutput, "]error") {
									return "error"
								}
								return ""
							}()}}
					}()
					c.Functions = append(c.Functions, &Function{
						Name:           fmt.Sprintf("__%v_setList", f.Name),
						Fullname:       fmt.Sprintf("%v::__%v_setList", c.Name, f.Name),
						Access:         "public",
						Virtual:        "non",
						Meta:           PLAIN,
						Output:         "void",
						Parameters:     params,
						Signature:      "()",
						Container:      strings.Split(f.Output, "<")[0],
						OverloadNumber: f.OverloadNumber,
						Overload:       f.Overload,
						NoMocDeduce:    true,
					})
				}

				if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_newList", f.Name), f.OverloadNumber) {
					c.Functions = append(c.Functions, &Function{
						Name:           fmt.Sprintf("__%v_newList", f.Name),
						Fullname:       fmt.Sprintf("%v::__%v_newList", c.Name, f.Name),
						Access:         "public",
						Virtual:        "non",
						Meta:           PLAIN,
						Output:         "void *",
						Signature:      "()",
						Container:      f.Output,
						OverloadNumber: f.OverloadNumber,
						Overload:       f.Overload,
						NoMocDeduce:    true,
					})
				}

				if IsPackedMap(CleanValue(f.Output)) {
					if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_keyList", f.Name), f.OverloadNumber) {
						var keyType, _ = UnpackedMapDirty(CleanValue(f.Output))
						c.Functions = append(c.Functions, &Function{
							Name:           fmt.Sprintf("__%v_keyList", f.Name),
							Fullname:       fmt.Sprintf("%v::__%v_keyList", c.Name, f.Name),
							Access:         "public",
							Virtual:        "non",
							Meta:           PLAIN,
							Output:         fmt.Sprintf("QList<%v>", keyType),
							Signature:      "()",
							OverloadNumber: f.OverloadNumber,
							Overload:       f.Overload,
							NoMocDeduce:    true,
							PureGoOutput: func() string {
								if strings.Contains(f.PureGoOutput, "[error]") {
									return "[]error"
								}
								return ""
							}(),
							Container: CleanValue(f.Output),
						})
						rec = true
					}
				}
			}

			for _, p := range f.Parameters {
				if IsPackedList(CleanValue(p.Value)) || IsPackedMap(CleanValue(p.Value)) {
					if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_%v_atList", f.Name, p.Name), f.OverloadNumber) {
						var key, output, isMap = func() (string, string, bool) {
							if IsPackedList(CleanValue(p.Value)) {
								return "int", strings.Split(strings.Split(p.Value, "<")[1], ">")[0], false
							}
							var key, value = UnpackedMapDirty(CleanValue(p.Value))
							return key, value, true
						}()

						params := []*Parameter{{Name: "i", Value: key}}
						if isMap {
							params[0].Name = "v"
							params = append(params, &Parameter{Name: "i", Value: "int"})
							params[0].PureGoType = func() string {
								if strings.Contains(p.PureGoType, "[error]") {
									return "error"
								}
								return ""
							}()
						} else {
							params[0].PureGoType = func() string {
								if strings.Contains(p.PureGoType, "]error") {
									return "error"
								}
								return ""
							}()
						}

						c.Functions = append(c.Functions, &Function{
							Name:           fmt.Sprintf("__%v_%v_atList", f.Name, p.Name),
							Fullname:       fmt.Sprintf("%v::__%v_%v_atList", c.Name, f.Name, p.Name),
							Access:         "public",
							Virtual:        "non",
							Meta:           PLAIN,
							Output:         output,
							Parameters:     params,
							Signature:      "()",
							Container:      strings.Split(p.Value, "<")[0],
							OverloadNumber: f.OverloadNumber,
							Overload:       f.Overload,
							NoMocDeduce:    true,
							PureGoOutput: func() string {
								if strings.Contains(p.PureGoType, "]error") {
									return "error"
								}
								return ""
							}(),
							IsMap: isMap,
						})
					}

					if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_%v_setList", f.Name, p.Name), f.OverloadNumber) {
						var params = func() []*Parameter {
							if IsPackedList(CleanValue(p.Value)) {
								return []*Parameter{{Name: "i", Value: strings.Split(strings.Split(p.Value, "<")[1], ">")[0],
									PureGoType: func() string {
										if strings.Contains(p.PureGoType, "]error") {
											return "error"
										}
										return ""
									}()}}
							}
							var key, value = UnpackedMapDirty(CleanValue(p.Value))
							return []*Parameter{{Name: "key", Value: key,
								PureGoType: func() string {
									if strings.Contains(p.PureGoType, "[error]") {
										return "error"
									}
									return ""
								}()}, {Name: "i", Value: value,
								PureGoType: func() string {
									if strings.Contains(p.PureGoType, "]error") {
										return "error"
									}
									return ""
								}()}}
						}()
						c.Functions = append(c.Functions, &Function{
							Name:           fmt.Sprintf("__%v_%v_setList", f.Name, p.Name),
							Fullname:       fmt.Sprintf("%v::__%v_%v_setList", c.Name, f.Name, p.Name),
							Access:         "public",
							Virtual:        "non",
							Meta:           PLAIN,
							Output:         "void",
							Parameters:     params,
							Signature:      "()",
							Container:      strings.Split(p.Value, "<")[0],
							OverloadNumber: f.OverloadNumber,
							Overload:       f.Overload,
							NoMocDeduce:    true,
						})
					}

					if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_%v_newList", f.Name, p.Name), f.OverloadNumber) {
						c.Functions = append(c.Functions, &Function{
							Name:           fmt.Sprintf("__%v_%v_newList", f.Name, p.Name),
							Fullname:       fmt.Sprintf("%v::__%v_%v_newList", c.Name, f.Name, p.Name),
							Access:         "public",
							Virtual:        "non",
							Meta:           PLAIN,
							Output:         "void *",
							Signature:      "()",
							Container:      p.Value,
							OverloadNumber: f.OverloadNumber,
							Overload:       f.Overload,
							NoMocDeduce:    true,
						})
					}

					if IsPackedMap(CleanValue(p.Value)) {
						if !c.HasFunctionWithNameAndOverloadNumber(fmt.Sprintf("__%v_%v_keyList", f.Name, p.Name), f.OverloadNumber) {
							var keyType, _ = UnpackedMapDirty(CleanValue(p.Value))
							c.Functions = append(c.Functions, &Function{
								Name:           fmt.Sprintf("__%v_%v_keyList", f.Name, p.Name),
								Fullname:       fmt.Sprintf("%v::__%v_%v_keyList", c.Name, f.Name, p.Name),
								Access:         "public",
								Virtual:        "non",
								Meta:           PLAIN,
								Output:         fmt.Sprintf("QList<%v>", keyType),
								Signature:      "()",
								OverloadNumber: f.OverloadNumber,
								Overload:       f.Overload,
								NoMocDeduce:    true,
								PureGoOutput: func() string {
									if strings.Contains(p.PureGoType, "[error]") {
										return "[]error"
									}
									return ""
								}(),
								Container: CleanValue(p.Value),
							})
							rec = true
						}
					}
				}
			}
		}
		if rec {
			c.FixGenericHelper()
		}
	}
}
