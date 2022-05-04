package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/josestg/go/flatten"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "flatten: encoding: %v\n", err)
			os.Exit(1)
		}

		cmd := scanner.Text()
		if cmd == ".exit" {
			return
		}

		var object flatten.Any
		// This Unmarshal supports JSON and Yaml.
		if err := yaml.Unmarshal(scanner.Bytes(), &object); err != nil {
			fmt.Fprintf(os.Stderr, "flatten: yaml: decoding: %v\n", err)
			os.Exit(1)
		}

		flatObject := flatten.Flatten(object, ".")
		if err := json.NewEncoder(os.Stdout).Encode(&flatObject); err != nil {
			fmt.Fprintf(os.Stderr, "flatten: encoding: %v\n", err)
			os.Exit(1)
		}
	}
}
