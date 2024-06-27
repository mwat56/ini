/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	inFileName  = `testIn.ini`
	outFilename = `testOut.ini`
)

func prepSectionList() *TSectionList {
	sl := NewSectionList()
	sl.AddSectionKey("", "key0", "")
	sl.AddSectionKey("s2", "float", "12345.6789")
	sl.AddSectionKey("s1", "bool", "nada")
	sl.AddSectionKey("s4", "uint", "1234567890")
	sl.AddSectionKey("s3", "int", "-12345")

	return sl
} // prepSectionList()

func Test_removeQuotes(t *testing.T) {
	si1, ws1 := "'this is a text'", "this is a text"
	si2, ws2 := " \" this is a text \" ", "this is a text"
	si3, ws3 := " \" this is a text ' ", "\" this is a text '"
	si4, ws4 := " this is a text ", "this is a text"
	si5, ws5 := " this is a text ' ", "this is a text '"
	tests := []struct {
		name        string
		args        string
		wantRString string
	}{
		{"1", si1, ws1},
		{"2", si2, ws2},
		{"3", si3, ws3},
		{"4", si4, ws4},
		{"5", si5, ws5},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRString := removeQuotes(tt.args); gotRString != tt.wantRString {
				t.Errorf("%q removeQuotes() = [%s], want [%s]",
					tt.name, gotRString, tt.wantRString)
			}
		})
	}
} // Test_removeQuotes

func TestTSectionList_addSection(t *testing.T) {
	sl := NewSectionList()
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"1", "SectTest1", true},
		{"2", "SectTest1", true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// sl := &tt.fields
			if got := sl.addSection(tt.args); got != tt.want {
				t.Errorf("%q TSectionList.addSection() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTIniList_addSection()

func TestTSectionList_AddSectionKey(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   string
	}
	sl := NewSectionList()
	tests := []struct {
		name    string
		args    tArgs
		wantROK bool
	}{
		{"1", tArgs{"", "k1", "v1"}, true},
		{"2", tArgs{"s2", "k2", "v2"}, true},
		{"3", tArgs{"s2", "k3", "v3"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// sl := &TSectionList{
			// 	defSect:  tt.fields.defSect,
			// 	fName:    tt.fields.fName,
			// 	secOrder: tt.fields.secOrder,
			// 	sections: tt.fields.sections,
			// }
			if gotROK := sl.AddSectionKey(tt.args.aSection, tt.args.aKey, tt.args.aValue); gotROK != tt.wantROK {
				t.Errorf("%q: TSectionList.AddSectionKey() = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSectionList_AddSectionKey()

//

func TestTSectionList_AsBool(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}
	sl := NewSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "1")
	_ = sl.AddSectionKey("", "key2", "2")
	_ = sl.AddSectionKey("", "key3", "funny")
	_ = sl.AddSectionKey("", "key4", "nightmare")
	_ = sl.AddSectionKey("", "key5", "talisman")
	tests := []struct {
		name  string
		args  tArgs
		want  bool
		want1 bool
	}{
		{"1", tArgs{"", ""}, false, false},
		{"2", tArgs{"", "key0"}, false, false},
		{"3", tArgs{"", "key1"}, true, true},
		{"4", tArgs{"", "key2"}, false, false},
		{"5", tArgs{"", "key3"}, false, true},
		{"6", tArgs{"", "key4"}, false, true},
		{"7", tArgs{"", "key5"}, true, true},
		{"8", tArgs{"", "n.a."}, false, false},
		{"9", tArgs{"n.a.", "-0"}, false, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsBool(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsBool() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsBool() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsBool()

//

func TestTSectionList_AsFloat32(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "123.456")
	_ = sl.AddSectionKey("", "key2", "0.0")
	_ = sl.AddSectionKey("", "key3", "-123.456")
	_ = sl.AddSectionKey("", "key4", "NaN")
	_ = sl.AddSectionKey("", "key5", "five dot five")
	tests := []struct {
		name     string
		args     tArgs
		wantRVal float32
		wantROK  bool
	}{
		{"1", tArgs{"", ""}, float32(0), false},
		{"2", tArgs{"", "key0"}, float32(0), false},
		{"3", tArgs{"", "key1"}, float32(123.456), true},
		{"4", tArgs{"", "key2"}, float32(0.0), true},
		{"5", tArgs{"", "key3"}, float32(-123.456), true},
		{"6", tArgs{"", "key4"}, float32(0), false},
		{"7", tArgs{"", "key5"}, float32(0), false},
		{"8", tArgs{"", "n.a."}, float32(0), false},
		{"9", tArgs{"n.a.", "-0.0"}, float32(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := sl.AsFloat32(tt.args.aSection, tt.args.aKey)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q: TSectionList.AsFloat32() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q: TSectionList.AsFloat32() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSectionList_AsFloat32()

func TestTSectionList_AsFloat64(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "1.7976931348623157e+308")
	_ = sl.AddSectionKey("", "key2", "0.0")
	_ = sl.AddSectionKey("", "key3", "-1.7976931348623157e+308")
	_ = sl.AddSectionKey("", "key4", "NaN")
	_ = sl.AddSectionKey("", "key5", "five dot five")

	tests := []struct {
		name     string
		args     tArgs
		wantRVal float64
		wantROK  bool
	}{
		{"1", tArgs{"", ""}, float64(0), false},
		{"2", tArgs{"", "key0"}, float64(0), false},
		{"3", tArgs{"", "key1"}, float64(1.7976931348623157e+308), true},
		{"4", tArgs{"", "key2"}, float64(0.0), true},
		{"5", tArgs{"", "key3"}, float64(-1.7976931348623157e+308), true},
		{"6", tArgs{"", "key4"}, float64(0), false},
		{"7", tArgs{"", "key5"}, float64(0), false},
		{"8", tArgs{"", "n.a."}, float64(0), false},
		{"9", tArgs{"n.a.", "-0.0"}, float64(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := sl.AsFloat64(tt.args.aSection, tt.args.aKey)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q: TSectionList.AsFloat64() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q: TSectionList.AsFloat64() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSectionList_AsFloat64()

//

func TestTSectionList_AsInt(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "-2147483648")
	_ = sl.AddSectionKey("", "key3", "2147483647")
	_ = sl.AddSectionKey("", "key4", "9223372036854775808") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")                  // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")                // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  int
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, int(0), false},
		{"2", tArgs{"", "key1"}, int(0), true},
		{"3", tArgs{"", "key2"}, int(-2147483648), true},
		{"4", tArgs{"", "key3"}, int(2147483647), true},
		{"5", tArgs{"", "key4"}, int(0), false},
		{"6", tArgs{"", "key5"}, int(-1), true},
		{"7", tArgs{"", "key6"}, int(0), false},
		{"8", tArgs{"", "n.a."}, int(0), false},
		{"9", tArgs{"n.a.", "-0"}, int(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsInt(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsInt() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsInt() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsInt()

func TestTSectionList_AsInt8(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "-128")
	_ = sl.AddSectionKey("", "key3", "127")
	_ = sl.AddSectionKey("", "key4", "128")  // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")   // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.") // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  int8
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, int8(0), false},
		{"2", tArgs{"", "key1"}, int8(0), true},
		{"3", tArgs{"", "key2"}, int8(-128), true},
		{"4", tArgs{"", "key3"}, int8(127), true},
		{"5", tArgs{"", "key4"}, int8(0), false},
		{"6", tArgs{"", "key5"}, int8(-1), true},
		{"7", tArgs{"", "key6"}, int8(0), false},
		{"8", tArgs{"", "n.a."}, int8(0), false},
		{"9", tArgs{"n.a.", "8"}, int8(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsInt8(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsInt8() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsInt8() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsInt8()

func TestTSectionList_AsInt16(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "-32768")
	_ = sl.AddSectionKey("", "key3", "32767")
	_ = sl.AddSectionKey("", "key4", "32768") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")    // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")  // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  int16
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, int16(0), false},
		{"2", tArgs{"", "key1"}, int16(0), true},
		{"3", tArgs{"", "key2"}, int16(-32768), true},
		{"4", tArgs{"", "key3"}, int16(32767), true},
		{"5", tArgs{"", "key4"}, int16(0), false},
		{"6", tArgs{"", "key5"}, int16(-1), true},
		{"7", tArgs{"", "key6"}, int16(0), false},
		{"8", tArgs{"", "n.a."}, int16(0), false},
		{"9", tArgs{"n.a.", "16"}, int16(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsInt16(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsInt16() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsInt16() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsInt16()

func TestTSectionList_AsInt32(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "-2147483648")
	_ = sl.AddSectionKey("", "key3", "2147483647")
	_ = sl.AddSectionKey("", "key4", "2147483648") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")         // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")       // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  int32
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, int32(0), false},
		{"2", tArgs{"", "key1"}, int32(0), true},
		{"3", tArgs{"", "key2"}, int32(-2147483648), true},
		{"4", tArgs{"", "key3"}, int32(2147483647), true},
		{"5", tArgs{"", "key4"}, int32(0), false},
		{"6", tArgs{"", "key5"}, int32(-1), true},
		{"7", tArgs{"", "key6"}, int32(0), false},
		{"8", tArgs{"", "n.a."}, int32(0), false},
		{"9", tArgs{"n.a.", "n.a."}, int32(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsInt32(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsInt32() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsInt32() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsInt32()

func TestTSectionList_AsInt64(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "-9223372036854775808")
	_ = sl.AddSectionKey("", "key3", "9223372036854775807")
	_ = sl.AddSectionKey("", "key4", "9223372036854775808") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")                  // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")                // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  int64
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, int64(0), false},
		{"2", tArgs{"", "key1"}, int64(0), true},
		{"3", tArgs{"", "key2"}, int64(-9223372036854775808), true},
		{"4", tArgs{"", "key3"}, int64(9223372036854775807), true},
		{"5", tArgs{"", "key4"}, int64(0), false},
		{"6", tArgs{"", "key5"}, int64(-1), true},
		{"7", tArgs{"", "key6"}, int64(0), false},
		{"8", tArgs{"", "n.a."}, int64(0), false},
		{"9", tArgs{"n.a.", "n.a."}, int64(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsInt64(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsInt64() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsInt64() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsInt64()

//

func TestTSectionList_AsString(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		// fields fields
		args  tArgs
		want  string
		want1 bool
	}{
		{"0", tArgs{"", ""}, "", false},
		{"1", tArgs{"s1", "bool"}, "nada", true},
		{"2", tArgs{"s2", "float"}, "12345.6789", true},
		{"3", tArgs{"s3", "int"}, "-12345", true},
		{"4", tArgs{"s4", "uint"}, "1234567890", true},
		{"5", tArgs{"", "n.a."}, "", false},
		{"6", tArgs{"n.a.", "n.a."}, "", false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsString(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsString() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsString() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsString()

//

func TestTSectionList_AsUInt(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "0")
	_ = sl.AddSectionKey("", "key3", "18446744073709551615") // amd64
	_ = sl.AddSectionKey("", "key4", "18446744073709551616") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")                   // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")                 // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  uint
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, uint(0), false},
		{"2", tArgs{"", "key1"}, uint(0), true},
		{"3", tArgs{"", "key2"}, uint(0), true},
		{"4", tArgs{"", "key3"}, uint(18446744073709551615), true},
		{"5", tArgs{"", "key4"}, uint(0), false},
		{"6", tArgs{"", "key5"}, uint(0), false},
		{"7", tArgs{"", "key6"}, uint(0), false},
		{"8", tArgs{"", "n.a."}, uint(0), false},
		{"9", tArgs{"n.a.", "-0"}, uint(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsUInt(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsUInt() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsUInt() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsUInt()

func TestTSectionList_AsUInt8(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "0")
	_ = sl.AddSectionKey("", "key3", "255")
	_ = sl.AddSectionKey("", "key4", "256")  // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")   // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.") // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  uint8
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, uint8(0), false},
		{"2", tArgs{"", "key1"}, uint8(0), true},
		{"3", tArgs{"", "key2"}, uint8(0), true},
		{"4", tArgs{"", "key3"}, uint8(255), true},
		{"5", tArgs{"", "key4"}, uint8(0), false},
		{"6", tArgs{"", "key5"}, uint8(0), false},
		{"7", tArgs{"", "key6"}, uint8(0), false},
		{"8", tArgs{"", "n.a."}, uint8(0), false},
		{"9", tArgs{"n.a.", "8"}, uint8(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsUInt8(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsUInt8() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsUInt8() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsUInt8()

func TestTSectionList_AsUInt16(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "0")
	_ = sl.AddSectionKey("", "key3", "65535")
	_ = sl.AddSectionKey("", "key4", "65536") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")    // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")  // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  uint16
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, uint16(0), false},
		{"2", tArgs{"", "key1"}, uint16(0), true},
		{"3", tArgs{"", "key2"}, uint16(0), true},
		{"4", tArgs{"", "key3"}, uint16(65535), true},
		{"5", tArgs{"", "key4"}, uint16(0), false},
		{"6", tArgs{"", "key5"}, uint16(0), false},
		{"7", tArgs{"", "key6"}, uint16(0), false},
		{"8", tArgs{"", "n.a."}, uint16(0), false},
		{"9", tArgs{"n.a.", "16"}, uint16(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsUInt16(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsUInt16() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsUInt16() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsUInt16()

func TestTSectionList_AsUInt32(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "0")
	_ = sl.AddSectionKey("", "key3", "4294967295")
	_ = sl.AddSectionKey("", "key4", "4294967296") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")         // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")       // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  uint32
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, uint32(0), false},
		{"2", tArgs{"", "key1"}, uint32(0), true},
		{"3", tArgs{"", "key2"}, uint32(0), true},
		{"4", tArgs{"", "key3"}, uint32(4294967295), true},
		{"5", tArgs{"", "key4"}, uint32(0), false},
		{"6", tArgs{"", "key5"}, uint32(0), false},
		{"7", tArgs{"", "key6"}, uint32(0), false},
		{"8", tArgs{"", "n.a."}, uint32(0), false},
		{"9", tArgs{"n.a.", "n.a."}, uint32(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsUInt32(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsUInt32() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsUInt32() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsUInt32()

func TestTSectionList_AsUInt64(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	_ = sl.AddSectionKey("", "key0", "")
	_ = sl.AddSectionKey("", "key1", "0")
	_ = sl.AddSectionKey("", "key2", "0")
	_ = sl.AddSectionKey("", "key3", "18446744073709551615")
	_ = sl.AddSectionKey("", "key4", "18446744073709551616") // Overflow
	_ = sl.AddSectionKey("", "key5", "-1")                   // Negative number
	_ = sl.AddSectionKey("", "key6", "n.a.")                 // Non-numeric string

	tests := []struct {
		name  string
		args  tArgs
		want  uint64
		want1 bool
	}{
		{"1", tArgs{"", "key0"}, uint64(0), false},
		{"2", tArgs{"", "key1"}, uint64(0), true},
		{"3", tArgs{"", "key2"}, uint64(0), true},
		{"4", tArgs{"", "key3"}, uint64(18446744073709551615), true},
		{"5", tArgs{"", "key4"}, uint64(0), false},
		{"6", tArgs{"", "key5"}, uint64(0), false},
		{"7", tArgs{"", "key6"}, uint64(0), false},
		{"8", tArgs{"", "n.a."}, uint64(0), false},
		{"9", tArgs{"n.a.", "n.a."}, uint64(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.AsUInt64(tt.args.aSection, tt.args.aKey)
			if got != tt.want {
				t.Errorf("%q: TSectionList.AsUInt64() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.AsUInt64() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_AsUInt64()

//

func TestTSectionList_Clear(t *testing.T) {
	sl := prepSectionList()
	cl := NewSectionList()

	tests := []struct {
		name string
		want *TSectionList
	}{
		{"1", cl},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.Clear(); !got.CompareTo(tt.want) {
				t.Errorf("%q: TSectionList.Clear() = {\n%v}, want '{\n%v}'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_Clear()

func TestTSectionList_CompareTo(t *testing.T) {
	sl := prepSectionList()
	sl1 := NewSectionList()
	sl2 := prepSectionList()
	sl3 := prepSectionList()
	sl3.AddSectionKey("", "key0", "null")
	sl4 := prepSectionList()
	sl4.RemoveSection("s4")
	sl4.AddSectionKey("s5", "whatever", "n.a.")

	tests := []struct {
		name string
		args *TSectionList
		want bool
	}{
		{"1", sl1, false},
		{"2", sl2, true},
		{"3", sl3, false},
		{"4", sl4, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.CompareTo(tt.args); got != tt.want {
				t.Errorf("%q: TSectionList.CompareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_CompareTo()

func TestTSectionList_GetSection(t *testing.T) {
	sl := prepSectionList()
	nl := &TSection{}
	kl2 := sl.sections[DefSection]
	kl3 := sl.sections["s3"]
	tests := []struct {
		name string
		args string
		want *TSection
	}{
		{"1", "n.a.", nl},
		{"2", "", kl2},
		{"3", "s3", kl3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.GetSection(tt.args); !got.CompareTo(tt.want) {
				t.Errorf("%q: TSectionList.GetSection() = {\n%v}, want {\n%v}",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_GetSection()

func TestTSectionList_HasSection(t *testing.T) {
	sl := prepSectionList()
	tests := []struct {
		name    string
		args    string
		wantROK bool
	}{
		{"0", "n.a.", false},
		{"1", "", true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotROK := sl.HasSection(tt.args); gotROK != tt.wantROK {
				t.Errorf("%q: TSectionList.HasSection() = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSectionList_HasSection()

func TestTSectionList_HasSectionKey(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", ""}, false},
		{"1", tArgs{"", "key0"}, true},
		{"2", tArgs{"s2", "float"}, true},
		{"3", tArgs{"s3", "n.a."}, false},
		{"3", tArgs{"sX", "n.a."}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.HasSectionKey(tt.args.aSection, tt.args.aKey); got != tt.want {
				t.Errorf("%q: TSectionList.HasSectionKey() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_HasSectionKey()

func TestTSectionList_RemoveSection(t *testing.T) {
	sl := prepSectionList()
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"1", "", true},     // first
		{"2", "s4", true},   // last
		{"3", "s2", true},   // middle
		{"4", "n.a.", true}, // n.a.
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.RemoveSection(tt.args); got != tt.want {
				t.Errorf("%q TSectionList.RemoveSection() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_RemoveSection()

func TestTSectionList_RemoveSectionKey(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		// fields fields
		args tArgs
		want bool
	}{
		{"0", tArgs{"", ""}, true},
		{"1", tArgs{"s1", "key0"}, true},
		{"2", tArgs{"s2", "float"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.RemoveSectionKey(tt.args.aSection, tt.args.aKey); got != tt.want {
				t.Errorf("%q: TSectionList.RemoveSectionKey() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_RemoveSectionKey()

func TestTSectionList_Sections(t *testing.T) {
	sl := prepSectionList()

	// sl.AddSectionKey("", "key0", "")
	// sl.AddSectionKey("s1", "bool", "nada")
	// sl.AddSectionKey("s2", "float", "12345.6789")
	// sl.AddSectionKey("s3", "int", "-12345")
	// sl.AddSectionKey("s4", "uint", "1234567890")

	tests := []struct {
		name  string
		want  []string
		want1 int
	}{
		{"1", []string{DefSection, "s2", "s1", "s4", "s3"}, 5},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sl.Sections()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: TSectionList.Sections() got = {\n%v},\nwant {\n%v}",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q: TSectionList.Sections() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_Sections()

func TestTSectionList_SetFilename(t *testing.T) {
	sl := NewSectionList()
	tests := []struct {
		name string
		args string
		want *TSectionList
	}{
		{"1", "/dev/null", sl},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.SetFilename(tt.args); got.Filename() != tt.args {
				t.Errorf("%q: TSectionList.SetFilename() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_SetFilename()

func TestTSectionList_Sort(t *testing.T) {
	sl := prepSectionList()
	wl := prepSectionList().Sort()
	tests := []struct {
		name string
		want *TSectionList
	}{
		{"1", wl},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.Sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: TSectionList.Sort() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_Sort()

func TestTSectionList_String(t *testing.T) {
	//NOTE: Since the order of the key/value pairs is not guaranteed
	// this test may occasionally fail.
	kld := NewSection()
	kld.AddKey("key1", "val1")
	kld.AddKey("key2", "val2")

	kl2 := NewSection()
	kl2.AddKey("key3", "val3")
	kl2.AddKey("key4", "")

	sl := &TSectionList{
		defSect: "Default",
		secOrder: tSectionOrder{
			"Default",
			"Sect2",
			"NOOP",
		},
		sections: tSections{
			"Sect2":   kl2,
			"Default": kld,
			"NOOP":    NewSection(),
		},
	}

	tests := []struct {
		name        string
		wantRString string
	}{
		{"1", "\n[Default]\nkey1 = val1\nkey2 = val2\n\n[Sect2]\nkey3 = val3\nkey4 =\n\n[NOOP]\n"},
		{"2", `
[Default]
key1 = val1
key2 = val2

[Sect2]
key3 = val3
key4 =

[NOOP]
`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRString := sl.String(); gotRString != tt.wantRString {
				t.Errorf("%q: TSectionList.String() = {\n%v},\nwant {\n%v}",
					tt.name, gotRString, tt.wantRString)
			}
		})
	}
} // TestTSectionList_String()

func Benchmark_TSectionList_String(b *testing.B) {
	sl, _ := NewIni(inFileName)
	for n := 0; n < b.N*8*4; n++ {
		if "" == sl.String() {
			continue
		}
	}
} // Benchmark_TSectionList_String()

// func Benchmark_TSectionList_String2(b *testing.B) {
// 	sl, _ := New(inFileName)
// 	for n := 0; n < b.N*8*4; n++ {
// 		if 0 == len(sl.String2()) {
// 			continue
// 		}
// 	}
// } // Benchmark_TSectionList_String()

func TestTSectionList_updateSectKey(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   string
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		// TODO: Add test cases.
		{"1", tArgs{"", "", ""}, false},
		{"2", tArgs{"general", "", ""}, false}, // empty key
		{"3", tArgs{"", "loglevel", ""}, true},
		{"4", tArgs{"general", "loglevel", ""}, true},
		{"5", tArgs{"general", "loglevel", "8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.updateSectKey(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSectionList.updateSectKey() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_updateSectKey()

func TestTSectionList_UpdateSectKeyBool(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   bool
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		// TODO: Add test cases.
		{"1", tArgs{"", "", false}, false},
		{"2", tArgs{"general", "", true}, false},
		{"3", tArgs{"", "loglevel", false}, true},
		{"4", tArgs{"general", "loglevel", true}, true},
		{"5", tArgs{"general", "loglevel", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.UpdateSectKeyBool(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSectionList.UpdateSectKeyBool() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_UpdateSectKeyBool()

func TestTSectionList_UpdateSectKeyFloat(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   float64
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", "", 0.0}, false},
		{"1", tArgs{"s1", "", 0.0}, false},
		{"2", tArgs{"s2", "float", 1.7976931348623157e+308}, true},
		{"3", tArgs{"s3", "int", -9223372036854775808}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.UpdateSectKeyFloat(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSectionList.UpdateSectKeyFloat() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_UpdateSectKeyFloat()

func TestTSectionList_UpdateSectKeyInt(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   int64
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", "", 0.0}, false},
		{"1", tArgs{"s1", "", 0.0}, false},
		{"2", tArgs{"s2", "float", -0}, true},
		{"3", tArgs{"s3", "int", 9223372036854775807}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.UpdateSectKeyInt(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSectionList.UpdateSectKeyInt() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_UpdateSectKeyUInt()

func TestTSectionList_UpdateSectKeyUInt(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   uint64
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", "", 0.0}, false},
		{"1", tArgs{"s1", "", 0.0}, false},
		{"2", tArgs{"s2", "float", 0}, true},
		{"3", tArgs{"s3", "int", 18446744073709551615}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.UpdateSectKeyUInt(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSectionList.UpdateSectKeyInt() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_UpdateSectKeyInt()

func TestTSectionList_UpdateSectKeyStr(t *testing.T) {
	type tArgs struct {
		aSection string
		aKey     string
		aValue   string
	}

	sl := prepSectionList()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", "", ""}, false},
		{"1", tArgs{"s1", "", "1"}, false},
		{"2", tArgs{"s2", "float", "float"}, true},
		{"3", tArgs{"s3", "int", "int"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.UpdateSectKeyStr(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSectionList.UpdateSectKeyStr() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_UpdateSectKeyStr()

func TestTSectionList_Merge(t *testing.T) {
	sl := prepSectionList()
	sl2 := prepSectionList()
	sl2.AddSectionKey("n.a.", "n.a.", "n.a.")
	tests := []struct {
		name string
		args *TSectionList
		want *TSectionList
	}{
		{"0", nil, sl},
		{"1", sl2, sl},
		{"2", sl2, sl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := sl.Merge(tt.args); !reflect.DeepEqual(got, tt.want) {
			if got := sl.Merge(tt.args); !got.CompareTo(tt.want) {
				t.Errorf("%q: TSectionList.Merge() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_Merge()

func TestTSectionList_WriteFile(t *testing.T) {
	ini, _ := NewIni(inFileName)
	ini.SetFilename(outFilename)

	tests := []struct {
		name       string
		id         *TSectionList
		wantRBytes int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{"1", ini, 4062, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRBytes, err := tt.id.Store()
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: TSectionList.WriteFile() error = {%v}, wantErr {%v}",
					tt.name, err, tt.wantErr)
				return
			}
			if gotRBytes != tt.wantRBytes {
				t.Errorf("%q: TSectionList.WriteFile() = '%v', want '%v'",
					tt.name, gotRBytes, tt.wantRBytes)
			}
		})
	}
} // TestTSectionList_WriteFile()

func walkFunc(aSect, aKey, aVal string) {
	fmt.Printf("\nSection: %s\nKey: %s\nValue: %s\n", aSect, aKey, aVal)
} // walkFunc()

func TestTSections_Walk(t *testing.T) {
	type tArgs struct {
		aFunc TIniWalkFunc
	}

	sl := prepSectionList()
	// sl.AddSectionKey("", "key0", "")
	// sl.AddSectionKey("s1", "bool", "nada")
	// sl.AddSectionKey("s2", "int", "12345.6789")
	// sl.AddSectionKey("s3", "int", "-12345")
	// sl.AddSectionKey("s4", "uint", "1234567890")
	tests := []struct {
		name   string
		fields TSectionList
		args   tArgs
	}{
		// TODO: Add test cases.
		{" 1", *sl, tArgs{walkFunc}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Walk(tt.args.aFunc)
		})
	}
} // TestTSections_Walk()

type tListWalk int

func (tw tListWalk) Walk(aSect, aKey, aVal string) {
	fmt.Fprintf(os.Stderr, "\nSection: %s\nKey: %s\nValue: %s\n",
		aSect, aKey, aVal)
} // walkFunc()

func TestTSectionList_Walker(t *testing.T) {
	var lw tListWalk

	sl := prepSectionList()
	tests := []struct {
		name string
		args TIniWalker
	}{
		{"1", lw},
		// TODO: add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl.Walker(tt.args)
		})
	}
} // TestTSectionList_Walker()

/* _EoF_ */
