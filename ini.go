// Package ini provides functions to read/write INI files from/to disc
// and methods to access the section's key/value pairs.
package ini

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type (
	// TKeyVal represents an INI key/value pair.
	TKeyVal struct {
		Key   string
		Value string
	}

	// TSection is a slice of key/value pairs.
	TSection []TKeyVal

	// `tIniSections` is a list (map) of INI sections.
	tIniSections map[string]*TSection

	// A helper slice of strings (i.e. section names)
	// used to preserve the order of INI sections.
	tOrder = []string

	// This opaque data structure is filled by e.g. `LoadFile(…)`.
	tSections struct {
		defSect  string       // name of default section
		fname    string       // name of the INI file to use
		secOrder tOrder       // slice containing the order of sections
		sections tIniSections // list of INI sections
	}
)

// TIniList is a list of INI sections.
//
// This opaque data structure is filled by e.g. `LoadFile(…)`.
//
// For accessing the sections and key/value pairs
// it provides the appropriate methods.
type TIniList tSections

const (
	ilDefCapacity = 16

	// DefSection is the name of the default section in the INI
	// file which is used when there are key/value pairs in the
	// file without a preceding section header like `[SectName]`.
	DefSection = "Default"
)

// Regular expressions to identify certain parts of an INI file.
var (
	// match: [section]
	ilSectionRE = regexp.MustCompile(`^\[\s*([^\]]*?)\s*]$`)
	// match: key = val
	ilKeyValRE = regexp.MustCompile(`^([^=]+?)\s*=\s*(.*)$`)
	// match: quoted ' " string " '
	ilQuotesRE = regexp.MustCompile(`^(['"])(.*)(['"])$`)
)

// `removeQuotes()` returns a quoted string w/o the quote characters.
func removeQuotes(aString string) (rString string) {
	// remove leading/trailing UTF whitespace:
	rString = strings.TrimSpace(aString)

	// get a slice of RegEx matches:
	matches := ilQuotesRE.FindStringSubmatch(rString)
	// we expect: (1) leading quote, (2) text, (3) trailing quote
	if (3 < len(matches)) && (matches[1] == matches[3]) {
		return matches[2]
	}

	return
} // removeQuotes()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// String returns a string representation of the key/value pair.
func (kv *TKeyVal) String() string {
	if 0 == len(kv.Value) {
		return kv.Key + " ="
	}

	return kv.Key + " = " + kv.Value
} // String()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// AddKey appends a new key/value pair returning
// `true` on success or `false` otherwise.
//
// If `aKey` is an empty string the method's result
// will be `false`.
//
//	`aKey` the key of the key/value pair to add.
//
//	`aValue` the value of the key/value pair to add.
func (cs *TSection) AddKey(aKey, aValue string) bool {
	if 0 < len(aKey) {
		idx := cs.IndexOf(aKey)
		if 0 > idx {
			*cs = append(*cs, TKeyVal{aKey, aValue})
		} else {
			(*cs)[idx].Value = aValue
		}

		if val, ok := cs.AsString(aKey); ok {
			return (val == aValue)
		}
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
// All other values will give `false` as the second (`rOK`) result value.
//
// This method actually checks only the first character of the key's
// value so one can write e.g. "false" or "NO" (for a `false` result),
// or "True" or "yes" (for a 'true' result).
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsBool(aKey string) (rVal, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		val = val[:1]
		switch val {
		case "0", "f", "F", "n", "N":
			rVal, rOK = false, true
		case "1", "t", "T", "y", "Y":
			rVal, rOK = true, true
		}
	}

	return
} // AsBool()

// AsFloat32 returns the value of `aKey` as a 32bit floating point.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsFloat32(aKey string) (rVal float32, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		if f64, err := strconv.ParseFloat(val, 32); nil == err {
			return float32(f64), true
		}
	}

	return
} // AsFloat32()

// AsFloat64 returns the value of `aKey` as a 64bit floating point.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
// 	aKey` the name of the key to lookup.
func (cs *TSection) AsFloat64(aKey string) (rVal float64, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		if f64, err := strconv.ParseFloat(val, 64); nil == err {
			rVal, rOK = f64, true
		}
	}

	return
} // AsFloat64()

// AsInt returns the value of `aKey` as an integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsInt(aKey string) (rVal int, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i, err := strconv.Atoi(val); nil == err {
			rVal, rOK = i, true
		}
	}

	return
} // AsInt()

// AsInt16 returns the value of `aKey` as a 16bit integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsInt16(aKey string) (rVal int16, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i64, err := strconv.ParseInt(val, 10, 16); nil == err {
			rVal, rOK = int16(i64), true
		}
	}

	return
} // AsInt16()

// AsInt32 returns the value of `aKey` as a 32bit integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsInt32(aKey string) (rVal int32, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i64, err := strconv.ParseInt(val, 10, 32); nil == err {
			rVal, rOK = int32(i64), true
		}
	}

	return
} // AsInt32()

// AsInt64 returns the value of `aKey` as a 64bit integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsInt64(aKey string) (rVal int64, rOK bool) {
	if val, exists := cs.AsString(aKey); exists {
		if i64, err := strconv.ParseInt(val, 10, 64); nil == err {
			rVal, rOK = i64, true
		}
	}

	return
} // AsInt64()

// AsString returns the value of `aKey` as a string.
//
// If the given `aKey` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aKey` the name of the key to lookup.
func (cs *TSection) AsString(aKey string) (rVal string, rOK bool) {
	for _, val := range *cs {
		if val.Key == aKey {
			rVal, rOK = val.Value, true
		}
	}

	return
} // AsString()

// Clear removes all entries in this INI section.
func (cs *TSection) Clear() *TSection {
	(*cs) = (*cs)[:0]

	return cs
} // Clear()

// HasKey returns whether `aKey` exists in this INI section.
//
//	`aKey` the key to lookup.
func (cs *TSection) HasKey(aKey string) bool {
	return (0 <= cs.IndexOf(aKey))
} // HasKey()

// IndexOf returns the index of `aKey` in this INI section
// or `-1` if not found.
//
//	`aKey` the key to lookup.
func (cs *TSection) IndexOf(aKey string) int {
	for result, kv := range *cs {
		if kv.Key == aKey {
			return result
		}
	}

	return -1
} // IndexOf()

// Len returns the number of key/value pairs in this section.
func (cs *TSection) Len() int {
	return len(*cs)
} // Len()

// RemoveKey removes `aKey` from this section.
//
// This method returns 'true' if `aKey` doesn't exist at all,
// or if `aKey` was successfully removed, or `false` otherwise.
//
//	`aKey` the name of the key/value pair to remove.
func (cs *TSection) RemoveKey(aKey string) bool {
	idx := cs.IndexOf(aKey)
	if 0 > idx {
		return true
	}
	slen := len(*cs) - 1 // new slice length (i.e. one shorter)
	(*cs)[idx] = TKeyVal{}
	switch idx {
	case 0:
		(*cs) = (*cs)[1:]
	case slen:
		(*cs) = (*cs)[:slen]
	default:
		(*cs) = append((*cs)[:idx], (*cs)[1+idx:]...)
	}

	return (0 > cs.IndexOf(aKey))
} // RemoveKey()

// String returns a string representation of an INI section.
//
// The single key/value pairs are delimited by a linefeed ('\n).
func (cs *TSection) String() (rString string) {
	for _, kv := range *cs {
		rString += kv.String() + "\n"
	}

	return
} // String()

// UpdateKey replaces the current value of `aKey`
// by the provided new `aValue`.
//
// In case `aKey` doesn't already exist in the list
// (and therefor can't be updated) it will be added
// by calling the `AddKey()` method.
//
// If `aKey` is an empty string the method's result
// will be `false`.
//
//	`aKey` the key of the key/value pair to update.
//
//	`aValue` the value of the key/value pair to update.
func (cs *TSection) UpdateKey(aKey, aValue string) bool {
	if 0 == len(aKey) {
		return false
	}
	for idx, val := range *cs {
		if val.Key == aKey {
			(*cs)[idx] = TKeyVal{Key: aKey, Value: aValue}
			return true
		}
	}

	// if aKey doesn't exist then create a new entry
	return cs.AddKey(aKey, aValue)
} // updateKey()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// addSection appends a new INI section returning `true` on success
// or `false` otherwise.
//
//	`aSection` name of the INI section to add.
func (il *TIniList) addSection(aSection string) bool {
	if _, exists := il.sections[aSection]; exists {
		return true // already there: nothing more to do
	}

	// we make room for initially 8 key/value pairs
	sect := make(TSection, 0, ilDefCapacity)
	il.sections[aSection] = &sect
	if _, ok := il.sections[aSection]; ok {
		// add new section name to order list
		il.secOrder = append(il.secOrder, aSection)

		return true
	}

	return false
} // addSection()

// AddSectionKey appends a new key/value pair to `aSection`
// returning `true` on success or `false` otherwise.
//
//	`aSection` name of the INI section to use.
//
//	`aKey` the key of the key/value pair to add.
//
//	`aValue` the value of the key/value pair to add.
func (il *TIniList) AddSectionKey(aSection, aKey, aValue string) bool {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if !il.addSection(aSection) {
		return false // can't find nor add the section
	}

	if cs, exists := il.sections[aSection]; exists {
		return cs.AddKey(aKey, aValue)
	}

	return false
} // AddSectionKey()

/*
 * Public methods to return INI values from a section as a certain data type.
 */

// AsBool returns the value of `aKey` in `aSection` as a boolean value.
//
// If the given aKey in `aSection` doesn't exist then the second
// (bool) return value will be `false`.
//
// "0", "f", "F", "n" and "N" are considered `false` while
// "1", "t", "T", "y" and "Y" are considered 'true';
// these values will be given in the first result value.
// All other values will give `false` as the second result value.
//
// This method actually checks only the first character of the key's
// value so one can write e.g. "false" or "NO" (for a `false`
// result), or "True" or "yes" (for a 'true' result).
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsBool(aSection, aKey string) (rVal, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsBool(aKey)
	}

	return
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
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsFloat32(aSection, aKey string) (rVal float32, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsFloat32(aKey)
	}

	return
} // AsFloat32()

// AsFloat64 returns the value of `aKey` in `aSection` as a 64bit
// floating point.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
// If s is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded
// using IEEE754 unbiased rounding.
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsFloat64(aSection, aKey string) (rVal float64, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsFloat64(aKey)
	}

	return
} // AsFloat64()

// AsInt returns the value of `aKey` in `aSection` as an integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsInt(aSection, aKey string) (rVal int, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt(aKey)
	}

	return
} // AsInt()

// AsInt16 return the value of `aKey` in `aSection` as a 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsInt16(aSection, aKey string) (rVal int16, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt16(aKey)
	}

	return
} // AsInt16()

// AsInt32 return the value of `aKey` in `aSection` as a 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsInt32(aSection, aKey string) (rVal int32, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt32(aKey)
	}

	return
} // AsInt32()

// AsInt64 return the value of `aKey` in `aSection` as a 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsInt64(aSection, aKey string) (rVal int64, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt64(aKey)
	}

	return
} // AsInt64()

// AsString returns the value of `aKey` in `aSection` as a string.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//
//	`aKey` the name of the key to lookup.
func (il *TIniList) AsString(aSection, aKey string) (rVal string, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsString(aKey)
	}

	return
} // AsString()

// Clear empties the internal data structures.
//
// This method can be called once the program has used the config values
// stored in the INI file to setup the application. Emptying these data
// structures should help the garbage collector do release the data not
// needed anymore.
func (il *TIniList) Clear() bool {
	// we leave defSect alone for now
	il.secOrder = make(tOrder, 0, ilDefCapacity)
	for name := range il.sections {
		if cs, exists := il.sections[name]; exists {
			cs.Clear()
		}
		delete(il.sections, name)
	}
	il.sections = make(tIniSections)

	return ((0 == len(il.sections)) && (0 == len(il.secOrder)))
} // Clear()

// Filename returns the configured filename of the INI file.
func (il *TIniList) Filename() string {
	return il.fname
} // Filename()

// GetSection returns the INI section named `aSection`,
// or an empty list if not found.
//
//	`aSection` is the name of the INI section to lookup.
func (il *TIniList) GetSection(aSection string) *TSection {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if result, ok := il.sections[aSection]; ok {
		return result
	}

	return &TSection{}
} // GetSection()

// HasSection checks whether the INI data contain `aSection`.
//
//	`aSection` is the name of the INI section to lookup.
func (il *TIniList) HasSection(aSection string) bool {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	_, result := il.sections[aSection]

	return result
} // HasSection()

// HasSectionKey checks whether the INI data contain `aSection`
// with `aKey` returning whether it exists at all.
//
//	`aSection` is the INI section to lookup.
//
//	`aKey` is the key name to lookup in `aSection`.
func (il *TIniList) HasSectionKey(aSection, aKey string) bool {
	if 0 == len(aKey) {
		return false
	}
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, ok := il.sections[aSection]; ok {
		return cs.HasKey(aKey)
	}

	return false
} // HasSectionKey()

// Len returns the number of INI sections.
func (il *TIniList) Len() int {
	return len(il.sections)
} // Len()

// Load reads the configured filename returning the data structure
// read from the INI file and a possible error condition.
//
// This method reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
func (il *TIniList) Load() (*TIniList, error) {
	file, rErr := os.Open(il.fname)
	if nil != rErr {
		return il, rErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err := il.read(scanner)
	return il, err
} // Load()

// Merge copies or merges all INI sections with all key/value pairs
// into this list.
//
//	`aINI` the INI list to merge with this one.
func (il *TIniList) Merge(aINI *TIniList) *TIniList {
	aINI.Walk(il.mergeWalker)

	return il
} // Merge()

// `mergeWalker()` inserts the given key/value pair in `aSection`.
func (il *TIniList) mergeWalker(aSection, aKey, aValue string) {
	il.AddSectionKey(aSection, aKey, aValue)
} // mergeWalker()

// `read()` parses the INI file returning the number of bytes read
// and a possible error.
//
// This method reads one line of the INI file at a time skipping both
// empty lines and comments (identified by '#' or ';' at line start).
func (il *TIniList) read(aScanner *bufio.Scanner) (rRead int, rErr error) {
	section := il.defSect
	var lastLine string

	for lineRead := aScanner.Scan(); lineRead; lineRead = aScanner.Scan() {
		line := aScanner.Text()
		rRead += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		lLen := len(line)
		if 0 == lLen {
			if 0 == len(lastLine) {
				// Skip blank lines
				continue
			}
			line = lastLine
			lastLine = ""
		}
		if ';' == line[0] || '#' == line[0] {
			if 0 == len(lastLine) {
				// Skip comment lines
				continue
			}
			line = lastLine
			lastLine = ""
		}
		if '\\' == line[lLen-1] {
			if (1 < lLen) && (' ' == line[lLen-2]) {
				lastLine += line[:lLen-1]
			} else {
				lastLine += line[:lLen-1] + " "
			}
			line = ""
			continue
		}
		if 0 < len(lastLine) {
			line = lastLine + line
			lastLine = ""
		}

		if matches := ilSectionRE.FindStringSubmatch(line); nil != matches {
			// update the current section name
			section = strings.TrimSpace(matches[1])
		} else if matches := ilKeyValRE.FindStringSubmatch(line); nil != matches {
			// get a slice of RegEx matches,
			// we expect (1) key, (2) value
			key := strings.TrimSpace(matches[1])
			val := removeQuotes(matches[2])

			il.AddSectionKey(section, key, val)
		} else {
			// ignore broken lines
			line = ""
		}
	}
	rErr = aScanner.Err()

	return
} // read()

// RemoveSection deletes `aSection` from the list of INI sections.
//
//	`aSection` is the name of the INI section to remove.
func (il *TIniList) RemoveSection(aSection string) bool {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	_, exists := il.sections[aSection]
	if !exists {
		// section doesn't exist
		return true
	}
	delete(il.sections, aSection)
	if 0 < len(il.sections) {
		if _, exists = il.sections[aSection]; exists {
			return false // this should never happen!
		}
	}

	// len - 1: because list is zero-based
	olen := len(il.secOrder) - 1
	if 0 > olen {
		// empty list
		return true
	}

	// remove secOrder entry:
	for idx, name := range il.secOrder {
		if name != aSection {
			continue
		}
		switch idx {
		case 0:
			if 0 == olen {
				// the only list entry: replace by an empty list
				il.secOrder = make(tOrder, 0, ilDefCapacity)
			} else {
				// first list entry: move the remaining data
				il.secOrder = il.secOrder[1:]
			}
		case olen:
			// last list entry
			il.secOrder = il.secOrder[:idx]
		default:
			il.secOrder = append(il.secOrder[:idx], il.secOrder[idx+1:]...)
		}
		return true
	}

	return false
} // RemoveSection()

// RemoveSectionKey removes aKey from aSection.
//
// This method returns 'true' if either `aSection` or `aKey`
// doesn't exist or if `aKey` in `aSection` was successfully
// removed, or `false` otherwise.
//
//	`aSection` is the name of the INI section to use.
//
//	`aKey` is the name of the key/value pair to remove.
func (il *TIniList) RemoveSectionKey(aSection, aKey string) bool {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	cs, exists := il.sections[aSection]
	if (!exists) || (0 == len(aKey)) {
		// section or key doesn't exist
		return true
	}

	return cs.RemoveKey(aKey)
} // RemoveSectionKey()

// SetFilename set the filename of the INI file to use.
//
//	`aFilename` is the name to use for the INI file.
func (il *TIniList) SetFilename(aFilename string) *TIniList {
	il.fname = aFilename

	return il
} // SetFilename()

// Store writes all INI data to the configured filename
// returning the number of bytes written and a possible error.
func (il *TIniList) Store() (int, error) {
	file, err := os.Create(il.fname)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(il.String()))
} // Store()

// String returns a string representation of an INI section list.
func (il *TIniList) String() (rString string) {
	// use the secOrder list to determine the order of sections
	for _, name := range il.secOrder {
		if 0 == len(name) {
			name = il.defSect
		}
		if cs, exists := il.sections[name]; exists {
			rString += "\n[" + name + "]\n" + cs.String()
		}
	}

	return
} // String()

// updateSectKey replaces the current value of `aKey` in `aSection`
// by the provided new `aValue`.
//
// Private method used by the UpdateSectKeyXXX() methods.
func (il *TIniList) updateSectKey(aSection, aKey, aValue string) bool {
	if 0 == len(aKey) {
		return false
	}
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.UpdateKey(aKey, aValue)
	}

	// if aSection or aKey doesn't exist then create a new entry
	return il.AddSectionKey(aSection, aKey, aValue)
} // updateSectKey()

// UpdateSectKeyBool replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` boolean.
//
// If the given `aValue` is 'true' the string "true" is used
// otherwise the string "false".
//
//	`aSection` is the name of the INI section to lookup.
//
//	`aKey` is the name of the key/value pair to use.
//
//	`aValue` is the boolean value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyBool(aSection, aKey string, aValue bool) bool {
	if aValue {
		return il.updateSectKey(aSection, aKey, "true")
	}

	return il.updateSectKey(aSection, aKey, "false")
} // UpdateSectKeyBool()

// UpdateSectKeyFloat replaces the current value of aKey in `aSection`
// by the provided new `aValue` float.
//
//	`aSection` is the name of the INI section to lookup.
//
//	`aKey` is the name of the key/value pair to use.
//
//	`aValue` is the float64 value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyFloat(aSection, aKey string, aValue float64) bool {
	return il.updateSectKey(aSection, aKey, fmt.Sprintf("%f", aValue))
} // UpdateSectKeyFloat()

// UpdateSectKeyInt replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` integer.
//
//	`aSection` is the name of the INI section to lookup.
//
//	`aKey` is the name of the key/value pair to use.
//
//	`aValue` is the int64 value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyInt(aSection, aKey string, aValue int64) bool {
	return il.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyInt()

// UpdateSectKeyStr replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` string.
//
//	`aSection` is the name of the INI section to lookup.
//
//	`aKey` is the name of the key/value pair to use.
//
//	`aValue` is the string value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyStr(aSection, aKey, aValue string) bool {
	return il.updateSectKey(aSection, aKey, strings.TrimSpace(aValue))
} // UpdateSectKeyStr()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TWalkFunc is used by `Walk()` when visiting an entry
	// in the INI list.
	//
	// see `Walk()`
	TWalkFunc func(aSection, aKey, aVal string)

	// TIniWalker is used by `Walker()` when visiting an entry
	// in the INI list.
	//
	// see `Walker()`
	TIniWalker interface {
		Walk(aSection, aKey, aVal string)
	}
)

// Walk traverses through all entries in the INI list sections
// calling `aFunc` for each entry.
//
//	`aFunc` is the function called for each key/value pair in all sections.
func (il *TIniList) Walk(aFunc TWalkFunc) {
	// we ignore the secOrder list because the
	// order of sections doesn't matter here.
	for section := range il.sections {
		if 0 == len(section) {
			section = il.defSect
		}
		cs := il.sections[section]
		for _, kv := range *cs {
			aFunc(section, kv.Key, kv.Value)
		}
	}
} // Walk()

// Walker traverses through all entries in the INI list sections
// calling `aWalker` for each entry.
//
//	`aWalker` is an object implementing the `TIniWalker` interface.
func (il *TIniList) Walker(aWalker TIniWalker) {
	il.Walk(aWalker.Walk)
} // Walker()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// New reads the given `aFilename` returning the data structure
// read from that INI file and a possible error condition.
//
// This function reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
//	`aFilename` is the name of the INI file to read.
func New(aFilename string) (*TIniList, error) {
	result := &TIniList{
		defSect:  DefSection,
		fname:    aFilename,
		secOrder: make(tOrder, 0, ilDefCapacity),
		sections: make(tIniSections),
	}
	return result.Load()
} // New()

/* _EoF_ */
