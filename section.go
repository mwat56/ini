/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	// Default list capacity.
	kvDefCapacity = 16
)

type (
	// `tKeyVal` represents an key/value pair.
	tKeyVal struct {
		Key   string
		Value string
	}
	// a list of key/value pairs
	tKeyValList []tKeyVal

	// `TSection` is a slice of sorted key/value pairs.
	TSection struct {
		data tKeyValList
		mtx  sync.RWMutex
	}

	// `TSectionWalkFunc()` is used by `Walk()` when visiting the entries
	// in the section.
	//
	// see `Walk()`
	TSectionWalkFunc func(aKey, aVal string)

	// A `TSectionWalker` is used by `Walker()` when visiting an entry
	// in the section's list.
	//
	// see `Walker()`
	TSectionWalker interface {
		Walk(aKey, aVal string)
	}
)

// --------------------------------------------------------------------------

// `compareTo()` returns whether the given list of key/value pairs is
// equal to this instance.
//
// It takes another `tKeyValList` instance as a parameter and returns a
// boolean value indicating whether the two sections are equal or not.
//
// If the method iterates over all key-value pairs without finding any
// differences, it returns `true` indicating that the two sections are equal.
// If either the respective section's lengths are different, or if at least
// one of the keyValue pairs are different, the method's result is `false`.
//
// Parameters:
// - `aList` The other `tKeyValList` instance to compare with.
//
// Returns:
// - `bool`: `true` if `aList` is equal to this instance, or `false` otherwise.
func (kvl tKeyValList) compareTo(aList *tKeyValList) bool {
	if len(kvl) != len(*aList) {
		return false
	}

	for _, kv := range kvl {
		val, exists := aList.value(kv.Key)
		if (!exists) || (val != kv.Value) {
			return false
		}
	}

	return true
} // compareTo()

// `copy()` returns a new `tKeyValList` that is a copy of the original list.
//
// This method creates a new slice and appends all elements from the original
// list to the new slice. It returns the new slice.
//
// Returns:
// - `tKeyValList`: A new slice that is a copy of the original list.
func (kvl tKeyValList) copy() *tKeyValList {
	twin := make(tKeyValList, 0, len(kvl))
	twin = append(twin, kvl...)

	return &twin
} // copy()

// `hasKey()` checks if a given key exists in the list.
//
// Parameters:
// - `aKey` string: The key to search for in the list.
//
// Returns:
// - `bool`: `true` if the key is found in the list, or `false` otherwise.
func (kvl tKeyValList) hasKey(aKey string) bool {
	for _, entry := range kvl {
		if aKey == entry.Key {
			return true
		}
	}

	return false
} // hasKey()

// `insert()` inserts a new key/value pair returning `true` on success or
// `false` otherwise.
//
// If `aKey` is an empty string the method's result will be `false`.
//
// Parameters:
// - `aKeyVal` The key/value pair to add.
//
// Returns:
// - `bool`: `true` if `aKeyVal` was added successfully, `false` otherwise.
func (kvl *tKeyValList) insert(aKeyVal tKeyVal) bool {
	if aKeyVal.Key = strings.TrimSpace(aKeyVal.Key); "" == aKeyVal.Key {
		return false
	}

	sLen := len(*kvl)
	idx := sort.Search(sLen, func(i int) bool {
		return (*kvl)[i].Key >= aKeyVal.Key
	})

	if sLen == idx { // key not found
		*kvl = append(*kvl, tKeyVal{})
		copy((*kvl)[idx+1:], (*kvl)[idx:])
	} else if (*kvl)[idx].Key != aKeyVal.Key { // it's a new key
		*kvl = append(*kvl, tKeyVal{})
		copy((*kvl)[idx+1:], (*kvl)[idx:])
	}
	(*kvl)[idx] = aKeyVal // update the vale

	return true
} // insert()

func (kvl *tKeyValList) merge(aList *tKeyValList) *tKeyValList {
	for _, kv := range *aList {
		kvl.insert(tKeyVal{kv.Key, kv.Value})
	}

	return kvl
} // merge()

// `remove()` deletes `aKey` in the list of key/value pairs.
//
// Parameters:
// - `aKey` string: The name of the key to remove.
//
// Returns:
// - `bool`: `true` if the `aKey` is found/removed, or `false` otherwise.
func (kvl *tKeyValList) remove(aKey string) bool {
	for idx, entry := range *kvl {
		if aKey == entry.Key {
			(*kvl) = append((*kvl)[:idx], (*kvl)[idx+1:]...)
			return true
		}
	}

	return false
} // remove()

// `String()` returns a string representation of the whole INI section.
//
// The single key/value pairs are delimited by a linefeed ('\n).
//
// Returns:
// - `string`: The string representation of the current section.
func (kvl tKeyValList) String() (rString string) {
	for _, kv := range kvl {
		if "" == kv.Value {
			rString += kv.Key + " =\n"
		} else {
			rString += kv.Key + " = " + kv.Value + "\n"
		}
	}

	return
} // String()

// `value()` returns the value of `aKey` as a string.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `string, bool`: The value associated with `aKey`.
// - `bool`: `true` if the aKey was found, `false` otherwise.
func (kvl tKeyValList) value(aKey string) (string, bool) {
	for _, kv := range kvl {
		if aKey == kv.Key {
			return kv.Value, true
		}
	}

	return "", false
} // value()

// --------------------------------------------------------------------------

// `AddKey()` appends a new key/value pair returning `true` on success or
// `false` otherwise.
//
// If `aKey` is an empty string the method's result will be `false`.
//
// Parameters:
// - `aKey` The key of the key/value pair to add.
// - `aValue` The value of the key/value pair to add.
//
// Returns:
// - `bool`: `true` if `aKey` was added successfully, `false` otherwise.
func (kl *TSection) AddKey(aKey, aValue string) bool {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false
	}
	kv := tKeyVal{aKey, strings.TrimSpace(aValue)}

	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	return kl.data.insert(kv)
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `bool`: The value of `aKey` a `bool`.
// - `bool`: `true` if the aKey was found, `false` otherwise.
// All other values will give `false` as the second return value.
func (kl *TSection) AsBool(aKey string) (bool, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false, false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		value += "\t" // in case of empty string: default FALSE
		// Since all values are TRIMed there can never be a TAB at the start.

		switch value[:1] {
		case `0`, `f`, `F`, `n`, `N`:
			return false, true

		case `1`, `t`, `T`, `y`, `Y`, `j`, `J`, `o`, `O`:
			// True, Yes (English), Ja (German), Oui (French)`
			return true, true
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
// -`aKey` The name of the key to lookup.
//
// Returns:
// - `float32`: The value of `aKey` as a 32bit floating point.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsFloat32(aKey string) (float32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return float32(0.0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if f64, err := strconv.ParseFloat(value, 32); (nil == err) && (f64 == f64) {
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
// - `aKey` the name of the key to lookup.
//
// Returns:
// - `float64`: The value of `aKey` as a 64bit floating point.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsFloat64(aKey string) (float64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return float64(0.0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if f64, err := strconv.ParseFloat(value, 64); (nil == err) && (f64 == f64) {
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
// - `aKey` The name of the key to lookup in the list.
//
// Returns:
// - `int`: The value of `aKey` as an integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt(aKey string) (int, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if i64, err := strconv.ParseInt(value, 10, 0); nil == err {
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
// - `aKey` The name of the key to lookup in the list.
//
// Returns:
// - `int8`: The value of `aKey` as an 8bit integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt8(aKey string) (int8, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int8(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if i64, err := strconv.ParseInt(value, 10, 8); nil == err {
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
// - `aKey` The name of the key to lookup in the list.
//
// Returns:
// - `int16`: The value of `aKey` as a 16bit integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt16(aKey string) (int16, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int16(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if i64, err := strconv.ParseInt(value, 10, 16); nil == err {
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int32`: The value of `aKey` as a 32bit integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt32(aKey string) (int32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int32(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if i64, err := strconv.ParseInt(value, 10, 32); nil == err {
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int64`: The value of `aKey` as a 64bit integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsInt64(aKey string) (int64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int64(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if i64, err := strconv.ParseInt(value, 10, 64); nil == err {
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `string`: The value of `aKey` as a string.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsString(aKey string) (string, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return "", false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		return value, true
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint`: The value of `aKey` as an unsigned integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt(aKey string) (uint, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint8`: The value of `aKey` as a 8bit unsigned integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt8(aKey string) (uint8, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint8(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if ui64, err := strconv.ParseUint(value, 10, 8); nil == err {
			return uint8(ui64), true
		}
	}

	return uint8(0), false
} // AsUInt8()

// `AsUInt16()` returns the value of `aKey` as an unsigned 16bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint16`: The value of `aKey` as a 16bit unsigned integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt16(aKey string) (uint16, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint16(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
		if ui64, err := strconv.ParseUint(value, 10, 16); nil == err {
			return uint16(ui64), true
		}
	}

	return uint16(0), false
} // AsUInt16()

// `AsUInt32()` returns the value of `aKey` as an unsigned 32bit integer.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
// Parameters:
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint32`: The value of `aKey` as a 32bit unsigned integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt32(aKey string) (uint32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint32(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
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
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint64`: The value of `aKey` as a 64bit unsigned integer.
// - `bool`: `true` if `aKey` was found, `false` otherwise.
func (kl *TSection) AsUInt64(aKey string) (uint64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint64(0), false
	}

	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	if value, exists := kl.data.value(aKey); exists {
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
// Returns:
// - `TSection`: The current section.
func (kl *TSection) Clear() *TSection {
	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	// replace the current list by fresh/empty one
	kl.data = make(tKeyValList, 0, kvDefCapacity)

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
// - `aSection` The `TSection` instance to compare.
//
// Returns:
// - `bool`: `true` if `aSection` is equal to this instance, `false` otherwise.
func (kl *TSection) CompareTo(aSection *TSection) bool {
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	kvl := aSection.data
	return kl.data.compareTo(&kvl)
} // compareTo()

// `Copy()` returns a copy of the current section.
//
// Returns:
// - `*TSection`: A copy of the current section.
func (kl *TSection) Copy() (rSection *TSection) {
	rSection = &TSection{
		data: make(tKeyValList, 0, kvDefCapacity),
	}
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	kvl := kl.data.copy()
	rSection.data = *kvl

	return
} // Copy()

// `HasKey()` returns whether `aKey` exists in this INI section.
//
// Parameters:
// - `aKey` The key to lookup.
//
// Returns:
// - `bool`: `true` if `aKey` exists in the section, `false` otherwise.
func (kl *TSection) HasKey(aKey string) bool {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false
	}
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	return kl.data.hasKey(aKey)
} // HasKey()

// `Len()` counts the number of key/value pairs in this section.
//
// Returns:
// - `int`: The number of key/value pairs in this section.
func (kl *TSection) Len() int {
	return len(kl.data)
} // Len()

// `Merge()` merges all section key/value pairs into this section.
//
// The method adds all non-existing key/value pairs from `aSection` to this
// section and updates all existing keys with the values from `aSection`
//
// Parameters:
// - `aSection`: The INI section to merge with this section.
//
// Returns:
// - `TSection`: This section added/updated from `aSection`.
func (kl *TSection) Merge(aSection *TSection) *TSection {
	if nil == aSection {
		return kl
	}
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	merged := kl.data.merge(&aSection.data)
	kl.data = *merged

	return kl
} // Merge()

// `RemoveKey()` removes `aKey` from this section.
//
// This method returns 'true' if `aKey` doesn't exist at all, or if
// `aKey` was successfully removed, and `false` otherwise.
//
// Parameters:
// - `aKey` The name of the key/value pair to remove.
//
// Returns:
// - `bool`: `true` if `aKey` was successfully removed, `false` otherwise.
func (kl *TSection) RemoveKey(aKey string) bool {
	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	if kl.data.remove(aKey) {
		return true
	}

	// `aKey` not found is considered removed
	return true
} // RemoveKey()

// `Sort()` sorts the key/value pairs in the section alphabetically by key.
//
// The original map is replaced with the new sorted map.
//
// The method returns a pointer to the same section, so that one can chain
// method calls like this:
//
//	kl.Sort().AddKey("key", "value")
//
// Returns:
// - *TSection: A pointer to the same section after sorting the key-value pairs.
func (kl *TSection) Sort() *TSection {
	kl.mtx.Lock()
	defer kl.mtx.Unlock()

	sort.Slice(kl.data, func(i, j int) bool {
		return kl.data[i].Key < kl.data[j].Key
	})

	return kl
} // Sort()

// `String()` returns a string representation of the whole INI section.
//
// The single key/value pairs are delimited by a linefeed ('\n).
//
// Returns:
// - `string`: The string representation of the current section.
func (kl *TSection) String() (rString string) {
	kl.mtx.RLock()
	defer kl.mtx.RUnlock()

	return kl.data.String()
} // String()

// `UpdateKey()` replaces the current value of `aKey` by the provided
// new `aValue`.
//
// If `aKey` is an empty string the method's result will be `false`.
//
// Parameters:
// - `aKey` The key of the key/value pair to update.
// - `aValue` The value of the key/value pair to update.
//
// Returns:
// - `bool`: `true` if `aKey` was updated successfully, `false` otherwise.
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
// - `aKey` The name of the key/value pair to use.
// - `aValue` The boolean value of the key/value pair to update.
//
// Returns:
// - `bool`: `true` if `aKey` was updated successfully, `false` otherwise.
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
// - `aKey` The name of the key/value pair to use.
// - `aValue` The float64 value of the key/value pair to update.
//
// Returns:
// - `bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateSectKeyFloat(aKey string, aValue float64) bool {
	return kl.UpdateKey(aKey, fmt.Sprintf("%f", aValue))
} // UpdateKeyFloat()

// `UpdateKeyInt()` replaces the current value of `aKey`
// by the provided new `aValue` integer.
//
// Parameters:
// - `aKey` The name of the key/value pair to use.
// - `aValue` The int64 value of the key/value pair to update.
//
// Returns:
// - `bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyInt(aKey string, aValue int64) bool {
	return kl.UpdateKey(aKey, fmt.Sprintf("%d", aValue))
} // UpdateKeyInt()

// `UpdateKeyUInt()` replaces the current value of `aKey`
// by the provided new `aValue` unsigned integer.
//
// Parameters:
// - `aKey` The name of the key/value pair to use.
// - `aValue` The int64 value of the key/value pair to update.
//
// Returns:
// - `bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyUInt(aKey string, aValue uint64) bool {
	return kl.UpdateKey(aKey, fmt.Sprintf("%d", aValue))
} // UpdateKeyUInt()

// `UpdateKeyStr` replaces the current value of `aKey`
// by the provided new `aValue` string.
//
// Parameters:
// - `aKey` The name of the key/value pair to use.
// - `aValue` The string value of the key/value pair to update.
//
// Returns:
// - `bool`: `true` if `aKey` was updated successfully, `false` otherwise.
func (kl *TSection) UpdateKeyStr(aKey, aValue string) bool {
	return kl.UpdateKey(aKey, aValue)
} // UpdateKeyStr()

// `Walk()` traverses through all entries in the section calling
// `aFunc` for each entry.
//
// Parameters:
// - `aFunc` The function called for each key/value pair in the sections.
func (kl *TSection) Walk(aFunc TSectionWalkFunc) {
	for _, kv := range kl.data {
		aFunc(kv.Key, kv.Value)
	}
} // Walk()

// `Walker()` traverses through all entries of current section
// calling `aWalker` for each entry.
//
// Parameters:
// - `aWalker` An object implementing the `TSectionWalker` interface.
func (kl *TSection) Walker(aWalker TSectionWalker) {
	kl.Walk(aWalker.Walk)
} // Walker()

// utility function

// `NewSection()` returns a new instance of `TSection`.
//
// Returns:
// - `*TSection`: A new instance of `TSection`.
func NewSection() *TSection {
	return &TSection{
		data: make(tKeyValList, 0, kvDefCapacity),
	}
} // NewSection()

/* _EoF_ */
