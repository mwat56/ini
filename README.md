# INI

[![GoDoc](https://godoc.org/github.com/mwat56/ini?status.svg)](https://godoc.org/github.com/mwat56/ini)
[![License](https://img.shields.io/eclipse-marketplace/l/notepad4e.svg)](https://github.com/mwat56/ini/blob/master/LICENSE)

- [INI](#ini)
	- [Purpose](#purpose)
	- [Installation](#installation)
	- [Usage](#usage)
	- [Licence](#licence)


## Purpose

Over the times several different file formats have been developed just for storing configuration data for some program.
While they all may have some merits, for me the two-dimensional INI file format – made popular by the DR-/MS-/PC-DOS and MS-Windows versions in the 80s of the last century – was always sufficient for my needs.
This package provides the `TSections` class to read/parse, modify, and write such INI files. It doesn't need any configuration but simply does what it's supposed to do.

## Installation

You can use `Go` to install this package for you:

    go get -u github.com/mwat56/ini

## Usage

An INI file usually looks like this:

    ; This is a comment

    [aSectionName]
        key1 = value 1
        key2 = value2
        …

    [anotherSection]
        key1 = value1
        key2 = value 2 is \
        really long and\
        spans several lines
        …

Leading whitespace is ignored, empty lines and those beginning with either a semicolon (`;`) or a number sign (`#`) are skipped (and not preserved when overwriting the file).
Lines that can't be identified as either a _section heading_ or a _key/value pair_ are silently ignored as well.
Quotes and whitespace surrounding a key or a value are ignored.

A line ending with a backslash (`\`) will be concatenated with the following line (unless that's a comment line).
By that mechanism you can use really long values spaning several lines.

You can create a `TSections` instance by either calling `ini.NewSections()` and then using the numerous methods (including `Load()` and `Store()`).
Or you simply call `ini.LoadFile(aFilename)` which does – as the name suggests – the loading for you.

_Note_ that both section and key names are _case sensitive_ to allow for the broadest possible range when naming them.
The same is true for the values which are, of course, case sensitive.
An application using this package, however, is free to interpret the values returned in any way they like.

Please look at the [source code documentation](https://godoc.org/github.com/mwat56/ini#TSections) to see the numerous methods provided to load, get, set, and update sections and key/value pairs.

## Licence

    Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                    All rights reserved
                EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.
