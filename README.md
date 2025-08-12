# gojsonfmt

[![PkgGoDev](https://pkg.go.dev/badge/github.com/chris-peterson444/gojsonfmt)](https://pkg.go.dev/github.com/chris-peterson444/gojsonfmt)

This package provides a way to format JSON into in a Go-like way. Opening and
closing braces will be compacted where possible, and indentation always uses
tabs.

## Usage

### As a Package

```bash
go get github.com/chris-peterson444/gojsonfmt
```

### Install the binary

```bash
go install github.com/chris-peterson444/gojsonfmt/cmd/gojsonfmt@latest
```

```bash
$ gojsonfmt --help
Usage of gojsonfmt:
  -file string
        Path to JSON data to format
  -stdin
        Read raw JSON from stdin and format it
```

## Examples

The following examples demonstrate the difference between `jq` formatting and 
what `gojsonfmt` will produce. 

### Lists of Objects

Objects in a list will have their final `}` character on the same line as the
first `{` character in the next object. The final `]` character of the list
will follow right after the last closing `}` of the final object in the list.

Original: 

```json
{
  "foo": [
    {
      "bar": 1,
      "baz": 2
    },
    {
      "bar": 3,
      "baz": 4
    }
  ]
}
```

Formatted:

```json
{
	"foo": [{
		"bar": 1,
		"baz": 2
	}, {
		"bar": 3,
		"baz": 4
	}]
}
```

### Lists of Lists

Lists of lists, and lists of alternating types are handled similarly.

Original:

```json
{
  "foo": [
    [
      1,
      2,
      3
    ],
    [
      {
        "bar": 1
      },
      {
        "bar": 2
      }
    ]
  ]
}
```

Formatted:

```json
{
	"foo": [[
		1,
		2,
		3
	], [{
		"bar": 1
	}, {
		"bar": 2
	}]]
}
```

### Empty Objects and Lists

Objects with no attributes and empty lists will be formatted in-line. Lists
with empty objects will attempt to compact as well.

Original:

```json
{
  "foo": {},
  "bar": [
    [
      {}
    ]
  ],
  "baz": {
    "qux": {}
  }
}
```

Formatted:

```json
{
	"foo": {},
	"bar": [[{}]],
	"baz": {
		"qux": {}
	}
}
```

### 
