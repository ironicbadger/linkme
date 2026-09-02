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

func TestLoadAnalyticsProviders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	data := []byte(`
analytics:
  umami:
    website_id: site-123
    script_url: https://umami.example/script.js
    host_url: https://collector.example
    respect_do_not_track: true
  matomo:
    url: https://matomo.example///
    site_id: "7"
    disable_cookies: true
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.Analytics.Umami.WebsiteID; got != "site-123" {
		t.Fatalf("Umami website ID = %q, want site-123", got)
	}
	if !cfg.Analytics.Umami.RespectDoNotTrack {
		t.Fatal("Umami respect_do_not_track was not loaded")
	}
	if got := cfg.Analytics.Matomo.BaseURL(); got != "https://matomo.example/" {
		t.Fatalf("Matomo base URL = %q, want https://matomo.example/", got)
	}
	if !cfg.Analytics.Matomo.DisableCookies {
		t.Fatal("Matomo disable_cookies was not loaded")
	}
}
