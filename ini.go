/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

//lint:file-ignore ST1017 - I prefer Yoda conditions

// `New()` reads the given `aFilename` returning the data structure read
// from that INI file and a possible error condition.
//
// This function reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
//	`aFilename` The name of the INI file to read.
func New(aFilename string) (*TSectionList, error) {
	result := &TSectionList{
		defSect:  DefSection,
		fName:    aFilename,
		secOrder: make(TSectionOrder, 0, slDefCapacity),
		sections: make(tSections),
	}
	return result.Load()
} // New()

/* _EoF_ */
