package parser

import (
	"strings"
)

func (c *Class) remove() {
	c.removeFunctions()
	c.removeFunctions_Version()

	c.removeEnums()
	c.removeEnums_Version()

	c.removeBases()
}

func (c *Class) removeFunctions() {
	for i := len(c.Functions) - 1; i >= 0; i-- {
		f := c.Functions[i]

		switch {
		case (f.Status == "obsolete" || f.Status == "compat") ||
			!(f.Access == "public" || f.Access == "protected") ||
			strings.ContainsAny(f.Name, "&<>=/!()[]{}^|*+-") ||
			strings.Contains(f.Name, "Operator"):
			{
				c.Functions = append(c.Functions[:i], c.Functions[i+1:]...)
			}

		case (f.Virtual == IMPURE || f.Virtual == PURE) && f.Meta == CONSTRUCTOR:
			{
				c.Functions = append(c.Functions[:i], c.Functions[i+1:]...)
			}
		}
	}
}

func (c *Class) removeFunctions_Version() {
	for i := len(c.Functions) - 1; i >= 0; i-- {
		switch c.Functions[i].Fullname {
		case "QTextBrowser::isModified", "QTextBrowser::setModified":
			{
				c.Functions = append(c.Functions[:i], c.Functions[i+1:]...)
			}

		case "QSemaphoreReleaser::QSemaphoreReleaser":
			{
				if c.Functions[i].OverloadNumber == "4" {
					c.Functions = append(c.Functions[:i], c.Functions[i+1:]...)
				}
			}
		}
	}
}

func (c *Class) removeEnums() {
	for i := len(c.Enums) - 1; i >= 0; i-- {
		if e := c.Enums[i]; (e.Status == "obsolete" || e.Status == "compat") ||
			!(e.Access == "public" || e.Access == "protected") {

			c.Enums = append(c.Enums[:i], c.Enums[i+1:]...)
		}
	}
}

func (c *Class) removeEnums_Version() {
	for i := len(c.Enums) - 1; i >= 0; i-- {
		switch c.Enums[i].ClassName() {
		case "QCss", "QScript", "Http2":
			{
				c.Enums = append(c.Enums[:i], c.Enums[i+1:]...)
				continue
			}
		}
		switch e := c.Enums[i]; e.Fullname {
		case "QTimeZone::anonymous":
			{
				c.Enums = append(c.Enums[:i], c.Enums[i+1:]...)
				continue
			}
		case "Qt::InputMethodQuery":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if e.Values[iv].Name == "ImQueryInput" {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
						break
					}
				}
				continue
			}
		case "QV4::PropertyFlag":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if e.Values[iv].Name == "Attr_ReadOnly" ||
						e.Values[iv].Name == "Attr_ReadOnly_ButConfigurable" {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
					}
				}
				continue
			}
		case "QDBusConnection::RegisterOption":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if strings.HasPrefix(e.Values[iv].Name, "ExportAll") {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
					}
				}
				continue
			}
		case "QDateTimeEdit::Section":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if e.Values[iv].Name == "TimeSections_Mask" ||
						e.Values[iv].Name == "DateSections_Mask" {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
					}
				}
				continue
			}
		case "QDockWidget::DockWidgetFeature":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if e.Values[iv].Name == "AllDockWidgetFeatures" {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
					}
				}
				continue
			}
		case "QSGNode::DirtyStateBit":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if e.Values[iv].Name == "DirtyPropagationMask" {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
					}
				}
				continue
			}
		case "QWebEnginePage::WebAction":
			{
				for iv := len(e.Values) - 1; iv >= 0; iv-- {
					if e.Values[iv].Name == "NoWebAction" {
						c.Enums[i].Values = append(c.Enums[i].Values[:iv], c.Enums[i].Values[iv+1:]...)
					}
				}
				continue
			}
		}
	}
}

func (c *Class) removeBases() {
	var bases = c.GetBases()
	for i := len(bases) - 1; i >= 0; i-- {
		if _, ok := State.ClassMap[bases[i]]; !ok {
			bases = append(bases[:i], bases[i+1:]...)
		}
	}
	c.Bases = strings.Join(bases, ",")
}
