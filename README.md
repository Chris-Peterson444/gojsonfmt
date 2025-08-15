# gojsonfmt

[![PkgGoDev](https://pkg.go.dev/badge/github.com/chris-peterson444/gojsonfmt)](https://pkg.go.dev/github.com/chris-peterson444/gojsonfmt)

This package provides a way to format JSON into in a Go-like way. Opening and
closing braces will be compacted where possible, and indentation always uses
tabs.

## Disclaimer

This is not a Canonical product. Just a program written by a Canonical employee
on company time. Assume no support, warranty, etc.

## Installation

### As a Package

```bash
go get github.com/chris-peterson444/gojsonfmt
```

### Get the binary

```bash
go install github.com/chris-peterson444/gojsonfmt/cmd/gojsonfmt@latest
```

```
$ gojsonfmt --help
Usage: gojsonfmt [JSON_TEXT]
       gojsonfmt --file <file-path>
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

  -file string
        Path to JSON data to format
  -version
        show version information
```

## Examples

The following examples demonstrate the difference between `jq` formatting and 
what `gojsonfmt` will produce. 

### Lists of Objects

Objects in a list will have their final `}` character on the same line as the
first `{` character in the next object. The final `]` character of the list
will follow right after the last closing `}` of the final object in the list.

jq: 

```
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

gojsonfmt:

```
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

jq:

```
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

gojsonfmt:

```
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

jq:

```
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

gojsonfmt:

```
{
	"foo": {},
	"bar": [[{}]],
	"baz": {
		"qux": {}
	}
}
```

### 
