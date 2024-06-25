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
