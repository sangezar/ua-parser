package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

	sc := bufio.NewScanner(os.Stdin)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxCapacity)
	enc := json.NewEncoder(os.Stdout)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			fmt.Fprintln(os.Stderr, "invalid json:", err)
			continue
		}

		uaStr, prefix, ok := internal.GetByDotPath(obj, fieldPath)
		if !ok || strings.TrimSpace(uaStr) == "" {
			_ = enc.Encode(obj) // return unchanged
			continue
		}

		internal.EnrichFlat(obj, prefix, uaStr)
		_ = enc.Encode(obj)
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
	}
}
