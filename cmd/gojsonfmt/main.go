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

	"github.com/chris-peterson444/gojsonfmt"
)

func main() {
	stdinFlag := flag.Bool("stdin", false, "Read raw JSON from stdin and format it")
	fileFlag := flag.String("file", "", "Path to JSON data to format")
	flag.Parse()

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
