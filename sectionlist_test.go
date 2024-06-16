/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"fmt"
	"reflect"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	inFileName  = `testIn.ini`
	outFilename = `testOut.ini`
)

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
	il1 := TSectionList{
		defSect:  DefSection,
		fName:    "/tmp/test1.ini",
		secOrder: tSectionOrder{},
		sections: tSections{},
	}
	tests := []struct {
		name   string
		fields TSectionList
		args   string
		want   bool
	}{
		{"1", il1, "SectTest1", true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := &tt.fields
			if got := sl.addSection(tt.args); got != tt.want {
				t.Errorf("%q TSectionList.addSection() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTIniList_addSection()

//

func TestTSectionList_Clear(t *testing.T) {
	cis, _ := New(inFileName)
	tests := []struct {
		name   string
		fields *TSectionList
		want   *TSectionList
	}{
		// TODO: Add test cases.
		{"1", cis, cis},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.fields
			if got := id.Clear(); got != tt.want {
				t.Errorf("%q TSectionList.Clear() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_Clear()

func TestTSectionList_RemoveSection(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TSectionList
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
			id := &TSectionList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.RemoveSection(tt.args.aSection); got != tt.want {
				t.Errorf("%q TSectionList.RemoveSection() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_RemoveSection()

func TestTSectionList_RemoveSectionKey(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TSectionList
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
			id := &TSectionList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.RemoveSectionKey(tt.args.aSection, tt.args.aKey); got != tt.want {
				t.Errorf("%q TSectionList.RemoveSectionKey() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_RemoveSectionKey()

func TestTSectionList_Sections(t *testing.T) {
	il0 := TSectionList{
		defSect:  DefSection,
		fName:    "/tmp/test1.ini",
		secOrder: tSectionOrder{},
		sections: tSections{},
	}

	il1 := TSectionList{
		defSect:  DefSection,
		fName:    "/tmp/test1.ini",
		secOrder: tSectionOrder{"One", "Two", "Three"},
		sections: tSections{},
	}
	tests := []struct {
		name   string
		fields TSectionList
		want   []string
		want1  int
	}{
		{"0", il0, il0.secOrder, 0},
		{"1", il1, il1.secOrder, 3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := &tt.fields
			got, got1 := sl.Sections()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q TSectionList.Sections() list = '%v', want '%v'",
					tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("%q TSectionList.Sections() int = %d, want %d",
					tt.name, got1, tt.want1)
			}
		})
	}
} // TestTSectionList_Sections()

func TestTSectionList_String(t *testing.T) {
	id1 := TSectionList{
		defSect: "Default",
		secOrder: tSectionOrder{
			"Default",
			"Sect2",
			"NOOP",
		},
		sections: tSections{
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
		fields TSectionList
		want   string
	}{
		// TODO: Add test cases.
		{" 1", id1, rl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TSectionList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.String(); got != tt.want {
				t.Errorf("%q TSectionList.String() =\n{%v}, want\n{%v}",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_String()

func Benchmark_TSectionList_String(b *testing.B) {
	sl, _ := New(inFileName)
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
	type kArgs struct {
		aSection string
		aKey     string
		aValue   string
	}

	sl, _ := New(inFileName)
	cs := *sl
	tests := []struct {
		name   string
		fields TSectionList
		args   kArgs
		want   bool
	}{
		// TODO: Add test cases.
		{"1", cs, kArgs{"", "", ""}, false},
		{"2", cs, kArgs{"general", "", ""}, false},
		{"3", cs, kArgs{"", "loglevel", ""}, true},
		{"4", cs, kArgs{"general", "loglevel", ""}, true},
		{"5", cs, kArgs{"general", "loglevel", "8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &TSectionList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.updateSectKey(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q TSectionList.updateSectKey() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_updateSectKey()

func TestTSectionList_UpdateSectKeyBool(t *testing.T) {
	cis, _ := New(inFileName)
	type fields TSectionList
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
			id := &TSectionList{
				defSect:  tt.fields.defSect,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			if got := id.UpdateSectKeyBool(tt.args.aSection, tt.args.aKey, tt.args.aValue); got != tt.want {
				t.Errorf("%q TSectionList.UpdateSectKeyBool() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTSectionList_UpdateSectKeyBool()

func TestTSectionList_WriteFile(t *testing.T) {
	ini, _ := New(inFileName)
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
				t.Errorf("%q TSectionList.WriteFile() error = {%v}, wantErr {%v}",
					tt.name, err, tt.wantErr)
				return
			}
			if gotRBytes != tt.wantRBytes {
				t.Errorf("%q TSectionList.WriteFile() = '%v', want '%v'",
					tt.name, gotRBytes, tt.wantRBytes)
			}
		})
	}
} // TestTSectionList_WriteFile()

func walkFunc(aSect, aKey, aVal string) {
	fmt.Printf("\nSection: %s\nKey: %s\nValue: %s\n", aSect, aKey, aVal)
} // walkFunc()

func TestTSections_Walk(t *testing.T) {
	sl, _ := New(inFileName)
	type args struct {
		aFunc TWalkFunc
	}
	tests := []struct {
		name   string
		fields TSectionList
		args   args
	}{
		// TODO: Add test cases.
		{" 1", *sl, args{walkFunc}},
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

func TestTSectionList_Walker(t *testing.T) {
	var walker tTestWalk
	sl, _ := New(inFileName)
	tests := []struct {
		name   string
		fields TSectionList
		args   TIniWalker
	}{
		{"1", *sl, walker},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := &TSectionList{
				defSect:  tt.fields.defSect,
				fName:    tt.fields.fName,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			sl.Walker(tt.args)
		})
	}
} // TestTSectionList_Walker()
