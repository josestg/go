package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/josestg/go/flatten"
	"io"
	"os"
	"strings"
)

const help = `Flattens the object - it'll return an object one level deep, 
regardless of how nested the original object was.

example:
  > echo '{"a":{"b":{"c":{"d":"value"}}}}' | flatten 
  {"a.b.c.d":"value"}

The flatten supports both JSON and YAML as input.
`

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		if strings.EqualFold(args[0], "help") {
			fmt.Fprint(os.Stderr, help)
			return
		}
	}

	var object flatten.Any
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "flatten: read input: %v\n", err)
		os.Exit(1)
	}
	// This Unmarshal supports JSON and Yaml.
	if err := yaml.Unmarshal(buf.Bytes(), &object); err != nil {
		fmt.Fprintf(os.Stderr, "flatten: decoding: %v\n", err)
		os.Exit(1)
	}

	flatObject := flatten.Flatten(object, ".")
	if err := json.NewEncoder(os.Stdout).Encode(&flatObject); err != nil {
		fmt.Fprintf(os.Stderr, "flatten: encoding: %v\n", err)
		os.Exit(1)
	}
}
