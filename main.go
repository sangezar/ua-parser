package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"ua-parser/internal"
)

func main() {
	fieldPath := os.Getenv("USER_AGENT_FIELD")
	if fieldPath == "" {
		fmt.Fprintln(os.Stderr, "USER_AGENT_FIELD env not set")
		os.Exit(1)
	}

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	for {
		var msg any
		if err := dec.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, "decode:", err)
			continue
		}

		obj, ok := msg.(map[string]any)
		if !ok {
			_ = enc.Encode(msg)
			continue
		}

		uaStr, prefix, ok := internal.GetByDotPath(obj, fieldPath)
		if !ok || strings.TrimSpace(uaStr) == "" {
			_ = enc.Encode(obj)
			continue
		}

		internal.EnrichFlat(obj, prefix, uaStr)
		_ = enc.Encode(obj)
	}
}
