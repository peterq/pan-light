package sub

import (
	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/subsub"

	"a" //vendor
	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/b"
	abc "github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/c"
	. "github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/d"

	"ago" //vendor
	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/bgo"
	abcgo "github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/cgo"
	. "github.com/peterq/pan-light/qt/tool-chain/cmd/moc/test/sub/dgo"
)

//TODO:  duplicate struct name issue

type SubTestStruct struct {
	subsubcustom.SubSubTestStruct

	_ func() `constructor:"Init"`

	_ func(string)        `signal:"someSignal"`
	_ func(string) string `slot:"someSlot"`
	_ string              `property:"someProperty"`

	_ func(StructSameGo)  `signal:"sameSignal"`
	_ func(*StructSameGo) `signal:"sameSignalR"`

	_ func(b.StructSubGoB)  `signal:"subSignalb1"`
	_ func(*b.StructSubGoB) `signal:"subSignalb1R"`

	_ func(b.StructSubMocB)  `signal:"subSignalb2"`
	_ func(*b.StructSubMocB) `signal:"subSignalb2R"`

	_ func(bgo.StructSubGoBGo)  `signal:"subSignalb3"`
	_ func(*bgo.StructSubGoBGo) `signal:"subSignalb3R"`

	_ func(abc.StructSubGoC)  `signal:"subSignalc1"`
	_ func(*abc.StructSubGoC) `signal:"subSignalc1R"`

	_ func(abc.StructSubMocC)  `signal:"subSignalc2"`
	_ func(*abc.StructSubMocC) `signal:"subSignalc2R"`

	_ func(abcgo.StructSubGoCGo)  `signal:"subSignalc3"`
	_ func(*abcgo.StructSubGoCGo) `signal:"subSignalc3R"`

	_ func(StructSubGoD)  `signal:"subSignald1"`
	_ func(*StructSubGoD) `signal:"subSignald1R"`

	_ func(StructSubMocD)  `signal:"subSignald2"`
	_ func(*StructSubMocD) `signal:"subSignald2R"`

	_ func(StructSubGoDGo)  `signal:"subSignald3"`
	_ func(*StructSubGoDGo) `signal:"subSignald3R"`
}

type SubTestStructNotUsed struct {
	subsubcustom.SubSubTestStruct

	_ func(a.StructSubGoA)  `signal:"vendorSignal1"`
	_ func(*a.StructSubGoA) `signal:"vendorSignal1R"`

	_ func(a.StructSubMocA)  `signal:"vendorSignal2"`
	_ func(*a.StructSubMocA) `signal:"vendorSignal2R"`

	_ func(ago.StructSubGoAGo)  `signal:"vendorSignal3"`
	_ func(*ago.StructSubGoAGo) `signal:"vendorSignal3R"`
}

func (s *SubTestStruct) Init() {
	s.SubSubConstructorProperty++
}

type StructSameGo struct {
	somestring string
}
