package qtmoc

import (
	"errors"
	. "math/rand"
	"os"
	"strings"
	"testing"
	tps "time"
	"unsafe"

	"github.com/peterq/pan-light/qt/core"
	"github.com/peterq/pan-light/qt/gui"
	"github.com/peterq/pan-light/qt/sql"
	"github.com/peterq/pan-light/qt/widgets"
	"github.com/peterq/pan-light/qt/xml"

	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub"
	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/subsub" //subsubcustom
)

type Application struct {
	widgets.QApplication
}

type testStruct struct {
	otherTestStruct
	testOther otherTestStruct

	_ bool             `property:"propBool"`
	_ int8             `property:"propInt8"`  // -> string
	_ uint8            `property:"propInt82"` // -> string
	_ int16            `property:"propInt16"`
	_ uint16           `property:"propInt162"`
	_ int32            `property:"propInt32"`  // -> int
	_ uint32           `property:"propInt322"` // -> int
	_ int              `property:"propInt"`
	_ uint             `property:"propInt2"`
	_ int64            `property:"propInt64"`
	_ uint64           `property:"propInt642"`
	_ float32          `property:"propFloat"`
	_ float64          `property:"propFloat2"`
	_ string           `property:"propString"`
	_ []string         `property:"propString2"`
	_ uintptr          `property:"propPointer"`
	_ unsafe.Pointer   `property:"propPointer2"`
	_ core.QVariant    `property:"propObject"`  // -> T (c++)
	_ *core.QObject    `property:"propObject2"` // -> *T
	_ *core.QVariant   `property:"propObject3"` // -> *T //TODO:
	_ core.Qt__Key     `property:"propEnum"`
	_ error            `property:"propError"`
	_ otherTestStruct  `property:"propReturnTest"`  // -> *T
	_ *otherTestStruct `property:"propReturnTest2"` // -> *T

	_ []bool             `property:"propListBool"`
	_ []int8             `property:"propListInt8"`  // -> string
	_ []uint8            `property:"propListInt82"` // -> string
	_ []int16            `property:"propListInt16"`
	_ []uint16           `property:"propListInt162"`
	_ []int32            `property:"propListInt32"`  // -> int
	_ []uint32           `property:"propListInt322"` // -> int
	_ []int              `property:"propListInt"`
	_ []uint             `property:"propListInt2"`
	_ []int64            `property:"propListInt64"`
	_ []uint64           `property:"propListInt642"`
	_ []float32          `property:"propListFloat"`
	_ []float64          `property:"propListFloat2"`
	_ []string           `property:"propListString"`
	_ []string           `property:"propListString2"`
	_ []uintptr          `property:"propListPointer"`
	_ []unsafe.Pointer   `property:"propListPointer2"`
	_ []core.QVariant    `property:"propListObject"`  // -> T (c++)
	_ []*core.QObject    `property:"propListObject2"` // -> *T
	_ []*core.QVariant   `property:"propListObject3"` // -> *T //TODO:
	_ []core.Qt__Key     `property:"propListEnum"`
	_ []error            `property:"propListError"`
	_ []otherTestStruct  `property:"propListReturnTest"`  // -> *T
	_ []*otherTestStruct `property:"propListReturnTest2"` // -> *T

	_ map[int]bool             `property:"propMapBool"`
	_ map[int]int8             `property:"propMapInt8"`  // -> string
	_ map[int]uint8            `property:"propMapInt82"` // -> string
	_ map[int]int16            `property:"propMapInt16"`
	_ map[int]uint16           `property:"propMapInt162"`
	_ map[int]int32            `property:"propMapInt32"`  // -> int
	_ map[int]uint32           `property:"propMapInt322"` // -> int
	_ map[int]int              `property:"propMapInt"`
	_ map[int]uint             `property:"propMapInt2"`
	_ map[int]int64            `property:"propMapInt64"`
	_ map[int]uint64           `property:"propMapInt642"`
	_ map[int]float32          `property:"propMapFloat"`
	_ map[int]float64          `property:"propMapFloat2"`
	_ map[int]string           `property:"propMapString"`
	_ map[int][]string         `property:"propMapString2"`
	_ map[int]uintptr          `property:"propMapPointer"`
	_ map[int]unsafe.Pointer   `property:"propMapPointer2"`
	_ map[int]core.QVariant    `property:"propMapObject"`  // -> T (c++)
	_ map[int]*core.QObject    `property:"propMapObject2"` // -> *T
	_ map[int]*core.QVariant   `property:"propMapObject3"` // -> *T //TODO:
	_ map[int]core.Qt__Key     `property:"propMapEnum"`
	_ map[int]error            `property:"propMapError"`
	_ map[int]otherTestStruct  `property:"propMapReturnTest"`  // -> *T
	_ map[int]*otherTestStruct `property:"propMapReturnTest2"` // -> *T

	_ map[bool]bool       `property:"propMapKeyBool"`
	_ map[int8]int8       `property:"propMapKeyInt8"`  // -> string
	_ map[uint8]uint8     `property:"propMapKeyInt82"` // -> string
	_ map[int16]int16     `property:"propMapKeyInt16"`
	_ map[uint16]uint16   `property:"propMapKeyInt162"`
	_ map[int32]int32     `property:"propMapKeyInt32"`  // -> int
	_ map[uint32]uint32   `property:"propMapKeyInt322"` // -> int
	_ map[int]int         `property:"propMapKeyInt"`
	_ map[uint]uint       `property:"propMapKeyInt2"`
	_ map[int64]int64     `property:"propMapKeyInt64"`
	_ map[uint64]uint64   `property:"propMapKeyInt642"`
	_ map[float32]float32 `property:"propMapKeyFloat"`
	_ map[float64]float64 `property:"propMapKeyFloat2"`
	_ map[string]string   `property:"propMapKeyString"`
	//is invalid in go _ map[[]string][]string             `property:"propMapKeyString2"`
	_ map[uintptr]uintptr               `property:"propMapKeyPointer"`
	_ map[unsafe.Pointer]unsafe.Pointer `property:"propMapKeyPointer2"`
	//will probably never work _ map[core.QVariant]core.QVariant       `property:"propMapKeyObject"`  // -> T (c++)
	_ map[*core.QObject]*core.QObject `property:"propMapKeyObject2"` // -> *T
	//will work once * is recognized (TODO) _ map[*core.QVariant]*core.QVariant `property:"propMapKeyObject3"` // -> *T //TODO:
	_ map[core.Qt__Key]core.Qt__Key `property:"propMapKeyEnum"`
	_ map[error]error               `property:"propMapKeyError"`
	//is invalid in go _ map[otherTestStruct]otherTestStruct   `property:"propMapKeyReturnTest"`  // -> *T
	_ map[*otherTestStruct]otherTestStruct  `property:"propMapKeyReturnTest"`  // -> *T
	_ map[*otherTestStruct]*otherTestStruct `property:"propMapKeyReturnTest2"` // -> *T

	a, b bool
	ab   func(bool) bool
	abc  func(bool) bool `test:"test"`
	abcd struct {
		a, b bool
		ab   func(bool) bool
		abc  func(bool) bool `test:"test"`
	}

	test widgets.QWidget
	core.QObject
	widgets.QWidget

	_ func(bool, bool)                   `signal:"BoolSignalInput"`
	_ func(int8, uint8)                  `signal:"Int8SignalInput"` // -> string
	_ func(int16, uint16)                `signal:"Int16SignalInput"`
	_ func(int32, uint32)                `signal:"Int32SignalInput"` // -> int
	_ func(int, uint)                    `signal:"IntSignalInput"`
	_ func(int64, uint64)                `signal:"Int64SignalInput"`
	_ func(float32, float64)             `signal:"FloatSignalInput"`
	_ func(string, []string)             `signal:"StringSignalInput"`
	_ func(uintptr, unsafe.Pointer)      `signal:"PointerSignalInput"`
	_ func(core.QVariant, *core.QObject) `signal:"ObjectSignalInput"` // -> T (c++) *T
	_ func(core.Qt__Key)                 `signal:"EnumSignalInput"`
	_ func(error)                        `signal:"ErrorSignalInput"`

	_ func([]bool, []bool)                   `signal:"BoolSignalListInput"`
	_ func([]int8, []uint8)                  `signal:"Int8SignalListInput"` // -> string
	_ func([]int16, []uint16)                `signal:"Int16SignalListInput"`
	_ func([]int32, []uint32)                `signal:"Int32SignalListInput"` // -> int
	_ func([]int, []uint)                    `signal:"IntSignalListInput"`
	_ func([]int64, []uint64)                `signal:"Int64SignalListInput"`
	_ func([]float32, []float64)             `signal:"FloatSignalListInput"`
	_ func([]string, []string)               `signal:"StringSignalListInput"`
	_ func([]uintptr, []unsafe.Pointer)      `signal:"PointerSignalListInput"`
	_ func([]core.QVariant, []*core.QObject) `signal:"ObjectSignalListInput"` // -> T (c++) *T
	_ func([]core.Qt__Key)                   `signal:"EnumSignalListInput"`
	_ func([]error)                          `signal:"ErrorSignalListInput"`

	_ func(map[int]bool, map[int]bool)                   `signal:"BoolSignalMapInput"`
	_ func(map[int]int8, map[int]uint8)                  `signal:"Int8SignalMapInput"` // -> string
	_ func(map[int]int16, map[int]uint16)                `signal:"Int16SignalMapInput"`
	_ func(map[int]int32, map[int]uint32)                `signal:"Int32SignalMapInput"` // -> int
	_ func(map[int]int, map[int]uint)                    `signal:"IntSignalMapInput"`
	_ func(map[int]int64, map[int]uint64)                `signal:"Int64SignalMapInput"`
	_ func(map[int]float32, map[int]float64)             `signal:"FloatSignalMapInput"`
	_ func(map[int]string, map[int]string)               `signal:"StringSignalMapInput"`
	_ func(map[int]uintptr, map[int]unsafe.Pointer)      `signal:"PointerSignalMapInput"`
	_ func(map[int]core.QVariant, map[int]*core.QObject) `signal:"ObjectSignalMapInput"` // -> T (c++) *T
	_ func(map[int]core.Qt__Key)                         `signal:"EnumSignalMapInput"`
	_ func(map[int]error)                                `signal:"ErrorSignalMapInput"`

	_ func(bool, bool)                   `slot:"BoolSlotInput"`
	_ func(int8, uint8)                  `slot:"Int8SlotInput"` // -> string
	_ func(int16, uint16)                `slot:"Int16SlotInput"`
	_ func(int32, uint32)                `slot:"Int32SlotInput"` // -> int
	_ func(int, uint)                    `slot:"IntSlotInput"`
	_ func(int64, uint64)                `slot:"Int64SlotInput"`
	_ func(float32, float64)             `slot:"FloatSlotInput"`
	_ func(string, []string)             `slot:"StringSlotInput"`
	_ func(uintptr, unsafe.Pointer)      `slot:"PointerSlotInput"`
	_ func(core.QVariant, *core.QObject) `slot:"ObjectSlotInput"` // -> T (c++) *T
	_ func(core.Qt__Key)                 `slot:"EnumSlotInput"`
	_ func(error)                        `slot:"ErrorSlotInput"`

	_ func([]bool, []bool)                   `slot:"BoolSlotListInput"`
	_ func([]int8, []uint8)                  `slot:"Int8SlotListInput"` // -> string
	_ func([]int16, []uint16)                `slot:"Int16SlotListInput"`
	_ func([]int32, []uint32)                `slot:"Int32SlotListInput"` // -> int
	_ func([]int, []uint)                    `slot:"IntSlotListInput"`
	_ func([]int64, []uint64)                `slot:"Int64SlotListInput"`
	_ func([]float32, []float64)             `slot:"FloatSlotListInput"`
	_ func([]string, []string)               `slot:"StringSlotListInput"`
	_ func([]uintptr, []unsafe.Pointer)      `slot:"PointerSlotListInput"`
	_ func([]core.QVariant, []*core.QObject) `slot:"ObjectSlotListInput"` // -> T (c++) *T
	_ func([]core.Qt__Key)                   `slot:"EnumSlotListInput"`
	_ func([]error)                          `slot:"ErrorSlotListInput"`

	_ func(map[int]bool, map[int]bool)                   `slot:"BoolSlotMapInput"`
	_ func(map[int]int8, map[int]uint8)                  `slot:"Int8SlotMapInput"` // -> string
	_ func(map[int]int16, map[int]uint16)                `slot:"Int16SlotMapInput"`
	_ func(map[int]int32, map[int]uint32)                `slot:"Int32SlotMapInput"` // -> int
	_ func(map[int]int, map[int]uint)                    `slot:"IntSlotMapInput"`
	_ func(map[int]int64, map[int]uint64)                `slot:"Int64SlotMapInput"`
	_ func(map[int]float32, map[int]float64)             `slot:"FloatSlotMapInput"`
	_ func(map[int]string, map[int]string)               `slot:"StringSlotMapInput"`
	_ func(map[int]uintptr, map[int]unsafe.Pointer)      `slot:"PointerSlotMapInput"`
	_ func(map[int]core.QVariant, map[int]*core.QObject) `slot:"ObjectSlotMapInput"` // -> T (c++) *T
	_ func(map[int]core.Qt__Key)                         `slot:"EnumSlotMapInput"`
	_ func(map[int]error)                                `slot:"ErrorSlotMapInput"`

	_ func(bool) bool                     `slot:"BoolSlotOutput"`
	_ func(bool) bool                     `slot:"BoolSlotOutput2"`
	_ func(int8) int8                     `slot:"Int8SlotOutput"`  // -> string
	_ func(uint8) uint8                   `slot:"Int8SlotOutput2"` // -> string
	_ func(int16) int16                   `slot:"Int16SlotOutput"`
	_ func(uint16) uint16                 `slot:"Int16SlotOutput2"`
	_ func(int32) int32                   `slot:"Int32SlotOutput"`  // -> int
	_ func(uint32) uint32                 `slot:"Int32SlotOutput2"` // -> int
	_ func(int) int                       `slot:"IntSlotOutput"`
	_ func(uint) uint                     `slot:"IntSlotOutput2"`
	_ func(int64) int64                   `slot:"Int64SlotOutput"`
	_ func(uint64) uint64                 `slot:"Int64SlotOutput2"`
	_ func(float32) float32               `slot:"FloatSlotOutput"`
	_ func(float64) float64               `slot:"FloatSlotOutput2"`
	_ func(string) string                 `slot:"StringSlotOutput"`
	_ func([]string) []string             `slot:"StringSlotOutput2"`
	_ func(uintptr) uintptr               `slot:"PointerSlotOutput"`
	_ func(unsafe.Pointer) unsafe.Pointer `slot:"PointerSlotOutput2"`
	_ func(core.QVariant) core.QVariant   `slot:"ObjectSlotOutput"`  // -> T (c++)
	_ func(*core.QObject) *core.QObject   `slot:"ObjectSlotOutput2"` // -> *T
	_ func(core.Qt__Key) core.Qt__Key     `slot:"EnumSlotOutput"`
	_ func(error) error                   `slot:"ErrorSlotOutput"`
	_ func(testStruct) testStruct         `slot:"returnTest"`  // -> *T
	_ func(*testStruct) *testStruct       `slot:"returnTest2"` // -> *T
	_ func(a0 string) (a1 string)         `slot:"returnName"`
	_ func(a0, a1 int) (a2 int)           `slot:"returnName2"`
	_ func()                              `slot:"other"`
	_ func() bool                         `slot:"other2"`

	_ func([]bool) []bool                     `slot:"BoolSlotListOutput"`
	_ func([]bool) []bool                     `slot:"BoolSlotListOutput2"`
	_ func([]int8) []int8                     `slot:"Int8SlotListOutput"`  // -> string
	_ func([]uint8) []uint8                   `slot:"Int8SlotListOutput2"` // -> string
	_ func([]int16) []int16                   `slot:"Int16SlotListOutput"`
	_ func([]uint16) []uint16                 `slot:"Int16SlotListOutput2"`
	_ func([]int32) []int32                   `slot:"Int32SlotListOutput"`  // -> int
	_ func([]uint32) []uint32                 `slot:"Int32SlotListOutput2"` // -> int
	_ func([]int) []int                       `slot:"IntSlotListOutput"`
	_ func([]uint) []uint                     `slot:"IntSlotListOutput2"`
	_ func([]int64) []int64                   `slot:"Int64SlotListOutput"`
	_ func([]uint64) []uint64                 `slot:"Int64SlotListOutput2"`
	_ func([]float32) []float32               `slot:"FloatSlotListOutput"`
	_ func([]float64) []float64               `slot:"FloatSlotListOutput2"`
	_ func([]string) []string                 `slot:"StringSlotListOutput"`
	_ func([]string) []string                 `slot:"StringSlotListOutput2"`
	_ func([]uintptr) []uintptr               `slot:"PointerSlotListOutput"`
	_ func([]unsafe.Pointer) []unsafe.Pointer `slot:"PointerSlotListOutput2"`
	_ func([]core.QVariant) []core.QVariant   `slot:"ObjectSlotListOutput"`  // -> T (c++)
	_ func([]*core.QObject) []*core.QObject   `slot:"ObjectSlotListOutput2"` // -> *T
	_ func([]core.Qt__Key) []core.Qt__Key     `slot:"EnumSlotListOutput"`
	_ func([]error) []error                   `slot:"ErrorSlotListOutput"`
	_ func([]testStruct) []testStruct         `slot:"returnListTest"`  // -> *T
	_ func([]*testStruct) []*testStruct       `slot:"returnListTest2"` // -> *T
	_ func(a0 []string) (a1 []string)         `slot:"returnListName"`
	_ func(a0, a1 []int) (a2 []int)           `slot:"returnListName2"`
	_ func() []bool                           `slot:"otherList"`

	_ func(map[int]bool) map[int]bool                     `slot:"BoolSlotMapOutput"`
	_ func(map[int]bool) map[int]bool                     `slot:"BoolSlotMapOutput2"`
	_ func(map[int]int8) map[int]int8                     `slot:"Int8SlotMapOutput"`  // -> string
	_ func(map[int]uint8) map[int]uint8                   `slot:"Int8SlotMapOutput2"` // -> string
	_ func(map[int]int16) map[int]int16                   `slot:"Int16SlotMapOutput"`
	_ func(map[int]uint16) map[int]uint16                 `slot:"Int16SlotMapOutput2"`
	_ func(map[int]int32) map[int]int32                   `slot:"Int32SlotMapOutput"`  // -> int
	_ func(map[int]uint32) map[int]uint32                 `slot:"Int32SlotMapOutput2"` // -> int
	_ func(map[int]int) map[int]int                       `slot:"IntSlotMapOutput"`
	_ func(map[int]uint) map[int]uint                     `slot:"IntSlotMapOutput2"`
	_ func(map[int]int64) map[int]int64                   `slot:"Int64SlotMapOutput"`
	_ func(map[int]uint64) map[int]uint64                 `slot:"Int64SlotMapOutput2"`
	_ func(map[int]float32) map[int]float32               `slot:"FloatSlotMapOutput"`
	_ func(map[int]float64) map[int]float64               `slot:"FloatSlotMapOutput2"`
	_ func(map[int]string) map[int]string                 `slot:"StringSlotMapOutput"`
	_ func(map[int]string) map[int]string                 `slot:"StringSlotMapOutput2"`
	_ func(map[int]uintptr) map[int]uintptr               `slot:"PointerSlotMapOutput"`
	_ func(map[int]unsafe.Pointer) map[int]unsafe.Pointer `slot:"PointerSlotMapOutput2"`
	_ func(map[int]core.QVariant) map[int]core.QVariant   `slot:"ObjectSlotMapOutput"`  // -> T (c++)
	_ func(map[int]*core.QObject) map[int]*core.QObject   `slot:"ObjectSlotMapOutput2"` // -> *T
	_ func(map[int]core.Qt__Key) map[int]core.Qt__Key     `slot:"EnumSlotMapOutput"`
	_ func(map[int]error) map[int]error                   `slot:"ErrorSlotMapOutput"`
	_ func(map[int]testStruct) map[int]testStruct         `slot:"returnMapTest"`  // -> *T
	_ func(map[int]*testStruct) map[int]*testStruct       `slot:"returnMapTest2"` // -> *T
	_ func(a0 map[int]string) (a1 map[int]string)         `slot:"returnMapName"`
	_ func(a0, a1 map[int]int) (a2 map[int]int)           `slot:"returnMapName2"`
	_ func() map[int]bool                                 `slot:"otherMap"`

	_ func(*core.QObject, []*core.QObject) `signal:"mixedSignal"`
	_ func([]*core.QObject, *core.QObject) `signal:"mixedSignal2"`
	_ func(*core.QObject, []*core.QObject) `signal:"mixedSlot"`
	_ func([]*core.QObject, *core.QObject) `signal:"mixedSlot2"`

	_ func(string, string, string) error `slot:"errorStringTest1,auto"`
	_ func(string, string, error) string `slot:"errorStringTest2,auto"`
	_ func(string, error, string) string `slot:"errorStringTest3,auto"`

	_ func(string, string, string) []error `slot:"errorStringTest4,auto"`
	_ func(string, string, []error) string `slot:"errorStringTest5,auto"`
	_ func(string, []error, string) string `slot:"errorStringTest6,auto"`

	_ func(string, string, string) map[error]error `slot:"errorStringTest7,auto"`
	_ func(string, string, map[error]error) string `slot:"errorStringTest8,auto"`
	_ func(string, map[error]error, string) string `slot:"errorStringTest9,auto"`

	_ func(string, map[error]error, []error) error `slot:"errorStringTest10,auto"`
	_ func(map[error]error, []error, error) string `slot:"errorStringTest11,auto"`
	_ func([]error, error, string) map[error]error `slot:"errorStringTest12,auto"`

	_ func(map[string]error, map[error]string, map[error]error) map[string]string          `slot:"errorStringTest13,auto"`
	_ func(map[error][]string, map[string][]string, map[error][]error, map[string][]error) `slot:"errorStringTest14,auto"`
}

type subTestStruct struct {
	testStruct
	_ func(*subTestStruct) *subTestStruct `slot:"returnTest3"`
}

func (s *subSubTestStruct) init() {
	s.SubSubConstructorProperty++
}

type (
	subSubTestStruct struct {
		sub.SubTestStruct

		_ func() `constructor:"init"`

		_ *sub.SubTestStruct                          `property:"StructPropertyTest"`
		_ func(*sub.SubTestStruct) *sub.SubTestStruct `slot:"StructSlotTest"`
		_ func(*sub.SubTestStruct)                    `signal:"StructSignalTest"`
	}

	otherTestStruct struct {
		core.QObject

		a, b bool
		ab   func(bool) bool
		abc  func(bool) bool `test:"test"`
		abcd struct {
			a, b bool
			ab   func(bool) bool
			abc  func(bool) bool `test:"test"`
		}

		_ bool `property:"propBoolSub"`
	}
)

type abstractTestStruct1 struct {
	core.QAbstractItemModel
}

type abstractTestStruct2 struct {
	core.QAbstractListModel
}

type abstractTestStruct3 struct {
	core.QStringListModel
}

type abstractTestStruct4 struct {
	core.QAbstractProxyModel
}

type abstractTestStruct5 struct {
	core.QAbstractTableModel
}

type abstractTestStruct6 struct {
	sql.QSqlQueryModel
}

type abstractTestStruct7 struct {
	gui.QStandardItemModel
}

type abstractTestStruct8 struct {
	widgets.QFileSystemModel
}

type pureGoTestStruct struct {
	core.QObject

	_ func(*bool)                         `signal:"goSignal1"`
	_ func(*int)                          `signal:"goSignal2"`
	_ func(*string)                       `signal:"goSignal3"`
	_ func(interface{})                   `signal:"goSignal4"`
	_ func(func())                        `signal:"goSignal5"`
	_ func(func(func()))                  `signal:"goSignal6"`
	_ func(func(func(func())))            `signal:"goSignal7"`
	_ func(func(func(func(interface{})))) `signal:"goSignal8"`
	_ func(chan<- bool)                   `signal:"goSignal9"`
	_ func(<-chan bool)                   `signal:"goSignal10"`
	_ func(chan bool)                     `signal:"goSignal11"`
	_ func(a chan<- bool)                 `signal:"goSignal12"`
	_ func(b <-chan bool)                 `signal:"goSignal13"`
	_ func(c chan bool)                   `signal:"goSignal14"`

	_ func(*bool)                         `slot:"goSlot1"`
	_ func(*int)                          `slot:"goSlot2"`
	_ func(*string)                       `slot:"goSlot3"`
	_ func(interface{})                   `slot:"goSlot4"`
	_ func(func())                        `slot:"goSlot5"`
	_ func(func(func()))                  `slot:"goSlot6"`
	_ func(func(func(func())))            `slot:"goSlot7"`
	_ func(func(func(func(interface{})))) `slot:"goSlot8"`
	_ func(chan<- bool)                   `slot:"goSlot9"`

	_ func(a *bool) *bool                                                 `slot:"goRSlot1"`
	_ func(a *int) *int                                                   `slot:"goRSlot2"`
	_ func(a *string) *string                                             `slot:"goRSlot3"`
	_ func(a interface{}) interface{}                                     `slot:"goRSlot4"`
	_ func(a func()) func()                                               `slot:"goRSlot5"`
	_ func(a func(func())) func(func())                                   `slot:"goRSlot6"`
	_ func(a func(func(func()))) func(func(func()))                       `slot:"goRSlot7"`
	_ func(a func(func(func(interface{})))) func(func(func(interface{}))) `slot:"goRSlot8"`
	_ func(a chan<- bool) chan<- bool                                     `slot:"goRSlot9"`
	_ func(a func(b func(c func()))) (x func(y func(z func())))           `slot:"goRSlot10"`
	_ func(a func(b func(c func())))                                      `slot:"goRSlot11"`
	//TODO: _ func(func(b func(c func())))    `slot:"goRSlot12"`
	_ func(a func(func(c func())))    `slot:"goRSlot13"`
	_ func(a func(b func(func())))    `slot:"goRSlot14"`
	_ func(a func(b func(func(int)))) `slot:"goRSlot15"`
	//TODO: _ func(a func(b func(func(d int)))) `slot:"goRSlot16"`
	_ func() (x func(y func(z func(int)))) `slot:"goRSlot17"`

	_ func([][]string)                   `slot:"goSlotArray"`
	_ func(map[string]map[string]string) `slot:"goSlotMap"`
	_ func(goStruct)                     `slot:"goStruct1"`
	_ func(*goStruct)                    `slot:"goStruct2"`
	_ func(func(goStruct)) goStruct      `slot:"goStruct3"`
	_ func(func(*goStruct)) *goStruct    `slot:"goStruct4"`

	_ func(tps.Duration)   `signal:"goStdSignal1"`
	_ func(strings.Reader) `signal:"goStdSignal2"`
	_ func(Rand)           `signal:"goStdSignal3"`
}

type goStruct struct {
}

var (
	b0, b1   bool         = false, true
	b0L, b1L []bool       = []bool{b0}, []bool{b1}
	b0M, b1M map[int]bool = map[int]bool{0: b0}, map[int]bool{0: b1}

	i0  int16          = 1
	i0L []int16        = []int16{i0}
	i0M map[int]int16  = map[int]int16{0: i0}
	i1  uint16         = 2
	i1L []uint16       = []uint16{i1}
	i1M map[int]uint16 = map[int]uint16{0: i1}

	i2  int          = 3
	i2L []int        = []int{i2}
	i2M map[int]int  = map[int]int{0: i2}
	i3  uint         = 4
	i3L []uint       = []uint{i3}
	i3M map[int]uint = map[int]uint{0: i3}

	i4  int          = 5
	i4L []int        = []int{i4}
	i4M map[int]int  = map[int]int{0: i4}
	i5  uint         = 6
	i5L []uint       = []uint{i5}
	i5M map[int]uint = map[int]uint{0: i5}

	i6  int64          = 7
	i6L []int64        = []int64{i6}
	i6M map[int]int64  = map[int]int64{0: i6}
	i7  uint64         = 8
	i7L []uint64       = []uint64{i7}
	i7M map[int]uint64 = map[int]uint64{0: i7}

	f0  float32         = 123.45
	f0L []float32       = []float32{f0}
	f0M map[int]float32 = map[int]float32{0: f0}
	f1  float64         = 678.91
	f1L []float64       = []float64{f1}
	f1M map[int]float64 = map[int]float64{0: f1}

	s0  string           = "test"
	s1  []string         = []string{"t", "e", "s", "t"}
	s0M map[int]string   = map[int]string{0: s0}
	s1M map[int][]string = map[int][]string{0: s1}

	p0  uintptr                = uintptr(12345)
	p0L []uintptr              = []uintptr{p0}
	p0M map[int]uintptr        = map[int]uintptr{0: p0}
	p1  unsafe.Pointer         = unsafe.Pointer(uintptr(67891))
	p1L []unsafe.Pointer       = []unsafe.Pointer{p1}
	p1M map[int]unsafe.Pointer = map[int]unsafe.Pointer{0: p1}

	o0 *core.QVariant
	o1 *core.QObject

	e0  core.Qt__Key         = core.Qt__Key_Z
	e0L []core.Qt__Key       = []core.Qt__Key{e0}
	e0M map[int]core.Qt__Key = map[int]core.Qt__Key{0: e0}
	e1  error                = errors.New("test")
	e1L []error              = []error{e1}
	e1M map[int]error        = map[int]error{0: e1}
)

func init() { gui.NewQGuiApplication(len(os.Args), os.Args) }

func TestGeneral(t *testing.T) {
	if false {
		NewTestStruct(nil)
		NewOtherTestStruct(nil)
		NewSubTestStruct(nil)
		sub.NewSubTestStruct(nil)
		subsubcustom.NewSubSubTestStruct(nil)
	}
	if res := NewSubSubTestStruct(nil).SubSubConstructorProperty; res != 3 {
		t.Fatal(res, "!=", 3)
	}
}

func TestProperties(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectPropBoolChanged(func(propBool bool) {
		if propBool != b1 {
			t.Fatal(propBool, b1)
		}
	})
	test.SetPropBool(b1)
	test.PropBoolChanged(b1)
	if test.IsPropBool() != b1 {
		t.Fatal("IsPropBool")
	}

	//_ int8             `property:"propInt8"`
	//_ uint8            `property:"propInt82"`

	test.ConnectPropInt16Changed(func(propInt16 int16) {
		if propInt16 != i0 {
			t.Fatal(propInt16, i0)
		}
	})
	test.SetPropInt16(i0)
	test.PropInt16Changed(i0)
	if test.PropInt16() != i0 {
		t.Fatal("PropInt16")
	}

	test.ConnectPropInt162Changed(func(propInt162 uint16) {
		if propInt162 != i1 {
			t.Fatal(propInt162, i1)
		}
	})
	test.SetPropInt162(i1)
	test.PropInt162Changed(i1)
	if test.PropInt162() != i1 {
		t.Fatal("PropInt162")
	}

	test.ConnectPropInt32Changed(func(propInt32 int) {
		if propInt32 != i2 {
			t.Fatal(propInt32, i2)
		}
	})
	test.SetPropInt32(i2)
	test.PropInt32Changed(i2)
	if test.PropInt32() != i2 {
		t.Fatal("PropInt32")
	}

	test.ConnectPropInt322Changed(func(propInt322 uint) {
		if propInt322 != i3 {
			t.Fatal(propInt322, i3)
		}
	})
	test.SetPropInt322(i3)
	test.PropInt322Changed(i3)
	if test.PropInt322() != i3 {
		t.Fatal("PropInt322")
	}

	test.ConnectPropIntChanged(func(propInt int) {
		if propInt != i2 {
			t.Fatal(propInt, i2)
		}
	})
	test.SetPropInt(i2)
	test.PropIntChanged(i2)
	if test.PropInt() != i2 {
		t.Fatal("PropInt")
	}

	test.ConnectPropInt2Changed(func(propInt2 uint) {
		if propInt2 != i3 {
			t.Fatal(propInt2, i3)
		}
	})
	test.SetPropInt2(i3)
	test.PropInt2Changed(i3)
	if test.PropInt2() != i3 {
		t.Fatal("PropInt2")
	}

	test.ConnectPropInt64Changed(func(propInt64 int64) {
		if propInt64 != i6 {
			t.Fatal(propInt64, i6)
		}
	})
	test.SetPropInt64(i6)
	test.PropInt64Changed(i6)
	if test.PropInt64() != i6 {
		t.Fatal("PropInt64")
	}

	test.ConnectPropInt642Changed(func(propInt642 uint64) {
		if propInt642 != i7 {
			t.Fatal(propInt642, i7)
		}
	})
	test.SetPropInt642(i7)
	test.PropInt642Changed(i7)
	if test.PropInt642() != i7 {
		t.Fatal("PropInt642")
	}

	test.ConnectPropFloatChanged(func(propFloat float32) {
		if propFloat != f0 {
			t.Fatal(propFloat, f0)
		}
	})
	test.SetPropFloat(f0)
	test.PropFloatChanged(f0)
	if test.PropFloat() != f0 {
		t.Fatal("PropFloat")
	}

	test.ConnectPropFloat2Changed(func(propFloat2 float64) {
		if propFloat2 != f1 {
			t.Fatal(propFloat2, f1)
		}
	})
	test.SetPropFloat2(f1)
	test.PropFloat2Changed(f1)
	if test.PropFloat2() != f1 {
		t.Fatal("PropFloat2")
	}

	test.ConnectPropStringChanged(func(propString string) {
		if propString != s0 {
			t.Fatal(propString, s0)
		}
	})
	test.SetPropString(s0)
	test.PropStringChanged(s0)
	if test.PropString() != s0 {
		t.Fatal("PropString")
	}

	test.ConnectPropString2Changed(func(propString2 []string) {
		if strings.Join(propString2, "") != strings.Join(s1, "") {
			t.Fatal(propString2, s1)
		}
	})
	test.SetPropString2(s1)
	test.PropString2Changed(s1)
	if strings.Join(test.PropString2(), "") != strings.Join(s1, "") {
		t.Fatal("PropString2")
	}

	test.ConnectPropPointerChanged(func(propPointer uintptr) {
		if int(propPointer) != int(p0) {
			t.Fatal(propPointer, p0, int(propPointer), int(p0))
		}
	})
	test.SetPropPointer(p0)
	test.PropPointerChanged(p0)
	if int(test.PropPointer()) != int(p0) {
		t.Fatal("PropPointer")
	}

	test.ConnectPropPointer2Changed(func(propPointer2 unsafe.Pointer) {
		if int(uintptr(propPointer2)) != int(uintptr(p1)) {
			t.Fatal(propPointer2, p1, int(uintptr(propPointer2)), int(uintptr(p1)))
		}
	})
	test.SetPropPointer2(p1)
	test.PropPointer2Changed(p1)
	if int(uintptr(test.PropPointer2())) != int(uintptr(p1)) {
		t.Fatal("PropPointer2")
	}

	test.ConnectPropObjectChanged(func(propObject *core.QVariant) {
		if propObject.ToString() != o0.ToString() { //TODO:
			t.Fatal(propObject, o0, propObject.ToString(), o0.ToString())
		}
	})
	o0 = core.NewQVariant14("test")
	test.SetPropObject(o0)
	test.PropObjectChanged(o0)
	if test.PropObject().ToString() != o0.ToString() {
		t.Fatal("PropObject")
	}

	test.ConnectPropObject2Changed(func(propObject2 *core.QObject) {
		if propObject2.ObjectName() != o1.ObjectName() {
			t.Fatal(propObject2, o1, propObject2.ObjectName(), o1.ObjectName())
		}
	})
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")
	test.SetPropObject2(o1)
	test.PropObject2Changed(o1)
	if test.PropObject2().ObjectName() != o1.ObjectName() {
		t.Fatal("PropObject2")
	}

	//TODO: ConnectPropObject3Changed *QVariant

	test.ConnectPropEnumChanged(func(propEnum core.Qt__Key) {
		if propEnum != e0 {
			t.Fatal(propEnum, e0)
		}
	})
	test.SetPropEnum(e0)
	test.PropEnumChanged(e0)
	if test.PropEnum() != e0 {
		t.Fatal("PropEnum")
	}

	test.ConnectPropErrorChanged(func(propError error) {
		if propError.Error() != e1.Error() {
			t.Fatal(propError, e1, propError.Error(), e1.Error())
		}
	})
	test.SetPropError(e1)
	test.PropErrorChanged(e1)
	if test.PropError().Error() != e1.Error() {
		t.Fatal("PropError")
	}

	sTest := NewOtherTestStruct(nil)
	sTest.SetObjectName("test")
	test.ConnectPropReturnTestChanged(func(propReturnTest *otherTestStruct) {
		if propReturnTest.ObjectName() != sTest.ObjectName() {
			t.Fatal(propReturnTest, sTest, propReturnTest.ObjectName(), sTest.ObjectName())
		}
	})
	test.SetPropReturnTest(sTest)
	test.PropReturnTestChanged(sTest)
	if test.PropReturnTest().ObjectName() != sTest.ObjectName() {
		t.Fatal("PropReturnTest")
	}

	test.ConnectPropReturnTest2Changed(func(propReturnTest2 *otherTestStruct) {
		if propReturnTest2.ObjectName() != sTest.ObjectName() {
			t.Fatal(propReturnTest2, sTest, propReturnTest2.ObjectName(), sTest.ObjectName())
		}
	})
	test.SetPropReturnTest2(sTest)
	test.PropReturnTest2Changed(sTest)
	if test.PropReturnTest2().ObjectName() != sTest.ObjectName() {
		t.Fatal("PropReturnTest2")
	}

	test.ConnectPropBoolSubChanged(func(propBoolSub bool) {
		if propBoolSub != b1 {
			t.Fatal(propBoolSub, b1)
		}
	})
	test.SetPropBoolSub(b1)
	test.PropBoolSubChanged(b1)
	if test.IsPropBoolSub() != b1 {
		t.Fatal("IsPropBoolSub")
	}
}

func TestPropertiesList(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectPropListBoolChanged(func(propBool []bool) {
		if propBool[0] != b1 {
			t.Fatal(propBool, b1)
		}
	})
	test.SetPropListBool(b1L)
	test.PropListBoolChanged(b1L)
	if test.PropListBool()[0] != b1 {
		t.Fatal("IsPropListBool")
	}

	//_ int8             `property:"propInt8"`
	//_ uint8            `property:"propInt82"`

	test.ConnectPropListInt16Changed(func(propInt16 []int16) {
		if propInt16[0] != i0 {
			t.Fatal(propInt16, i0)
		}
	})
	test.SetPropListInt16(i0L)
	test.PropListInt16Changed(i0L)
	if test.PropListInt16()[0] != i0 {
		t.Fatal("PropListInt16")
	}

	test.ConnectPropListInt162Changed(func(propInt162 []uint16) {
		if propInt162[0] != i1 {
			t.Fatal(propInt162, i1)
		}
	})
	test.SetPropListInt162(i1L)
	test.PropListInt162Changed(i1L)
	if test.PropListInt162()[0] != i1 {
		t.Fatal("PropListInt162")
	}

	test.ConnectPropListInt32Changed(func(propInt32 []int) {
		if propInt32[0] != i2 {
			t.Fatal(propInt32, i2)
		}
	})
	test.SetPropListInt32(i2L)
	test.PropListInt32Changed(i2L)
	if test.PropListInt32()[0] != i2 {
		t.Fatal("PropListInt32")
	}

	test.ConnectPropListInt322Changed(func(propInt322 []uint) {
		if propInt322[0] != i3 {
			t.Fatal(propInt322, i3)
		}
	})
	test.SetPropListInt322(i3L)
	test.PropListInt322Changed(i3L)
	if test.PropListInt322()[0] != i3 {
		t.Fatal("PropListInt322")
	}

	test.ConnectPropListIntChanged(func(propInt []int) {
		if propInt[0] != i2 {
			t.Fatal(propInt, i2)
		}
	})
	test.SetPropListInt(i2L)
	test.PropListIntChanged(i2L)
	if test.PropListInt()[0] != i2 {
		t.Fatal("PropListInt")
	}

	test.ConnectPropListInt2Changed(func(propInt2 []uint) {
		if propInt2[0] != i3 {
			t.Fatal(propInt2, i3)
		}
	})
	test.SetPropListInt2(i3L)
	test.PropListInt2Changed(i3L)
	if test.PropListInt2()[0] != i3 {
		t.Fatal("PropListInt2")
	}

	test.ConnectPropListInt64Changed(func(propInt64 []int64) {
		if propInt64[0] != i6 {
			t.Fatal(propInt64, i6)
		}
	})
	test.SetPropListInt64(i6L)
	test.PropListInt64Changed(i6L)
	if test.PropListInt64()[0] != i6 {
		t.Fatal("PropListInt64")
	}

	test.ConnectPropListInt642Changed(func(propInt642 []uint64) {
		if propInt642[0] != i7 {
			t.Fatal(propInt642, i7)
		}
	})
	test.SetPropListInt642(i7L)
	test.PropListInt642Changed(i7L)
	if test.PropListInt642()[0] != i7 {
		t.Fatal("PropListInt642")
	}

	test.ConnectPropListFloatChanged(func(propFloat []float32) {
		if propFloat[0] != f0 {
			t.Fatal(propFloat, f0)
		}
	})
	test.SetPropListFloat(f0L)
	test.PropListFloatChanged(f0L)
	if test.PropListFloat()[0] != f0 {
		t.Fatal("PropListFloat")
	}

	test.ConnectPropListFloat2Changed(func(propFloat2 []float64) {
		if propFloat2[0] != f1 {
			t.Fatal(propFloat2, f1)
		}
	})
	test.SetPropListFloat2(f1L)
	test.PropListFloat2Changed(f1L)
	if test.PropListFloat2()[0] != f1 {
		t.Fatal("PropListFloat2")
	}

	test.ConnectPropListPointerChanged(func(propPointer []uintptr) {
		if int(propPointer[0]) != int(p0) {
			t.Fatal(propPointer, p0, int(propPointer[0]), int(p0))
		}
	})
	test.SetPropListPointer(p0L)
	test.PropListPointerChanged(p0L)
	if int(test.PropListPointer()[0]) != int(p0) {
		t.Fatal("PropListPointer")
	}

	test.ConnectPropListPointer2Changed(func(propPointer2 []unsafe.Pointer) {
		if int(uintptr(propPointer2[0])) != int(uintptr(p1)) {
			t.Fatal(propPointer2, p1, int(uintptr(propPointer2[0])), int(uintptr(p1)))
		}
	})
	test.SetPropListPointer2(p1L)
	test.PropListPointer2Changed(p1L)
	if int(uintptr(test.PropListPointer2()[0])) != int(uintptr(p1)) {
		t.Fatal("PropListPointer2")
	}

	test.ConnectPropListObjectChanged(func(propObject []*core.QVariant) {
		if propObject[0].ToString() != o0.ToString() { //TODO:
			t.Fatal(propObject, o0, propObject[0].ToString(), o0.ToString())
		}
	})
	o0 = core.NewQVariant14("test")
	test.SetPropListObject([]*core.QVariant{o0})
	test.PropListObjectChanged([]*core.QVariant{o0})
	if test.PropListObject()[0].ToString() != o0.ToString() {
		t.Fatal("PropListObject")
	}

	test.ConnectPropListObject2Changed(func(propObject2 []*core.QObject) {
		if propObject2[0].ObjectName() != o1.ObjectName() {
			t.Fatal(propObject2, o1, propObject2[0].ObjectName(), o1.ObjectName())
		}
	})
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")
	test.SetPropListObject2([]*core.QObject{o1})
	test.PropListObject2Changed([]*core.QObject{o1})
	if test.PropListObject2()[0].ObjectName() != o1.ObjectName() {
		t.Fatal("PropListObject2")
	}

	//TODO: ConnectPropListObject3Changed *QVariant

	test.ConnectPropListEnumChanged(func(propEnum []core.Qt__Key) {
		if propEnum[0] != e0 {
			t.Fatal(propEnum, e0)
		}
	})
	test.SetPropListEnum(e0L)
	test.PropListEnumChanged(e0L)
	if test.PropListEnum()[0] != e0 {
		t.Fatal("PropListEnum")
	}

	test.ConnectPropListErrorChanged(func(propError []error) {
		if propError[0].Error() != e1.Error() {
			t.Fatal(propError, e1, propError[0].Error(), e1.Error())
		}
	})
	test.SetPropListError(e1L)
	test.PropListErrorChanged(e1L)
	if test.PropListError()[0].Error() != e1.Error() {
		t.Fatal("PropListError")
	}

	sTest := NewOtherTestStruct(nil)
	sTest.SetObjectName("test")
	test.ConnectPropListReturnTestChanged(func(propReturnTest []*otherTestStruct) {
		if propReturnTest[0].ObjectName() != sTest.ObjectName() {
			t.Fatal(propReturnTest, sTest, propReturnTest[0].ObjectName(), sTest.ObjectName())
		}
	})
	test.SetPropListReturnTest([]*otherTestStruct{sTest})
	test.PropListReturnTestChanged([]*otherTestStruct{sTest})
	if test.PropListReturnTest()[0].ObjectName() != sTest.ObjectName() {
		t.Fatal("PropListReturnTest")
	}

	test.ConnectPropListReturnTest2Changed(func(propReturnTest2 []*otherTestStruct) {
		if propReturnTest2[0].ObjectName() != sTest.ObjectName() {
			t.Fatal(propReturnTest2, sTest, propReturnTest2[0].ObjectName(), sTest.ObjectName())
		}
	})
	test.SetPropListReturnTest2([]*otherTestStruct{sTest})
	test.PropListReturnTest2Changed([]*otherTestStruct{sTest})
	if test.PropListReturnTest2()[0].ObjectName() != sTest.ObjectName() {
		t.Fatal("PropListReturnTest2")
	}
}

func TestPropertiesMap(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectPropMapBoolChanged(func(propBool map[int]bool) {
		if propBool[0] != b1 {
			t.Fatal(propBool, b1)
		}
	})
	test.SetPropMapBool(b1M)
	test.PropMapBoolChanged(b1M)
	if test.PropMapBool()[0] != b1 {
		t.Fatal("IsPropMapBool")
	}

	//_ int8             `property:"propInt8"`
	//_ uint8            `property:"propInt82"`

	test.ConnectPropMapInt16Changed(func(propInt16 map[int]int16) {
		if propInt16[0] != i0 {
			t.Fatal(propInt16, i0)
		}
	})
	test.SetPropMapInt16(i0M)
	test.PropMapInt16Changed(i0M)
	if test.PropMapInt16()[0] != i0 {
		t.Fatal("PropMapInt16")
	}

	test.ConnectPropMapInt162Changed(func(propInt162 map[int]uint16) {
		if propInt162[0] != i1 {
			t.Fatal(propInt162, i1)
		}
	})
	test.SetPropMapInt162(i1M)
	test.PropMapInt162Changed(i1M)
	if test.PropMapInt162()[0] != i1 {
		t.Fatal("PropMapInt162")
	}

	test.ConnectPropMapInt32Changed(func(propInt32 map[int]int) {
		if propInt32[0] != i2 {
			t.Fatal(propInt32, i2)
		}
	})
	test.SetPropMapInt32(i2M)
	test.PropMapInt32Changed(i2M)
	if test.PropMapInt32()[0] != i2 {
		t.Fatal("PropMapInt32")
	}

	test.ConnectPropMapInt322Changed(func(propInt322 map[int]uint) {
		if propInt322[0] != i3 {
			t.Fatal(propInt322, i3)
		}
	})
	test.SetPropMapInt322(i3M)
	test.PropMapInt322Changed(i3M)
	if test.PropMapInt322()[0] != i3 {
		t.Fatal("PropMapInt322")
	}

	test.ConnectPropMapIntChanged(func(propInt map[int]int) {
		if propInt[0] != i2 {
			t.Fatal(propInt, i2)
		}
	})
	test.SetPropMapInt(i2M)
	test.PropMapIntChanged(i2M)
	if test.PropMapInt()[0] != i2 {
		t.Fatal("PropMapInt")
	}

	test.ConnectPropMapInt2Changed(func(propInt2 map[int]uint) {
		if propInt2[0] != i3 {
			t.Fatal(propInt2, i3)
		}
	})
	test.SetPropMapInt2(i3M)
	test.PropMapInt2Changed(i3M)
	if test.PropMapInt2()[0] != i3 {
		t.Fatal("PropMapInt2")
	}

	test.ConnectPropMapInt64Changed(func(propInt64 map[int]int64) {
		if propInt64[0] != i6 {
			t.Fatal(propInt64, i6)
		}
	})
	test.SetPropMapInt64(i6M)
	test.PropMapInt64Changed(i6M)
	if test.PropMapInt64()[0] != i6 {
		t.Fatal("PropMapInt64")
	}

	test.ConnectPropMapInt642Changed(func(propInt642 map[int]uint64) {
		if propInt642[0] != i7 {
			t.Fatal(propInt642, i7)
		}
	})
	test.SetPropMapInt642(i7M)
	test.PropMapInt642Changed(i7M)
	if test.PropMapInt642()[0] != i7 {
		t.Fatal("PropMapInt642")
	}

	test.ConnectPropMapFloatChanged(func(propFloat map[int]float32) {
		if propFloat[0] != f0 {
			t.Fatal(propFloat, f0)
		}
	})
	test.SetPropMapFloat(f0M)
	test.PropMapFloatChanged(f0M)
	if test.PropMapFloat()[0] != f0 {
		t.Fatal("PropMapFloat")
	}

	test.ConnectPropMapFloat2Changed(func(propFloat2 map[int]float64) {
		if propFloat2[0] != f1 {
			t.Fatal(propFloat2, f1)
		}
	})
	test.SetPropMapFloat2(f1M)
	test.PropMapFloat2Changed(f1M)
	if test.PropMapFloat2()[0] != f1 {
		t.Fatal("PropMapFloat2")
	}

	test.ConnectPropMapStringChanged(func(propMapString map[int]string) {
		if propMapString[0] != s0 {
			t.Fatal(propMapString, s0)
		}
	})
	test.SetPropMapString(s0M)
	test.PropMapStringChanged(s0M)
	if test.PropMapString()[0] != s0 {
		t.Fatal("PropMapString")
	}

	test.ConnectPropMapString2Changed(func(propMapString2 map[int][]string) {
		if strings.Join(propMapString2[0], "") != strings.Join(s1, "") {
			t.Fatal(propMapString2, s1, strings.Join(propMapString2[0], ""), strings.Join(s1, ""))
		}
	})
	test.SetPropMapString2(s1M)
	test.PropMapString2Changed(s1M)
	if strings.Join(test.PropMapString2()[0], "") != strings.Join(s1, "") {
		t.Fatal("PropMapString2")
	}

	test.ConnectPropMapPointerChanged(func(propPointer map[int]uintptr) {
		if int(propPointer[0]) != int(p0) {
			t.Fatal(propPointer, p0, int(propPointer[0]), int(p0))
		}
	})
	test.SetPropMapPointer(p0M)
	test.PropMapPointerChanged(p0M)
	if int(test.PropMapPointer()[0]) != int(p0) {
		t.Fatal("PropMapPointer")
	}

	test.ConnectPropMapPointer2Changed(func(propPointer2 map[int]unsafe.Pointer) {
		if int(uintptr(propPointer2[0])) != int(uintptr(p1)) {
			t.Fatal(propPointer2, p1, int(uintptr(propPointer2[0])), int(uintptr(p1)))
		}
	})
	test.SetPropMapPointer2(p1M)
	test.PropMapPointer2Changed(p1M)
	if int(uintptr(test.PropMapPointer2()[0])) != int(uintptr(p1)) {
		t.Fatal("PropMapPointer2")
	}

	test.ConnectPropMapObjectChanged(func(propObject map[int]*core.QVariant) {
		if propObject[0].ToString() != o0.ToString() { //TODO:
			t.Fatal(propObject, o0, propObject[0].ToString(), o0.ToString())
		}
	})
	o0 = core.NewQVariant14("test")
	test.SetPropMapObject(map[int]*core.QVariant{0: o0})
	test.PropMapObjectChanged(map[int]*core.QVariant{0: o0})
	if test.PropMapObject()[0].ToString() != o0.ToString() {
		t.Fatal("PropMapObject")
	}

	test.ConnectPropMapObject2Changed(func(propObject2 map[int]*core.QObject) {
		if propObject2[0].ObjectName() != o1.ObjectName() {
			t.Fatal(propObject2, o1, propObject2[0].ObjectName(), o1.ObjectName())
		}
	})
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")
	test.SetPropMapObject2(map[int]*core.QObject{0: o1})
	test.PropMapObject2Changed(map[int]*core.QObject{0: o1})
	if test.PropMapObject2()[0].ObjectName() != o1.ObjectName() {
		t.Fatal("PropMapObject2")
	}

	//TODO: ConnectPropMapObject3Changed *QVariant

	test.ConnectPropMapEnumChanged(func(propEnum map[int]core.Qt__Key) {
		if propEnum[0] != e0 {
			t.Fatal(propEnum, e0)
		}
	})
	test.SetPropMapEnum(e0M)
	test.PropMapEnumChanged(e0M)
	if test.PropMapEnum()[0] != e0 {
		t.Fatal("PropMapEnum")
	}

	test.ConnectPropMapErrorChanged(func(propError map[int]error) {
		if propError[0].Error() != e1.Error() {
			t.Fatal(propError, e1, propError[0].Error(), e1.Error())
		}
	})
	test.SetPropMapError(e1M)
	test.PropMapErrorChanged(e1M)
	if test.PropMapError()[0].Error() != e1.Error() {
		t.Fatal("PropMapError")
	}

	sTest := NewOtherTestStruct(nil)
	sTest.SetObjectName("test")
	test.ConnectPropMapReturnTestChanged(func(propReturnTest map[int]*otherTestStruct) {
		if propReturnTest[0].ObjectName() != sTest.ObjectName() {
			t.Fatal(propReturnTest, sTest, propReturnTest[0].ObjectName(), sTest.ObjectName())
		}
	})
	test.SetPropMapReturnTest(map[int]*otherTestStruct{0: sTest})
	test.PropMapReturnTestChanged(map[int]*otherTestStruct{0: sTest})
	if test.PropMapReturnTest()[0].ObjectName() != sTest.ObjectName() {
		t.Fatal("PropMapReturnTest")
	}

	test.ConnectPropMapReturnTest2Changed(func(propReturnTest2 map[int]*otherTestStruct) {
		if propReturnTest2[0].ObjectName() != sTest.ObjectName() {
			t.Fatal(propReturnTest2, sTest, propReturnTest2[0].ObjectName(), sTest.ObjectName())
		}
	})
	test.SetPropMapReturnTest2(map[int]*otherTestStruct{0: sTest})
	test.PropMapReturnTest2Changed(map[int]*otherTestStruct{0: sTest})
	if test.PropMapReturnTest2()[0].ObjectName() != sTest.ObjectName() {
		t.Fatal("PropMapReturnTest2")
	}
}

func TestSignalInput(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectBoolSignalInput(func(v0, v1 bool) {
		if v0 != b0 || v1 != b1 {
			t.Fatal(v0, b0, v1, b1)
		}
	})

	test.ConnectInt8SignalInput(func(v0 string, v1 string) {
		if v0 != s0 || v1 != s0 {
			t.Fatal(v0, s0, v1, s0)
		}
	})

	test.ConnectInt16SignalInput(func(v0 int16, v1 uint16) {
		if v0 != i0 || v1 != i1 {
			t.Fatal(v0, i0, v1, i1)
		}
	})

	test.ConnectInt32SignalInput(func(v0 int, v1 uint) {
		if v0 != i2 || v1 != i3 {
			t.Fatal(v0, i2, v1, i3)
		}
	})

	test.ConnectIntSignalInput(func(v0 int, v1 uint) {
		if v0 != i4 || v1 != i5 {
			t.Fatal(v0, i4, v1, i5)
		}
	})

	test.ConnectInt64SignalInput(func(v0 int64, v1 uint64) {
		if v0 != i6 || v1 != i7 {
			t.Fatal(v0, i6, v1, i7)
		}
	})

	test.ConnectFloatSignalInput(func(v0 float32, v1 float64) {
		if v0 != f0 || v1 != f1 {
			t.Fatal(v0, f0, v1, f1)
		}
	})

	test.ConnectStringSignalInput(func(v0 string, v1 []string) {
		if v0 != s0 || strings.Join(v1, "") != strings.Join(s1, "") {
			t.Fatal(v0, s0, v1, s1)
		}
	})

	test.ConnectPointerSignalInput(func(v0 uintptr, v1 unsafe.Pointer) {
		if int(v0) != int(p0) || int(uintptr(v1)) != int(uintptr(p1)) {
			t.Fatal(v0, p0, v1, p1, int(v0), int(p0), int(uintptr(v1)), int(uintptr(p1)))
		}
	})

	test.ConnectObjectSignalInput(func(v0 *core.QVariant, v1 *core.QObject) {
		if v0.ToString() != o0.ToString() || v1.ObjectName() != o1.ObjectName() {
			t.Fatal(v0, o0, v1, o1)
		}
	})

	test.ConnectEnumSignalInput(func(v0 core.Qt__Key) {
		if v0 != e0 {
			t.Fatal(v0, e0)
		}
	})

	test.ConnectErrorSignalInput(func(v0 error) {
		if v0.Error() != e1.Error() {
			t.Fatal(v0, e1)
		}
	})

	test.BoolSignalInput(b0, b1)
	//test.Int8SignalInput(s0, s0)
	test.Int16SignalInput(i0, i1)
	test.Int32SignalInput(i2, i3)
	test.IntSignalInput(i4, i5)
	test.Int64SignalInput(i6, i7)
	test.FloatSignalInput(f0, f1)
	test.StringSignalInput(s0, s1)
	test.PointerSignalInput(p0, p1)

	o0 = core.NewQVariant14("test")
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")
	test.ObjectSignalInput(o0, o1)

	test.EnumSignalInput(e0)
	test.ErrorSignalInput(e1)
}

func TestSignalListInput(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectBoolSignalListInput(func(v0, v1 []bool) {
		if v0[0] != b0 || v1[0] != b1 {
			t.Fatal(v0, b0, v1, b1)
		}
	})

	test.ConnectInt16SignalListInput(func(v0 []int16, v1 []uint16) {
		if v0[0] != i0 || v1[0] != i1 {
			t.Fatal(v0, i0, v1, i1)
		}
	})

	test.ConnectInt32SignalListInput(func(v0 []int, v1 []uint) {
		if v0[0] != i2 || v1[0] != i3 {
			t.Fatal(v0, i2, v1, i3)
		}
	})

	test.ConnectIntSignalListInput(func(v0 []int, v1 []uint) {
		if v0[0] != i4 || v1[0] != i5 {
			t.Fatal(v0, i4, v1, i5)
		}
	})

	test.ConnectInt64SignalListInput(func(v0 []int64, v1 []uint64) {
		if v0[0] != i6 || v1[0] != i7 {
			t.Fatal(v0, i6, v1, i7)
		}
	})

	test.ConnectFloatSignalListInput(func(v0 []float32, v1 []float64) {
		if v0[0] != f0 || v1[0] != f1 {
			t.Fatal(v0, f0, v1, f1)
		}
	})

	test.ConnectPointerSignalListInput(func(v0 []uintptr, v1 []unsafe.Pointer) {
		if int(v0[0]) != int(p0) || int(uintptr(v1[0])) != int(uintptr(p1)) {
			t.Fatal(v0, p0, v1, p1, int(v0[0]), int(p0), int(uintptr(v1[0])), int(uintptr(p1)))
		}
	})

	test.ConnectObjectSignalListInput(func(v0 []*core.QVariant, v1 []*core.QObject) {
		if v0[0].ToString() != o0.ToString() || v1[0].ObjectName() != o1.ObjectName() {
			t.Fatal(v0, o0, v1, o1)
		}
	})

	test.ConnectEnumSignalListInput(func(v0 []core.Qt__Key) {
		if v0[0] != e0 {
			t.Fatal(v0, e0)
		}
	})

	test.ConnectErrorSignalListInput(func(v0 []error) {
		if v0[0].Error() != e1.Error() {
			t.Fatal(v0, e1)
		}
	})

	test.BoolSignalListInput(b0L, b1L)
	//test.Int8SignalListInput(s0, s0)
	test.Int16SignalListInput(i0L, i1L)
	test.Int32SignalListInput(i2L, i3L)
	test.IntSignalListInput(i4L, i5L)
	test.Int64SignalListInput(i6L, i7L)
	test.FloatSignalListInput(f0L, f1L)
	test.PointerSignalListInput(p0L, p1L)

	o0 = core.NewQVariant14("test")
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")
	test.ObjectSignalListInput([]*core.QVariant{o0}, []*core.QObject{o1})

	test.EnumSignalListInput(e0L)
	test.ErrorSignalListInput(e1L)
}

func TestSlotInput(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectBoolSlotInput(func(v0, v1 bool) {
		if v0 != b0 || v1 != b1 {
			t.Fatal(v0, b0, v1, b1)
		}
	})

	test.ConnectInt8SlotInput(func(v0 string, v1 string) {
		if v0 != s0 || v1 != s0 {
			t.Fatal(v0, s0, v1, s0)
		}
	})

	test.ConnectInt16SlotInput(func(v0 int16, v1 uint16) {
		if v0 != i0 || v1 != i1 {
			t.Fatal(v0, i0, v1, i1)
		}
	})

	test.ConnectInt32SlotInput(func(v0 int, v1 uint) {
		if v0 != i2 || v1 != i3 {
			t.Fatal(v0, i2, v1, i3)
		}
	})

	test.ConnectIntSlotInput(func(v0 int, v1 uint) {
		if v0 != i4 || v1 != i5 {
			t.Fatal(v0, i4, v1, i5)
		}
	})

	test.ConnectInt64SlotInput(func(v0 int64, v1 uint64) {
		if v0 != i6 || v1 != i7 {
			t.Fatal(v0, i6, v1, i7)
		}
	})

	test.ConnectFloatSlotInput(func(v0 float32, v1 float64) {
		if v0 != f0 || v1 != f1 {
			t.Fatal(v0, f0, v1, f1)
		}
	})

	test.ConnectStringSlotInput(func(v0 string, v1 []string) {
		if v0 != s0 || strings.Join(v1, "") != strings.Join(s1, "") {
			t.Fatal(v0, s0, v1, s1)
		}
	})

	test.ConnectPointerSlotInput(func(v0 uintptr, v1 unsafe.Pointer) {
		if int(v0) != int(p0) || int(uintptr(v1)) != int(uintptr(p1)) {
			t.Fatal(v0, p0, v1, p1, int(v0), int(p0), int(uintptr(v1)), int(uintptr(p1)))
		}
	})

	test.ConnectObjectSlotInput(func(v0 *core.QVariant, v1 *core.QObject) {
		if v0.ToString() != o0.ToString() || v1.ObjectName() != o1.ObjectName() {
			t.Fatal(v0, o0, v1, o1)
		}
	})

	test.ConnectEnumSlotInput(func(v0 core.Qt__Key) {
		if v0 != e0 {
			t.Fatal(v0, e0)
		}
	})

	test.ConnectErrorSlotInput(func(v0 error) {
		if v0.Error() != e1.Error() {
			t.Fatal(v0, e1)
		}
	})

	test.BoolSlotInput(b0, b1)
	//test.Int8SlotInput(s0, s0)
	test.Int16SlotInput(i0, i1)
	test.Int32SlotInput(i2, i3)
	test.IntSlotInput(i4, i5)
	test.Int64SlotInput(i6, i7)
	test.FloatSlotInput(f0, f1)
	test.StringSlotInput(s0, s1)
	test.PointerSlotInput(p0, p1)

	o0 = core.NewQVariant14("test")
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")
	test.ObjectSlotInput(o0, o1)

	test.EnumSlotInput(e0)
	test.ErrorSlotInput(e1)
}

func TestSlotOutput(t *testing.T) {
	test := NewTestStruct(nil)

	test.ConnectBoolSlotOutput(func(v0 bool) bool { return v0 })
	test.ConnectBoolSlotOutput2(func(v0 bool) bool { return v0 })

	test.ConnectInt8SlotOutput(func(v0 string) string { return v0 })
	test.ConnectInt8SlotOutput2(func(v0 string) string { return v0 })

	test.ConnectInt16SlotOutput(func(v0 int16) int16 { return v0 })
	test.ConnectInt16SlotOutput2(func(v0 uint16) uint16 { return v0 })

	test.ConnectInt32SlotOutput(func(v0 int) int { return v0 })
	test.ConnectInt32SlotOutput2(func(v0 uint) uint { return v0 })

	test.ConnectIntSlotOutput(func(v0 int) int { return v0 })
	test.ConnectIntSlotOutput2(func(v0 uint) uint { return v0 })

	test.ConnectInt64SlotOutput(func(v0 int64) int64 { return v0 })
	test.ConnectInt64SlotOutput2(func(v0 uint64) uint64 { return v0 })

	test.ConnectFloatSlotOutput(func(v0 float32) float32 { return v0 })
	test.ConnectFloatSlotOutput2(func(v0 float64) float64 { return v0 })

	test.ConnectStringSlotOutput(func(v0 string) string { return v0 })
	test.ConnectStringSlotOutput2(func(v0 []string) []string { return v0 })

	test.ConnectPointerSlotOutput(func(v0 uintptr) uintptr { return v0 })
	test.ConnectPointerSlotOutput2(func(v0 unsafe.Pointer) unsafe.Pointer { return v0 })

	test.ConnectObjectSlotOutput(func(v0 *core.QVariant) *core.QVariant { return v0 })
	test.ConnectObjectSlotOutput2(func(v0 *core.QObject) *core.QObject { return v0 })

	test.ConnectEnumSlotOutput(func(v0 core.Qt__Key) core.Qt__Key { return v0 })
	test.ConnectErrorSlotOutput(func(v0 error) error { return v0 })

	test.ConnectReturnTest(func(v0 *testStruct) *testStruct { return v0 })
	test.ConnectReturnTest2(func(v0 *testStruct) *testStruct { return v0 })

	test.ConnectReturnName(func(a0 string) string { return a0 })
	test.ConnectReturnName2(func(a0 int, a1 int) int { return a0 + a1 })

	if test.BoolSlotOutput(b0) != b0 {
		t.Fatal("BoolSlotOutput")
	}
	if test.BoolSlotOutput2(b1) != b1 {
		t.Fatal("BoolSlotOutput2")
	}

	/*
		if test.Int8SlotOutput(s0) != s0 {
			t.Fatal("Int8SlotOutput")
		}
		if test.Int8SlotOutput2(s0) != s0 {
			t.Fatal("Int8SlotOutput2")
		}
	*/

	if test.Int16SlotOutput(i0) != i0 {
		t.Fatal("Int16SlotOutput")
	}
	if test.Int16SlotOutput2(i1) != i1 {
		t.Fatal("Int16SlotOutput2")
	}

	if test.Int32SlotOutput(i2) != i2 {
		t.Fatal("Int32SlotOutput")
	}
	if test.Int32SlotOutput2(i3) != i3 {
		t.Fatal("Int32SlotOutput2")
	}

	if test.IntSlotOutput(i4) != i4 {
		t.Fatal("IntSlotOutput")
	}
	if test.IntSlotOutput2(i5) != i5 {
		t.Fatal("IntSlotOutput2")
	}

	if test.Int64SlotOutput(i6) != i6 {
		t.Fatal("Int64SlotOutput")
	}
	if test.Int64SlotOutput2(i7) != i7 {
		t.Fatal("Int64SlotOutput2")
	}

	if test.FloatSlotOutput(f0) != f0 {
		t.Fatal("FloatSlotOutput")
	}
	if test.FloatSlotOutput2(f1) != f1 {
		t.Fatal("FloatSlotOutput2")
	}

	if test.StringSlotOutput(s0) != s0 {
		t.Fatal("StringSlotOutput")
	}
	if strings.Join(test.StringSlotOutput2(s1), "") != strings.Join(s1, "") {
		t.Fatal("StringSlotOutput2")
	}

	if int(test.PointerSlotOutput(p0)) != int(p0) {
		t.Fatal("PointerSlotOutput")
	}
	if int(uintptr(test.PointerSlotOutput2(p1))) != int(uintptr(p1)) {
		t.Fatal("PointerSlotOutput2")
	}

	o0 = core.NewQVariant14("test")
	o1 = core.NewQObject(nil)
	o1.SetObjectName("test")

	if test.ObjectSlotOutput(o0).ToString() != o0.ToString() {
		t.Fatal("ObjectSlotOutput")
	}
	if test.ObjectSlotOutput2(o1).ObjectName() != o1.ObjectName() {
		t.Fatal("ObjectSlotOutput2")
	}

	if test.EnumSlotOutput(e0) != e0 {
		t.Fatal("EnumSlotOutput")
	}

	if test.ErrorSlotOutput(e1).Error() != e1.Error() {
		t.Fatal("ErrorSlotOutput")
	}

	test.SetObjectName("testName")
	if test.ReturnTest(test.ReturnTest(test)).ReturnTest(test).ObjectName() != test.ObjectName() {
		t.Fatal("ReturnTest")
	}
	if test.ReturnTest2(test.ReturnTest2(test)).ReturnTest2(test).ObjectName() != test.ObjectName() {
		t.Fatal("ReturnTest2")
	}

	if test.ReturnName(s0) != s0 {
		t.Fatal("ReturnName")
	}
	if test.ReturnName2(i2, i2) != i2*2 {
		t.Fatal("ReturnName2")
	}

	sTest := NewSubTestStruct(nil)
	sTest.ConnectReturnTest3(func(v0 *subTestStruct) *subTestStruct { return v0 })
	sTest.SetObjectName("testSubName")
	if sTest.ReturnTest3(sTest.ReturnTest3(sTest)).ReturnTest3(sTest).ObjectName() != sTest.ObjectName() {
		t.Fatal("ReturnName3")
	}
}

func TestSubSubStructs(t *testing.T) {
	sTest := NewSubSubTestStruct(nil)
	sTest.SetObjectName("testObjectName")

	sTest.ConnectStructPropertyTestChanged(func(StructPropertyTest *sub.SubTestStruct) {
		if StructPropertyTest.ObjectName() != sTest.ObjectName() {
			t.Fatal(StructPropertyTest.ObjectName(), "!=", sTest.ObjectName())
		}
	})

	sTest.ConnectStructSignalTest(func(v0 *sub.SubTestStruct) {
		if v0.ObjectName() != sTest.ObjectName() {
			t.Fatal(v0.ObjectName(), "!=", sTest.ObjectName())
		}
	})

	sTest.ConnectStructSlotTest(func(v0 *sub.SubTestStruct) *sub.SubTestStruct {
		if v0.ObjectName() != sTest.ObjectName() {
			t.Fatal(v0.ObjectName(), "!=", sTest.ObjectName())
		}
		return v0
	})

	sTest.SetStructPropertyTest(sTest)
	sTest.StructSignalTest(sTest)
	sTest.StructSlotTest(sTest)

	//

	sTest.SetSomeProperty(s0)
	if o := sTest.SomeProperty(); o != s0 {
		t.Fatal(sTest.SomeProperty(), "!=", s0)
	}
	/*
		if o := sTest.SubTestStruct_PTR().SomeProperty(); o != s0 {
			t.Fatal(o, len(o), "!=", s0, len(s0))
		}
	*/

	sTest.ConnectSomeSignal(func(v0 string) {
		if v0 != s0 {
			t.Fatal(v0, "!=", s0)
		}
	})
	sTest.SomeSignal(s0)
	//sTest.SubTestStruct_PTR().SomeSignal(s0)

	sTest.ConnectSomeSlot(func(v0 string) string {
		if v0 != s0 {
			t.Fatal(v0, "!=", s0)
		}
		return v0
	})
	if o := sTest.SomeSlot(s0); o != s0 {
		t.Fatal(sTest.SomeSlot(s0), "!=", s0)
	}

	/*
		if o := sTest.SubTestStruct_PTR().SomeSlot(s0); o != s0 {
			t.Fatal(sTest.SomeSlot(s0), "!=", s0)
		}
	*/
}

func (*testStruct) errorStringTest1(string, string, string) error { return errors.New("") }
func (*testStruct) errorStringTest2(string, string, error) string { return "" }
func (*testStruct) errorStringTest3(string, error, string) string { return "" }

func (*testStruct) errorStringTest4(string, string, string) []error { return []error{} }
func (*testStruct) errorStringTest5(string, string, []error) string { return "" }
func (*testStruct) errorStringTest6(string, []error, string) string { return "" }

func (*testStruct) errorStringTest7(string, string, string) map[error]error {
	return map[error]error{}
}
func (*testStruct) errorStringTest8(string, string, map[error]error) string { return "" }
func (*testStruct) errorStringTest9(string, map[error]error, string) string { return "" }

func (*testStruct) errorStringTest10(string, map[error]error, []error) error { return errors.New("") }
func (*testStruct) errorStringTest11(map[error]error, []error, error) string { return "" }
func (*testStruct) errorStringTest12([]error, error, string) map[error]error {
	return map[error]error{}
}

func (*testStruct) errorStringTest13(map[string]error, map[error]string, map[error]error) map[string]string {
	return map[string]string{}
}

func (*testStruct) errorStringTest14(map[error][]string, map[string][]string, map[error][]error, map[string][]error) {
}

func TestBoolPointer(t *testing.T) {
	r := xml.NewQXmlSimpleReader()
	r.ConnectFeature(func(name string, ok *bool) bool { *ok = true; return *ok })
	var ok bool
	if r := r.Feature("", &ok); r == false || ok == false {
		t.Fatal()
	}
}
