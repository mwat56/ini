/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"reflect"
	"runtime"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func prepSection() *TSection {
	kl := NewSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("bool", "")
	_ = kl.AddKey("float", "")
	_ = kl.AddKey("int", "")
	_ = kl.AddKey("uint", "")

	return kl
} // prepSection()

func Benchmark_merge1(b *testing.B) {
	runtime.GOMAXPROCS(1)

	kl := NewSection()
	kl2 := prepSection()
	kl3 := prepSection()
	_ = kl3.AddKey("key0", "0")
	_ = kl3.AddKey("key1", "1")
	_ = kl3.AddKey("key2", "2")
	kl4 := prepSection()
	_ = kl4.AddKey("key3", "funny")
	_ = kl4.AddKey("key4", "nightmare")
	_ = kl4.AddKey("key5", "talisman")
	for n := 0; n < b.N<<5; n++ {
		kl5 := kl.merge(kl2)
		kl6 := kl5.merge(kl3)
		_ = kl6.merge(kl3)
	}
} // Benchmark_merge1()

func Benchmark_merge2(b *testing.B) {
	runtime.GOMAXPROCS(1)

	kl := NewSection()
	kl2 := prepSection()
	kl3 := prepSection()
	_ = kl3.AddKey("key0", "0")
	_ = kl3.AddKey("key1", "1")
	_ = kl3.AddKey("key2", "2")
	kl4 := prepSection()
	_ = kl4.AddKey("key3", "funny")
	_ = kl4.AddKey("key4", "nightmare")
	_ = kl4.AddKey("key5", "talisman")

	for n := 0; n < b.N<<5; n++ {
		kl5 := kl.merge2(kl2)
		kl6 := kl5.merge2(kl3)
		_ = kl6.merge2(kl3)
	}
} // Benchmark_merge2()

func Benchmark_Merge3(b *testing.B) {
	runtime.GOMAXPROCS(1)

	kl := NewSection()
	kl2 := prepSection()
	kl3 := prepSection()
	_ = kl3.AddKey("key0", "0")
	_ = kl3.AddKey("key1", "1")
	_ = kl3.AddKey("key2", "2")
	kl4 := prepSection()
	_ = kl4.AddKey("key3", "funny")
	_ = kl4.AddKey("key4", "nightmare")
	_ = kl4.AddKey("key5", "talisman")

	for n := 0; n < b.N<<5; n++ {
		kl5 := kl.Merge(kl2)
		kl6 := kl5.Merge(kl3)
		_ = kl6.Merge(kl3)
	}
} // Benchmark_Merge()

func TestNewSection(t *testing.T) {
	kl := &TSection{
		data: make(map[string]string),
	}
	tests := []struct {
		name string
		want *TSection
	}{
		{"1", kl},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: NewSection() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestNewSection()

func TestTSection_AddKey(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue string
	}
	kl := prepSection()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"1", tArgs{"", ""}, false},
		{"2", tArgs{"key2", ""}, true},
		{"3", tArgs{"key3", "val3"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.AddKey(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSection.AddKey() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_AddKey()

func TestTSection_AsBool(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "1")
	_ = kl.AddKey("key2", "2")
	_ = kl.AddKey("key3", "funny")
	_ = kl.AddKey("key4", "nightmare")
	_ = kl.AddKey("key5", "talisman")
	tests := []struct {
		args  string
		want  bool
		want1 bool
	}{
		{"bool", false, false},
		{"key0", false, true},
		{"key1", true, true},
		{"key2", false, false},
		{"key3", false, true},
		{"key4", false, true},
		{"key5", true, true},
		{"600", false, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsBool(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsBool(%q) got = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsBool(%q) got1 = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsBool()

func TestTSection_AsFloat32(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "123.45")
	_ = kl.AddKey("key2", "0.0")
	_ = kl.AddKey("key3", "-123.45")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "five")
	tests := []struct {
		args  string
		want  float32
		want1 bool
	}{
		{"", float32(0), false},
		{"key1", float32(123.45), true},
		{"key2", float32(0.0), true},
		{"key3", float32(-123.45), true},
		{"key4", float32(0.0), false},
		{"key5", float32(0.0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsFloat32(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsFloat32(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsFloat32(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsFloat32()

func TestTSection_AsFloat64(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "123.45")
	_ = kl.AddKey("key2", "0.0")
	_ = kl.AddKey("key3", "-123.456")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "five")
	tests := []struct {
		args  string
		want  float64
		want1 bool
	}{
		{"", float64(0), false},
		{"key1", float64(123.45), true},
		{"key2", float64(0.0), true},
		{"key3", float64(-123.456), true},
		{"key4", float64(0.0), false},
		{"key5", float64(0.0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsFloat64(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsFloat64(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsFloat64(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsFloat64()

func TestTSection_AsInt(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "123.45")
	_ = kl.AddKey("key2", "-456")
	_ = kl.AddKey("key3", "-123.456")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "123456789")
	tests := []struct {
		args  string
		want  int
		want1 bool
	}{
		{"", 0, false},
		{"key1", 0, false},
		{"key2", -456, true},
		{"key3", 0, false},
		{"key4", 0, false},
		{"key5", 123456789, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsInt(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsInt(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsInt(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsInt()

func TestTSection_AsInt8(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "")
	_ = kl.AddKey("key2", "128")
	_ = kl.AddKey("key3", "-127")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "127")
	tests := []struct {
		args  string
		want  int8
		want1 bool
	}{
		{"", 0, false},
		{"key1", int8(0), false},
		{"key2", int8(0), false},
		{"key3", int8(-127), true},
		{"key4", int8(0), false},
		{"key5", int8(127), true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsInt8(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsInt8(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsInt8(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsInt8()

func TestTSection_AsInt16(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "")
	_ = kl.AddKey("key2", "32768")
	_ = kl.AddKey("key3", "-32768")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "32767")

	tests := []struct {
		args  string
		want  int16
		want1 bool
	}{
		{"", 0, false},
		{"key1", int16(0), false},
		{"key2", int16(0), false},
		{"key3", int16(-32768), true},
		{"key4", int16(0), false},
		{"key5", int16(32767), true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsInt16(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsInt16(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsInt16(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsInt16()

func TestTSection_AsInt32(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "")
	_ = kl.AddKey("key2", "2147483648")
	_ = kl.AddKey("key3", "-2147483648")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "2147483647")

	tests := []struct {
		args  string
		want  int32
		want1 bool
	}{
		{"", 0, false},
		{"key1", int32(0), false},
		{"key2", int32(0), false},
		{"key3", int32(-2147483648), true},
		{"key4", int32(0), false},
		{"key5", int32(2147483647), true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsInt32(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsInt32(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsInt32(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsInt32()

func TestTSection_AsInt64(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "0")
	_ = kl.AddKey("key1", "")
	_ = kl.AddKey("key2", "9223372036854775808")
	_ = kl.AddKey("key3", "-9223372036854775808")
	_ = kl.AddKey("key4", "nan")
	_ = kl.AddKey("key5", "9223372036854775807")

	tests := []struct {
		args  string
		want  int64
		want1 bool
	}{
		{"", 0, false},
		{"key1", int64(0), false},
		{"key2", int64(0), false},
		{"key3", int64(-9223372036854775808), true},
		{"key4", int64(0), false},
		{"key5", int64(9223372036854775807), true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsInt64(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsInt64(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsInt64(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsInt32()

func TestTSection_AsString(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "1")
	_ = kl.AddKey("key2", "funny")
	_ = kl.AddKey("key3", "nightmare")
	_ = kl.AddKey("key4", "tailor")
	tests := []struct {
		args  string
		want  string
		want1 bool
	}{
		{"", "", false},
		{"key0", "", true},
		{"key1", "1", true},
		{"key2", "funny", true},
		{"key3", "nightmare", true},
		{"key4", "tailor", true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {

			got, got1 := kl.AsString(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsString(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsString(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
}

func TestTSection_AsUInt(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "0")
	_ = kl.AddKey("key2", "123456789")
	_ = kl.AddKey("key3", "9223372036854775807")
	_ = kl.AddKey("key4", "18446744073709551616") // Overflow
	_ = kl.AddKey("key5", "-1")                   // Negative number
	_ = kl.AddKey("key6", "abc")                  // Non-numeric string

	tests := []struct {
		args  string
		want  uint
		want1 bool
	}{
		{"", uint(0), false},
		{"key1", uint(0), true},
		{"key2", uint(123456789), true},
		{"key3", uint(9223372036854775807), true},
		{"key4", uint(0), false},
		{"key5", uint(0), false},
		{"key6", uint(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsUInt(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsUInt(%q) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsUInt(%q) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsUInt()

func TestTSection_AsUInt8(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "0")
	_ = kl.AddKey("key2", "123")
	_ = kl.AddKey("key3", "255")
	_ = kl.AddKey("key4", "256") // Overflow
	_ = kl.AddKey("key5", "-1")  // Negative number
	_ = kl.AddKey("key6", "abc") // Non-numeric string

	tests := []struct {
		args  string
		want  uint8
		want1 bool
	}{
		{"", uint8(0), false},
		{"key1", uint8(0), true},
		{"key2", uint8(123), true},
		{"key3", uint8(255), true},
		{"key4", uint8(0), false},
		{"key5", uint8(0), false},
		{"key6", uint8(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsUInt8(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsUInt8(%q)) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsUInt8(%q)) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsUInt8

func TestTSection_AsUInt16(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "0")
	_ = kl.AddKey("key2", "256")
	_ = kl.AddKey("key3", "65535")
	_ = kl.AddKey("key4", "65536") // Overflow
	_ = kl.AddKey("key5", "-1")    // Negative number
	_ = kl.AddKey("key6", "abc")   // Non-numeric string

	tests := []struct {
		args  string
		want  uint16
		want1 bool
	}{
		{"", uint16(0), false},
		{"key1", uint16(0), true},
		{"key2", uint16(256), true},
		{"key3", uint16(65535), true},
		{"key4", uint16(0), false},
		{"key5", uint16(0), false},
		{"key6", uint16(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsUInt16(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsUInt16(%q)) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsUInt16(%q)) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsUInt16

func TestTSection_AsUInt32(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "0")
	_ = kl.AddKey("key2", "65536")
	_ = kl.AddKey("key3", "4294967295")
	_ = kl.AddKey("key4", "4294967296") // Overflow
	_ = kl.AddKey("key5", "-1")         // Negative number
	_ = kl.AddKey("key6", "abc")        // Non-numeric string

	tests := []struct {
		args  string
		want  uint32
		want1 bool
	}{
		{"", uint32(0), false},
		{"key1", uint32(0), true},
		{"key2", uint32(65536), true},
		{"key3", uint32(4294967295), true},
		{"key4", uint32(0), false},
		{"key5", uint32(0), false},
		{"key6", uint32(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsUInt32(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsUInt32(%q)) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsUInt32(%q)) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsUInt32

func TestTSection_AsUInt64(t *testing.T) {
	kl := prepSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "0")
	_ = kl.AddKey("key2", "4294967296")
	_ = kl.AddKey("key3", "18446744073709551615")
	_ = kl.AddKey("key4", "18446744073709551616") // Overflow
	_ = kl.AddKey("key5", "-1")                   // Negative number
	_ = kl.AddKey("key6", "abc")                  // Non-numeric string

	tests := []struct {
		args  string
		want  uint64
		want1 bool
	}{
		{"", uint64(0), false},
		{"key1", uint64(0), true},
		{"key2", uint64(4294967296), true},
		{"key3", uint64(18446744073709551615), true},
		{"key4", uint64(0), false},
		{"key5", uint64(0), false},
		{"key6", uint64(0), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, got1 := kl.AsUInt64(tt.args)
			if got != tt.want {
				t.Errorf("TSection.AsUInt64(%q)) val = %v, want %v",
					tt.args, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSection.AsUInt64(%q)) ok = %v, want %v",
					tt.args, got1, tt.want1)
			}
		})
	}
} // TestTSection_AsUInt64

func TestTSection_Clear(t *testing.T) {
	kl := prepSection()
	klw := NewSection()

	tests := []struct {
		name   string
		fields *TSection
		want   *TSection
	}{
		{"1", kl, klw},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.Clear(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: TSection.Clear() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_Clear()

func TestTSection_CompareTo(t *testing.T) {
	kl := prepSection()
	kl2 := NewSection()
	kl3 := prepSection()
	_ = kl3.AddKey("key0", "val0")
	kl4 := &TSection{}

	tests := []struct {
		name string
		args *TSection
		want bool
	}{
		{"1", kl, true},
		{"2", kl2, false},
		{"3", kl3, false},
		{"4", kl4, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.CompareTo(tt.args); got != tt.want {
				t.Errorf("%q: TSection.compareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_CompareTo()

func TestTSection_HasKey(t *testing.T) {
	kl := prepSection()

	tests := []struct {
		args    string
		wantROK bool
	}{
		{"", false},
		{"key0", true},
		{"key1", false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if gotROK := kl.HasKey(tt.args); gotROK != tt.wantROK {
				t.Errorf("TSection.HasKey(%q) = %v, want %v",
					tt.args, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_HasKey()

func TestTSection_Len(t *testing.T) {
	kl1 := prepSection()
	kl2 := NewSection()
	kl3 := NewSection()
	_ = kl3.AddKey("key0", "")
	_ = kl3.AddKey("key1", "0")
	tests := []struct {
		name   string
		fields *TSection
		want   int
	}{
		{"1", kl1, 5},
		{"2", kl2, 0},
		{"3", kl3, 2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kl := tt.fields
			if got := kl.Len(); got != tt.want {
				t.Errorf("%q: TSection.Len() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_Len()

func TestTSection_Merge(t *testing.T) {
	kl1 := prepSection()
	kl2 := prepSection()
	kl3 := NewSection()
	_ = kl3.AddKey("key0", "")
	_ = kl3.AddKey("bool", "")
	kl4 := prepSection()
	_ = kl4.AddKey("401", "")
	_ = kl4.AddKey("402", "")

	tests := []struct {
		name   string
		fields *TSection
		args   *TSection
		want   *TSection
	}{
		{"0", kl1, kl1, kl1},
		{"1", kl1, kl2, kl1},
		{"2", kl2, kl3, kl2},
		{"3", kl3, kl4, kl4},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kl := &TSection{
				data: tt.fields.data,
				// mtx:  tt.fields.mtx,
			}
			if got := kl.Merge(tt.args); !kl.CompareTo(got) {
				t.Errorf("%q: TSection.Merge() = {\n%v},\nwant {\n%v}",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_Merge()

func Benchmark_TSection_Merge(b *testing.B) {
	runtime.GOMAXPROCS(1)

	kl1 := prepSection()
	kl2 := NewSection()
	_ = kl2.AddKey("key2", "")
	_ = kl2.AddKey("bool", "")

	for n := 0; n < b.N<<5; n++ {
		if nil == kl1.Merge(kl2) {
			continue
		}
	}
} // Benchmark_TSection_Merge()

func TestTSection_merge(t *testing.T) {
	kl1 := NewSection()
	kl2 := prepSection()
	kl3 := NewSection()
	_ = kl3.AddKey("301", "")
	_ = kl3.AddKey("302", "")

	tests := []struct {
		name   string
		fields *TSection
		args   *TSection
		want   *TSection
	}{
		{"1", kl1, kl2, kl2},
		{"2", kl2, kl3, kl3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kl := &TSection{
				data: tt.fields.data,
			}
			if got := kl.merge(tt.args); !kl.CompareTo(got) {
				t.Errorf("%q: TSection.merge() = {\n%v}, want {\n%v}",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_merge()

func Benchmark_TSection_merge(b *testing.B) {
	runtime.GOMAXPROCS(1)

	kl1 := prepSection()
	kl2 := NewSection()
	_ = kl2.AddKey("key2", "")
	_ = kl2.AddKey("bool", "")

	for n := 0; n < b.N<<5; n++ {
		if nil == kl1.merge(kl2) {
			continue
		}
	}
} // Benchmark_TSection_merge()

func TestTSection_RemoveKey(t *testing.T) {
	kl := prepSection()

	tests := []struct {
		args string
		want bool
	}{
		{"", true},
		{"key0", true},  // first
		{"uint", true},  // last
		{"float", true}, // middle
		{"n.a.", true},  // not available
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := kl.RemoveKey(tt.args); got != tt.want {
				t.Errorf("TSection.RemoveKey(%q) = %v, want %v",
					tt.args, got, tt.want)
			}
		})
	}
} // TestTSection_RemoveKey()

func TestTSection_String(t *testing.T) {
	//NOTE: Since the order of the key/value pairs is not guaranteed
	// this test may occasionally fail.
	kl := NewSection()
	_ = kl.AddKey("key0", "")
	_ = kl.AddKey("key1", "1")
	_ = kl.AddKey("key2", "two")

	tests := []struct {
		name        string
		wantRString string
	}{
		// NOTE: The order of the key/value pairs is not guaranteed.
		{"1", "key0 =\nkey1 = 1\nkey2 = two\n"},
		{"2", `key0 =
key1 = 1
key2 = two
`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRString := kl.String(); gotRString != tt.wantRString {
				t.Errorf("%q: TSection.String() = %q,\nwant %q",
					tt.name, gotRString, tt.wantRString)
			}
		})
	}
} // TestTSection_String()

func TestTSection_UpdateKey(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue string
	}

	kl := prepSection()
	tests := []struct {
		name    string
		args    tArgs
		wantROK bool
	}{
		{"0", tArgs{"", ""}, false},
		{"1", tArgs{"", "value2"}, false},
		{"2", tArgs{"key0", "value3"}, true},
		{"3", tArgs{"n.a.", "value4"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotROK := kl.UpdateKey(tt.args.aKey, tt.args.aValue); gotROK != tt.wantROK {
				t.Errorf("%q: TSection.UpdateKey() = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTSection_UpdateKey()

func TestTSection_UpdateKeyBool(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue bool
	}

	kl := prepSection()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", false}, false},
		{"1", tArgs{"", true}, false},
		{"2", tArgs{"key0", false}, true},
		{"3", tArgs{"n.a.", true}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.UpdateKeyBool(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSection.UpdateKeyBool() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_UpdateKeyBool()

func TestTSection_UpdateSectKeyFloat(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue float64
	}

	kl := prepSection()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", 0.0}, false},
		{"1", tArgs{"", 1.1}, false},
		{"2", tArgs{"key0", 2.2}, true},
		{"3", tArgs{"n.a.", 3.3}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.UpdateSectKeyFloat(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSection.UpdateSectKeyFloat() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_UpdateSectKeyFloat()

func TestTSection_UpdateKeyInt(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue int64
	}

	kl := prepSection()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", -0}, false},
		{"1", tArgs{"", -1}, false},
		{"2", tArgs{"key0", -2}, true},
		{"3", tArgs{"n.a.", -3}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.UpdateKeyInt(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSection.UpdateKeyInt() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_UpdateKeyInt()

func TestTSection_UpdateKeyUInt(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue uint64
	}

	kl := prepSection()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", 0}, false},
		{"1", tArgs{"", 1}, false},
		{"2", tArgs{"key0", 2}, true},
		{"3", tArgs{"n.a.", 3}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.UpdateKeyUInt(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSection.UpdateKeyUInt() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_UpdateKeyUInt()

func TestTSection_UpdateKeyStr(t *testing.T) {
	type tArgs struct {
		aKey   string
		aValue string
	}

	kl := prepSection()
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{"", "0"}, false},
		{"1", tArgs{"", "1"}, false},
		{"2", tArgs{"key0", "2"}, true},
		{"3", tArgs{"n.a.", "3"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kl.UpdateKeyStr(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q: TSection.UpdateKeyStr() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSection_UpdateKeyStr()
