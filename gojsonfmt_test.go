package gojsonfmt_test

import (
	"testing"

	"github.com/chris-peterson444/gojsonfmt"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		summary  string
		input    string
		expected string
	}{{
		summary: "Compact JSON object",
		input:   `{"foo": [{"bar": 1}, {"bar": 2}]}`,
		expected: `{
	"foo": [{
		"bar": 1
	}, {
		"bar": 2
	}]
}`,
	}, {
		summary:  "empty is empty",
		input:    "",
		expected: "",
	}, {
		summary: "Long example formatted with jq",
		input: `{
  "single-attr": {
    "one-attr-obj": {
      "thing": "thing"
    },
    "multiple-attrs-obj": {
      "some-string": "foo-bar-baz",
      "some-number": 123.123,
      "some-bool": true,
      "some-null": null,
      "empty-obj": {},
      "empty-list": []
    },
    "non-empty-list-of-objects": [
      {
        "some-string": "foo-bar-baz",
        "some-number": 123
      },
      {
        "some-bool": false,
        "some-list": []
      }
    ],
    "object-with-empty-objs": {
      "some-obj-foo": {},
      "some-obj-bar": {}
    }
  },
  "list-of-objs": [
    {
      "some-string": "foo-bar-baz",
      "inner-list-of-objects": [
        {
          "some-string": "foo"
        }
      ],
      "inner-list-of-objects-two": [
        {
          "some-numer": 123
        },
        {
          "some-bool": true
        },
        {
          "some-empty-object": {}
        }
      ]
    },
    {
      "some-empty-list": []
    }
  ],
  "list-of-lists": [
    [
      1,
      2,
      3
    ],
    [
      "foo",
      "bar",
      "bar"
    ],
    [
      null
    ],
    [
      true,
      false
    ],
    [
      {
        "some-empty-list": []
      },
      {
        "some-empty-obj": {}
      }
    ]
  ]
}
`,
		expected: `{
	"single-attr": {
		"one-attr-obj": {
			"thing": "thing"
		},
		"multiple-attrs-obj": {
			"some-string": "foo-bar-baz",
			"some-number": 123.123,
			"some-bool": true,
			"some-null": null,
			"empty-obj": {},
			"empty-list": []
		},
		"non-empty-list-of-objects": [{
			"some-string": "foo-bar-baz",
			"some-number": 123
		}, {
			"some-bool": false,
			"some-list": []
		}],
		"object-with-empty-objs": {
			"some-obj-foo": {},
			"some-obj-bar": {}
		}
	},
	"list-of-objs": [{
		"some-string": "foo-bar-baz",
		"inner-list-of-objects": [{
			"some-string": "foo"
		}],
		"inner-list-of-objects-two": [{
			"some-numer": 123
		}, {
			"some-bool": true
		}, {
			"some-empty-object": {}
		}]
	}, {
		"some-empty-list": []
	}],
	"list-of-lists": [[
		1,
		2,
		3
	], [
		"foo",
		"bar",
		"bar"
	], [
		null
	], [
		true,
		false
	], [{
		"some-empty-list": []
	}, {
		"some-empty-obj": {}
	}]]
}`,
	}, {
		summary: "empty list and object edge cases",
		input: `{
			  "list-of-single-empty-list": [
			    []
			  ],
			  "list-of-multiple-empty-lists": [
			    [],
			    [],
			    []
			  ],
			  "list-of-empty-object": [
			    {}
			  ],
			  "list-of-multiple-empty-object": [
			    {},
			    {}
			  ]
}`,
		expected: `{
	"list-of-single-empty-list": [[]],
	"list-of-multiple-empty-lists": [[], [], []],
	"list-of-empty-object": [{}],
	"list-of-multiple-empty-object": [{}, {}]
}`,
	}, {
		summary: "lists and objects of various completeness and nesting",
		input: `{
  "foo": [
    {
      "foo": 1
    },
    {},
    [
      {
        "foo": 1
      }
    ],
    [],
    [
      {
        "foo": [
          {
            "bar": 1,
	    "baz": [1,2,3]
          }
        ]
      }
    ],
    [],
    {}
  ]
}`,
		expected: `{
	"foo": [{
		"foo": 1
	}, {}, [{
		"foo": 1
	}], [], [{
		"foo": [{
			"bar": 1,
			"baz": [
				1,
				2,
				3
			]
		}]
	}], [], {}]
}`,
	}, {
		summary: "list of lists example",
		input: `{
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
}`,
		expected: `{
	"foo": [[
		1,
		2,
		3
	], [{
		"bar": 1
	}, {
		"bar": 2
	}]]
}`,
	}, {
		summary: "empty obj examples",
		input: `{
		  "foo": {},
		  "bar": [
			[
			  {}
			]
			],
		  "baz": {
		    "qux": {}
		  }
		}`,
		expected: `{
	"foo": {},
	"bar": [[{}]],
	"baz": {
		"qux": {}
	}
}`,
	}}

	for _, test := range tests {
		t.Logf("summary: %s", test.summary)
		generated, err := gojsonfmt.FormatJSONString(test.input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if generated != test.expected {
			t.Fatalf("expected:\n %s\nbut got:\n%s\n", test.expected, generated)
		}
	}
}
