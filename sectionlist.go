/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

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

	// `TSectionList` is a list of INI sections.
	//
	// This opaque data structure is filled by e.g. `load()`.
	//
	// For accessing the sections and key/value pairs it provides
	// the appropriate methods.
	TSectionList struct {
		defSect  string        // name of default section
		fName    string        // name of the INI file to use
		secOrder tSectionOrder // slice containing the order of sections
		sections tSections     // map of INI sections
	}

	// `TIniWalkFunc()` is used by `Walk()` when visiting an entry
	// in the INI list.
	//
	// see `Walk()`
	TIniWalkFunc func(aSection, aKey, aVal string)

	// A `TIniWalker` is used by `Walker()` when visiting an entry
	// in the INI list.
	//
	// see `Walker()`
	TIniWalker interface {
		Walk(aSection, aKey, aVal string)
	}
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
// Parameters:
// - `aString` The quoted string to process.
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
// Parameters:
// - `aSection` The name of the INI section to add.
//
// Returns:
// - bool: `true` if the section list was successfully updated,
// `false` otherwise.
func (sl *TSectionList) addSection(aSection string) (rOK bool) {
	if _, rOK = sl.sections[aSection]; rOK {
		return // already there: nothing more to do
	}

	sl.sections[aSection] = NewSection()
	if _, rOK = sl.sections[aSection]; rOK {
		// add new section name to order list
		sl.secOrder = append(sl.secOrder, aSection)

		// just to be safe:
		_, rOK = sl.sections[aSection]
	}

	return
} // addSection()

// `AddSectionKey()` appends a new key/value pair to `aSection`
// returning `true` on success or `false` otherwise.
//
// Parameters:
// - `aSection` The name of the INI section to use.
// - `aKey` The key of the key/value pair to add.
// - `aValue` The value of the key/value pair to add.
//
// Returns:
// - `bool`: `true` on success, of `false` if either `aKey` is empty, or
// `aSection` can't be found or added.
func (sl *TSectionList) AddSectionKey(aSection, aKey, aValue string) (rOK bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if rOK = sl.addSection(aSection); !rOK {
		return // can't find nor add the section
	}

	if kl, exists := sl.sections[aSection]; exists {
		rOK = kl.AddKey(aKey, aValue)
	}

	return
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
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `bool`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsBool(aSection, aKey string) (bool, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false, false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsBool(aKey)
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
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `float32`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsFloat32(aSection, aKey string) (float32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return float32(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsFloat32(aKey)
	}

	return float32(0), false
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
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `float64`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsFloat64(aSection, aKey string) (float64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return float64(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsFloat64(aKey)
	}

	return float64(0), false
} // AsFloat64()

// Int

// `AsInt()` returns the value of `aKey` in `aSection` as an integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsInt(aSection, aKey string) (int, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsInt(aKey)
	}

	return int(0), false
} // AsInt()

// `AsInt8()` returns the value of `aKey` in `aSection` as a 8bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int8`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsInt8(aSection, aKey string) (int8, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int8(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsInt8(aKey)
	}

	return int8(0), false
} // AsInt8()

// `AsInt16()` return the value of `aKey` in `aSection` as a 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int16`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsInt16(aSection, aKey string) (int16, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int16(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsInt16(aKey)
	}

	return int16(0), false
} // AsInt16()

// `AsInt32()` return the value of `aKey` in `aSection` as a 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int32`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsInt32(aSection, aKey string) (int32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int32(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsInt32(aKey)
	}

	return int32(0), false
} // AsInt32()

// `AsInt64()` return the value of `aKey` in `aSection` as a 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second return
// value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `int64`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsInt64(aSection, aKey string) (int64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return int64(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsInt64(aKey)
	}

	return int64(0), false
} // AsInt64()

//

// `AsString()` returns the value of `aKey` in `aSection` as a string.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `string`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsString(aSection, aKey string) (string, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return "", false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsString(aKey)
	}

	return "", false
} // AsString()

// Uint

// `AsUInt()` returns the value of `aKey` in `aSection` as an
// unsigned integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsUInt(aSection, aKey string) (uint, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsUInt(aKey)
	}

	return uint(0), false
} // AsUInt()

// `AsUInt8()` returns the value of `aKey` in `aSection` as an
// unsigned 8bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint8`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsUInt8(aSection, aKey string) (uint8, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint8(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsUInt8(aKey)
	}

	return uint8(0), false
} // AsUInt8()

// `AsUInt16()` return the value of `aKey` in `aSection` as an
// unsigned 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint16`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsUInt16(aSection, aKey string) (uint16, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint16(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsUInt16(aKey)
	}

	return uint16(0), false
} // AsUInt16()

// `AsUInt32()` return the value of `aKey` in `aSection` as an
// unsigned 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// return value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint32`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsUInt32(aSection, aKey string) (uint32, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint32(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsUInt32(aKey)
	}

	return uint32(0), false
} // AsUInt32()

// `AsUInt64()` return the value of `aKey` in `aSection` as an unsigned
// 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second return
// value will be `false`.
//
// Parameters:
// - `aSection` the name of the INI section to lookup.
// - `aKey` The name of the key to lookup.
//
// Returns:
// - `uint64`: The value associated with `aKey`.
// - `bool`: `true` if `aKey` was found, or false otherwise.
func (sl *TSectionList) AsUInt64(aSection, aKey string) (uint64, bool) {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return uint64(0), false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.AsUInt64(aKey)
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
// Returns:
// - `*TSectionList`: The return value is the cleared list.
func (sl *TSectionList) Clear() *TSectionList {
	// we leave `defSect` alone for now
	sl.secOrder = make(tSectionOrder, 0, slDefCapacity)
	for name := range sl.sections {
		if kl, exists := sl.sections[name]; exists {
			kl.Clear()
		}
		delete(sl.sections, name)
	}
	sl.sections = make(tSections)

	return sl
} // Clear()

// `CompareTo()` compares the current `TSectionList` with another
// `TSectionList`.
// It checks whether both lists have the same number of sections and
// whether each section in the current list has the same keys and values
// as in the other list.
//
// Parameters:
// - `aINI`: The `TSectionList` to compare with the current one.
//
// Returns:
// - `bool`: `true` if both lists are equal, `false` otherwise.
func (sl *TSectionList) CompareTo(aINI *TSectionList) bool {
	// Check if both lists have the same number of sections
	if len(sl.sections) != len(aINI.sections) {
		return false
	}

	// Iterate over each section in the current list
	for name, kl := range sl.sections {
		// Check if the other list has the same section
		section, exists := aINI.sections[name]
		if !exists {
			return false
		}
		// Compare the keys and values of the sections
		if !kl.CompareTo(section) {
			return false
		}
	}

	// If all checks passed, the lists are equal
	return true
} // CompareTo()

// `Filename()` returns the configured filename of the INI file.
func (sl *TSectionList) Filename() string {
	return sl.fName
} // Filename()

// `GetSection()` returns the INI section named `aSection`, or an empty list
// if not found.
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
//
// Returns:
// - `*TSection`: The requested section or an empty if not found.
func (sl *TSectionList) GetSection(aSection string) *TSection {
	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if result, ok := sl.sections[aSection]; ok {
		return result
	}

	return &TSection{}
} // GetSection()

// `HasSection()` checks whether the INI data contain `aSection`.
//
// Parameters:
// - `aSection` is the name of the INI section to lookup.
//
// Returns:
// - `bool`: `true` if `aSection` is found, or `false` otherwise.
func (sl *TSectionList) HasSection(aSection string) (rOK bool) {
	if aSection = strings.TrimSpace(aSection); "" == aSection {
		_, rOK = sl.sections[sl.defSect]
	} else {
		_, rOK = sl.sections[aSection]
	}

	return
} // HasSection()

// `HasSectionKey()` checks whether the INI data contain `aSection` with
// `aKey` returning whether it exists at all.
//
// Parameters:
// - `aSection` The INI section to lookup.
// - `aKey` The key name to lookup in `aSection`.
//
// Returns:
// - `bool`: `true` if `aKey` exists, or `false` otherwise.
func (sl *TSectionList) HasSectionKey(aSection, aKey string) bool {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, ok := sl.sections[aSection]; ok {
		return kl.HasKey(aKey)
	}

	return false
} // HasSectionKey()

// `Len()` returns the number of INI sections.
//
// It is used to determine the size of the list of sections.
//
// Returns:
// - `int`: The number of sections in the INI file.
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
// - `*TSectionList`: The loaded INI list.
// - `error`: A possible error condition.
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
// - `aSection`: The name of the INI section to lookup.
// - `aKey` The name of the key/value pair to use.
// - `aValue`: The value of the key/value pair to update.
func (sl *TSectionList) mergeWalker(aSection, aKey, aValue string) {
	sl.AddSectionKey(aSection, aKey, aValue) // ignore the return value
} // mergeWalker()

// `Merge()` copies or merges all INI sections with all key/value pairs
// into this list.
//
// Parameters:
// - `aINI` The INI sections to merge with this list.
//
// Returns:
// - `TSectionList` This sections list merged with the other one.
func (sl *TSectionList) Merge(aINI *TSectionList) *TSectionList {
	if nil != aINI {
		aINI.Walk(sl.mergeWalker)
	}

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
// - `aScanner`: A bufio.Scanner instance that reads from the INI file.
//
// Returns:
// - `int`: The number of bytes read from the INI file.
// - `error`: A possible error condition.
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

// `RemoveSection()` deletes `aSection` from the list of sections.
//
// Parameters:
// - `aSection` The name of the INI section to remove.
//
// Returns:
// - `bool`: `true` on success, `false` on failure.
func (sl *TSectionList) RemoveSection(aSection string) bool {
	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}
	if _, exists := sl.sections[aSection]; !exists {
		// section doesn't exist which satisfies the removal request
		return true
	}

	delete(sl.sections, aSection)
	if _, exists := sl.sections[aSection]; exists {
		return false // this should never happen!
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
// - `aSection` is the name of the INI section to use.
// - `aKey` The name of the key/value pair to remove.
//
// Returns:
// - `bool`: `true` on success, `false` on failure.
func (sl *TSectionList) RemoveSectionKey(aSection, aKey string) bool {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return true
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.RemoveKey(aKey)
	}

	// section or key doesn't exist: assume successful removal
	return true
} // RemoveSectionKey()

// `Sections()` returns a list of section names in the order they
// appear in the INI file.
//
// The returned list is a slice of strings. The length of the slice
// is the number of sections in the INI file.
//
// Returns:
// - `[]string`: A list of section names
// - `int`: The number of sections in the returned list.
func (sl *TSectionList) Sections() ([]string, int) {
	dest := make([]string, len(sl.secOrder))
	len := copy(dest, sl.secOrder)

	return dest, len
} // Sections()

// `SetFilename()` sets the filename of the INI file to use.
//
// Parameters:
// - `aFilename` The name to use for the INI file.
func (sl *TSectionList) SetFilename(aFilename string) *TSectionList {
	sl.fName = strings.TrimSpace(aFilename)

	return sl
} // SetFilename()

// `Sort()` sorts the sections in the order they appear in the INI file.
//
// This method sorts the key/value pairs in each section.
//
// Returns:
// - `*TSectionList`: The sorted instance of the `TSectionList`.
func (sl *TSectionList) Sort() *TSectionList {
	// use the secOrder list to determine the order of sections
	for _, name := range sl.secOrder {
		if kl, exists := sl.sections[name]; exists {
			sl.sections[name] = kl.Sort()
		}
	}

	return sl
} // Sort()

// `Store()` writes all INI data to the configured filename.
//
// Returns:
// - `int`: The number of bytes written.
// - `error`: An possible error during writing the data to file.
func (sl *TSectionList) Store() (int, error) {
	file, err := os.Create(sl.fName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(sl.String()))
} // Store()

// `String()` returns a string representation of the INI section list.
//
// Returns:
// - `string`: The string representation of the INI section list.
func (sl *TSectionList) String() (rString string) {
	// use the secOrder list to determine the order of sections
	for _, name := range sl.secOrder {
		if kl, exists := sl.sections[name]; exists {
			// ensure that all sections are sorted internally
			sl.sections[name] = kl.Sort()

			rString += "\n[" + name + "]\n" + kl.String()
		}
	}

	return
} // String()

// `updateSectKey()` updates the current value of `aKey` in `aSection`
// by the provided new `aValue`.
//
// Private method used by the UpdateSectKeyXXX() methods.
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
// - `aKey`     The name of the key/value pair to use.
// - `aValue`   The value of the key/value pair to update.
//
// Returns:
// - bool: `true` if the key/value pair was successfully updated,
// `false` otherwise.
func (sl *TSectionList) updateSectKey(aSection, aKey, aValue string) bool {
	if aKey = strings.TrimSpace(aKey); "" == aKey {
		return false
	}

	if aSection = strings.TrimSpace(aSection); "" == aSection {
		aSection = sl.defSect
	}

	if kl, exists := sl.sections[aSection]; exists {
		return kl.UpdateKey(aKey, aValue)
	}

	// if `aSection` doesn't exist we create a new entry
	return sl.AddSectionKey(aSection, aKey, aValue)
} // updateSectKey()

// `UpdateSectKeyBool()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` boolean.
//
// If the given `aValue` is `true` then the string "True" is used
// otherwise the string "False".
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
// - `aKey` The name of the key/value pair to use.
// - `aValue` The boolean value of the key/value pair to update.
//
// Returns:
// - bool: `true` if the key/value pair was successfully updated,
// or `false` otherwise.
func (sl *TSectionList) UpdateSectKeyBool(aSection, aKey string, aValue bool) bool {
	if aValue {
		return sl.updateSectKey(aSection, aKey, `True`)
	}

	return sl.updateSectKey(aSection, aKey, `False`)
} // UpdateSectKeyBool()

// `UpdateSectKeyFloat()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` float.
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
// - `aKey` The name of the key/value pair to use.
// - `aValue` The float64 value of the key/value pair to update.
//
// Returns:
// - bool: `true` if the key/value pair was successfully updated,
// or `false` otherwise.
func (sl *TSectionList) UpdateSectKeyFloat(aSection, aKey string, aValue float64) bool {
	return sl.updateSectKey(aSection, aKey, fmt.Sprintf("%f", aValue))
} // UpdateSectKeyFloat()

// `UpdateSectKeyInt()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` integer.
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
// - `aKey` The name of the key/value pair to use.
// - `aValue` The int64 value of the key/value pair to update.
//
// Returns:
// - bool: `true` if the key/value pair was successfully updated,
// or `false` otherwise.
func (sl *TSectionList) UpdateSectKeyInt(aSection, aKey string, aValue int64) bool {
	return sl.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyInt()

// `UpdateSectKeyUInt()` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` unsigned integer.
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
// - `aKey` The name of the key/value pair to use.
// - `aValue` The int64 value of the key/value pair to update.
//
// Returns:
// - bool: `true` if the key/value pair was successfully updated,
// or `false` otherwise.
func (sl *TSectionList) UpdateSectKeyUInt(aSection, aKey string, aValue uint64) bool {
	return sl.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyUInt()

// `UpdateSectKeyStr` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` string.
//
// Parameters:
// - `aSection` The name of the INI section to lookup.
// - `aKey` The name of the key/value pair to use.
// - `aValue` The string value of the key/value pair to update.
//
// Returns:
// - bool: `true` if the key/value pair was successfully updated,
// or `false` otherwise.
func (sl *TSectionList) UpdateSectKeyStr(aSection, aKey, aValue string) bool {
	return sl.updateSectKey(aSection, aKey, aValue)
} // UpdateSectKeyStr()

// ----------------------------------------------------------------

// `Walk()` traverses through all entries in the INI list sections calling
// `aFunc` for each entry.
//
// Parameters:
// - `aFunc` The function called for each key/value pair in all sections.
func (sl *TSectionList) Walk(aFunc TIniWalkFunc) {
	// We ignore the `secOrder` list because the
	// order of sections doesn't matter here.
	for name, kl := range sl.sections {
		for _, kv := range kl.data {
			aFunc(name, kv.Key, kv.Value)
		}
	}

} // Walk()

// `Walker()` traverses through all entries in all INI sections
// calling `aWalker` for each entry.
//
// Parameters:
// - `aWalker` An object implementing the `TIniWalker` interface.
func (sl *TSectionList) Walker(aWalker TIniWalker) {
	sl.Walk(aWalker.Walk)
} // Walker()

// ----------------------------------------------------------------

// `NewSectionList()` creates a new instance of the `TSectionList`.
//
// This method initializes a new `TSectionList` instance with the default section name.
//
// Returns:
// - *TSectionList: A new instance of the `TSectionList`.
func NewSectionList() *TSectionList {
	return &TSectionList{
		defSect:  DefSection,
		secOrder: make(tSectionOrder, 0, slDefCapacity),
		sections: make(tSections),
	}
} // NewSectionList()

/* _EoF_ */
