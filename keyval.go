/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import "strconv"

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// TKeyVal represents an INI key/value pair.
	TKeyVal struct {
		Key   string
		Value string
	}
)

// `AsBool` returns the value as a boolean value.
//
// `0`, `f`, `F`, `n`, and `N` are considered `false` while
// `1`, `t`, `T`, `y`, and `Y` are considered `true`;
// these values will be given in the result value.
//
// This method actually checks only the first character of the key's
// value so one can write e.g. "false" or "NO" (for a `false` result),
// or "True" or "yes" (for a `true` result).
func (kv *TKeyVal) AsBool() (rVal bool) {
	val := kv.Value + `0` // in case of empty string: default FALSE
	switch val[:1] {
	case `0`, `f`, `F`, `n`, `N`:
		return false
	case `1`, `t`, `T`, `y`, `Y`:
		return true
	default:
		return false
	}
} // AsBool()

// `AsFloat32` returns the value of the key/value pair as
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

	return 0.0, false
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

// `AsInt` returns the value of the key/value pair as an integer.
func (kv *TKeyVal) AsInt() (int, bool) {
	if i, err := strconv.Atoi(kv.Value); nil == err {
		return i, true
	}

	return 0, false
} // AsInt()

// `AsInt16` returns the value of the key/value pair as a 16bit integer.
func (kv *TKeyVal) AsInt16() (int16, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 16); nil == err {
		return int16(i64), true
	}

	return 0, false
} // AsInt16()

// `AsInt32` returns the value of the key/value pair as a 32bit integer.
func (kv *TKeyVal) AsInt32() (int32, bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 32); nil == err {
		return int32(i64), true
	}

	return 0, false
} // AsInt32()

// `AsInt64` returns the value of the key/value pair as a 64bit integer.
func (kv *TKeyVal) AsInt64() (rVal int64, rOK bool) {
	if i64, err := strconv.ParseInt(kv.Value, 10, 64); nil == err {
		return i64, true
	}

	return 0, false
} // AsInt64()

// `AsString` returns the value of the key/value pair as a string.
//
// If the key's value is empty then the second return value
// will be `false`.
func (kv *TKeyVal) AsString() (string, bool) {
	return kv.Value, (0 < len(kv.Value))
} // AsString()

// `String` returns a string representation of the key/value pair.
//
// The returned string follows the pattern `Key = value`.
func (kv *TKeyVal) String() string {
	if 0 == len(kv.Value) {
		return kv.Key + ` =`
	}

	return kv.Key + ` = ` + kv.Value
} // String()

// `UpdateValue` replaces the current value of the key/value pair by
// the provided new `aValue`.
//
// If `aValue` is an empty string the method's result will be `false`.
//
//	`aValue` The value of the key/value pair to update.
func (kv *TKeyVal) UpdateValue(aValue string) bool {
	kv.Value = aValue
	return (0 < len(kv.Value))
} // UpdateValue()

/* _EoF_ */
