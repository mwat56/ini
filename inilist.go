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

type (
	// `tSectionList` is a list (map) of INI sections.
	tSectionList map[string]*TSection

	// A helper slice of strings (i.e. section names)
	// used to preserve the order of INI sections.
	tSectionOrder = []string

	// This opaque data structure is filled by e.g. `LoadFile(…)`.
	tIniSectionsList struct {
		defSect  string        // name of default section
		fName    string        // name of the INI file to use
		secOrder tSectionOrder // slice containing the order of sections
		sections tSectionList  // list of INI sections
	}
)

// TIniList is a list of INI sections.
//
// This opaque data structure is filled by e.g. `LoadFile(…)`.
//
// For accessing the sections and key/value pairs it provides
// the appropriate methods.
type TIniList tIniSectionsList

const (
	// DefSection is the name of the default section in the INI file
	// which is used when there are key/value pairs in the file
	// without a preceding section header like `[SectionName]`.
	DefSection = `Default`

	// Default list capacity.
	ilDefCapacity = 16
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
//
//	`aString` The quoted string to process.
func removeQuotes(aString string) (rString string) {
	// remove leading/trailing UTF whitespace:
	rString = strings.TrimSpace(aString)

	// get a slice of RegEx matches:
	matches := ilQuotesRE.FindStringSubmatch(rString)
	// we expect: (1) leading quote, (2) text, (3) trailing quote
	if (3 < len(matches)) && (matches[1] == matches[3]) {
		rString = matches[2]
	}

	return
} // removeQuotes()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `addSection` appends a new INI section returning `true` on success or
// `false` otherwise.
//
//	`aSection` The name of the INI section to add.
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

// `AddSectionKey` appends a new key/value pair to `aSection`
// returning `true` on success or `false` otherwise.
//
//	`aSection` The name of the INI section to use.
//	`aKey` The key of the key/value pair to add.
//	`aValue` The value of the key/value pair to add.
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

// `AsBool` returns the value of `aKey` in `aSection` as a boolean value.
//
// If the given aKey in `aSection` doesn't exist then the second (bool)
// return value will be `false`.
//
// `0`, `f`, `F`, `n`, and `N` are considered `false` while
// `1`, `t`, `T`, `y`, and `Y` are considered `true`;
// these values will be given in the first result value.
// All other values will give `false` as the second result value.
//
// This method actually checks only the first character of the key's value
// so one can write e.g. "false" or "NO" (for a `false` result), or "True" or
// "yes" (for a `true` result).
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (il *TIniList) AsBool(aSection, aKey string) (rVal, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsBool(aKey)
	}

	return
} // AsBool()

// `AsFloat32` returns the value of `aKey` in `aSection` as a 32bit floating
// point.
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
func (il *TIniList) AsFloat32(aSection, aKey string) (rVal float32, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
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
func (il *TIniList) AsFloat64(aSection, aKey string) (rVal float64, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsFloat64(aKey)
	}

	return
} // AsFloat64()

// `AsInt` returns the value of `aKey` in `aSection` as an integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (il *TIniList) AsInt(aSection, aKey string) (rVal int, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt(aKey)
	}

	return
} // AsInt()

// `AsInt16` return the value of `aKey` in `aSection` as a 16bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second
// (`rOK`) return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (il *TIniList) AsInt16(aSection, aKey string) (rVal int16, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt16(aKey)
	}

	return
} // AsInt16()

// `AsInt32` return the value of `aKey` in `aSection` as a 32bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (il *TIniList) AsInt32(aSection, aKey string) (rVal int32, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt32(aKey)
	}

	return
} // AsInt32()

// `AsInt64` return the value of `aKey` in `aSection` as a 64bit integer.
//
// If the given `aKey` in `aSection` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (il *TIniList) AsInt64(aSection, aKey string) (rVal int64, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsInt64(aKey)
	}

	return
} // AsInt64()

// `AsString` returns the value of `aKey` in `aSection` as a string.
//
// If the given `aKey` in `aSection` doesn't exist then the second (`rOK`)
// return value will be `false`.
//
//	`aSection` the name of the INI section to lookup.
//	`aKey` The name of the key to lookup.
func (il *TIniList) AsString(aSection, aKey string) (rVal string, rOK bool) {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if cs, exists := il.sections[aSection]; exists {
		return cs.AsString(aKey)
	}

	return
} // AsString()

// `Clear` empties the internal data structures.
//
// This method can be called once the program has used the config values
// stored in the INI file to setup the application. Emptying these data
// structures should help the garbage collector to release the data not
// needed anymore.
//
// The return value is the cleared list.
func (il *TIniList) Clear() *TIniList {
	// we leave `defSect` alone for now
	il.secOrder = make(tSectionOrder, 0, ilDefCapacity)
	for name := range il.sections {
		if cs, exists := il.sections[name]; exists {
			cs.Clear()
		}
		delete(il.sections, name)
	}
	il.sections = make(tSectionList)

	return il
} // Clear()

// `Filename` returns the configured filename of the INI file.
func (il *TIniList) Filename() string {
	return il.fName
} // Filename()

// `GetSection` returns the INI section named `aSection`, or an empty list
// if not found.
//
//	`aSection` The name of the INI section to lookup.
func (il *TIniList) GetSection(aSection string) *TSection {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	if result, ok := il.sections[aSection]; ok {
		return result
	}

	return &TSection{}
} // GetSection()

// `HasSection` checks whether the INI data contain `aSection`.
//
//	`aSection` is the name of the INI section to lookup.
func (il *TIniList) HasSection(aSection string) (rOK bool) {
	if 0 == len(aSection) {
		_, rOK = il.sections[il.defSect]
	} else {
		_, rOK = il.sections[aSection]
	}
	return
} // HasSection()

// `HasSectionKey` checks whether the INI data contain `aSection` with
// `aKey` returning whether it exists at all.
//
//	`aSection` The INI section to lookup.
//	`aKey` The key name to lookup in `aSection`.
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

// `Len` returns the number of INI sections.
func (il *TIniList) Len() int {
	return len(il.sections)
} // Len()

// `Load` reads the configured filename returning the data structure
// read from the INI file and a possible error condition.
//
// This method reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
func (il *TIniList) Load() (*TIniList, error) {
	file, rErr := os.Open(il.fName)
	if nil != rErr {
		return il, rErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err := il.read(scanner)

	return il, err
} // Load()

// `Merge` copies or merges all INI sections with all key/value pairs
// into this list.
//
//	`aINI` The INI list to merge with this one.
func (il *TIniList) Merge(aINI *TIniList) *TIniList {
	aINI.Walk(il.mergeWalker)

	return il
} // Merge()

// `mergeWalker` inserts the given key/value pair in `aSection`.
func (il *TIniList) mergeWalker(aSection, aKey, aValue string) {
	il.AddSectionKey(aSection, aKey, aValue)
} // mergeWalker()

// `read` parses the INI file returning the number of bytes read
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
		lineLen := len(line)
		if 0 == lineLen {
			if 0 == len(lastLine) {
				continue // Skip blank lines
			}
			line, lastLine = lastLine, ``
		}
		if ';' == line[0] || '#' == line[0] {
			if 0 == len(lastLine) {
				continue // Skip comment lines
			}
			line, lastLine = lastLine, ``
		}
		if '\\' == line[lineLen-1] {
			if (1 < lineLen) && (' ' == line[lineLen-2]) {
				lastLine += line[:lineLen-1]
			} else {
				lastLine += line[:lineLen-1] + ` `
			}
			line = ``
			continue // concatenation handled
		}
		if 0 < len(lastLine) {
			line, lastLine = lastLine+line, ``
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
			line = `` // ignore broken lines
		}
	}
	rErr = aScanner.Err()

	return
} // read()

// `RemoveSection` deletes `aSection` from the list of INI sections.
//
//	`aSection` The name of the INI section to remove.
func (il *TIniList) RemoveSection(aSection string) bool {
	if 0 == len(aSection) {
		aSection = il.defSect
	}
	_, exists := il.sections[aSection]
	if !exists {
		// section doesn't exist which satisfies the removal wish
		return true
	}

	delete(il.sections, aSection)
	if 0 < len(il.sections) {
		if _, exists = il.sections[aSection]; exists {
			return false // this should never happen!
		}
	}

	// len - 1: because list is zero-based
	oLen := len(il.secOrder) - 1
	if 0 > oLen {
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
			if 0 == oLen {
				// the only list entry: replace by an empty list
				il.secOrder = make(tSectionOrder, 0, ilDefCapacity)
			} else {
				// first list entry: move the remaining data
				il.secOrder = il.secOrder[1:]
			}

		case oLen:
			// last list entry
			il.secOrder = il.secOrder[:idx]

		default:
			il.secOrder = append(il.secOrder[:idx], il.secOrder[idx+1:]...)
		}

		return true
	}

	return false
} // RemoveSection()

// `RemoveSectionKey` removes aKey from aSection.
//
// This method returns 'true' if either `aSection` or `aKey` doesn't exist
// or if `aKey` in `aSection` was successfully removed, or `false` otherwise.
//
//	`aSection` is the name of the INI section to use.
//	`aKey` The name of the key/value pair to remove.
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

// `SetFilename` sets the filename of the INI file to use.
//
//	`aFilename` The name to use for the INI file.
func (il *TIniList) SetFilename(aFilename string) *TIniList {
	il.fName = strings.TrimSpace(aFilename)

	return il
} // SetFilename()

// `Store` writes all INI data to the configured filename returning
// the number of bytes written and a possible error.
func (il *TIniList) Store() (int, error) {
	file, err := os.Create(il.fName)
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

// `updateSectKey` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue`.
//
// Private method used by the UpdateSectKeyXXX() methods.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The value of the key/value pair to update.
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

	// if `aSection` or `aKey` doesn't exist then create a new entry
	return il.AddSectionKey(aSection, aKey, aValue)
} // updateSectKey()

// `UpdateSectKeyBool` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` boolean.
//
// If the given `aValue` is `true` then the string "true" is used
// otherwise the string "false".
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The boolean value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyBool(aSection, aKey string, aValue bool) bool {
	if aValue {
		return il.updateSectKey(aSection, aKey, `true`)
	}

	return il.updateSectKey(aSection, aKey, `false`)
} // UpdateSectKeyBool()

// `UpdateSectKeyFloat` replaces the current value of aKey in `aSection`
// by the provided new `aValue` float.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The float64 value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyFloat(aSection, aKey string, aValue float64) bool {
	return il.updateSectKey(aSection, aKey, fmt.Sprintf("%f", aValue))
} // UpdateSectKeyFloat()

// `UpdateSectKeyInt` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` integer.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The int64 value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyInt(aSection, aKey string, aValue int64) bool {
	return il.updateSectKey(aSection, aKey, fmt.Sprintf("%d", aValue))
} // UpdateSectKeyInt()

// `UpdateSectKeyStr` replaces the current value of `aKey` in `aSection`
// by the provided new `aValue` string.
//
//	`aSection` The name of the INI section to lookup.
//	`aKey` The name of the key/value pair to use.
//	`aValue` The string value of the key/value pair to update.
func (il *TIniList) UpdateSectKeyStr(aSection, aKey, aValue string) bool {
	return il.updateSectKey(aSection, aKey, strings.TrimSpace(aValue))
} // UpdateSectKeyStr()

/* _EoF_ */
