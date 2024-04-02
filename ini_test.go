/*
Copyright © 2019, 2022 M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"fmt"
	"testing"
)

const (
	inFileName  = `testIn.ini`
	outFilename = `testOut.ini`
)

func TestTIniList_Clear(t *testing.T) {
	cis, _ := New(inFileName)
	tests := []struct {
		name   string
		fields *TIniList
		want   *TIniList
	}{
		// TODO: Add test cases.
		{"1", cis, cis},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.fields
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
		secOrder: tSectionOrder{
			"Default",
			"Sect2",
			"NOOP",
		},
		sections: tSectionList{
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
	type kArgs struct {
		aSection string
		aKey     string
		aValue   string
	}

	il, _ := New(inFileName)
	cs := TIniList(*il)
	tests := []struct {
		name   string
		fields TIniList
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
	ini, _ := New(inFileName)
	ini.SetFilename(outFilename)

	tests := []struct {
		name       string
		id         *TIniList
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

func TestTIniList_Walker(t *testing.T) {
	type args struct {
		aWalker TIniWalker
	}
	var walker tTestWalk
	il, _ := New(inFileName)
	tests := []struct {
		name   string
		fields TIniList
		args   args
	}{
		{"1", *il, args{walker}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := &TIniList{
				defSect:  tt.fields.defSect,
				fName:    tt.fields.fName,
				secOrder: tt.fields.secOrder,
				sections: tt.fields.sections,
			}
			il.Walker(tt.args.aWalker)
		})
	}
} // TestTIniList_Walker()

/* _EoF_ */
