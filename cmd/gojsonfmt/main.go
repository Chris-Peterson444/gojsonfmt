// Copyright (C) 2025 Canonical Ltd.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/chris-peterson444/gojsonfmt"
)

// Update for releases.
const VERSION = "devel"

func getVersion() string {
	// Get exact commit version if possible.
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return VERSION
	}
	if VERSION != "devel" {
		return VERSION
	}
	var version, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			version = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}

	if modified == "true" {
		version += "+dirty"
	}
	// return version
	return VERSION + "+" + version
}

func customUsage() {
	summary := `
Format JSON in a Go-like way. Opening and closing braces will be compacted where
possible, and indentation always uses tabs. 

By default it will read in raw JSON from stdin in, but you can pass --file
to read the text in a specified file.

Example:

	$ echo '{"foo": [{"bar": 1}, {"bar": 2}]}' | gojsonfmt
	{
		"foo": [{
			"bar": 1
		}, {
			"bar": 2
		}]
	}
`

	fmt.Printf("Usage: %s [JSON_TEXT]\n", os.Args[0])
	fmt.Printf("       %s --file <file-path>", os.Args[0])
	fmt.Println(summary)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = customUsage
	fileFlag := flag.String("file", "", "Path to JSON data to format")
	versionFlag := flag.Bool("version", false, "show version information")

	flag.Parse()

	if *versionFlag {
		version := getVersion()
		fmt.Printf("version: %q\n", version)
		return
	}

	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read file: %v\n", err)
		}
		formatJSON(data)
		return
	} else {
		data, err := io.ReadAll(bufio.NewReader(os.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read stdin: %v\n", err)
			os.Exit(1)
		}
		formatJSON(data)
		return
	}
}

func formatJSON(data []byte) {
	formatted, err := gojsonfmt.FormatJSONBytes(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to format JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(formatted))
}
