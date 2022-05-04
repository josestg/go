package flatten

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
	"strings"
	"testing"
)

type any = interface{}

func TestFlatten(t *testing.T) {

	testCases := []struct {
		desc string
		inp  string
		sep  string
		out  any
	}{
		{
			desc: "a simple nested json; expecting seperated by '.'",
			inp:  `{"a":{"b":{"c":{"d":"value"}}}}`,
			sep:  ".",
			out:  map[string]any{"a.b.c.d": "value"},
		},
		{
			desc: "a simple nested json; expecting seperated by '_'",
			inp:  `{"a":{"b":{"c":{"d":"value"}}}}`,
			sep:  "_",
			out:  map[string]any{"a_b_c_d": "value"},
		},
		{
			desc: "a nested json with branches",
			inp:  `{"a":{"b":{"c1":{"d":"value 1"},"c2":{"d":"value 2"}}}}`,
			sep:  ".",
			out:  map[string]any{"a.b.c1.d": "value 1", "a.b.c2.d": "value 2"},
		},
		{
			desc: "a nested json with recursive branches",
			inp:  `{"a1":{"b1":{"c1":{"d":"value 1"},"c2":{"d":"value 2"}},"b2":{"c1":{"d":"value 3"},"c2":{"d":"value 4"}}},"a2":{"b1":{"c1":{"d":"value 5"},"c2":{"d":"value 6"}},"b2":{"c1":{"d":"value 7"},"c2":{"d":"value 8"}}}}`,
			sep:  ".",
			out: map[string]any{
				"a1.b1.c1.d": "value 1",
				"a1.b1.c2.d": "value 2",
				"a1.b2.c1.d": "value 3",
				"a1.b2.c2.d": "value 4",
				"a2.b1.c1.d": "value 5",
				"a2.b1.c2.d": "value 6",
				"a2.b2.c1.d": "value 7",
				"a2.b2.c2.d": "value 8",
			},
		},
		{
			desc: "a nested json with array of objects",
			inp:  `{"a":{"b":[{"c":"1","d":"2"},{"e":"3","f":"4"}]}}`,
			sep:  ".",
			out: map[string]any{
				"a.b.0.c": "1",
				"a.b.0.d": "2",
				"a.b.1.e": "3",
				"a.b.1.f": "4",
			},
		},
		{
			desc: "a nested json with array of objects and scalars",
			inp:  `{"a":{"b":[{"c":"1","d":"2"},{"e":"3","f":"4"},"5","6"]}}`,
			sep:  ".",
			out: map[string]any{
				"a.b.0.c": "1",
				"a.b.0.d": "2",
				"a.b.1.e": "3",
				"a.b.1.f": "4",
				"a.b.2":   "5",
				"a.b.3":   "6",
			},
		},
		{
			desc: "a nested arrays",
			inp:  `[[["1","2","3"]]]`,
			sep:  ".",
			out: map[string]any{
				"0.0.0": "1",
				"0.0.1": "2",
				"0.0.2": "3",
			},
		},
		{
			desc: "a nested array of objects",
			inp:  `[[[{"a":"1"},{"b":"2"},{"c":"3"}]]]`,
			sep:  ".",
			out: map[string]any{
				"0.0.0.a": "1",
				"0.0.1.b": "2",
				"0.0.2.c": "3",
			},
		},
		{
			desc: "a complex json structure",
			inp:  `{"a":["1",{"b":{"c":"2","d":["3",["4",{"e":"5","f":"6"},"7",["8","9"]]]},"g":"10","h":{"i":{"j":"11"}}},"12",{"k":"13","l":[{"m":"14"},{"n":"15"}]}],"o":{"p":["16",{"q":"17"}]}}`,
			sep:  ".",
			out: map[string]any{
				"a.0":           "1",
				"a.1.b.c":       "2",
				"a.1.b.d.0":     "3",
				"a.1.b.d.1.0":   "4",
				"a.1.b.d.1.1.e": "5",
				"a.1.b.d.1.1.f": "6",
				"a.1.b.d.1.2":   "7",
				"a.1.b.d.1.3.0": "8",
				"a.1.b.d.1.3.1": "9",
				"a.1.g":         "10",
				"a.1.h.i.j":     "11",
				"a.2":           "12",
				"a.3.k":         "13",
				"a.3.l.0.m":     "14",
				"a.3.l.1.n":     "15",
				"o.p.0":         "16",
				"o.p.1.q":       "17",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {

			if testing.Verbose() {
				fmt.Fprintf(os.Stdout, "=== TEST %s ===\n", strings.ToTitle(tc.desc))
				inp, _ := json.MarshalIndent(json.RawMessage(tc.inp), "", "  ")
				out, _ := json.MarshalIndent(tc.out, "", "  ")
				fmt.Fprintf(os.Stdout, "=== INPUT\n%s\n", string(inp))
				fmt.Fprintf(os.Stdout, "=== EXPECTED\n%s\n", string(out))
			}

			var obj any
			if err := json.NewDecoder(strings.NewReader(tc.inp)).Decode(&obj); err != nil {
				t.Fatalf("creating obj: err: %v", err)
			}

			flat := Flatten(obj, tc.sep)
			if diff := cmp.Diff(flat, tc.out); len(diff) != 0 {
				t.Fatal(diff)
			}
		})
	}
}

func ExampleFlatten() {
	mapObj := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 1,
				"d": 2,
			},
		},
	}

	sliceObj := []interface{}{"a", []string{"b", "c"}}

	flattenMapObj := Flatten(mapObj, ".")
	flattenSliceObj := Flatten(sliceObj, ".")
	fmt.Println(flattenMapObj)
	fmt.Println(flattenSliceObj)
	// Output:
	// map[a.b.c:1 a.b.d:2]
	// map[0:a 1.0:b 1.1:c]
}
