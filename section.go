/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// internally used string list
	tStringMap map[string]string

	// `TSection` is a slice of key/value pairs.
	TSection struct {
		data tStringMap
		mtx  sync.RWMutex
	}

	// `TSectionWalkFunc()` is used by `Walk()` when visiting the entries
	// in the INI list.
	//
	// see `Walk()`
	TSectionWalkFunc func(aKey, aVal string)
)

// `AddKey()` appends a new key/value pair returning `true` on success or
// `false` otherwise.
//
// If `aKey` is an empty string the method's result will be `false`.
//
// Parameters:
//
//	`aKey` The key of the key/value pair to add.
//	`aValue` The value of the key/value pair to add.
//
// Returns:
//
//	`bool`: `true` if `aKey` was added successfully, `false` otherwise.
func (kl *TSection) AddKey(aKey, aValue string) (rVal bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return
	}

	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	kl.data[aKey] = strings.TrimSpace(aValue)
	_, rVal = kl.data[aKey]

	return
} // AddKey()

// Bool

// `AsBool()` returns the value of `aKey` as a boolean value.
//
// If the given `aKey` doesn't exist then the second (bool) return value
// will be `false`.
//
// `0`, `f`, `F`, `n`, and `N` are considered `false` while
// `1`, `t`, `T`, `y`, `Y`, `j`, `J`, `o`, `O` are considered `true`;
// these values will be given in the first return value with the second
// being `true`.
// This method actually checks only the first character of the key's value
// so one can write e.g. "false" or "NO" (for a `false` result), or "True"
// or "yes" (for a `true` result).
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`bool`: The value of `aKey` a `bool`.
//	`bool`: `true` if the aKey was found, `false` otherwise.
//
// All other values will give `false` as the second return value.
func (kl *TSection) AsBool(aKey string) (bool, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false, false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data[aKey]; exists {
		value += "\t" // in case of empty string: default FALSE
		// Since all values are TRIMed there can never be a TAB at the start.

		switch value[:1] {
		case `0`, `f`, `F`, `n`, `N`:
			return false, true

		case `1`, `t`, `T`, `y`, `Y`, `j`, `J`, `o`, `O`:
			// True, Yes (English), Ja (German), Oui (French)`
			return true, true

			// default:
			// 	return false, false
		}
	}

	return false, false
} // AsBool()

// Float

// `AsFloat32()` returns the value of `aKey` as a 32bit floating point.
//
// If the given `aKey` doesn't exist then the second return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded using
// IEEE754 unbiased rounding.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`float32`: The value of `aKey` as a 32bit floating point.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsFloat32(aKey string) (float32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return float32(0.0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if f64, err := strconv.ParseFloat(val, 32); (nil == err) && (f64 == f64) {
			// for NaN the inequality comparison with itself returns true
			return float32(f64), true
		}
	}

	return float32(0.0), false
} // AsFloat32()

// `AsFloat64()` returns the value of `aKey` as a 64bit floating point.
//
// If the given `aKey` doesn't exist then the second return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded using
// IEEE754 unbiased rounding.
//
// Parameters:
//
//	aKey` the name of the key to lookup.
//
// Returns:
//
//	`float64`: The value of `aKey` as a 64bit floating point.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsFloat64(aKey string) (float64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return float64(0.0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if f64, err := strconv.ParseFloat(val, 64); (nil == err) && (f64 == f64) {
			// for NaN the inequality comparison with itself returns true
			return f64, true
		}
	}

	return float64(0.0), false
} // AsFloat64()

// Int

// `AsInt()` returns the value of `aKey` as an integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup in the list.
//
// Returns:
//
//	`int`: The value of `aKey` as an integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt(aKey string) (int, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if i64, err := strconv.ParseInt(val, 10, 0); nil == err {
			return int(i64), true
		}
	}

	return int(0), false
} // AsInt()

// `AsInt8()` returns the value of `aKey` as an 8bit integer.
//
// If the given `aKey` doesn't exist or the key's value is empty then
// the second return value will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup in the list.
//
// Returns:
//
//	`int8`: The value of `aKey` as an 8bit integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt8(aKey string) (int8, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int8(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if i64, err := strconv.ParseInt(val, 10, 8); nil == err {
			return int8(i64), true
		}
	}

	return int8(0), false
} // AsInt()

// `AsInt16()` returns the value of `aKey` as a 16bit integer.
//
// If the given `aKey` doesn't exist or the key's value is empty then
// the second return value will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup in the list.
//
// Returns:
//
//	`int16`: The value of `aKey` as a 16bit integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt16(aKey string) (int16, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int16(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if i64, err := strconv.ParseInt(val, 10, 16); nil == err {
			return int16(i64), true
		}
	}

	return int16(0), false
} // AsInt16()

// `AsInt32()` returns the value of `aKey` as a 32bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`int32`: The value of `aKey` as a 32bit integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt32(aKey string) (int32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int32(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if i64, err := strconv.ParseInt(val, 10, 32); nil == err {
			return int32(i64), true
		}
	}

	return int32(0), false
} // AsInt32()

// `AsInt64()` returns the value of `aKey` as a 64bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`int64`: The value of `aKey` as a 64bit integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt64(aKey string) (int64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int64(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		if i64, err := strconv.ParseInt(val, 10, 64); nil == err {
			return i64, true
		}
	}

	return int64(0), false
} // AsInt64()

// String

// `AsString()` returns the value of `aKey` as a string.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`string`: The value of `aKey` as a string.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsString(aKey string) (string, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return "", false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if val, exists := kl.data[aKey]; exists {
		return val, true
	}

	return "", false
} // AsString()

// UInt

// `AsUInt()` returns the value of `aKey` as an unsigned integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`uint`: The value of `aKey` as an unsigned integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt(aKey string) (uint, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data[aKey]; exists {
		if ui64, err := strconv.ParseUint(value, 10, 0); nil == err {
			return uint(ui64), true
		}
	}

	return uint(0), false
} // AsUInt()

// `AsUInt8()` returns the value of `aKey` as an unsigned 8bit integer.
//
// If the given `aKey` doesn't exist then the second  return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`uint8`: The value of `aKey` as a 8bit unsigned integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt8(aKey string) (uint8, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint8(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data[aKey]; exists {
		if ui64, err := strconv.ParseUint(value, 10, 8); nil == err {
			return uint8(ui64), true
		}
	}

	return uint8(0), false
} // AsUInt()

// `AsInt16()` returns the value of `aKey` as an unsigned 16bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`uint16`: The value of `aKey` as a 16bit unsigned integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt16(aKey string) (uint16, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint16(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data[aKey]; exists {
		if ui64, err := strconv.ParseUint(value, 10, 16); nil == err {
			return uint16(ui64), true
		}
	}

	return uint16(0), false
} // AsUInt16()

// `AsInt32()` returns the value of `aKey` as an unsigned 32bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`uint32`: The value of `aKey` as a 32bit unsigned integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt32(aKey string) (uint32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint32(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data[aKey]; exists {
		if ui64, err := strconv.ParseUint(value, 10, 32); nil == err {
			return uint32(ui64), true
		}
	}

	return uint32(0), false
} // AsUInt32()

// `AsUInt64()` returns the value of `aKey` as an unsigned 64bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
//
//	`aKey` The name of the key to lookup.
//
// Returns:
//
//	`uint64`: The value of `aKey` as a 64bit unsigned integer.
//	`bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt64(aKey string) (uint64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint64(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data[aKey]; exists {
		if ui64, err := strconv.ParseUint(value, 10, 64); nil == err {
			return ui64, true
		}
	}

	return uint64(0), false
} // AsUInt64()

//

// `Clear()` removes all entries in this INI section.
//
// It returns a pointer to the same section, so that you can chain
// method calls like this:
//
//	kl.Clear().AddKey("key", "value")
//
// This method does not return any error, because it does not
// perform any I/O operation.
func (kl *TSection) Clear() *TSection {
	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	// replace the current list by fresh one
	kl.data = make(map[string]string)

	return kl
} // Clear()

// `CompareTo()` is used to compare two sections for equality.
// It takes another `TSection` instance as a parameter and returns a
// boolean value indicating whether the two sections are equal or not.
//
// If the method iterates over all key-value pairs without finding any
// differences, it returns `true` indicating that the two sections are equal.
// If either the respective section's lengths are different, or if at least
// one of the keyValue pairs are different, the method's result is `false`.
//
// Parameters:
//
//	`aSection` The `TSection` instance to compare.
//
// Returns:
//
//	`bool`: `true` if `aSection` is equal to this instance, `false` otherwise.
func (kl *TSection) CompareTo(aSection *TSection) bool {
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if len(kl.data) != len(aSection.data) {
		return false
	}

	for key, value := range kl.data {
		if val, exists := aSection.data[key]; (!exists) || (val != value) {
			return false
		}
	}

	return true
} // compareTo()

// `HasKey()` returns whether `aKey` exists in this INI section.
//
// Parameters:
//
//	`aKey` The key to lookup.
//
// Returns:
//
//	`bool`: `true` if `aKey` exists in the section, `false` otherwise.
func (kl *TSection) HasKey(aKey string) (rOK bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return
	}

	_, rOK = kl.data[aKey]

	return
} // HasKey()

// `Len()` counts the number of key/value pairs in this section.
//
// Returns:
//
//	`int`: The number of key/value pairs in this section.
func (kl *TSection) Len() int {
	return len(kl.data)
} // Len()

// `Merge()` merges all section key/value pairs into this section.
//
// The method adds all non-existing key/value pairs from `aSection` to this
// section and updates all existing keys with the values from `aSection`
//
// Parameters:
//
//	`aSection`: The INI section to merge with this section.
//
// Returns:
//
//	`TSection`: This INI section including all added/updated key/value pairs.
func (kl *TSection) merge(aSection *TSection) *TSection {
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	for key, value := range aSection.data {
		if val, exists := kl.data[key]; (!exists) || (val != value) {
			kl.data[key] = value
		}
	}

	return kl
} // merge()

func (kl *TSection) merge2(aSection *TSection) *TSection {
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	// Strangely enough despite having no tests this implementation
	// seems _slower_ then the one with at least three tests:
	//
	// goos: linux
	// goarch: amd64
	// pkg: github.com/mwat56/ini
	// cpu: AMD Ryzen 9 5950X 16-Core Processor
	// Benchmark_merge1-32    	  108182	     11898 ns/op	       0 B/op	       0 allocs/op
	// testing: Benchmark_merge1-32 left GOMAXPROCS set to 1
	// Benchmark_merge2-32    	  101205	     11839 ns/op	       0 B/op	       0 allocs/op
	// testing: Benchmark_merge2-32 left GOMAXPROCS set to 1

	for key, value := range aSection.data {
		kl.data[key] = value
	}

	return kl
} // merge()

// `RemoveKey()` removes `aKey` from this section.
//
// This method returns 'true' if `aKey` doesn't exist at all, or if
// `aKey` was successfully removed, and `false` otherwise.
//
// Parameters:
//
//	`aKey` The name of the key/value pair to remove.
//
// Returns:
//
//	`bool`: `true` if `aKey` was successfully removed, `false` otherwise.
func (kl *TSection) RemoveKey(aKey string) bool {
	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	if _, exists := kl.data[aKey]; !exists {
		// if aKey doesn't exist we consider the removal successful
		return true
	}

	delete(kl.data, aKey)
	if _, exists := kl.data[aKey]; !exists {
		// aKey successfully removed
		return true
	}

	return false
} // RemoveKey()

// `String()` returns a string representation of the whole INI section.
//
// The single key/value pairs are delimited by a linefeed ('\n).
//
// NOTE: The order of the key/value pairs is not guaranteed.
//
// Returns:
//
//	`string`: The string representation of the current section.
func (kl *TSection) String() (rString string) {
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	for key, value := range kl.data {
		if "" == value {
			rString += key + " =\n"
		} else {
			rString += key + " = " + value + "\n"
		}
	}

	return
} // String()

// `UpdateKey()` replaces the current value of `aKey` by the provided
// new `aValue`.
//
// If `aKey` is an empty string the method's result will be `false`.
//
// Parameters:
//
//	`aKey` The key of the key/value pair to update.
//	`aValue` The value of the key/value pair to update.
//
// Returns:
//
//	`bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKey(aKey, aValue string) (rOK bool) {
	// if "" == aKey {
	// 	return
	// }
	// kl.mtx.Lock()
	// defer kl.mtx.Unlock()

	// kl.data[aKey] = strings.TrimSpace(aValue)
	// _, rOK = kl.data[aKey]

	return kl.AddKey(aKey, aValue)
} // UpdateKey()

// `UpdateKeyBool()` replaces the current value of `aKey`
// by the provided new `aValue` boolean.
//
// If the given `aValue` is `true` then the string "true" is used
// otherwise the string "false".
//
// Parameters:
//
//	`aKey` The name of the key/value pair to use.
//	`aValue` The boolean value of the key/value pair to update.
//
// Returns:
//
//	`bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyBool(aKey string, aValue bool) bool {
	if aValue {
		return kl.UpdateKey(aKey, `true`)
	}

	return kl.UpdateKey(aKey, `false`)
} // UpdateKeyBool()

// `UpdateKeyFloat()` replaces the current value of `aKey`
// by the provided new `aValue` float.
//
// Parameters:
//
//	`aKey` The name of the key/value pair to use.
//	`aValue` The float64 value of the key/value pair to update.
//
// Returns:
//
//	`bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateSectKeyFloat(aKey string, aValue float64) bool {
	return kl.UpdateKey(aKey, fmt.Sprintf("%f", aValue))
} // UpdateKeyFloat()

// `UpdateKeyInt()` replaces the current value of `aKey`
// by the provided new `aValue` integer.
//
// Parameters:
//
//	`aKey` The name of the key/value pair to use.
//	`aValue` The int64 value of the key/value pair to update.
//
// Returns:
//
//	`bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyInt(aKey string, aValue int64) bool {
	return kl.UpdateKey(aKey, fmt.Sprintf("%d", aValue))
} // UpdateKeyInt()

// `UpdateKeyUInt()` replaces the current value of `aKey`
// by the provided new `aValue` unsigned integer.
//
// Parameters:
//
//	`aKey` The name of the key/value pair to use.
//	`aValue` The int64 value of the key/value pair to update.
//
// Returns:
//
//	`bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyUInt(aKey string, aValue uint64) bool {
	return kl.UpdateKey(aKey, fmt.Sprintf("%d", aValue))
} // UpdateKeyUInt()

// `UpdateKeyStr` replaces the current value of `aKey`
// by the provided new `aValue` string.
//
// Parameters:
//
//	`aKey` The name of the key/value pair to use.
//	`aValue` The string value of the key/value pair to update.
//
// Returns:
//
//	`bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyStr(aKey, aValue string) bool {
	return kl.UpdateKey(aKey, aValue)
} // UpdateKeyStr()

// `mergeWalker()` inserts the given key/value pair in this `TSection`.
//
// This method is called by the `Merge()` method.
//
// Parameters:
//
//	`aKey` The name of the key/value pair to use.
//	`aValue` The value of the key/value pair to update.
func (kl *TSection) mergeWalker(aKey, aValue string) {
	kl.AddKey(aKey, aValue) // ignore the method's result
} // mergeWalker

// `Merge()` copies or merges all section key/value pairs
// into this section.
//
// Parameters:
//
//	`aSection` The INI section to merge with this section.
//
// Returns:
//
//	`TSection` This INI section that got merged with the other one.
func (kl *TSection) Merge(aSection *TSection) *TSection {
	aSection.Walk(kl.mergeWalker)

	return kl
} // Merge()

// `Walk()` traverses through all entries in the section calling
// `aFunc` for each entry.
//
// Parameters:
//
//	`aFunc` The function called for each key/value pair in the sections.
func (kl *TSection) Walk(aFunc TSectionWalkFunc) {
	for key, value := range kl.data {
		aFunc(key, value)
	}
} // Walk()

// utility function

// `NewSection()` returns a new instance of `TSection`.
//
// Returns:
//
//	`*TSection`: A new instance of `TSection`.
func NewSection() *TSection {
	return &TSection{
		data: make(tStringMap),
	}
} // NewSection()

/* _EoF_ */
