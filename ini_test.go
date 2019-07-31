/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package ini

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"fmt"
	"testing"
)

const (
	inFileName  = "testin.ini"
	outFileName = "testout.ini"
)

func Test_tKeyVal_String(t *testing.T) {
	kv1 := TKeyVal{"key1", "val1"}
	rs1 := "key1 = val1"
	kv2 := TKeyVal{"key2", ""}
	rs2 := "key2 ="
	kv3 := TKeyVal{"", ""}
	rs3 := " ="
	tests := []struct {
		name string
		kv   *TKeyVal
		want string
	}{
		// TODO: Add test cases.
		{" 1", &kv1, rs1},
		{" 2", &kv2, rs2},
		{" 3", &kv3, rs3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := TKeyVal{
				Key:   tt.kv.Key,
				Value: tt.kv.Value,
			}
			if got := kv.String(); got != tt.want {
				t.Errorf("tKeyVal.String() = {%v}, want {%v}", got, tt.want)
			}
		})
	}
} // Test_tKeyVal_String()

func Benchmark_tKeyVal_String(b *testing.B) {
	kv1 := TKeyVal{"key1", "val1"}
	for n := 0; n < b.N; n++ {
		if 0 > len(kv1.String()) {
			continue
		}
	}
} // Benchmark_tKeyVal_String()

func Benchmark_tKeyVal_string0(b *testing.B) {
	kv1 := TKeyVal{"key1", "val1"}
	for n := 0; n < b.N; n++ {
		if 0 > len(kv1.string0()) {
			continue
		}
	}
} // Benchmark_tKeyVal_string0()

func Test_removeQuotes(t *testing.T) {
	si1 := "'this is a text'"
	so1 := "this is a text"
	si2 := " \" this is a text \" "
	ws2 := " this is a text "
	si3 := " \" this is a text ' "
	ws3 := "\" this is a text '"
	type args struct {
		aString string
	}
	tests := []struct {
		name        string
		args        args
		wantRString string
	}{
		// TODO: Add test cases.
		{" 1", args{si1}, so1},
		{" 2", args{si2}, ws2},
		{" 3", args{si3}, ws3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRString := removeQuotes(tt.args.aString); gotRString != tt.wantRString {
				t.Errorf("removeQuotes() = %v, want %v", gotRString, tt.wantRString)
			}
		})
	}
} // Test_removeQuotes

func Test_tSection_String(t *testing.T) {
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
				t.Errorf("tSection.String() = {%v}, want {%v}", got, tt.want)
			}
		})
	}
} // Test_tSection_String()

func Test_tSection_UpdateKey(t *testing.T) {
	cs1 := make(TSection, 0, defCapacity)
	cs1.AddKey("Key1", "Value1")
	cs1.AddKey("Key2", "Value2")
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
		{" 0", &cs1, args{"", ""}, false},
		{" 1", &cs1, args{"Key1", "Value 1 (new)"}, true},
		{" 2", &cs1, args{"Key 2", "Value 2 (new)"}, true},
		{" 3", &cs1, args{"Key2", "Value 2 (new)"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.UpdateKey(tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("TSection.UpdateKey() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSection_UpdateKey()

func Benchmark_TSection_String(b *testing.B) {
	sl1 := TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "val2"},
		TKeyVal{"key3", "val3"},
		TKeyVal{"key4", ""},
	}
	for n := 0; n < b.N; n++ {
		if 0 > len(sl1.String()) {
			continue
		}
	}
} // Benchmark_TSection_String()

func Benchmark_TSection_string0(b *testing.B) {
	sl1 := TSection{
		TKeyVal{"key1", "val1"},
		TKeyVal{"key2", "val2"},
		TKeyVal{"key3", "val3"},
		TKeyVal{"key4", ""},
	}
	for n := 0; n < b.N; n++ {
		if 0 > len(sl1.string0()) {
			continue
		}
	}
} // Benchmark_TSection_string0()

func TestTIniList_Clear(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TIniList
	cs := fields(*cis)
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{"1", cs, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TIniList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.Clear(); got != tt.want {
				t.Errorf("TIniList.Clear() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTIniList_Clear()

func TestTIniList_RemoveSection(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TIniList
	type args struct {
		aSection string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{"1", fields(*cis), args{"Default"}, true}, // first
		{"2", fields(*cis), args{"port3"}, true},   // last
		{"3", fields(*cis), args{"sql3"}, true},    // middle
		{"4", fields(*cis), args{"nichda"}, true},  // n.a.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TIniList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.RemoveSection(tt.args.aSection); got != tt.want {
				t.Errorf("TIniList.RemoveSection() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTIniList_RemoveSection()

func TestTIniList_RemoveSectionKey(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TIniList
	type args struct {
		aSection string
		aKey     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{"1", fields(*cis), args{"", "ach jeh"}, true},
		{"2", fields(*cis), args{"Default", "ach jeh"}, true},
		{"3", fields(*cis), args{"sql3", "password"}, true},
		{"4", fields(*cis), args{"sql3", "port"}, true},
		{"5", fields(*cis), args{"port0", ""}, true},
		{"6", fields(*cis), args{"", ""}, true},
		{"7", fields(*cis), args{"general", "n.a."}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TIniList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.RemoveSectionKey(tt.args.aSection, tt.args.aKey); got != tt.want {
				t.Errorf("TIniList.RemoveSectionKey() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTIniList_RemoveSectionKey()

func TestTIniList_String(t *testing.T) {
	id1 := TIniList{
		defSect: "Default",
		secOrder: tOrder{
			"Default",
			"Sect2",
			"NOOP",
		},
		sections: tIniSections{
			"Sect2": &TSection{
				TKeyVal{"key3", "val3"},
				TKeyVal{"key4", ""},
			},
			"Default": &TSection{
				TKeyVal{"key1", "val1"},
				TKeyVal{"key2", "val2"},
			},
		},
	}
	rl1 := "\n[Default]\nkey1 = val1\nkey2 = val2\n\n[Sect2]\nkey3 = val3\nkey4 =\n"
	tests := []struct {
		name   string
		fields TIniList
		want   string
	}{
		// TODO: Add test cases.
		{" 1", id1, rl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TIniList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.String(); got != tt.want {
				t.Errorf("TIniList.String() =\n{%v}, want\n{%v}", got, tt.want)
			}
		})
	}
} // TestTIniList_String()

func Benchmark_TSections_String(b *testing.B) {
	id1 := TIniList{
		defSect: "Default",
		secOrder: tOrder{
			"Default",
			"Sect2",
			"NOOP",
		},
		sections: tIniSections{
			"Sect2": &TSection{
				TKeyVal{"key3", "val3"},
				TKeyVal{"key4", ""},
			},
			"Default": &TSection{
				TKeyVal{"key1", "val1"},
				TKeyVal{"key2", "val2"},
			},
		},
	}
	for n := 0; n < b.N; n++ {
		if 0 > len(id1.String()) {
			continue
		}
	}
} // Benchmark_TSections_String()

func Benchmark_TSections_string0(b *testing.B) {
	id1 := TIniList{
		defSect: "Default",
		secOrder: tOrder{
			"Default",
			"Sect2",
			"NOOP",
		},
		sections: tIniSections{
			"Sect2": &TSection{
				TKeyVal{"key3", "val3"},
				TKeyVal{"key4", ""},
			},
			"Default": &TSection{
				TKeyVal{"key1", "val1"},
				TKeyVal{"key2", "val2"},
			},
		},
	}
	for n := 0; n < b.N; n++ {
		if 0 > len(id1.string0()) {
			continue
		}
	}
} // Benchmark_TSections_string0()

func compare1(aString string) {
	if "" == aString {
		return
	}
} // compare1()

func compare2(aString string) {
	if 0 == len(aString) {
		return
	}
} // compare2()

func Benchmark_compare1(b *testing.B) {
	for n := 0; n < b.N*8; n++ {
		compare1("qwertzuiopü+#äölkjhgfdsa<yxcvbnm,.-^1234567890ß´qwertzuiop")
	}
} // Benchmark_compare1()

func Benchmark_compare2(b *testing.B) {
	for n := 0; n < b.N*8; n++ {
		compare2("qwertzuiopü+#äölkjhgfdsa<yxcvbnm,.-^1234567890ß´qwertzuiop")
	}
} // Benchmark_compare2()

func TestTIniList_updateSectKey(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TIniList
	cs := fields(*cis)
	type args struct {
		aSection string
		aKey     string
		aValue   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{"1", cs, args{"", "", ""}, false},
		{"2", cs, args{"general", "", ""}, false},
		{"3", cs, args{"", "loglevel", ""}, true},
		{"4", cs, args{"general", "loglevel", ""}, true},
		{"5", cs, args{"general", "loglevel", "8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TIniList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.updateSectKey(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("TIniList.updateSectKey() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTIniList_updateSectKey()

func TestTIniList_UpdateSectKeyBool(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TIniList
	cs := fields(*cis)
	type args struct {
		aSection string
		aKey     string
		aValue   bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{"1", cs, args{"", "", false}, false},
		{"2", cs, args{"general", "", true}, false},
		{"3", cs, args{"", "loglevel", false}, true},
		{"4", cs, args{"general", "loglevel", true}, true},
		{"5", cs, args{"general", "loglevel", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TIniList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.UpdateSectKeyBool(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("TIniList.UpdateSectKeyBool() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTIniList_UpdateSectKeyBool()

func TestTIniList_WriteFile(t *testing.T) {
	id, _ := New(inFileName)
	tests := []struct {
		name       string
		id         *TIniList
		wantRBytes int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{"1", id, 4253, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRBytes, err := tt.id.Store()
			if (err != nil) != tt.wantErr {
				t.Errorf("TIniList.WriteFile() error = {%v}, wantErr {%v}", err, tt.wantErr)
				return
			}
			if gotRBytes != tt.wantRBytes {
				t.Errorf("TIniList.WriteFile() = '%v', want '%v'", gotRBytes, tt.wantRBytes)
			}
		})
	}
} // TestTIniList_WriteFile()

func walkFunc(aSect, aKey, aVal string) {
	fmt.Printf("\nSection: %s\nKey: %s\nValue: %s\n", aSect, aKey, aVal)
} // walkFunc()

func TestTSections_Walk(t *testing.T) {
	il, _ := New(inFileName)
	type args struct {
		aFunc TWalkFunc
	}
	tests := []struct {
		name   string
		fields TIniList
		args   args
	}{
		// TODO: Add test cases.
		{" 1", *il, args{walkFunc}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Walk(tt.args.aFunc)
		})
	}
} // TestTSections_Walk()

type tTestWalk int

func (tw tTestWalk) Walk(aSect, aKey, aVal string) {
	fmt.Printf("\nSection: %s\nKey: %s\nValue: %s\n", aSect, aKey, aVal)
} // walkFunc()

func TestTSections_Walker(t *testing.T) {
	il, _ := New(inFileName)
	type args struct {
		aWalker TIniWalker
	}
	var walker tTestWalk
	tests := []struct {
		name   string
		fields *TIniList
		args   args
	}{
		// TODO: Add test cases.
		{" 1", il, args{walker}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Walker(tt.args.aWalker)
		})
	}
} // TestTSections_Walker()
