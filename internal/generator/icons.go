package generator

import (
	"embed"
	"encoding/json"
	"strings"

	"github.com/kaugesaar/lucide-go"
)

//go:embed icons.json
var iconsFS embed.FS

var iconsData map[string]IconData

type IconData struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Hex   string `json:"hex"`
	Path  string `json:"path"`
}

func init() {
	data, err := iconsFS.ReadFile("icons.json")
	if err != nil {
		// Will be populated at build time
		iconsData = make(map[string]IconData)
		return
	}

	var icons []IconData
	if err := json.Unmarshal(data, &icons); err != nil {
		iconsData = make(map[string]IconData)
		return
	}

	iconsData = make(map[string]IconData)
	for _, icon := range icons {
		// to avoid slug mismatch when looking up in simpleicons
		iconsData[strings.ToLower(icon.Slug)] = icon
	}
}

const placeholderSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="12" r="10"/></svg>`

func GetIconSVG(iconName, provider string) string {
	iconName = normalizeSlug(iconName)
	iconProvider := strings.TrimSpace(strings.ToLower(provider))

	switch iconProvider {
	case "", "simpleicon":
		return GetSimpleIcon(iconName)

	case "lucide":
		return getLucideIconSVG(iconName)

	default:
		return placeholderSVG
	}
}

func normalizeSlug(slug string) string {
	return strings.ToLower(strings.TrimSpace(slug))
}

func getLucideIconSVG(slug string) string {
	icon := lucide.Icon(slug, map[string]any{
		"size":        24,
		"strokeWidth": 2,
		"color":       "currentColor",
	})

	// lucide.Icon returns template.HTML
	if icon == "" {
		return placeholderSVG
	}

	return string(icon)
}

// GetSimpleIcon returns an SVG string for the given icon slug
func GetSimpleIcon(slug string) string {
	slug = strings.ToLower(slug)
	icon, ok := iconsData[slug]
	if !ok {
		// Return a placeholder circle if icon not found
		return placeholderSVG
	}

	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="` + icon.Path + `"/></svg>`
}

// GetIconColor returns the brand color for the given icon slug
func GetIconColor(slug string) string {
	slug = strings.ToLower(slug)
	icon, ok := iconsData[slug]
	if !ok {
		return "#ffffff"
	}
	return "#" + icon.Hex
}
