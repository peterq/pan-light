package sub

import (
	"github.com/peterq/pan-light/qt/core"

	_ "github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/subsub"
)

var SomeType *someType

type someType struct {
	core.QObject

	_ bool   `property:"someBool,auto"`
	_ string `property:"someString,auto"`

	_ func(string) `signal:"someSignal,auto"`
	_ func(string) `slot:"someSlot,auto"`

	_ bool   `property:"someBoolB"`
	_ string `property:"someStringB"`

	_ func(string) `signal:"someSignalB"`
	_ func(string) `slot:"someSlotB"`

	_ string `property:"someSubProp,->(subsubcustom.SubSubTestStructInstance.subsubProperty)"`
	_ string `property:"someSubPropA,<-(subsubcustom.SubSubTestStructInstance.subsubProperty)"`

	_ func(string) `signal:"someSignalC,->(subsubcustom.SubSubTestStructInstance.subPropertySignal)"`
	_ func(string) `slot:"someSlotC,->(subsubcustom.SubSubTestStructInstance.subPropertySlot)"`

	_ func(string) `signal:"SomeSignalD,<-(subsubcustom.SubSubTestStructInstance.subPropertySignal)"`
	_ func(string) `slot:"SomeSlotD,<-(subsubcustom.SubSubTestStructInstance.subPropertySlot)"`
}

func (t *someType) isSomeBool() bool     { return true }
func (t *someType) setSomeBool(bool)     {}
func (t *someType) someBoolChanged(bool) {}

func (t *someType) someString() string       { return "test" }
func (t *someType) setSomeString(string)     {}
func (t *someType) someStringChanged(string) {}

func (t *someType) someSignal(string) {}
func (t *someType) someSlot(string)   {}

type someOtherType struct {
	core.QObject

	_ func() `constructor:"init"`

	_ bool   `property:"someBool,->(SomeType)"`
	_ string `property:"someString,->(SomeType)"`

	_ bool   `property:"someBoolA,->(SomeType.someBool)"`
	_ string `property:"someStringA,->(SomeType.someString)"`

	_ bool   `property:"someBoolB,<-(SomeType)"`
	_ string `property:"someStringB,<-(SomeType)"`

	_ bool   `property:"someBoolC,<-(SomeType.someBool)"`
	_ string `property:"someStringC,<-(SomeType.someString)"`

	_ func(string) `signal:"someSignal,->(SomeType)"`
	_ func(string) `signal:"someSlot,->(SomeType)"`

	_ func(string) `signal:"someSignalA,->(SomeType.someSignal)"`
	_ func(string) `signal:"someSlotA,->(SomeType.someSlot)"`

	_ func(string) `signal:"someSignalB,<-(SomeType)"`
	_ func(string) `signal:"someSlotB,<-(SomeType)"`

	_ func(string) `signal:"someSignalC,<-(SomeType.someSignal)"`
	_ func(string) `signal:"someSlotC,<-(SomeType.someSlot)"`
}

func (t *someOtherType) init() { SomeType = NewSomeType(nil) }

func (t *someOtherType) someSignalB(string) {}
func (t *someOtherType) someSlotB(string)   {}

func (t *someOtherType) someSignalC(string) {}
func (t *someOtherType) someSlotC(string)   {}

type someOtherTypeOut struct {
	core.QObject

	_ func() `constructor:"init"`

	_ bool   `property:"someBool,->(this.c)"`
	_ string `property:"someString,->(this.c)"`

	_ bool   `property:"someBoolA,->(this.c.someBool)"`
	_ string `property:"someStringA,->(this.c.someString)"`

	_ func(string) `signal:"someSignal,->(this.c)"`
	_ func(string) `signal:"someSlot,->(this.c)"`

	_ func(string) `signal:"someSignalA,->(this.c.someSignal)"`
	_ func(string) `signal:"someSlotA,->(this.c.someSlot)"`

	c *someType
}

func (t *someOtherTypeOut) init() { t.c = NewSomeType(nil) }

type someOtherTypeIn struct {
	core.QObject

	_ func() `constructor:"init"`

	_ bool   `property:"someBool,<-(this.c)"`
	_ string `property:"someString,<-(this.c)"`

	_ bool   `property:"someBoolA,<-(this.c.someBool)"`
	_ string `property:"someStringA,<-(this.c.someString)"`

	_ func(string) `signal:"someSignal,<-(this.c)"`
	_ func(string) `slot:"someSlot,<-(this.c)"`

	_ func(string) `signal:"someSignalA,<-(this.c.someSignal)"`
	_ func(string) `slot:"someSlotA,<-(this.c.someSlot)"`

	_ func(string) `signal:"SomeSignalB,<-(this.c.someSignal)"`
	_ func(string) `slot:"SomeSlotB,<-(this.c.someSlot)"`

	c *someType
}

func (t *someOtherTypeIn) init() { t.c = NewSomeType(nil) }

func (t *someOtherTypeIn) someSignal(string) {}
func (t *someOtherTypeIn) someSlot(string)   {}

func (t *someOtherTypeIn) someSignalA(string) {}
func (t *someOtherTypeIn) someSlotA(string)   {}
