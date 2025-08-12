# gojsonfmt

This package provides a way to format JSON into in a Go-like way. Opening and
closing braces will be compacted where possible, and indentation always uses
tabs.

## Examples:

### Lists of Objects

The following JSON document, formatted by `jq`, contains a single key, which
holds a list of objects.

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

will be formatted into:

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

Lists of lists, and lists of alternating types are handled as well:

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

formatted:

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
