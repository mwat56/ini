/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"os"
	"path/filepath"
)

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
		secOrder: make(tSectionOrder, 0, slDefCapacity),
		sections: make(tSections),
	}
	return result.Load()
} // New()

// `ReadIniData()` returns the config values read from INI file(s).
//
//	The steps here are:
//	(1) read the local `./.aName.ini`,
//	(2) read the global `/etc/aName.ini`,
//	(3) read the user-local `~/.aName.ini`,
//	(4) read the user-local `~/.config/aName.ini`,
//	(5) read the `-ini` commandline argument.
//
// This function considers only the `Default` section of the INI files.
//
// Example:
//
//	iniData, _ := ReadIniData("myApp")
//	fmt.Println(iniData.AsString("myKey"))
//
// The function returns a pointer to the 'Default' section
// of the first INI file that contains it.
//
//	`aName` The application's name used as the INI file name (without extension).
func ReadIniData(aName string) *TSection {
	var (
		confDir    string
		err        error
		ini1, ini2 *TSectionList
	)
	// (1) ./
	fName, _ := filepath.Abs(`./` + aName + `.ini`)
	if ini1, err = New(fName); nil == err {
		ini1.AddSectionKey(ini1.defSect, `iniFile`, fName)
	}

	// (2) /etc/
	fName = `/etc/` + aName + `.ini`
	if ini2, err = New(fName); nil == err {
		ini1.Merge(ini2)
		ini1.AddSectionKey(ini1.defSect, `iniFile`, fName)
	}

	// (3) ~user/
	fName, err = os.UserHomeDir()
	if (nil == err) && (0 < len(fName)) {
		fName, _ = filepath.Abs(filepath.Join(fName, `.`+aName+`.ini`))
		if ini2, err = New(fName); nil == err {
			ini1.Merge(ini2)
			ini1.AddSectionKey(ini1.defSect, `iniFile`, fName)
		}
	}

	// (4) ~/.config/
	if confDir, err = os.UserConfigDir(); nil == err {
		fName, _ = filepath.Abs(filepath.Join(confDir, aName+`.ini`))
		if ini2, err = New(fName); nil == err {
			ini1.Merge(ini2)
			ini1.AddSectionKey(ini1.defSect, `iniFile`, fName)
		}
	}

	// (5) cmdline
	aLen := len(os.Args)
	for i := 1; i < aLen; i++ {
		if `-ini` == os.Args[i] {
			//XXX Note that this works only if `-ini` and
			// filename are two separate arguments. It will
			// fail if it's given in the form `-ini=filename`.
			i++
			if i < aLen {
				fName, _ = filepath.Abs(os.Args[i])
				if ini2, err = New(fName); nil == err {
					ini1.Merge(ini2)
					ini1.AddSectionKey(ini1.defSect, `iniFile`, fName)
				}
			}
			break
		}
	}

	return ini1.GetSection(ini1.defSect)
} // ReadIniData()

/* _EoF_ */
