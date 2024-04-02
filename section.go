/*
Copyright © 2019, 204  M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package ini

import (
	"os"
	"path/filepath"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// TSection is a slice of key/value pairs.
	TSection []TKeyVal
)

// `AddKey` appends a new key/value pair returning `true` on success or
// `false` otherwise.
//
// If `aKey` is an empty string the method's result will be the result
// of `ks.RemoveKey(aKey)` i.e. usually `true`.
// If `aKey` already exist its value will be updated by `aValue`.
// If `aKey` doesn't exist in the section's list a new key/value pair
// will be appended.
//
//	`aKey` The key of the key/value pair to add.
//	`aValue` The value of the key/value pair to add.
func (ks *TSection) AddKey(aKey, aValue string) bool {
	if 0 < len(aKey) {
		idx := ks.IndexOf(aKey)
		if 0 > idx {
			*ks = append(*ks, TKeyVal{aKey, aValue})
		} else {
			// key already exists: update
			(*ks)[idx].Value = aValue
		}

		if val, ok := ks.AsString(aKey); ok {
			return (val == aValue)
		}
	} else {
		return ks.RemoveKey(aKey)
	}

	return false
} // AddKey()

// `AsBool` returns the value of `aKey` as a boolean value.
//
// If the given `aKey` doesn't exist then the second (bool) return value
// will be `false`.
//
// `0`, `f`, `F`, `n`, and `N` are considered `false` while
// `1`, `t`, `T`, `y`, and `Y` are considered `true`;
// these values will be given in the first result value.
// All other values will give `false` as the second (`rOK`) result value.
//
// This method actually checks only the first character of the key's value
// so one can write e.g. "false" or "NO" (for a `false` result), or "True"
// or "yes" (for a `true` result).
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsBool(aKey string) (rVal, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal = kv.AsBool()
	}

	return
} // AsBool()

// * returns the value of `aKey` as a 32bit floating point.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat32` returns the nearest floating point number rounded using
// IEEE754 unbiased rounding.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsFloat32(aKey string) (rVal float32, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsFloat32()
	}

	return
} // AsFloat32()

// `AsFloat64` returns the value of `aKey` as a 64bit floating point.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return
// value will be `false`.
//
// If the string is well-formed and near a valid floating point number,
// `AsFloat64` returns the nearest floating point number rounded using
// IEEE754 unbiased rounding.
//
//	aKey` the name of the key to lookup.
func (ks *TSection) AsFloat64(aKey string) (rVal float64, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsFloat64()
	}

	return
} // AsFloat64()

// `AsInt` returns the value of `aKey` as an integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return value
// will be `false`.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsInt(aKey string) (rVal int, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsInt()
	}

	return
} // AsInt()

// `AsInt16` returns the value of `aKey` as a 16bit integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return value
// will be `false`.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsInt16(aKey string) (rVal int16, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsInt16()
	}

	return
} // AsInt16()

// `AsInt32` returns the value of `aKey` as a 32bit integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return value
// will be `false`.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsInt32(aKey string) (rVal int32, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsInt32()
	}

	return
} // AsInt32()

// `AsInt64` returns the value of `aKey` as a 64bit integer.
//
// If the given `aKey` doesn't exist then the second (`rOK`) return value
// will be `false`.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsInt64(aKey string) (rVal int64, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsInt64()
	}

	return
} // AsInt64()

// `asKeyVal` returns the value of `aKey` as a key/value pair.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) asKeyVal(aKey string) (TKeyVal, bool) {
	var kv TKeyVal
	for _, kv = range *ks {
		if kv.Key == aKey {
			return kv, true
		}
	}

	return kv, false
} // asKeyVal()

// `AsString` returns the value of `aKey` as a string.
//
// If the given `aKey` doesn't exist then the second return value
// will be `false`.
//
//	`aKey` The name of the key to lookup.
func (ks *TSection) AsString(aKey string) (rVal string, rOK bool) {
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		rVal, rOK = kv.AsString()
	}

	return
} // AsString()

// `Clear` removes all entries in this INI section.
func (ks *TSection) Clear() *TSection {
	(*ks) = (*ks)[:0]

	return ks
} // Clear()

// `HasKey` returns whether `aKey` exists in this INI section.
//
//	`aKey` The key to lookup.
func (ks *TSection) HasKey(aKey string) bool {
	return (0 <= ks.IndexOf(aKey))
} // HasKey()

// `IndexOf` returns the index of `aKey` in this INI section or `-1`
// if not found.
//
//	`aKey` The key to lookup.
func (ks *TSection) IndexOf(aKey string) int {
	for result, kv := range *ks {
		if kv.Key == aKey {
			return result
		}
	}

	return -1
} // IndexOf()

// `Len` returns the number of key/value pairs in this section.
func (ks *TSection) Len() int {
	return len(*ks)
} // Len()

// `RemoveKey` removes `aKey` from this section.
//
// This method returns 'true' if `aKey` doesn't exist at all, or if `aKey`
// was successfully removed, or `false` otherwise.
//
//	`aKey` The name of the key/value pair to remove.
func (ks *TSection) RemoveKey(aKey string) bool {
	idx := ks.IndexOf(aKey)
	if 0 > idx {
		return true
	}
	sLen := len(*ks) - 1 // new slice length (i.e. one shorter)
	(*ks)[idx] = TKeyVal{}
	switch idx {
	case 0:
		(*ks) = (*ks)[1:]
	case sLen:
		(*ks) = (*ks)[:sLen]
	default:
		(*ks) = append((*ks)[:idx], (*ks)[1+idx:]...)
	}

	return (0 > ks.IndexOf(aKey))
} // RemoveKey()

// `String` returns a string representation of an INI section.
//
// The single key/value pairs are delimited by a linefeed ('\n).
func (ks *TSection) String() (rString string) {
	for _, kv := range *ks {
		rString += kv.String() + "\n"
	}

	return
} // String()

// `UpdateKey` replaces the current value of `aKey` by the provided
// new `aValue`.
//
// In case `aKey` doesn't already exist in the list (and therefore can't
// be updated) it will be added by calling the `AddKey()` method.
//
// If `aKey` is an empty string the method's result will be `false`.
//
//	`aKey` The key of the key/value pair to update.
//	`aValue` The value of the key/value pair to update.
func (ks *TSection) UpdateKey(aKey, aValue string) (rOK bool) {
	if 0 == len(aKey) {
		return false
	}
	var kv TKeyVal
	if kv, rOK = ks.asKeyVal(aKey); rOK {
		if kv.UpdateValue(aValue) {
			return
		}
	}

	// if aKey doesn't exist then create a new entry
	return ks.AddKey(aKey, aValue)
} // UpdateKey()

// ReadIniData returns the config values read from INI file(s).
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
//	`aName` The application's name used in the INI file name
//
// (without `.ini` extension).
func ReadIniData(aName string) *TSection {
	var (
		confDir    string
		err        error
		ini1, ini2 *TIniList
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
