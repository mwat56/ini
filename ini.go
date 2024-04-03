/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

//lint:file-ignore ST1017 - I prefer Yoda conditions

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

// `Walk` traverses through all entries in the INI list sections calling
// `aFunc` for each entry.
//
//	`aFunc` The function called for each key/value pair in all sections.
func (il *TIniList) Walk(aFunc TWalkFunc) {
	// We ignore the `secOrder` list because the
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

// `Walker` traverses through all entries in the INI list sections
// calling `aWalker` for each entry.
//
//	`aWalker` An object implementing the `TIniWalker` interface.
func (il *TIniList) Walker(aWalker TIniWalker) {
	il.Walk(aWalker.Walk)
} // Walker()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `New` reads the given `aFilename` returning the data structure read
// from that INI file and a possible error condition.
//
// This function reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
//	`aFilename` The name of the INI file to read.
func New(aFilename string) (*TIniList, error) {
	result := &TIniList{
		defSect:  DefSection,
		fName:    aFilename,
		secOrder: make(tSectionOrder, 0, ilDefCapacity),
		sections: make(tSectionList),
	}
	return result.Load()
} // New()

/* _EoF_ */
