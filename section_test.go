/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"reflect"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestTSection_AddKey(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "val2"},
		TKeyVal{"key3", "val3"},
		TKeyVal{"key4", ""},
	}
	tests := []struct {
		name string
		cs   *TSection
		args TKeyVal
		want bool
	}{
		{"5", ks, TKeyVal{"key5", "5.5"}, true},
		{"6", ks, TKeyVal{"key5", "6.6"}, true},
		{"7", ks, TKeyVal{"", "7.7"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.AddKey(tt.args.Key, tt.args.Value); got != tt.want {
				t.Errorf("%q TSection.AddKey() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_AddKey()

func TestTSection_AsBool(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "temp"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", ""},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal bool
		wantROK  bool
	}{
		{"1", ks, "key1", false, true},
		{"2", ks, "key2", true, true},
		{"3", ks, "key3", false, true},
		{"4", ks, "key4", false, true},
		{"5", ks, "key5", false, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsBool(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsBool() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsBool() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsBool()

func TestTSection_AsFloat32(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal float32
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 12.34, true},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsFloat32(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsFloat32() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsFloat32() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsFloat32()

func TestTSection_AsFloat64(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal float64
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 12.34, true},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsFloat64(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsFloat64() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsFloat64() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsFloat64()

func TestTSection_AsInt(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal int
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsInt(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsInt() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsInt() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsInt()

func TestTSection_AsInt8(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal int8
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsInt8(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsInt8() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsInt8() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsInt()

func TestTSection_AsInt16(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "76543"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal int16
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsInt16(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsInt16() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsInt16() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsInt16()

func TestTSection_AsInt32(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "76543"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal int32
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 76543, true},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsInt32(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsInt32() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsInt32() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsInt32()

func TestTSection_AsInt64(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "76543"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal int64
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 76543, true},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsInt64(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsInt64() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsInt64() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsInt64()

func TestTSection_AsUInt(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal uint
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsUInt(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsUInt() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsUInt() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsUInt()

func TestTSection_AsUInt8(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal int8
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsInt8(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsInt8() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsInt8() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsUInt8()

func TestTSection_AsUInt16(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "76543"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal uint16
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 0, false},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsUInt16(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsUInt16() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsUInt16() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsUInt16()

func TestTSection_AsUInt32(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "76543"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal uint32
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 76543, true},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsUInt32(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsUInt32() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsUInt32() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsUInt32()

func TestTSection_AsUInt64(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "76543"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "12.34"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     string
		wantRVal uint64
		wantROK  bool
	}{
		{"1", ks, "key1", 0, false},
		{"2", ks, "key2", 76543, true},
		{"3", ks, "key3", 0, true},
		{"4", ks, "key4", 0, false},
		{"5", ks, "key5", 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsUInt64(tt.args)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsUInt64() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsUInt64() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsUInt64()

func TestTSection_AsString(t *testing.T) {
	type kArgs struct {
		aKey string
	}
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "hello"},
	}
	tests := []struct {
		name     string
		cs       *TSection
		args     kArgs
		wantRVal string
		wantROK  bool
	}{
		{"1", ks, kArgs{"key1"}, "1st", true},
		{"2", ks, kArgs{"key2"}, "2nd", true},
		{"3", ks, kArgs{"key3"}, "0", true},
		{"4", ks, kArgs{"key4"}, "hello", true},
		{"5", ks, kArgs{"key5"}, "", false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRVal, gotROK := tt.cs.AsString(tt.args.aKey)
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TSection.AsString() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TSection.AsString() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_AsString()

func TestTSection_Clear(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "hello"},
	}
	ws1 := &TSection{}
	ws2 := &TSection{}
	tests := []struct {
		name string
		cs   *TSection
		want *TSection
	}{
		{"1", ks, ws1},
		{"2", ws1, ws2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.Clear(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q TSection.Clear() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_Clear()

func TestTSection_HasKey(t *testing.T) {
	type kArgs struct {
		aKey string
	}
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "hello"},
	}
	tests := []struct {
		name string
		cs   *TSection
		args kArgs
		want bool
	}{
		{"1", ks, kArgs{"key1"}, true},
		{"5", ks, kArgs{"key5"}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.HasKey(tt.args.aKey); got != tt.want {
				t.Errorf("%q TSection.HasKey() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_HasKey()

func TestTSection_IndexOf(t *testing.T) {
	type kArgs struct {
		aKey string
	}
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "hello"},
	}
	tests := []struct {
		name string
		cs   *TSection
		args kArgs
		want int
	}{
		{"1", ks, kArgs{"key1"}, 0},
		{"2", ks, kArgs{"key2"}, 1},
		{"5", ks, kArgs{"key5"}, -1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.IndexOf(tt.args.aKey); got != tt.want {
				t.Errorf("%q TSection.IndexOf() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_IndexOf()

func TestTSection_Len(t *testing.T) {
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "hello"},
	}
	tests := []struct {
		name string
		cs   *TSection
		want int
	}{
		{"1", ks, 4},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.Len(); got != tt.want {
				t.Errorf("%q TSection.Len() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_Len()

func TestTSection_RemoveKey(t *testing.T) {
	type kArgs struct {
		aKey string
	}
	ks := &TSection{
		TKeyVal{"key1", "1st"},
		TKeyVal{"key2", "2nd"},
		TKeyVal{"key3", "0"},
		TKeyVal{"key4", "hello"},
	}
	tests := []struct {
		name string
		cs   *TSection
		args kArgs
		want bool
	}{
		{"1", ks, kArgs{"key1"}, true},
		{"2", ks, kArgs{"key1"}, true},
		{"3", ks, kArgs{"n.a."}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.RemoveKey(tt.args.aKey); got != tt.want {
				t.Errorf("%q TSection.RemoveKey() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_RemoveKey()

func TestTSection_String(t *testing.T) {
	sl1 := TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "val2"},
		TKeyVal{"key3", "val3"},
		TKeyVal{"key4", ""},
	}
	rl1 := "key1 = val1\nkey2 = val2\nkey3 = val3\nkey4 =\n"
	tests := []struct {
		name string
		cs   *TSection
		want string
	}{
		// TODO: Add test cases.
		{" 1", &sl1, rl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.String(); got != tt.want {
				t.Errorf("%q TSection.String() = {%v}, want {%v}",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_String()

func Benchmark_TSection_String(b *testing.B) {
	sl1 := TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "val2"},
		TKeyVal{"key3", "val3"},
		TKeyVal{"key4", ""},
	}
	for n := 0; n < b.N*8*8; n++ {
		if "" == sl1.String() {
			continue
		}
	}
} // Benchmark_TSection_String()

// func Benchmark_TSections_String2(b *testing.B) {
// 	id1 := TIniList{
// 		defSect: "Default",
// 		secOrder: tSectionOrder{
// 			"Default",
// 			"Sect2",
// 			"NOOP",
// 		},
// 		sections: tSectionList{
// 			"Sect2": &TSection{
// 				TKeyVal{"key3", "val3"},
// 				TKeyVal{"key4", ""},
// 			},
// 			"Default": &TSection{
// 				TKeyVal{"key1", "val1"},
// 				TKeyVal{"key2", "val2"},
// 			},
// 		},
// 	}
// 	for n := 0; n < b.N; n++ {
// 		if 0 == len(id1.String()) {
// 			continue
// 		}
// 	}
// } /

func TestTSection_UpdateKey(t *testing.T) {
	ks := make(TSection, 0, slDefCapacity)
	ks.AddKey("Key1", "Value1")
	ks.AddKey("Key2", "Value2")
	type args struct {
		aKey   string
		aValue string
	}
	tests := []struct {
		name string
		cs   *TSection
		args args
		want bool
	}{
		// TODO: Add test cases.
		{" 0", &ks, args{"", ""}, false},
		{" 1", &ks, args{"Key1", "Value 1 (new)"}, true},
		{" 2", &ks, args{"Key 2", "Value 2 (new)"}, true},
		{" 3", &ks, args{"Key2", "Value 2 (new)"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.UpdateKey(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q TSection.UpdateKey() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_UpdateKey()

/* _EoF_ */
