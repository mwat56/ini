/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `tSections` is a list (map) of INI sections.
	tSections map[string]*TSection

	// A helper slice of strings (i.e. section names)
	// used to preserve the order of INI sections.
	tSectionOrder = []string

	// This opaque data structure is filled by e.g. `load()`.
	tIniSectionsList struct {
		defSect  string        // name of default section
		fName    string        // name of the INI file to use
		secOrder tSectionOrder // slice containing the order of sections
		sections tSections     // list of INI sections
	}

	// `TWalkFunc()` is used by `Walk()` when visiting an entry
	// in the INI list.
	//
	// see `Walk()`
	TWalkFunc func(aSection, aKey, aVal string)

	// A `TIniWalker` is used by `Walker()` when visiting an entry
	// in the INI list.
	//
	// see `Walker()`
	TIniWalker interface {
		Walk(aSection, aKey, aVal string)
	}

	// `TSectionList` is a list of INI sections.
	//
	// This opaque data structure is filled by e.g. `load()`.
	//
	// For accessing the sections and key/value pairs it provides
	// the appropriate methods.
	TSectionList tIniSectionsList
)

const (
	// `DefSection` is the name of the default section in the INI file
	// which is used when there are key/value pairs in the file
	// without a preceding section header like `[SectionName]`.
	DefSection = `Default`

	// Default list capacity.
	slDefCapacity = 16
)

// Regular expressions to identify certain parts of an INI file.
var (
	// match: [section]
	isSectionRE = regexp.MustCompile(`^\[\s*([^\]]*?)\s*]$`)

	// match: key = val
	isKeyValRE = regexp.MustCompile(`^([^=]+?)\s*=\s*(.*)$`)

	// match: quoted ' " string " '
	isQuotesRE = regexp.MustCompile(`^\s*(['"])\s*(.*?)\s*(['"])\s*$`)
)

// `removeQuotes()` returns a quoted string w/o the quote characters.
//
//	`aString` The quoted string to process.
func removeQuotes(aString string) (rString string) {
	// remove leading/trailing UTF whitespace:
	rString = strings.TrimSpace(aString)

	// get a slice of RegEx matches:
	matches := isQuotesRE.FindStringSubmatch(rString)
	// we expect: (1) leading quote, (2) text, (3) trailing quote
	if (3 < len(matches)) && (matches[1] == matches[3]) {
		rString = matches[2]
	}

	return
} // removeQuotes()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `addSection()` appends a new INI section returning `true` on success or
// `false` otherwise.
//
//	`aSection` The name of the INI section to add.
func (sl *TSectionList) addSection(aSection string) bool {
	if _, exists := sl.sections[aSection]; exists {
		return true // already there: nothing more to do
	}

	// we make room for initially 16 key/value pairs
	sect := make(TSection, 0, slDefCapacity)
	sl.sections[aSection] = &sect
	if _, ok := sl.sections[aSection]; ok {
		// add new section name to order list
		sl.secOrder = append(sl.secOrder, aSection)

		return true
	}

	return false
} // addSection()

// `AddSectionKey()` appends a new key/value pair to `aSection`
// returning `true` on success or `false` otherwise.
//
//	`aSection` The name of the INI section to use.
//	`aKey` The key of the key/value pair to add.
//	`aValue` The value of the key/value pair to add.
func (sl *TSectionList) AddSectionKey(aSection, aKey, aValue string) bool {
	if "" == aSection {
		aSection = sl.defSect
	}
	if ok := sl.addSection(aSection); !ok {
		return false // can't find nor add the section
	}

	if cs, exists := sl.sections[aSection]; exists {
		return cs.AddKey(aKey, aValue)
	}

	return false
} // AddSectionKey()

/*
 * Public methods to return INI values from a section as a certain data type.
 */

// `AsBool()` returns the value of `aKey` in `aSection` as a boolean value.
//
// If the given aKey in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// `0`, `f`, `F`, `n`, and `N` are considered `false` while
// `1`, `t`, `T`, `y`, `Y`, `j`, `J`, `o`, `O` are considered `true`;
// these values will be given in the first result value.
// All other values will give `false` as the second result value.
//
// This method actually checks only the first character of the key's value
// so one can write e.g. "false" or "NO" (for a `false` result), or "True" or
// "yes" (for a `true` result).
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsBool(aSection, aKey string) (bool, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsBool(aKey)
	}

	return false, false
} // AsBool()

// Float

// `AsFloat32` returns the value of `aKey` in `aSection` as a 32bit
// floating point.
//
// If the given `aKey` in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// If the key's value is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded using IEEE754
// unbiased rounding.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsFloat32(aSection, aKey string) (rVal float32, rOK bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsFloat32(aKey)
	}

	return
} // AsFloat32()

// `AsFloat64` returns the value of `aKey` in `aSection` as a 64bit
// floating point.
//
// If the given `aKey` in `aSection` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
// If the key's value is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded using IEEE754
// unbiased rounding.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsFloat64(aSection, aKey string) (rVal float64, rOK bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsFloat64(aKey)
	}

	return
} // AsFloat64()

// Int

// `AsInt()` returns the value of `aKey` in `aSection` as an integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsInt(aSection, aKey string) (int, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsInt(aKey)
	}

	return int(0), false
} // AsInt()

// `AsInt8()` returns the value of `aKey` in `aSection` as a 8bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsInt8(aSection, aKey string) (int8, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsInt8(aKey)
	}

	return int8(0), false
} // AsInt8()

// `AsInt16()` return the value of `aKey` in `aSection` as a 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsInt16(aSection, aKey string) (int16, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsInt16(aKey)
	}

	return int16(0), false
} // AsInt16()

// `AsInt32()` return the value of `aKey` in `aSection` as a 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsInt32(aSection, aKey string) (int32, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsInt32(aKey)
	}

	return int32(0), false
} // AsInt32()

// `AsInt64()` return the value of `aKey` in `aSection` as a 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second return
// value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsInt64(aSection, aKey string) (int64, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsInt64(aKey)
	}

	return int64(0), false
} // AsInt64()

//

// `AsString()` returns the value of `aKey` in `aSection` as a string.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsString(aSection, aKey string) (string, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsString(aKey)
	}

	return "", false
} // AsString()

// Uint

// `AsUInt()` returns the value of `aKey` in `aSection` as an integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsUInt(aSection, aKey string) (uint, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsUInt(aKey)
	}

	return uint(0), false
} // AsUInt()

// `AsUInt8()` returns the value of `aKey` in `aSection` as a 8bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsUInt8(aSection, aKey string) (uint8, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsUInt8(aKey)
	}

	return uint8(0), false
} // AsUInt8()

// `AsUInt16()` return the value of `aKey` in `aSection` as a 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsUInt16(aSection, aKey string) (uint16, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsUInt16(aKey)
	}

	return uint16(0), false
} // AsUInt16()

// `AsUInt32()` return the value of `aKey` in `aSection` as a 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AUInt32(aSection, aKey string) (uint32, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsUInt32(aKey)
	}

	return uint32(0), false
} // AsUInt32()

// `AsUInt64()` return the value of `aKey` in `aSection` as an unsigned
// 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second return
// value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (sl *TSectionList) AsUInt64(aSection, aKey string) (uint64, bool) {
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.AsUInt64(aKey)
	}

	return uint64(0), false
} // AsUInt64()

//

// `Clear()` empties the internal data structures.
//
// This method can be called once the program has used the config values
// stored in the INI file to setup the application. Emptying these data
// structures should help the garbage collector to release the data not
// needed anymore.
//
// The return value is the cleared list.
func (sl *TSectionList) Clear() *TSectionList {
	// we leave `defSect` alone for now
	sl.secOrder = make(tSectionOrder, 0, slDefCapacity)
	for name := range sl.sections {
		if cs, exists := sl.sections[name]; exists {
			cs.Clear()
		}
		delete(sl.sections, name)
	}
	sl.sections = make(tSections)

	return sl
} // Clear()

// `Filename()` returns the configured filename of the INI file.
func (sl *TSectionList) Filename() string {
	return sl.fName
} // Filename()

// `GetSection()` returns the INI section named `aSection`, or an empty list
// if not found.
//
//	`aSection` The name of the INI section to lookup.
func (sl *TSectionList) GetSection(aSection string) *TSection {
	if "" == aSection {
		aSection = sl.defSect
	}
	if result, ok := sl.sections[aSection]; ok {
		return result
	}

	return &TSection{}
} // GetSection()

// `HasSection()` checks whether the INI data contain `aSection`.
//
//	`aSection` is the name of the INI section to lookup.
func (sl *TSectionList) HasSection(aSection string) (rOK bool) {
	if "" == aSection {
		_, rOK = sl.sections[sl.defSect]
	} else {
		_, rOK = sl.sections[aSection]
	}

	return
} // HasSection()

// `HasSectionKey()` checks whether the INI data contain `aSection` with
// `aKey` returning whether it exists at all.
//
//	`aSection` The INI section to lookup.
//	`aKey` The key name to lookup in `aSection`.
func (sl *TSectionList) HasSectionKey(aSection, aKey string) bool {
	if "" == aKey {
		return false
	}
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, ok := sl.sections[aSection]; ok {
		return cs.HasKey(aKey)
	}

	return false
} // HasSectionKey()

// `Len()` returns the number of INI sections.
//
// It is used to determine the size of the list of sections.
//
// Returns:
//
//	rOK: The number of sections in the INI file.
func (sl *TSectionList) Len() int {
	return len(sl.sections)
} // Len()

// `load()` reads the configured filename returning the data structure
// read from the INI file and a possible error condition.
//
// This method reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
// Returns:
//
//	*TSectionList: The loaded INI list.
//	error: A possible error condition.
func (sl *TSectionList) load() (*TSectionList, error) {
	file, rErr := os.Open(sl.fName)
	if nil != rErr {
		return sl, rErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err := sl.read(scanner)

	return sl, err
} // load()

// `mergeWalker()` inserts the given key/value pair in `aSection`.
//
// This method is called by the `Merge()` method.
//
// Parameters:
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The value of the key/value pair to update.
func (sl *TSectionList) mergeWalker(aSection, aKey, aValue string) {
	sl.AddSectionKey(aSection, aKey, aValue)
} // mergeWalker()

// `Merge()` copies or merges all INI sections with all key/value pairs
// into this list.
//
// Parameters:
//
//	`aINI` The INI sections to merge with this list.
//
// Returns:
//
//	`aINI` The INI list to merge with this one.
func (sl *TSectionList) Merge(aINI *TSectionList) *TSectionList {
	aINI.Walk(sl.mergeWalker)

	return sl
} // Merge()

// `read()` reads/parses the INI file data returning the number of bytes
// read and a possible error.
//
// This method reads one line of the INI file at a time skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
// The method updates the current section name and adds new key/value
// pairs to the list of sections.
//
// This method is called by the `load()` method.
//
// Parameters:
//
//	aScanner: A bufio.Scanner instance that reads from the INI file.
//
// Returns:
//
//	rRead: The number of bytes read from the INI file.
//	rErr:  A possible error condition.
func (sl *TSectionList) read(aScanner *bufio.Scanner) (rRead int, rErr error) {
	var lastLine string
	section := sl.defSect

	for lineRead := aScanner.Scan(); lineRead; lineRead = aScanner.Scan() {
		line := aScanner.Text()
		rRead += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		lineLen := len(line)
		if 0 == lineLen {
			if "" == lastLine {
				continue // Skip blank lines
			}
			line, lastLine = lastLine, ""
		}
		if ';' == line[0] || '#' == line[0] { // comment indicators
			if "" == lastLine {
				continue // Skip comment lines
			}
			line, lastLine = lastLine, ""
		}
		if '\\' == line[lineLen-1] { // possible value concatenation
			if (1 < lineLen) && (' ' == line[lineLen-2]) {
				lastLine += line[:lineLen-1]
			} else {
				lastLine += line[:lineLen-1] + " "
			}
			line = ``
			continue // concatenation handled
		}
		if 0 < len(lastLine) {
			line, lastLine = lastLine+line, ""
		}

		if matches := isSectionRE.FindStringSubmatch(line); nil != matches {
			// update the current section name
			section = strings.TrimSpace(matches[1])
		} else if matches := isKeyValRE.FindStringSubmatch(line); nil != matches {
			// get a slice of RegEx matches,
			// we expect (1) key, (2) value
			key := strings.TrimSpace(matches[1])
			val := removeQuotes(matches[2])

			sl.AddSectionKey(section, key, val) // ignore return value
		} else {
			line = "" // ignore broken lines
		}
	}
	rErr = aScanner.Err()

	return
} // read()

// `RemoveSection()` deletes `aSection` from the list of INI sections.
//
// Parameters:
//
//	`aSection` The name of the INI section to remove.
//
// Returns:
//
//	bool: `true` on success, `false` on failure.
func (sl *TSectionList) RemoveSection(aSection string) bool {
	if "" == aSection {
		aSection = sl.defSect
	}
	_, exists := sl.sections[aSection]
	if !exists {
		// section doesn't exist which satisfies the removal wish
		return true
	}

	delete(sl.sections, aSection)
	if 0 < len(sl.sections) {
		if _, exists = sl.sections[aSection]; exists {
			return false // this should never happen!
		}
	}

	// len - 1: because list is zero-based
	oLen := len(sl.secOrder) - 1
	if 0 > oLen {
		// empty list
		return true
	}

	// remove secOrder entry:
	for idx, name := range sl.secOrder {
		if name != aSection {
			continue
		}
		switch idx {
		case 0:
			if 0 == oLen {
				// the only list entry: replace by an empty list
				sl.secOrder = make(tSectionOrder, 0, slDefCapacity)
			} else {
				// first list entry: move the remaining data
				sl.secOrder = sl.secOrder[1:]
			}

		case oLen:
			// last list entry
			sl.secOrder = sl.secOrder[:idx]

		default:
			sl.secOrder = append(sl.secOrder[:idx], sl.secOrder[idx+1:]...)
		}

		return true
	}

	return false
} // RemoveSection()

// `RemoveSectionKey()` removes aKey from aSection.
//
// This method returns 'true' if either `aSection` or `aKey` doesn't exist
// or if `aKey` in `aSection` was successfully removed, or `false` otherwise.
//
// Parameters:
//
//	`aSection` is the name of the INI section to use.
//	`aKey` The name of the key/value pair to remove.
//
// Returns:
//
//	bool: `true` on success, `false` on failure.
func (sl *TSectionList) RemoveSectionKey(aSection, aKey string) bool {
	if "" == aSection {
		aSection = sl.defSect
	}
	cs, exists := sl.sections[aSection]
	if (!exists) || ("" == aKey) {
		// section or key doesn't exist
		return true
	}

	return cs.RemoveKey(aKey)
} // RemoveSectionKey()

// `Sections()` returns a list of section names in the order they
// appear in the INI file.
//
// The returned list is a slice of strings. The length of the slice
// is the number of sections in the INI file.
//
// Returns:
//
//	The number of sections in the returned list.
func (sl *TSectionList) Sections() ([]string, int) {
	dst := make([]string, len(sl.secOrder))
	len := copy(dst, sl.secOrder)

	return dst, len
} // Sections()

// `SetFilename()` sets the filename of the INI file to use.
//
// Parameters:
//
//	`aFilename` The name to use for the INI file.
func (sl *TSectionList) SetFilename(aFilename string) *TSectionList {
	sl.fName = strings.TrimSpace(aFilename)

	return sl
} // SetFilename()

// `Store()` writes all INI data to the configured filename returning
// the number of bytes written and a possible error.
func (sl *TSectionList) Store() (int, error) {
	file, err := os.Create(sl.fName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(sl.String()))
} // Store()

// String() returns a string representation of the INI section list.
func (sl *TSectionList) String() (rString string) {
	// use the secOrder list to determine the order of sections
	for _, name := range sl.secOrder {
		if "" == name {
			name = sl.defSect
		}
		if cs, exists := sl.sections[name]; exists {
			rString += "\n[" + name + "]\n" + cs.String()
		}
	}

	return
} // String()

// // String() returns a string representation of the INI section list.
// func (sl *TSectionList) String2() (rString string) {
// 	var sb strings.Builder
// 	// use the secOrder list to determine the order of sections
// 	for _, name := range sl.secOrder {
// 		if "" == name) {
// 			name = sl.defSect
// 		}
// 		if cs, exists := sl.sections[name]; exists {
// 			// ignore return values
// 			_, _ = sb.WriteString(fmt.Sprintf("\n[%s]\n%s",
// 				name, cs.String2()))
// 		}
// 	}

// 	return sb.String()
// } // String2()

// `updateSectKey()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue`.
//
// Private method used by the UpdateSectKeyXXX() methods.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The value of the key/value pair to update.
func (sl *TSectionList) updateSectKey(aSection, aKey, aValue string) bool {
	if "" == aKey {
		return false
	}
	if "" == aSection {
		aSection = sl.defSect
	}
	if cs, exists := sl.sections[aSection]; exists {
		return cs.UpdateKey(aKey, aValue)
	}

	// if `aSection` or `aKey` doesn't exist we create a new entry
	return sl.AddSectionKey(aSection, aKey, aValue)
} // updateSectKey()

// `UpdateSectKeyBool()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` boolean.
//
// If the given `aValue` is `true` then the string "true" is used
// otherwise the string "false".
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The boolean value of the key/value pair to update.
func (sl *TSectionList) UpdateSectKeyBool(aSection, aKey string, aValue bool) bool {
	if aValue {
		return sl.updateSectKey(aSection, aKey, `true`)
	}

	return sl.updateSectKey(aSection, aKey, `false`)
} // UpdateSectKeyBool()

// `UpdateSectKeyFloat()` replaces the current value of aKey in `aSection`
// by the provided new `aValue` float.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The float64 value of the key/value pair to update.
func (sl *TSectionList) UpdateSectKeyFloat(aSection, aKey string, aValue float64) bool {
	return sl.updateSectKey(aSection, aKey, fmt.Sprintf("%f", aValue))
} // UpdateSectKeyFloat()

// `UpdateSectKeyInt()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` integer.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The int64 value of the key/value pair to update.
func (sl *TSectionList) UpdateSectKeyInt(aSection, aKey string, aValue int64) bool {
	return sl.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyInt()

// `UpdateSectKeyUInt()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` unsigned integer.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The int64 value of the key/value pair to update.
func (sl *TSectionList) UpdateSectKeyUInt(aSection, aKey string, aValue uint64) bool {
	return sl.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyUInt()

// `UpdateSectKeyStr` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` string.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The string value of the key/value pair to update.
func (sl *TSectionList) UpdateSectKeyStr(aSection, aKey, aValue string) bool {
	return sl.updateSectKey(aSection, aKey, strings.TrimSpace(aValue))
} // UpdateSectKeyStr()

// `Walk()` traverses through all entries in the INI list sections calling
// `aFunc` for each entry.
//
//	`aFunc` The function called for each key/value pair in all sections.
func (sl *TSectionList) Walk(aFunc TWalkFunc) {
	// We ignore the `secOrder` list because the
	// order of sections doesn't matter here.
	for section := range sl.sections {
		if "" == section {
			section = sl.defSect
		}
		cs := sl.sections[section]
		for _, kv := range *cs {
			aFunc(section, kv.Key, kv.Value)
		}
	}
} // Walk()

// `Walker()` traverses through all entries in the INI list sections
// calling `aWalker` for each entry.
//
//	`aWalker` An object implementing the `TIniWalker` interface.
func (sl *TSectionList) Walker(aWalker TIniWalker) {
	sl.Walk(aWalker.Walk)
} // Walker()

/* _EoF_ */
