package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"ua-parser/internal"
)

func main() {
	uaFieldFlag := flag.String("user-agent-field", "", "Dot-path to the User-Agent field")
	flag.Parse()

	fieldPath := strings.TrimSpace(*uaFieldFlag)
	if fieldPath == "" {
		fieldPath = os.Getenv("USER_AGENT_FIELD")
	}
	if fieldPath == "" {
		fmt.Fprintln(os.Stderr, "missing user agent field: pass --user-agent-field or set USER_AGENT_FIELD env")
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
