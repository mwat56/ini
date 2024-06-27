/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"runtime"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestNewIni(t *testing.T) {
	fName := "testIn.ini"

	tests := []struct {
		name     string
		filename string
		// want     *TSectionList
		wantErr bool
	}{
		{"0", "", true},
		{"1", fName, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewIni(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: NewIni() error = %q, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
		})
	}
} // TestNewIni()

const cmpstring = "qwertzuiopü+#äölkjhgfdsa<yxcvbnm,.-^1234567890ß´qwertzuiop"

func Benchmark_compare1(b *testing.B) {
	runtime.GOMAXPROCS(1)

	for n := 0; n < b.N<<5; n++ {
		if "" == cmpstring {
			continue
		}
	}
} // Benchmark_compare1()

func Benchmark_compare2(b *testing.B) {
	runtime.GOMAXPROCS(1)

	for n := 0; n < b.N<<5; n++ {
		if 0 == len(cmpstring) {
			continue
		}
	}
} // Benchmark_compare2()

/* _EoF_ */
