/*
Copyright © 2019, 2023 M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"testing"
)

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

func TestTIniList_addSection(t *testing.T) {
	// type fields struct {
	// 	defSect  string
	// 	fName    string
	// 	secOrder tSectionOrder
	// 	sections tSectionList
	// }
	type args struct {
		aSection string
	}
	il1 := TIniList{
		defSect:  DefSection,
		fName:    "/tmp/test1.ini",
		secOrder: tSectionOrder{},
		sections: tSectionList{},
	}
	tests := []struct {
		name   string
		fields TIniList
		args   args
		want   bool
	}{
		{"1", il1, args{"SectTest1"}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := &tt.fields
			// il := &TIniList{
			// 	defSect:  tt.fields.defSect,
			// 	fName:    tt.fields.fName,
			// 	secOrder: tt.fields.secOrder,
			// 	sections: tt.fields.sections,
			// }
			if got := il.addSection(tt.args.aSection); got != tt.want {
				t.Errorf("TIniList.addSection() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTIniList_addSection()
