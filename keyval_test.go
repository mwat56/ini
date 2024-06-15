/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestTKeyVal_AsBool(t *testing.T) {
	tests := []struct {
		name     string
		fields   TKeyVal
		wantRVal bool
	}{
		{"1", TKeyVal{"key1", "val1"}, false},
		{"2", TKeyVal{"key2", `True`}, true},
		{"3", TKeyVal{"key3", `0`}, false},
		{"4", TKeyVal{"key4", ``}, false},
		{"5", TKeyVal{"key2", "ja"}, true},
		{"6", TKeyVal{"key2", "Oui"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if gotRVal := kv.AsBool(); gotRVal != tt.wantRVal {
				t.Errorf("%q TKeyVal.AsBool() = '%v', want '%v'",
					tt.name, gotRVal, tt.wantRVal)
			}
		})
	}
} // TestTKeyVal_AsBool()

func TestTKeyVal_AsFloat32(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   float32
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 5.5, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsFloat32()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsFloat32() got = '%v', want '%v'",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsFloat32() got1 = '%v', want '%v'",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsFloat32()

func TestTKeyVal_AsFloat64(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   float64
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 5.5, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsFloat64()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsFloat64() got = '%v', want '%v'",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsFloat64() got1 = '%v', want '%v'",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsFloat64()

func TestTKeyVal_AsInt(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   int
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsInt()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsInt() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsInt() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsInt()

func TestTKeyVal_AsInt8(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   int8
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsInt8()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsInt() got = %d, want %d",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsInt() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsInt()

func TestTKeyVal_AsInt16(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   int16
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsInt16()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsInt16() got = '%v', want '%v'",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("T%q KeyVal.AsInt16() got1 = '%v', want '%v'",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsInt16()

func TestTKeyVal_AsInt32(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   int32
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsInt32()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsInt32() got = %v, want %v",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsInt32() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsInt32()

func TestTKeyVal_AsInt64(t *testing.T) {
	tests := []struct {
		name     string
		fields   TKeyVal
		wantRVal int64
		wantROK  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, 0, false},
		{"2", TKeyVal{"key2", "0"}, 0, true},
		{"3", TKeyVal{"key3", "3"}, 3, true},
		{"4", TKeyVal{"key4", ""}, 0, false},
		{"5", TKeyVal{"key5", "5.5"}, 0, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			gotRVal, gotROK := kv.AsInt64()
			if gotRVal != tt.wantRVal {
				t.Errorf("%q TKeyVal.AsInt64() gotRVal = %v, want %v",
					tt.name, gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("%q TKeyVal.AsInt64() gotROK = %v, want %v",
					tt.name, gotROK, tt.wantROK)
			}
		})
	}
} // TestTKeyVal_AsInt64()

func TestTKeyVal_AsUInt(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   uint
		want1  bool
	}{
		{"1", TKeyVal{"key", ""}, uint(0), false},
		{"2", TKeyVal{"key", "0"}, uint(0), true},
		{"3", TKeyVal{"key", "9223372036854775807"}, uint(9223372036854775807), true},
		{"4", TKeyVal{"key", "92233720368547758088"}, uint(0), false}, // Overflow

		{"5", TKeyVal{"key", "-1"}, uint(0), false},  // Negative number
		{"6", TKeyVal{"key", "abc"}, uint(0), false}, // Non-numeric string
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			got, got1 := kv.AsUInt()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsUInt16() got = %d, want %d",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsUInt16() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsUInt()

func TestTKeyVal_AsUInt8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint8
		ok       bool
	}{
		{"Valid uint8 value", "123", 123, true},
		{"Invalid uint8 value", "256", 0, false},
		{"Negative value", "-10", 0, false},
		{"Leading whitespace", "  123", 123, true},
		{"Trailing whitespace", "123  ", 123, true},
		{"Leading and trailing whitespace", "  123  ", 123, true},
		{"Empty string", "", 0, false},
		{"Non-numeric string", "abc", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{Value: tt.input}
			result, ok := kv.AsUInt8()
			if result != tt.expected || ok != tt.ok {
				t.Errorf("%q AsUInt8() = (%v, %v), want (%v, %v)",
					tt.name, result, ok, tt.expected, tt.ok)
			}
		})
	}
} // TestTKeyVal_AsUInt8()

func TestTKeyVal_AsUInt16(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   uint16
		want1  bool
	}{
		{"1", TKeyVal{"key", ""}, uint16(0), false},
		{"2", TKeyVal{"key", "0"}, uint16(0), true},
		{"3", TKeyVal{"key", "65535"}, uint16(65535), true},
		{"4", TKeyVal{"key", "65536"}, uint16(0), false}, // Overflow
		{"5", TKeyVal{"key", "-1"}, uint16(0), false},    // Negative number
		{"6", TKeyVal{"key", "abc"}, uint16(0), false},   // Non-numeric string
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			got, got1 := kv.AsUInt16()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsUInt16() got = %d, want %d",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsUInt16() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsUInt16()

func TestTKeyVal_AsUInt32(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   uint32
		want1  bool
	}{
		{"1", TKeyVal{"key", ""}, uint32(0), false},
		{"2", TKeyVal{"key", "0"}, uint32(0), true},
		{"3", TKeyVal{"key", "4294967295"}, uint32(4294967295), true},
		{"4", TKeyVal{"key", "4294967296"}, uint32(0), false}, // Overflow
		{"5", TKeyVal{"key", "-1"}, uint32(0), false},         // Negative number
		{"6", TKeyVal{"key", "abc"}, uint32(0), false},        // Non-numeric string
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			got, got1 := kv.AsUInt32()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsUInt32() got = %d, want %d",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsUInt32() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsUInt32()

func TestTKeyVal_AsUInt64(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   uint64
		want1  bool
	}{
		{"1", TKeyVal{"key", ""}, uint64(0), false},
		{"2", TKeyVal{"key", "0"}, uint64(0), true},
		{"3", TKeyVal{"key", "18446744073709551615"}, uint64(18446744073709551615), true},
		{"4", TKeyVal{"key", "18446744073709551616"}, uint64(0), false}, // Overflow
		{"5", TKeyVal{"key", "-1"}, uint64(0), false},                   // Negative number
		{"6", TKeyVal{"key", "abc"}, uint64(0), false},                  // Non-numeric string
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			got, got1 := kv.AsUInt64()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsUInt64() got = %d, want %d",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsUInt64() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsUInt64()

func TestTKeyVal_AsString(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   string
		want1  bool
	}{
		{"1", TKeyVal{"key1", "val1"}, "val1", true},
		{"2", TKeyVal{"key2", "0"}, "0", true},
		{"3", TKeyVal{"key3", "3"}, "3", true},
		{"4", TKeyVal{"key4", ""}, "", true},
		{"5", TKeyVal{"key5", "5.5"}, "5.5", true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			got, got1 := kv.AsString()
			if got != tt.want {
				t.Errorf("%q TKeyVal.AsString() got = %q, want %q",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TKeyVal.AsString() got1 = %v, want %v",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTKeyVal_AsString()

func TestTKeyVal_String(t *testing.T) {
	tests := []struct {
		name   string
		fields TKeyVal
		want   string
	}{
		{"1", TKeyVal{"key1", "val1"}, "key1 = val1"},
		{"2", TKeyVal{"key2", "0"}, "key2 = 0"},
		{"3", TKeyVal{"key3", "3"}, "key3 = 3"},
		{"4", TKeyVal{"key4", ""}, "key4 ="},
		{"5", TKeyVal{"key5", "5.5"}, "key5 = 5.5"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := kv.String(); got != tt.want {
				t.Errorf("%q TKeyVal.String() = %q, want %q",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTKeyVal_String()

func Benchmark_TKeyVal_String(b *testing.B) {
	kv1 := TKeyVal{"key1", "val1"}
	for n := 0; n < b.N*8*8; n++ {
		if "" == kv1.String() {
			continue
		}
	}
} // Benchmark_TKeyVal_String()

// func Benchmark_TKeyVal_String2(b *testing.B) {
// 	kv1 := TKeyVal{"key1", "val1"}
// 	for n := 0; n < b.N*8*8; n++ {
// 		if 0 == len(kv1.String2()) {
// 			continue
// 		}
// 	}
// } // Benchmark_TKeyVal_String()

func TestTKeyVal_UpdateValue(t *testing.T) {
	type args struct {
		aValue string
	}
	tests := []struct {
		name   string
		fields TKeyVal
		args   args
		want   bool
	}{
		{"1", TKeyVal{"key1", "val1"}, args{"val-11"}, true},
		{"2", TKeyVal{"key2", "val2"}, args{""}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := kv.UpdateValue(tt.args.aValue); got != tt.want {
				t.Errorf("%q TKeyVal.UpdateValue() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTKeyVal_UpdateValue()

/* _EoF_ */
