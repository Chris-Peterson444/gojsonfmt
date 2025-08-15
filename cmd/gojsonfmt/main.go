// Copyright 2025 Chris Peterson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = customUsage
	stdinFlag := flag.Bool("stdin", true, "Read raw JSON from stdin and format it")
	fileFlag := flag.String("file", "", "Path to JSON data to format")
	versionFlag := flag.Bool("version", false, "show version information")

	flag.Parse()

	if *versionFlag {
		version := getVersion()
		fmt.Printf("version: %q\n", version)
		return
	}

	if *stdinFlag && *fileFlag != "" {
		fmt.Fprintln(os.Stderr, "Error: --stdin and --file cannot be used together")
		os.Exit(1)
	}

	if *stdinFlag {
		data, err := io.ReadAll(bufio.NewReader(os.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read stdin: %v\n", err)
			os.Exit(1)
		}
		formatJSON(data)
		return
	}

	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read file: %v\n", err)
		}
		formatJSON(data)
		return
	}

	fmt.Fprintln(os.Stderr, "error: You must specify either --stdin or --file")
	os.Exit(1)
}

func formatJSON(data []byte) {
	formatted, err := gojsonfmt.FormatJSONBytes(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to format JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(formatted))
}
