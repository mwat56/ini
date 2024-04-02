/*
Copyright © 2019, 204 M.Watermann, 10247 Berlin, Germany

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if gotRVal := kv.AsBool(); gotRVal != tt.wantRVal {
				t.Errorf("TKeyVal.AsBool() = '%v', want '%v'", gotRVal, tt.wantRVal)
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
				t.Errorf("TKeyVal.AsFloat32() got = '%v', want '%v'", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TKeyVal.AsFloat32() got1 = '%v', want '%v'", got1, tt.want1)
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
				t.Errorf("TKeyVal.AsFloat64() got = '%v', want '%v'", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TKeyVal.AsFloat64() got1 = '%v', want '%v'", got1, tt.want1)
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
				t.Errorf("TKeyVal.AsInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TKeyVal.AsInt() got1 = %v, want %v", got1, tt.want1)
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
				t.Errorf("TKeyVal.AsInt16() got = '%v', want '%v'", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TKeyVal.AsInt16() got1 = '%v', want '%v'", got1, tt.want1)
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
				t.Errorf("TKeyVal.AsInt32() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TKeyVal.AsInt32() got1 = %v, want %v", got1, tt.want1)
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
				t.Errorf("TKeyVal.AsInt64() gotRVal = %v, want %v", gotRVal, tt.wantRVal)
			}
			if gotROK != tt.wantROK {
				t.Errorf("TKeyVal.AsInt64() gotROK = %v, want %v", gotROK, tt.wantROK)
			}
		})
	}
} // TestTKeyVal_AsInt64()

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
		{"4", TKeyVal{"key4", ""}, "", false},
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
				t.Errorf("TKeyVal.AsString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TKeyVal.AsString() got1 = %v, want %v", got1, tt.want1)
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
				t.Errorf("TKeyVal.String() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTKeyVal_String()

func Benchmark_TKeyVal_String(b *testing.B) {
	kv1 := TKeyVal{"key1", "val1"}
	for n := 0; n < b.N; n++ {
		if 0 == len(kv1.String()) {
			continue
		}
	}
} // Benchmark_TKeyVal_String()

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
		{"2", TKeyVal{"key2", "val2"}, args{""}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &TKeyVal{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := kv.UpdateValue(tt.args.aValue); got != tt.want {
				t.Errorf("TKeyVal.UpdateValue() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTKeyVal_UpdateValue()

/* _EoF_ */
