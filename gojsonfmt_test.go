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

package gojsonfmt_test

import (
	"testing"

	"github.com/chris-peterson444/gojsonfmt"
)

var conversionTests = []struct {
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
			  ],
			 "object-with-multiple-empty-arrays": {
				"foo": [],
				"bar": []
			}
}`,
	expected: `{
	"list-of-single-empty-list": [[]],
	"list-of-multiple-empty-lists": [[], [], []],
	"list-of-empty-object": [{}],
	"list-of-multiple-empty-object": [{}, {}],
	"object-with-multiple-empty-arrays": {
		"foo": [],
		"bar": []
	}
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
	summary: "deeply nested lists with varying types",
	input: `{
  "foo": [
    [
      1,
      2,
      [
        [
          3,
          4,
          [
            {
              "bar": [
                5,
                6
              ]
            },
            7,
            8,
	    [ 8.1, 8.2],
            {
              "baz": 1
            },
            {},
            9,
            {},
            [],
            10
          ]
        ]
      ]
    ]
  ]
}`,
	expected: `{
	"foo": [[
		1,
		2,
		[[
			3,
			4,
			[{
				"bar": [
					5,
					6
				]
			},
			7,
			8,
			[
				8.1,
				8.2
			], {
				"baz": 1
			}, {},
			9,
			{}, [],
			10
	]]]]]
}`,
}, {
	summary: "deeply nested list inside of simple list",
	input: `{
  "foo": [
    1,
    2,
    3,
    [
      [
        {
          "bar": "baz",
          "baz": [
            []
          ]
        }
      ]
    ],
    4,
    5,
    6
  ]
}`,
	expected: `{
	"foo": [
		1,
		2,
		3,
		[[{
			"bar": "baz",
			"baz": [[]]
		}]],
		4,
		5,
		6
	]
}`,
}, {
	summary: "list that starts with object, then simple values",
	input: `{
		  "foo": [
		    {
		      "bar": 1
		    },
		    2,
		    3
		  ]
		}`,
	expected: `{
	"foo": [{
		"bar": 1
	},
	2,
	3
	]
}`,
}, {
	summary: "list that starts with simple values, then objects",
	input: `{
		  "foo": [
		    1,
		    {
		      "bar": 2
		    },
		    3,
		    {
		      "bar": 4
		    },
		    5
		  ]
		}`,
	expected: `{
	"foo": [
		1,
		{
			"bar": 2
		},
		3,
		{
			"bar": 4
		},
		5
	]
}`,
}}

func TestFormatJSONString(t *testing.T) {
	for _, test := range conversionTests {
		t.Run(test.summary, func(t *testing.T) {
			generated, err := gojsonfmt.FormatJSONString(test.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if generated != test.expected {
				t.Fatalf("expected:\n %s\nbut got:\n%s\n", test.expected, generated)
			}
		})
	}
}

func TestFormatJSONBytes(t *testing.T) {
	for _, test := range conversionTests {
		t.Run(test.summary, func(t *testing.T) {
			generated, err := gojsonfmt.FormatJSONBytes([]byte(test.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(generated) != test.expected {
				t.Fatalf("expected:\n %s\nbut got:\n%s\n", test.expected, generated)
			}
		})
	}
}
