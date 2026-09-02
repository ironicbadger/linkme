package generator

import (
	"strings"
	"testing"
)

func TestGetIconSVG(t *testing.T) {
	tests := []struct {
		name        string
		icon        string
		provider    string
		placeholder bool
		contains    string
	}{
		{name: "default simple icon", icon: "github", contains: "<path"},
		{name: "explicit simple icon", icon: " GitHub ", provider: " SIMPLEICON ", contains: "<path"},
		{name: "lucide icon", icon: " Broom ", provider: " LUCIDE ", contains: `stroke="currentColor"`},
		{name: "unknown simple icon", icon: "not-a-real-icon", placeholder: true},
		{name: "unknown lucide icon", icon: "not-a-real-icon", provider: "lucide", placeholder: true},
		{name: "unknown provider", icon: "github", provider: "other", placeholder: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIconSVG(tt.icon, tt.provider)
			if (got == placeholderSVG) != tt.placeholder {
				t.Fatalf("placeholder = %v, want %v; SVG: %s", got == placeholderSVG, tt.placeholder, got)
			}
			if tt.contains != "" && !strings.Contains(got, tt.contains) {
				t.Fatalf("SVG does not contain %q: %s", tt.contains, got)
			}
		})
	}
}

func TestGetSimpleIconNormalizesCase(t *testing.T) {
	if got, want := GetSimpleIcon("GitHub"), GetSimpleIcon("github"); got != want {
		t.Fatalf("mixed-case slug returned different SVG")
	}
}
