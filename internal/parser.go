package internal

import (
	"encoding/json"
	"regexp"
	"strings"

	ua "github.com/mileusna/useragent"
)

var nonWord = regexp.MustCompile(`[^A-Za-z0-9]+`)

func normPrefix(s string) string {
	s = strings.ToLower(s)
	s = nonWord.ReplaceAllString(s, "_")
	return strings.Trim(s, "_")
}

// GetByDotPath fetches a value from a map by a dot-path and returns:
// (value, last segment as prefix, ok)
func GetByDotPath(m map[string]any, path string) (string, string, bool) {
	parts := strings.Split(path, ".")
	var cur any = m
	for i, p := range parts {
		asMap, ok := cur.(map[string]any)
		if !ok {
			return "", "", false
		}
		var next any
		var found bool
		for k, v := range asMap {
			kl := strings.ToLower(k)
			pl := strings.ToLower(p)
			if kl == pl || kl == strings.ReplaceAll(pl, "_", "-") || kl == strings.ReplaceAll(pl, "-", "_") {
				next = v
				p = k // зберігаємо реальне ім'я ключа
				found = true
				break
			}
		}
		if !found {
			return "", "", false
		}
		if i == len(parts)-1 {
			if s, ok := next.(string); ok && strings.TrimSpace(s) != "" {
				return s, p, true
			}
			b, err := json.Marshal(next)
			if err != nil {
				return "", "", false
			}
			val := string(b)
			if strings.TrimSpace(val) == "" || val == "null" {
				return "", "", false
			}
			return val, p, true
		}
		cur = next
	}
	return "", "", false
}

// EnrichFlat adds flat UA fields with a prefix and no nesting.
func EnrichFlat(obj map[string]any, prefix string, uaStr string) {
	u := ua.Parse(uaStr)

	deviceType := "desktop"
	switch {
	case u.Bot:
		deviceType = "bot"
	case u.Tablet:
		deviceType = "tablet"
	case u.Mobile:
		deviceType = "mobile"
	case u.Desktop:
		deviceType = "desktop"
	}

	p := normPrefix(prefix)
	set := func(k string, v any) { obj[p+"_"+k] = v }

	// core fields
	set("browser_name", u.Name)
	set("browser_version", u.Version)
	set("os_name", u.OS)
	set("os_version", u.OSVersion)

	// synthetic/compatibility
	set("device_type", deviceType)
	set("is_mobile", u.Mobile)
	set("is_tablet", u.Tablet)
	set("is_desktop", u.Desktop && !u.Bot)
	set("is_bot", u.Bot)

	// additional fields from mileusna/useragent
	set("device_name", u.Device) // e.g., iPhone, iPad, Huawei...
	set("bot_url", u.URL)        // for bots; empty for regular agents
	set("is_unknown", u.IsUnknown())
}
