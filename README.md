UA Parser – enrich JSON lines with parsed User-Agent fields

UA Parser is a small, stateless CLI that parses User-Agent strings and enriches JSON lines with flat fields. It works in streaming mode: reads JSON from stdin and writes enriched JSON to stdout.

You can use it either:
- as a simple CLI filter in the terminal; or
- as a subprocess inside data pipelines (e.g., Redpanda Connect via `subprocess`), where it continuously processes messages.

Key points:
- Single binary, no runtime dependencies
- Stateless streaming via JSON decoder (no need for newline delimiters)
- Configure UA field via `--user-agent-field` flag (env `USER_AGENT_FIELD` as fallback)

### Build

- Build for both Linux amd64 and macOS arm64:
```bash
make
```

Artifacts are in `bin/`:
- `bin/ua-parser-linux-amd64`
- `bin/ua-parser-darwin-arm64`

### Console usage

The tool reads JSON objects from stdin and writes enriched JSON objects to stdout. Provide the User-Agent field path via flag (or env fallback).

Examples:
```bash
# Using flag
echo '{"user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36"}' \
| ./bin/ua-parser-darwin-arm64 --user-agent-field user_agent | jq .

# Using env (fallback)
echo '{"user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36"}' \
| USER_AGENT_FIELD=user_agent ./bin/ua-parser-darwin-arm64 | jq .
```

If a line does not contain the field (or it is empty), the line is returned unchanged.

### Output fields

Given a prefix derived from the last path segment (normalized to lowercase and underscores), the following fields are added:
- `<prefix>_browser_name`
- `<prefix>_browser_version`
- `<prefix>_os_name`
- `<prefix>_os_version`
- `<prefix>_device_type` (desktop|mobile|tablet|bot)
- `<prefix>_is_mobile`
- `<prefix>_is_tablet`
- `<prefix>_is_desktop`
- `<prefix>_is_bot`
- `<prefix>_device_name`
- `<prefix>_bot_url`
- `<prefix>_is_unknown`

### Redpanda Connect (rpk) example

Use the `subprocess` processor to call the binary and enrich messages. Below is a simple YAML pipeline that reads from stdin and writes to stdout, invoking the parser as a subprocess.

```yaml
input:
  label: stdin
  stdin: {}

pipeline:
  processors:
    - subprocess:
        name: /plugins/ua_parser
        args: ["--user-agent-field", "headers.user_agent"]

output:
  stdout: {}
```

Notes:
- Ensure the binary path is correct and executable.
- For macOS on Apple Silicon, use `ua-parser-darwin-arm64` in the command.
