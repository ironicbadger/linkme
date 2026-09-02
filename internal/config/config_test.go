package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadIconProviders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	data := []byte(`
links:
  - title: Blog
    icon: rss
    icon-provider: lucide
sections:
  - title: Elsewhere
    links:
      - title: Code
        icon: github
        icon-provider: simpleicon
socials:
  - icon: message-circle
    icon-provider: lucide
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.Links[0].IconProvider; got != "lucide" {
		t.Fatalf("link provider = %q, want lucide", got)
	}
	if got := cfg.Sections[0].Links[0].IconProvider; got != "simpleicon" {
		t.Fatalf("section link provider = %q, want simpleicon", got)
	}
	if got := cfg.Socials[0].IconProvider; got != "lucide" {
		t.Fatalf("social provider = %q, want lucide", got)
	}
}
