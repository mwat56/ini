/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"strconv"
	"strings"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// TKeyVal represents an INI key/value pair.
	TKeyVal struct {
		Key   string
		Value string
	}
)

// `AsBool()` returns the value as a boolean value.
//
// `0`, `f`, `F`, `n`, and `N` are considered `false` while
// `1`, `t`, `T`, `y`, and `Y` are considered `true`;
// these values will be given in the result value.
//
// This method actually checks only the first character of the key's
// value so one can write e.g. "false" or "NO" (for a `false` result),
// or "True" or "yes" (for a `true` result).
func (kv *TKeyVal) AsBool() (rVal bool) {
	val := kv.Value + `@` // in case of empty string: default FALSE
	switch val[:1] {
	case `0`, `f`, `F`, `n`, `N`:
		return false
	case `1`, `t`, `T`, `y`, `Y`, `j`, `J`, `o`, `O`:
		// True, Yes, Ja (German), Oui (French)`
		return true
	default:
		return false
	}
} // AsBool()

//

// `AsFloat32()` returns the value of the key/value pair as
// a 32bit floating point.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded using
// IEEE754 unbiased rounding.
func (kv *TKeyVal) AsFloat32() (float32, bool) {
	if f64, err := strconv.ParseFloat(kv.Value, 32); nil == err {
		return float32(f64), true
	}

	return float32(0.0), false
} // AsFloat32()

// `AsFloat64` returns the value of the key/value pair as
// a 64bit floating point.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded using
// IEEE754 unbiased rounding.
func (kv *TKeyVal) AsFloat64() (float64, bool) {
	if f64, err := strconv.ParseFloat(kv.Value, 64); nil == err {
		return f64, true
	}

	return 0.0, false
} // AsFloat64()

//

// `AsInt()` returns the value of the key/value pair as an integer.
func (kv *TKeyVal) AsInt() (int, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 0); nil == err {
		return int(i64), true
	}

	return int(0), false
} // AsInt()

// `AsInt8()` returns the value of the key/value pair as an integer.
func (kv *TKeyVal) AsInt8() (int8, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 8); nil == err {
		return int8(i64), true
	}

	return int8(0), false
} // AsInt8()

// `AsInt16()` returns the value of the key/value pair as a 16bit integer.
func (kv *TKeyVal) AsInt16() (int16, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 16); nil == err {
		return int16(i64), true
	}

	return int16(0), false
} // AsInt16()

// `AsInt32()` returns the value of the key/value pair as a 32bit integer.
func (kv *TKeyVal) AsInt32() (int32, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 32); nil == err {
		return int32(i64), true
	}

	return int32(0), false
} // AsInt32()

// `AsInt64()` returns the value of the key/value pair as a 64bit integer.
func (kv *TKeyVal) AsInt64() (int64, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 64); nil == err {
		return i64, true
	}

	return 0, false
} // AsInt64()

//

// `AsString()` returns the value of the key's value as a string.
//
// The second return value will be `true` (it exists for symmetry only).
func (kv *TKeyVal) AsString() (string, bool) {
	return kv.Value, true
} // AsString()

//

// `AsUInt()` returns the value of the key/value pair as an unsigned integer.
func (kv *TKeyVal) AsUInt() (uint, bool) {
	if ui64, err := strconv.ParseUint(kv.Value, 10, 0); nil == err {
		return uint(ui64), true
	}
	return uint(0), false
} // AsUInt()

// `AsUInt8()` returns the value of the key/value pair as an unsigned
// 8bit integer.
func (kv *TKeyVal) AsUInt8() (uint8, bool) {
	if ui64, err := strconv.ParseUint(strings.TrimSpace(kv.Value), 10, 8); nil == err {
		return uint8(ui64), true
	}

	return uint8(0), false
} // AsUInt8()

// `AsUInt16()` returns the value of the key/value pair as an unsigned
// 16bit integer.
func (kv *TKeyVal) AsUInt16() (uint16, bool) {
	if ui64, err := strconv.ParseUint(kv.Value, 10, 16); nil == err {
		return uint16(ui64), true
	}
	return uint16(0), false
} // AsUInt16()

// `AsUInt32()` returns the value of the key/value pair as an unsigned
// 32bit integer.
func (kv *TKeyVal) AsUInt32() (uint32, bool) {
	if ui64, err := strconv.ParseUint(kv.Value, 10, 32); nil == err {
		return uint32(ui64), true
	}
	return uint32(0), false
} // AsUInt32()

// `AsUInt64()` returns the value of the key/value pair as an unsigned
// 64bit integer.
func (kv *TKeyVal) AsUInt64() (uint64, bool) {
	if ui64, err := strconv.ParseUint(kv.Value, 10, 64); nil == err {
		return ui64, true
	}
	return uint64(0), false
} // AsUInt64()

// `String()` returns a string representation of the key/value pair.
//
// The returned string follows the pattern `Key = value`.
func (kv *TKeyVal) String() string {
	if "" == kv.Value {
		return kv.Key + ` =`
	}

	return kv.Key + ` = ` + kv.Value
} // String()

// func (kv *TKeyVal) String2() string {
// 	var sb strings.Builder

// 	if 0 == len(kv.Value) {
// 		_, _ = sb.WriteString(fmt.Sprintf("%s =", kv.Key))
// 	} else {
// 		_, _ = sb.WriteString(fmt.Sprintf("%s = %s", kv.Key, kv.Value))
// 	}
// 	return sb.String()
// } // String2()

// `UpdateValue()` replaces the current value of the key/value pair by
// the provided new `aValue`.
//
//	`aValue` The value of the key/value pair to update.
func (kv *TKeyVal) UpdateValue(aValue string) bool {
	kv.Value = strings.TrimSpace(aValue)

	return true
} // UpdateValue()

/* _EoF_ */
