// Package ini provides functions to read/write INI files from/to disc
// and methods to access the section's key-value pairs.
package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type (
	// TKeyVal represents an INI key-value pair.
	TKeyVal struct {
		Key   string
		Value string
	}

	// TSection is a slice of key-value pairs.
	TSection []TKeyVal

	// `tIniSections` is a list (map) of INI sections.
	tIniSections map[string]*TSection

	// A helper slice of strings (i.e. section names)
	// used to preserve the order of INI sections.
	tOrder = []string

	// This opaque data structure is filled by e.g. `LoadFile(…)`.
	tSections struct {
		// name of default section:
		defSect string
		// slice containing the order of sections:
		secOrder tOrder
		// list of INI sections:
		sections tIniSections
	}
)

// TSections is a list of INI sections.
//
// This opaque data structure is filled by e.g. `LoadFile(…)`.
//
// For accessing the sections and key-value pairs
// it provides the appropriate methods.
type TSections tSections

const (
	defCapacity = 16

	// DefSection is the name of the default section in the INI file.
	DefSection = "Default"
)

// Regular expressions to identify certain parts of an INI file.
var (
	// match: [section]
	sectionRE = regexp.MustCompile(`^\[\s*([^\]]*?)\s*]$`)
	// match: key = val
	keyValRE = regexp.MustCompile(`^([^=]+?)\s*=\s*(.*)$`)
	// quoted ' " string " '
	quotesRE = regexp.MustCompile(`^(['"])(.*)(['"])$`)
)

// trimRemoveQuotes returns a quoted string w/o the quote characters.
func trimRemoveQuotes(aString string) string {
	// remove leading/trailing UTF whitespace:
	result := strings.TrimSpace(aString)
	// get a slice of RegEx matches:
	matches := quotesRE.FindStringSubmatch(result)

	// we expect: (1) leading quote, (2) text, (3) trailing quote
	if (2 < len(matches)) && (matches[1] == matches[3]) {
		// rFiltered = strings.TrimSpace(matches[2])
		result = matches[2]
	}

	return result
} // trimRemoveQuotes()

// String returns a string representation of a key-value pair.
func (kv *TKeyVal) String() string {
	if 0 == len(kv.Value) {
		return kv.Key + " ="
	}

	return kv.Key + " = " + kv.Value
} // String()

// string0 is the initial but slower implementation.
func (kv *TKeyVal) string0() string {
	// NOTE: this implementation is ~7 times slower than the one above
	if 0 == len(kv.Value) {
		return fmt.Sprintf("%s =", kv.Key)
	}

	return fmt.Sprintf("%s = %s", kv.Key, kv.Value)
} // string0()

/*
 * Methods of TSection objects.
 */

// AddKey appends a new key-value pair
// returning `true` on success or `false` otherwise.
//
// `aKey` the key of the key-value pair to add.
//
// `aValue` the value of the key-value pair to add.
func (cs *TSection) AddKey(aKey, aValue string) bool {
	*cs = append(*cs, TKeyVal{aKey, aValue})

	if val, ok := cs.AsString(aKey); ok {
		return (val == aValue)
	}

	return false
} // AddKey()

// AsBool returns the value of `aKey` in `aSection` as a boolean value.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// "0", "f", "F", "n" and "N" are considered `false` while
// "1", "t", "T", "y" and "Y" are considered 'true';
// these values will be given in the first result value.
// All other values will give `false` as the second result value.
//
// This method actually checks only the first character of the key's
// value so one can write e.g. "false" or "NO" (for a `false` result),
// or "True" or "yes" (for a 'true' result).
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsBool(aKey string) (bool, bool) {
	if val, exists := cs.AsString(aKey); exists {
		val = val[:1]
		switch val {
		case "0", "f", "F", "n", "N":
			return false, true
		case "1", "t", "T", "y", "Y":
			return true, true
		}
	}

	return false, false
} // AsBool()

// AsFloat32 returns the value of `aKey` as a 32bit floating point.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// If s is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsFloat32(aKey string) (float32, bool) {
	if val, exists := cs.AsString(aKey); exists {
		if f64, err := strconv.ParseFloat(val, 32); nil == err {
			return float32(f64), true
		}
	}

	return 0, false
} // AsFloat32()

// AsFloat64 returns the value of `aKey` as a 64bit floating point.
//
// If the given `aKey` doesn't exist then the second (bool) return
// value will be `false`.
//
// If s is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsFloat64(aKey string) (float64, bool) {
	if val, exists := cs.AsString(aKey); exists {
		if f64, err := strconv.ParseFloat(val, 64); nil == err {
			return f64, true
		}
	}

	return 0, false
} // AsFloat64()

// AsInt returns the value of `aKey` as an integer.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsInt(aKey string) (int, bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i, err := strconv.Atoi(val); nil == err {
			return i, true
		}
	}

	return 0, false
} // AsInt()

// AsInt16 return the value of `aKey` as a 16bit integer.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsInt16(aKey string) (int16, bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i64, err := strconv.ParseInt(val, 10, 16); nil == err {
			return int16(i64), true
		}
	}

	return 0, false
} // AsInt16()

// AsInt32 return the value of `aKey` as a 32bit integer.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsInt32(aKey string) (int32, bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i64, err := strconv.ParseInt(val, 10, 32); nil == err {
			return int32(i64), true
		}
	}

	return 0, false
} // AsInt32()

// AsInt64 return the value of `aKey` as a 64bit integer.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsInt64(aKey string) (int64, bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i64, err := strconv.ParseInt(val, 10, 64); nil == err {
			return i64, true
		}
	}

	return 0, false
} // AsInt64()

// AsString returns the value of `aKey` as a string.
//
// If the given `aKey` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aKey` the name of the key to lookup.
func (cs *TSection) AsString(aKey string) (string, bool) {
	for _, val := range *cs {
		if val.Key == aKey {
			return val.Value, true
		}
	}
	return "", false
} // AsString()

// Clear removes all entries in this INI section.
func (cs *TSection) Clear() *TSection {
	(*cs) = (*cs)[:0]

	return cs
} // Clear()

// HasKey returns whether `aKey` exists in this INI section.
//
// `aKey` the key to lookup.
func (cs *TSection) HasKey(aKey string) bool {
	for _, kv := range *cs {
		if kv.Key == aKey {
			return true
		}
	}

	return false
} // HasKey()

// RemoveKey removes `aKey` from this section.
//
// This method returns 'true' if `aKey` doesn't exist at all,
// or if `aKey` was successfully removed, or `false` otherwise.
//
// `aKey` the name of the key-value pair to remove.
func (cs *TSection) RemoveKey(aKey string) bool {
	slen := len(*cs) - 1 // new slice length (i.e. one shorter)
	for idx, kv := range *cs {
		if kv.Key != aKey {
			continue
		}
		(*cs)[idx] = TKeyVal{}
		if 0 == idx {
			(*cs) = (*cs)[1:]
		} else if idx == slen {
			(*cs) = (*cs)[:slen]
		} else {
			(*cs) = append((*cs)[:idx], (*cs)[1+idx:]...)
		}
		return true
	}

	return (!cs.HasKey(aKey))
} // RemoveKey()

// String returns a string representation of an INI section.
//
// The single key-value pairs are delimited by a linefeed ('\n).
func (cs *TSection) String() (rString string) {
	for _, kv := range *cs {
		rString += kv.String() + "\n"
	}

	return
} // String()

// string0 is the initial but slower implementation.
func (cs *TSection) string0() string {
	// NOTE: this implementation is ~3 times slower than the one above.
	var result bytes.Buffer

	for _, kv := range *cs {
		result.WriteString(kv.string0() + "\n")
	}

	return result.String()
} // string0()

// UpdateKey replaces the current value of `aKey`
// by the provided new `aValue`.
func (cs *TSection) UpdateKey(aKey, aValue string) bool {
	if 0 == len(aKey) {
		return false
	}
	for idx, val := range *cs {
		if val.Key == aKey {
			kv := TKeyVal{Key: aKey, Value: aValue}
			(*cs)[idx] = kv
			return true
		}
	}

	// if aKey doesn't exist then create a new entry
	return cs.AddKey(aKey, aValue)
} // updateKey()

// addSection appends a new INI section returning `true` on success
// or `false` otherwise.
//
// `aSection` name of the INI section to add.
func (id *TSections) addSection(aSection string) bool {
	if _, exists := id.sections[aSection]; exists {
		return true // already there: nothing more to do
	}

	// we make room for initially 8 key-value pairs
	sect := make(TSection, 0, defCapacity)
	id.sections[aSection] = &sect
	if _, ok := id.sections[aSection]; ok {
		// add new section name to order list
		id.secOrder = append(id.secOrder, aSection)

		return true
	}

	return false
} // addSection()

// AddSectionKey appends a new key-value pair to `aSection`
// returning `true` on success or `false` otherwise.
//
// `aSection` name of the INI section to use.
//
// `aKey` the key of the key-value pair to add.
//
// `aValue` the value of the key-value pair to add.
func (id *TSections) AddSectionKey(aSection, aKey, aValue string) bool {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if !id.addSection(aSection) {
		return false // can't find nor add the section
	}

	if cs, exists := id.sections[aSection]; exists {
		return cs.AddKey(aKey, aValue)
	}

	return false
} // AddSectionKey()

/*
 * Public methods to return INI values from a section as a certain data type.
 */

// AsBool returns the value of `aKey` in `aSection` as a boolean value.
//
// If the given aKey in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// "0", "f", "F", "n" and "N" are considered `false` while
// "1", "t", "T", "y" and "Y" are considered 'true';
// these values will be given in the first result value.
// All other values will give `false` as the second result value.
//
// This method actually checks only the first character of the key's
// value so one can write e.g. "false" or "NO" (for a `false` result),
// or "True" or "yes" (for a 'true' result).
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsBool(aSection, aKey string) (bool, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsBool(aKey)
	}

	return false, false
} // AsBool()

// AsFloat32 returns the value of `aKey` in `aSection` as
// a 32bit floating point.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (bool) return value will be `false`.
//
// If s is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsFloat32(aSection, aKey string) (float32, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsFloat32(aKey)
	}

	return 0, false
} // AsFloat32()

// AsFloat64 returns the value of `aKey` in `aSection` as a 64bit floating point.
//
// If the given `aKey` in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// If s is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsFloat64(aSection, aKey string) (float64, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsFloat64(aKey)
	}

	return 0, false
} // AsFloat64()

// AsInt returns the value of `aKey` in `aSection` as an integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsInt(aSection, aKey string) (int, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsInt(aKey)
	}

	return 0, false
} // AsInt()

// AsInt16 return the value of `aKey` in `aSection` as a 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsInt16(aSection, aKey string) (int16, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsInt16(aKey)
	}

	return 0, false
} // AsInt16()

// AsInt32 return the value of `aKey` in `aSection` as a 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsInt32(aSection, aKey string) (int32, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsInt32(aKey)
	}

	return 0, false
} // AsInt32()

// AsInt64 return the value of `aKey` in `aSection` as a 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsInt64(aSection, aKey string) (int64, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsInt64(aKey)
	}

	return 0, false
} // AsInt64()

// AsString returns the value of `aKey` in `aSection` as a string.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (bool) return value will be `false`.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key to lookup.
func (id *TSections) AsString(aSection, aKey string) (string, bool) {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.AsString(aKey)
	}

	return "", false
} // AsString()

// Clear empties the internal data structures.
//
// This method can be called once the program has used the config values stored
// in the INI file to setup the application. Emptying these data structures
// helps the garbage collector do release the data not needed anymore.
func (id *TSections) Clear() bool {
	// we leave defSect alone for now
	id.secOrder = make(tOrder, 0, defCapacity)
	for name := range id.sections {
		if cs, exists := id.sections[name]; exists {
			cs.Clear()
		}
		delete(id.sections, name)
	}
	id.sections = make(tIniSections)

	return ((0 == len(id.sections)) && (0 == len(id.secOrder)))
} // Clear()

// GetSection returns the INI section named `aSection`,
// or `nil` if not found.
func (id *TSections) GetSection(aSection string) *TSection {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if result, ok := id.sections[aSection]; ok {
		return result
	}

	return nil
} // GetSection()

// HasSection checks whether the INI data contain `aSection`.
func (id *TSections) HasSection(aSection string) bool {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	_, result := id.sections[aSection]

	return result
} // HasSection()

// HasSectionKey checks whether the INI data contain `aSection`
// with `aKey` returning whether it exists at all.
//
// `aSection` the INI section to lookup.
//
// `aKey` is the key name to lookup in `aSection`.
func (id *TSections) HasSectionKey(aSection, aKey string) bool {
	if 0 == len(aKey) {
		return false
	}
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, ok := id.sections[aSection]; ok {
		return cs.HasKey(aKey)
		// for _, kv := range cs {
		// 	if kv.Key == aKey {
		// 		return true
		// 	}
		// }
	}

	return false
} // HasSectionKey()

// Len returns the number of INI sections and
// whether the IniSections structure is consistent.
func (id *TSections) Len() int {
	return len(id.sections)
} // Len()

// Load reads the given `aFilename` returning the data structure
// read from the INI file and a possible error condition.
//
// This method reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
// `aFilename` is the name of the INI file to read.
func (id *TSections) Load(aFilename string) (*TSections, error) {
	file, rErr := os.Open(aFilename)
	if nil != rErr {
		return nil, rErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err := id.read(scanner)
	return id, err
} // Load()

// read parses the INI file
// returning the number of bytes read and a possible error.
//
// This method reads one line of the INI file at a time skipping both
// empty lines and comments (identified by '#' or ';' at line start).
func (id *TSections) read(aScanner *bufio.Scanner) (int, error) {
	section := id.defSect
	var result int

	for lineRead := aScanner.Scan(); lineRead; lineRead = aScanner.Scan() {
		line := aScanner.Text()
		result += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		if 0 == len(line) {
			// Skip blank lines
			continue
		}
		if ';' == line[0] || '#' == line[0] {
			// Skip comment lines
			continue
		}

		if matches := sectionRE.FindStringSubmatch(line); nil != matches {
			// update the current section name
			section = strings.TrimSpace(matches[1])
		} else if matches := keyValRE.FindStringSubmatch(line); nil != matches {
			// get a slice of RegEx matches,
			// we expect (1) key, (2) value:
			key := strings.TrimSpace(matches[1])
			val := trimRemoveQuotes(matches[2])

			id.AddSectionKey(section, key, val)
			// } else {
			// 	// ignore broken lines
			// 	continue
		}
	}

	return result, aScanner.Err()
} // read()

// RemoveSection deletes `aSection` from the list of INI sections.
//
// `aSection` the name of the INI section to remove.
func (id *TSections) RemoveSection(aSection string) bool {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	_, exists := id.sections[aSection]
	if !exists {
		// section doesn't exist
		return true
	}
	delete(id.sections, aSection)
	if 0 == len(id.sections) {
		exists = false
	} else if _, exists = id.sections[aSection]; exists {
		return false
	}

	// len - 1: because list is zero-based
	olen := len(id.secOrder) - 1
	if 0 > olen {
		// empty list
		return true
	}

	// remove secOrder entry:
	for idx, name := range id.secOrder {
		if name != aSection {
			continue
		}
		if 0 == idx {
			if 0 == olen {
				// the only list entry: replace by an empty list
				id.secOrder = make(tOrder, 0, defCapacity)
			} else {
				// first list entry: move the remaining data
				id.secOrder = id.secOrder[1:]
			}
		} else if idx == olen {
			// last list entry
			id.secOrder = id.secOrder[:idx]
		} else {
			id.secOrder = append(id.secOrder[:idx], id.secOrder[idx+1:]...)
		}
		return true
	}

	return false
} // RemoveSection()

// RemoveSectionKey removes aKey from aSection.
//
// This method returns 'true' if either `aSection` or `aKey` doesn't exist
// or if `aKey` in `aSection` was successfully removed, or `false` otherwise.
//
// `aSection` the name of the INI section to use.
//
// `aKey` the name of the key-value pair to remove.
func (id *TSections) RemoveSectionKey(aSection, aKey string) bool {
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	cs, exists := id.sections[aSection]
	if (!exists) || (0 == len(aKey)) {
		// section or key doesn't exist
		return true
	}

	return cs.RemoveKey(aKey)
} // RemoveSectionKey()

// Store writes all INI data to `aFilename`
// returning the number of bytes written and a possible error.
//
// `aFilename` is the name of the INI file to write.
func (id *TSections) Store(aFilename string) (int, error) {
	file, err := os.Create(aFilename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	s := id.String()
	return file.Write([]byte(s))
} // Store()

// String returns a string representation of an INI section list.
func (id *TSections) String() (rString string) {
	// use the secOrder list to determine the order of sections
	for _, name := range id.secOrder {
		if 0 == len(name) {
			name = id.defSect
		}
		if cs, exists := id.sections[name]; exists {
			rString += "\n[" + name + "]\n" + cs.String()
		}
	}

	return
} // String()

// string0 is the initial but slower implementation.
func (id *TSections) string0() string {
	// NOTE: this implementation is ~3 times slower than the one above.
	var result bytes.Buffer

	// use the secOrder list to determine the order of sections
	for _, name := range id.secOrder {
		if 0 == len(name) {
			name = id.defSect
		}
		if cs, exists := id.sections[name]; exists {
			result.WriteString(fmt.Sprintf("\n[%s]\n%s", name, cs.string0()))
		}
	}

	return result.String()
} // String0()

// updateSectKey replaces the current value of `aKey` in `aSection`
// by the provided new `aValue`.
//
// Private method used by the UpdateSectKeyXXX() methods.
func (id *TSections) updateSectKey(aSection, aKey, aValue string) bool {
	if 0 == len(aKey) {
		return false
	}
	if 0 == len(aSection) {
		aSection = id.defSect
	}
	if cs, exists := id.sections[aSection]; exists {
		return cs.UpdateKey(aKey, aValue)
	}

	// if aSection or aKey doesn't exist then create a new entry
	return id.AddSectionKey(aSection, aKey, aValue)
} // updateSectKey()

// UpdateSectKeyBool replaces the current value of `aKey` in `aSection` by
// the provided new `aValue` boolean.
//
// If the given `aValue` is 'true' the string "true" is used
// otherwise the string "false".
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key-value pair to use.
//
// `aValue` the boolean value of the key-value pair to update.
func (id *TSections) UpdateSectKeyBool(aSection, aKey string, aValue bool) bool {
	if aValue {
		return id.updateSectKey(aSection, aKey, "true")
	}

	return id.updateSectKey(aSection, aKey, "false")
} // UpdateSectKeyBool()

// UpdateSectKeyFloat replaces the current value of aKey in `aSection` by
// the provided new `aValue` float.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key-value pair to use.
//
// `aValue` the float64 value of the key-value pair to update.
func (id *TSections) UpdateSectKeyFloat(aSection, aKey string, aValue float64) bool {
	return id.updateSectKey(aSection, aKey, fmt.Sprintf("%f", aValue))
} // UpdateSectKeyFloat()

// UpdateSectKeyInt replaces the current value of `aKey` in `aSection` by
// the provided new `aValue` integer.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key-value pair to use.
//
// `aValue` the int64 value of the key-value pair to update.
func (id *TSections) UpdateSectKeyInt(aSection, aKey string, aValue int64) bool {
	return id.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyInt()

// UpdateSectKeyStr replaces the current value of `aKey` in `aSection` by
// the provided new `aValue` string.
//
// `aSection` the name of the INI section to lookup.
//
// `aKey` the name of the key-value pair to use.
//
// `aValue` the string value of the key-value pair to update.
func (id *TSections) UpdateSectKeyStr(aSection, aKey, aValue string) bool {
	return id.updateSectKey(aSection, aKey, strings.TrimSpace(aValue))
} // UpdateSectKeyStr()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TWalkFunc is used by `Walk()` when visiting an entry in the INI list.
	//
	// see `Walk()`
	TWalkFunc func(aSect, aKey, aVal string)

	// TIniWalker is used by `Walker()` when visiting an entry in the INI list.
	//
	// see `Walker()`
	TIniWalker interface {
		Walk(aSect, aKey, aVal string)
	}
)

// Walk traverses through all entries in the INI list sections
// calling `aFunc` for each entry.
//
// `aFunc` is the function called for each key-value pair in all INI sections.
func (id *TSections) Walk(aFunc TWalkFunc) {
	// we ignore the secOrder list because the
	// order of sections doesn't matter here.
	for name := range id.sections {
		if 0 == len(name) {
			name = id.defSect
		}
		cs, _ := id.sections[name]
		for _, kv := range *cs {
			aFunc(name, kv.Key, kv.Value)
		}
	}
} // Walk()

// Walker traverses through all entries in the INI list sections
// calling `aWalker` for each entry.
//
// `aWalker` is an object implementing the `TIniWalker` interface.
func (id *TSections) Walker(aWalker TIniWalker) {
	id.Walk(aWalker.Walk)
} // Walker()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// LoadFile reads the given `aFilename` returning the data structure
// read from the INI file and a possible error condition.
//
// This function reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
// `aFilename` is the name of the INI file to read.
func LoadFile(aFilename string) (*TSections, error) {
	result := NewSections()

	return result.Load(aFilename)
} // LoadFile()

// NewSections creates a new/empty `IniSections` structure.
func NewSections() *TSections {
	return &TSections{
		defSect:  DefSection,
		secOrder: make(tOrder, 0, defCapacity),
		sections: make(tIniSections),
	}
} // NewIniSections()

/* _EoF_ */
