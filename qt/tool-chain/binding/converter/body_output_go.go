package converter

import "github.com/peterq/pan-light/qt/tool-chain/binding/parser"

func GoOutputParametersFromC(function *parser.Function, name string) string {
	if function.Meta == parser.CONSTRUCTOR {
		return goOutput(name, function.Name, function, function.PureGoOutput)
	}
	return goOutput(name, function.Output, function, function.PureGoOutput)
}

func GoJSOutputParametersFromC(function *parser.Function, name string) string {
	if function.Meta == parser.CONSTRUCTOR {
		return goOutputJS(name, function.Name, function, function.PureGoOutput)
	}
	return goOutputJS(name, function.Output, function, function.PureGoOutput)
}

func GoOutputParametersFromCFailed(function *parser.Function) string {
	if function.Meta == parser.CONSTRUCTOR {
		return goOutputFailed(function.Name, function, function.PureGoOutput)
	}
	return goOutputFailed(function.Output, function, function.PureGoOutput)
}
