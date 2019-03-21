package templater

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/converter"
	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
)

func cppFunctionCallback(function *parser.Function) string {
	var output = fmt.Sprintf("%v { %v };", cppFunctionCallbackHeader(function), cppFunctionCallbackBody(function))
	if function.IsSupported() {
		return cppFunctionCallbackWithGuards(function, output)
	}
	return ""
}

func cppFunctionCallbackWithGuards(function *parser.Function, output string) string {
	switch {
	case
		function.Fullname == "QProcess::nativeArguments", function.Fullname == "QProcess::setNativeArguments",
		function.Fullname == "QAbstractEventDispatcher::registerEventNotifier", function.Fullname == "QAbstractEventDispatcher::unregisterEventNotifier":
		{
			return fmt.Sprintf("#ifdef Q_OS_WIN\n\t\t%v\n\t#endif", output)
		}

	case
		function.Fullname == "QSensorGesture::detected":
		{
			return fmt.Sprintf("#ifdef Q_QDOC\n\t\t%v\n\t#endif", output)
		}
	}

	return output
}

func cppFunctionCallbackHeader(function *parser.Function) string {
	return fmt.Sprintf("%v %v%v(%v)%v",

		func() string {
			var c, _ = function.Class()
			if parser.IsPackedMap(function.Output) && c.Module == parser.MOC && function.IsMocFunction {
				var tHash = sha1.New()
				tHash.Write([]byte(function.Output))
				return fmt.Sprintf("type%v", hex.EncodeToString(tHash.Sum(nil)[:3]))
			}
			return function.Output
		}(),

		func() string {
			if function.Meta == parser.SIGNAL {
				return fmt.Sprintf("Signal_%v", strings.Title(function.Name))
			}
			var c, _ = function.Class()
			if strings.HasPrefix(function.Name, parser.TILDE) && c.Module != parser.MOC {
				return strings.Replace(function.Name, parser.TILDE, fmt.Sprintf("%vMy", parser.TILDE), -1)
			}
			return function.Name
		}(),

		func() string {
			if function.Meta == parser.SIGNAL {
				return function.OverloadNumber
			}
			return ""
		}(),

		converter.CppInputParametersForCallbackHeader(function),

		func() string {
			if strings.Contains(function.Signature, ") const") {
				return " const"
			}
			return ""
		}(),
	)
}

func cppFunctionCallbackBody(function *parser.Function) string {
	out := fmt.Sprintf("%v%v%v;",

		converter.CppInputParametersForCallbackBodyPrePack(function),

		func() string {
			if converter.CppHeaderOutput(function) != parser.VOID {
				return "return "
			}
			return ""
		}(),

		func() string {
			var output string
			if UseJs() {
				output = fmt.Sprintf("emscripten::val::global(\"Module\").call<%v>(\"_callback%v_%v%v\", %v)", converter.CppOutputTemplateJS(function), function.ClassName(), strings.Replace(strings.Title(function.Name), parser.TILDE, "Destroy", -1), function.OverloadNumber, converter.CppInputParametersForCallbackBody(function))
			} else {
				output = fmt.Sprintf("callback%v_%v%v(%v)", function.ClassName(), strings.Replace(strings.Title(function.Name), parser.TILDE, "Destroy", -1), function.OverloadNumber, converter.CppInputParametersForCallbackBody(function))
			}
			if converter.CppHeaderOutput(function) != parser.VOID {
				output = converter.CppInput(output, function.Output, function)
			}
			if UseJs() {
				output = strings.Replace(output, "static_cast", "reinterpret_cast", -1)
				output = strings.Replace(output, "enum_cast", "static_cast", -1)
			}
			return output
		}(),
	)
	return out
}

func cppFunction(function *parser.Function) string {
	var output = fmt.Sprintf("%v\n{\n%v\n}", cppFunctionHeader(function), cppFunctionUnused(function, cppFunctionBodyWithGuards(function)))
	if UseJs() {
		if !strings.Contains(output, "_Packed") && !strings.Contains(output, "emscripten::val") {
			output = strings.Replace(output, converter.CppHeaderName(function), "_KEEPALIVE_"+converter.CppHeaderName(function), -1)
			output = "EMSCRIPTEN_KEEPALIVE\n" + output
		} else {
			output = strings.Replace(output, "static_cast", "reinterpret_cast", -1)
			exportedFunctions = append(exportedFunctions, converter.CppHeaderName(function))
		}
		output = strings.Replace(output, "enum_cast", "static_cast", -1)
	}
	if function.IsSupported() {
		return output
	}
	return ""
}

func cppFunctionHeader(function *parser.Function) string {
	var output = fmt.Sprintf("%v %v(%v)", converter.CppHeaderOutput(function), converter.CppHeaderName(function), converter.CppHeaderInput(function))
	if UseJs() {
		if strings.Contains(output, "_Packed") || strings.Contains(output, "emscripten::val") {
			output = strings.Replace(output, "void*", "uintptr_t", -1)
			function.BoundByEmscripten = true
		}
	}
	if function.IsSupported() {
		return output
	}
	return ""
}

//TODO:
func cppFunctionUnused(function *parser.Function, body string) string {

	var tmp = make([]string, 0)
	if !(function.Static || function.Meta == parser.CONSTRUCTOR) {
		tmp = append(tmp, "ptr")
	}
	if function.Meta != parser.SIGNAL {
		for _, p := range function.Parameters {
			tmp = append(tmp, parser.CleanName(p.Name, p.Value))
		}
	}

	bb := new(bytes.Buffer)
	defer bb.Reset()
	for _, p := range tmp {
		if !strings.Contains(body, p) {
			fmt.Fprintf(bb, "\tQ_UNUSED(%v);\n", p)
		}
	}
	bb.WriteString(body)
	return bb.String()
}

func cppFunctionBodyWithGuards(function *parser.Function) string {

	if function.Default {
		switch {
		case
			strings.HasPrefix(function.ClassName(), "QMac") && !strings.HasPrefix(parser.State.ClassMap[function.ClassName()].Module, "QtMac"):
			{
				return fmt.Sprintf("#ifdef Q_OS_OSX\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}
		}
	} else {
		switch {
		case
			function.Fullname == "QMenu::setAsDockMenu", function.Fullname == "QSysInfo::macVersion", function.Fullname == "QSysInfo::MacintoshVersion":
			{
				return fmt.Sprintf("#ifdef Q_OS_OSX\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}

		case
			function.Fullname == "QProcess::nativeArguments", function.Fullname == "QProcess::setNativeArguments",
			function.Fullname == "QAbstractEventDispatcher::registerEventNotifier", function.Fullname == "QAbstractEventDispatcher::unregisterEventNotifier",
			function.Fullname == "QSysInfo::windowsVersion":
			{
				return fmt.Sprintf("#ifdef Q_OS_WIN\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}

		case
			function.Fullname == "QScreen::model":
			{
				return fmt.Sprintf("#ifndef Q_OS_WIN\n%v\n#endif", cppFunctionBody(function))
			}

		case
			function.Fullname == "QApplication::navigationMode", function.Fullname == "QApplication::setNavigationMode",
			function.Fullname == "QWidget::hasEditFocus", function.Fullname == "QWidget::setEditFocus":
			{
				return fmt.Sprintf("#ifdef QT_KEYPAD_NAVIGATION\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}

		case
			function.Fullname == "QMenuBar::defaultAction", function.Fullname == "QMenuBar::setDefaultAction":
			{
				return fmt.Sprintf("#ifdef Q_OS_WINCE\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}

		case
			function.Fullname == "QWidget::setupUi", function.Fullname == "QSensorGesture::detected":
			{
				return fmt.Sprintf("#ifdef Q_QDOC\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}

		case
			function.Fullname == "QTextDocument::print", function.Fullname == "QPlainTextEdit::print", function.Fullname == "QTextEdit::print":
			{
				return fmt.Sprintf("#ifndef Q_OS_IOS\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}

		case function.Name == "qmlRegisterType" && function.TemplateModeGo != "":
			{
				return fmt.Sprintf("#ifdef QT_QML_LIB\n%v%v\n#endif", cppFunctionBody(function), cppFunctionBodyFailed(function))
			}
		}
	}

	return cppFunctionBody(function)
}

func cppFunctionBodyFailed(function *parser.Function) string {
	if converter.CppHeaderOutput(function) != parser.VOID {
		return fmt.Sprintf("\n#else\n\treturn %v;", converter.CppOutputParametersFailed(function))
	}
	return ""
}

func cppFunctionBody(function *parser.Function) string {

	var fakeDefault bool
	if UseJs() && function.Meta == parser.SLOT {
		defer func() { function.Meta = parser.SLOT; function.Default = false }()
		function.Meta = parser.PLAIN
		function.Default = true
		fakeDefault = true
	}

	var polyinputs, polyName = function.PossiblePolymorphicDerivations(false)

	var polyinputsSelf []string

	var c, _ = function.Class()
	if function.Default && c.Module != parser.MOC {
		polyinputsSelf, _ = function.PossibleDerivationsReversedAndRemovedPure(true)
	} else {
		polyinputsSelf, _ = function.PossiblePolymorphicDerivations(true)
	}

	if c.Module == parser.MOC {
		var fc, ok = function.Class()
		if !ok {
			return ""
		}

		for _, bcn := range append([]string{function.ClassName()}, fc.GetBases()...) {
			var bc, ok = parser.State.ClassMap[bcn]
			if !ok {
				continue
			}
			var f = *function
			f.Fullname = fmt.Sprintf("%v::%v", bcn, f.Name)

			var ff = bc.GetFunction(f.Name)
			for _, fb := range parser.IsBlockedDefault() {
				if f.Fullname == fb || (ff != nil && ff.Virtual == parser.PURE && (bc.Module != parser.MOC && bc.Pkg == "")) {
					return ""
				}
			}
		}
	}

	if (len(polyinputsSelf) == 0 && len(polyinputs) == 0) ||
		function.SignalMode == parser.CONNECT || function.SignalMode == parser.DISCONNECT ||
		(len(polyinputsSelf) != 0 && function.Meta == parser.CONSTRUCTOR) || (function.Meta == parser.DESTRUCTOR || strings.HasPrefix(function.Name, parser.TILDE)) {
		out := cppFunctionBodyInternal(function)
		if fakeDefault {
			out = strings.Replace(out, "->"+parser.State.ClassMap[function.ClassName()].GetBases()[0]+"::", "->", -1)
		}
		return out
	}

	bb := new(bytes.Buffer)
	defer bb.Reset()

	bb.WriteString("\t")

	var deduce = func(input []string, polyName string, inner bool, body string) string {
		bbi := new(bytes.Buffer)
		defer bbi.Reset()
		for _, polyType := range input {
			if polyType == "QObject" || polyType == input[len(input)-1] {
				continue
			}

			if strings.HasPrefix(polyType, "QMac") {
				fmt.Fprint(bbi, "\n\t#ifdef Q_OS_OSX\n\t\t")
			}

			base := input[len(input)-1]
			if parser.State.ClassMap[polyType].IsSubClassOfQObject() {
				base = "QObject"
			}

			fmt.Fprintf(bbi, "if (dynamic_cast<%v*>(static_cast<%v*>(%v))) {\n", polyType, base, polyName)

			if strings.HasPrefix(polyType, "QMac") {
				fmt.Fprint(bbi, "\t#else\n\t\tif (false) {\n\t#endif\n")
			}

			fmt.Fprintf(bbi, "\t%v\n", func() string {
				var ibody string
				if function.Default && polyName == "ptr" {
					if fakeDefault && !inner {
						ibody = strings.Replace(body, "static_cast<"+input[len(input)-1]+"*>("+polyName+")->"+input[len(input)-1]+"::", "static_cast<My"+polyType+"*>("+polyName+")->My"+polyType+"::", -1)

						//TODO: only temporary until invoke works ->
						for _, s := range append(parser.State.ClassMap["QAbstractItemView"].GetAllDerivations(), "QAbstractItemView") {
							if strings.Contains(ibody, "static_cast<My"+s+"*>(ptr)->My"+s+"::update()") {
								ibody = ""
								break
							}
						}
						//<-
					} else {
						ibody = strings.Replace(body, "static_cast<"+input[len(input)-1]+"*>("+polyName+")->"+input[len(input)-1]+"::", "static_cast<"+polyType+"*>("+polyName+")->"+polyType+"::", -1)
					}
				} else {
					ibody = strings.Replace(body, "static_cast<"+input[len(input)-1]+"*>("+polyName+")", "static_cast<"+polyType+"*>("+polyName+")", -1)
				}

				if strings.HasPrefix(polyType, "QMac") {
					ibody = fmt.Sprintf("#ifdef Q_OS_OSX\n\t%v\n\t#endif", ibody)
				}

				if inner {
					return ibody
				}
				if strings.Count(ibody, "\n") > 1 {
					return "\t" + strings.Replace(ibody, "\n", "\n\t", -1)
				}
				return ibody
			}())
			fmt.Fprint(bbi, "\t} else ")
		}

		if len(input) > 0 {
			var _, ok = parser.State.ClassMap[input[len(input)-1]]
			if ok {
				var f = *function
				f.Fullname = fmt.Sprintf("%v::%v", input[len(input)-1], f.Name)

				for _, fb := range parser.IsBlockedDefault() {
					if f.Fullname == fb {
						body = ""
					}
				}
			}
		}

		if bbi.String() == "" {
			return body
		}
		fmt.Fprintf(bbi, "{\n\t%v\n\t}", func() string {
			if fakeDefault && !inner {
				body = strings.Replace(body, "static_cast<"+function.ClassName()+"*>(ptr)->"+function.ClassName()+"::", "static_cast<My"+function.ClassName()+"*>(ptr)->My"+function.ClassName()+"::", -1)
			}
			if strings.Count(body, "\n") > 1 {
				return "\t" + strings.Replace(body, "\n", "\n\t", -1)
			}
			return body
		}())
		return bbi.String()
	}

	if function.Static {
		fmt.Fprint(bb, deduce(polyinputs, polyName, true, cppFunctionBodyInternal(function)))
	} else if function.Meta == parser.GETTER || function.Meta == parser.SETTER || function.Meta == parser.SLOT {
		fmt.Fprint(bb, deduce(polyinputsSelf, "ptr", false, cppFunctionBodyInternal(function)))
	} else {
		fmt.Fprint(bb, deduce(polyinputsSelf, "ptr", false, deduce(polyinputs, polyName, true, cppFunctionBodyInternal(function))))
	}

	return bb.String()
}

func cppFunctionBodyInternal(function *parser.Function) string {

	switch function.Meta {
	case parser.CONSTRUCTOR:
		{
			return fmt.Sprintf("%v\treturn %vnew %v%v(%v)%v;",

				func() string {
					if parser.State.ClassMap[function.ClassName()].IsSubClassOf("QCoreApplication") || function.Name == "QAndroidService" {
						if UseJs() {
							return `	static int argcs = argc;
	static char** argvs = static_cast<char**>(malloc(argcs * sizeof(char*)));

	QList<QByteArray> aList = QString::fromStdString(argv["data"].as<std::string>()).toUtf8().split('|');
	for (int i = 0; i < argcs; i++)
		argvs[i] = (new QByteArray(aList.at(i)))->data();

`
						}
						return `	static int argcs = argc;
	static char** argvs = static_cast<char**>(malloc(argcs * sizeof(char*)));

	QList<QByteArray> aList = QByteArray(argv).split('|');
	for (int i = 0; i < argcs; i++)
		argvs[i] = (new QByteArray(aList.at(i)))->data();

`
					}
					return ""
				}(),

				func() string {
					if UseJs() && function.BoundByEmscripten {
						return "reinterpret_cast<uintptr_t>("
					}
					return ""
				}(),

				func() string {
					var class, _ = function.Class()
					if class.Module != parser.MOC {
						if class.HasCallbackFunctions() {
							return "My"
						}
					}
					return ""
				}(),

				func() string {
					if c := parser.State.ClassMap[function.ClassName()]; c != nil && c.Fullname != "" {
						return c.Fullname
					}
					return function.ClassName()
				}(),

				converter.CppInputParameters(function),

				func() string {
					if UseJs() && function.BoundByEmscripten {
						return ")"
					}
					return ""
				}(),
			)
		}

	case parser.SLOT:
		{
			var (
				functionOutputType string
				bb                 = new(bytes.Buffer)
			)
			defer bb.Reset()

			if reg := converter.CppRegisterMetaType(function); reg != "" {
				bb.WriteString(reg + "\n")
			}

			fmt.Fprint(bb, "\t")

			if converter.CppHeaderOutput(function) != parser.VOID {
				functionOutputType = converter.CppInputParametersForSlotArguments(function, &parser.Parameter{Name: "returnArg", Value: function.Output})
				if function.Output != "void*" && !parser.State.ClassMap[strings.TrimSuffix(functionOutputType, "*")].IsSubClassOfQObject() {
					functionOutputType = strings.TrimSuffix(functionOutputType, "*")
				}
				fmt.Fprintf(bb, "%v returnArg;\n\t", functionOutputType)
			}

			fmt.Fprintf(bb, "QMetaObject::invokeMethod(static_cast<%v*>(ptr), \"%v\"%v%v);",

				function.ClassName(),

				function.Name,

				func() string {
					if converter.CppHeaderOutput(function) != parser.VOID {

						if c, _ := function.Class(); c.Module == parser.MOC && parser.IsPackedMap(function.Output) && function.IsMocFunction {
							var tHash = sha1.New()
							tHash.Write([]byte(function.Output))
							return fmt.Sprintf(", Q_RETURN_ARG(%v, returnArg)", strings.Replace(functionOutputType, parser.CleanValue(function.Output), fmt.Sprintf("type%v", hex.EncodeToString(tHash.Sum(nil)[:3])), -1))
						}

						return fmt.Sprintf(", Q_RETURN_ARG(%v, returnArg)", functionOutputType)
					}
					return ""
				}(),

				converter.CppInputParametersForSlotInvoke(function),
			)

			if converter.CppHeaderOutput(function) != parser.VOID {
				fmt.Fprintf(bb, "\n\treturn %v;", converter.CppOutput("returnArg", functionOutputType, function))
			}

			return bb.String()
		}

	case parser.PLAIN, parser.DESTRUCTOR:
		{
			if (function.Meta == parser.DESTRUCTOR || strings.HasPrefix(function.Name, parser.TILDE)) && function.Default {
				return ""
			}

			if function.Fullname == "SailfishApp::application" || function.Fullname == "SailfishApp::main" {
				return fmt.Sprintf(`	static int argcs = argc;
	static char** argvs = static_cast<char**>(malloc(argcs * sizeof(char*)));

	QList<QByteArray> aList = QByteArray(argv).split('|');
	for (int i = 0; i < argcs; i++)
		argvs[i] = (new QByteArray(aList.at(i)))->data();

	return %v(%v);`,

					function.Fullname,

					converter.CppInputParameters(function),
				)
			}

			return fmt.Sprintf("\t%v%v;",

				func() string {
					if converter.CppHeaderOutput(function) != parser.VOID {
						return "return "
					}
					return ""
				}(),

				converter.CppOutputParameters(function, fmt.Sprintf("%v%v%v(%v)%v",

					func() string {
						var c, _ = function.Class()
						//TODO:
						if c.Name == "QAndroidJniEnvironment" && function.Meta == parser.PLAIN && strings.HasPrefix(function.Name, "Exception") {
							return "({ QAndroidJniEnvironment env; env->"
						}
						if function.NonMember {
							return ""
						}
						if function.Static {
							return fmt.Sprintf("%v::", function.ClassName())
						}
						return fmt.Sprintf("static_cast<%v*>(ptr)->",
							func() string {
								if c.Fullname != "" {
									return c.Fullname
								}
								if strings.HasSuffix(function.Name, "_atList") {
									if function.IsMap {
										return fmt.Sprintf("%v<%v,%v>", parser.CleanValue(function.Container), function.Parameters[0].Value, strings.TrimPrefix(function.Output, "const "))
									}
									return fmt.Sprintf("%v<%v>", parser.CleanValue(function.Container), strings.TrimPrefix(function.Output, "const "))
								}
								if strings.HasSuffix(function.Name, "_setList") {
									if len(function.Parameters) == 2 {
										return fmt.Sprintf("%v<%v,%v>", parser.CleanValue(function.Container), function.Parameters[0].Value, strings.TrimPrefix(function.Parameters[1].Value, "const "))
									}
									return fmt.Sprintf("%v<%v>", parser.CleanValue(function.Container), strings.TrimPrefix(function.Parameters[0].Value, "const "))
								}
								if strings.HasSuffix(function.Name, "_newList") {
									//will be overriden
								}
								if strings.HasSuffix(function.Name, "_keyList") {
									//will be overriden
								}
								return function.ClassName()
							}(),
						)
					}(),

					func() string {
						if function.Default {
							var c, _ = function.Class()
							if c.Module == parser.MOC {
								if function.IsMocProperty {
									return fmt.Sprintf("%vDefault", function.Name)
								}
								return fmt.Sprintf("%v::%v", parser.State.ClassMap[function.ClassName()].GetBases()[0], function.Name)
							} else {
								return fmt.Sprintf("%v::%v", function.ClassName(), function.Name)
							}
						}
						return function.Name
					}(),

					converter.CppOutputParametersDeducedFromGeneric(function), converter.CppInputParameters(function),
					//TODO:
					func() string {
						var c, _ = function.Class()
						if c.Name == "QAndroidJniEnvironment" && function.Meta == parser.PLAIN && strings.HasPrefix(function.Name, "Exception") {
							return "; })"
						}
						return ""
					}())))
		}

	case parser.GETTER:
		{
			return fmt.Sprintf("\treturn %v;", converter.CppOutputParameters(function,
				func() string {
					if function.Static {
						return function.Fullname
					}
					return fmt.Sprintf("static_cast<%v*>(ptr)->%v", func() string {
						if c := parser.State.ClassMap[function.ClassName()]; c != nil && c.Fullname != "" {
							return c.Fullname
						}
						return function.ClassName()
					}(), function.Name)
				}()))
		}

	case parser.SETTER:
		{
			var function = *function
			function.Name = function.TmpName
			function.Fullname = fmt.Sprintf("%v::%v", function.ClassName(), function.Name)

			return fmt.Sprintf("\t%v = %v;", converter.CppOutputParameters(&function,
				func() string {
					if function.Static {
						return function.Fullname
					}
					return fmt.Sprintf("static_cast<%v*>(ptr)->%v", func() string {
						if c := parser.State.ClassMap[function.ClassName()]; c != nil && c.Fullname != "" {
							return c.Fullname
						}
						return function.ClassName()
					}(), function.Name)
				}()),

				converter.CppInputParameters(&function),
			)
		}

	case parser.SIGNAL:
		{
			var bb = new(bytes.Buffer)
			defer bb.Reset()

			if function.SignalMode == parser.CONNECT {
				if reg := converter.CppRegisterMetaType(function); reg != "" {
					bb.WriteString(reg + "\n")
				}
			}

			var my string
			var c, _ = function.Class()
			if c.Module != parser.MOC {
				my = "My"
			}
			if converter.IsPrivateSignal(function) {
				fmt.Fprintf(bb, "\tQObject::%v(static_cast<%v*>(ptr), &%v::%v, static_cast<%v%v*>(ptr), static_cast<%v (%v%v::*)(%v)>(&%v%v::Signal_%v%v));", strings.ToLower(function.SignalMode), function.ClassName(), function.ClassName(), function.Name, my, function.ClassName(), function.Output, my, function.ClassName(), converter.CppInputParametersForSignalConnect(function), my, function.ClassName(), strings.Title(function.Name), function.OverloadNumber)
			} else {
				fmt.Fprintf(bb, "\tQObject::%v(static_cast<%v*>(ptr), static_cast<%v (%v::*)(%v)>(&%v::%v), static_cast<%v%v*>(ptr), static_cast<%v (%v%v::*)(%v)>(&%v%v::Signal_%v%v));",
					strings.ToLower(function.SignalMode),

					function.ClassName(), function.Output, function.ClassName(), converter.CppInputParametersForSignalConnect(function), function.ClassName(), function.Name,

					my, function.ClassName(), function.Output, my, function.ClassName(), converter.CppInputParametersForSignalConnect(function), my, function.ClassName(), strings.Title(function.Name), function.OverloadNumber)
			}
			return bb.String()
		}
	}

	function.Access = "unsupported_cppFunctionBody"
	return function.Access
}
