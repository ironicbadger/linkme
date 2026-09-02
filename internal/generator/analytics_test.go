package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ironicbadger/linkme/internal/config"
)

func TestGenerateAnalyticsForBuiltInThemes(t *testing.T) {
	analytics := config.Analytics{
		Google:      &config.GoogleAnalytics{ID: "G-TEST"},
		GoatCounter: &config.GoatCounter{ID: "stats"},
		Plausible:   &config.Plausible{Domain: "example.com", ScriptURL: "https://plausible.example/script.js"},
		Umami: &config.Umami{
			WebsiteID:         "site-123",
			ScriptURL:         "https://umami.example/script.js",
			HostURL:           "https://collector.example",
			RespectDoNotTrack: true,
		},
		Matomo: &config.Matomo{
			URL:            "https://matomo.example",
			SiteID:         "7",
			DisableCookies: true,
		},
	}
	markers := []string{
		"googletagmanager.com/gtag/js?id=G-TEST",
		`data-goatcounter="https://stats.goatcounter.com/count"`,
		`data-domain="example.com"`,
		`data-website-id="site-123"`,
		`data-host-url="https://collector.example"`,
		`data-do-not-track="true"`,
		"_paq.push(['disableCookies'])",
		`var u = "https:\/\/matomo.example\/"`,
		"_paq.push(['setSiteId', '7'])",
	}

	for _, themeName := range []string{"default", "polysleek"} {
		t.Run(themeName, func(t *testing.T) {
			html := generateTestSite(t, themeName, analytics)
			for _, marker := range markers {
				if count := strings.Count(html, marker); count != 1 {
					t.Errorf("marker %q count = %d, want 1", marker, count)
				}
			}
		})
	}
}

func TestGenerateOmitsUnconfiguredAnalytics(t *testing.T) {
	analytics := config.Analytics{
		Umami:  &config.Umami{WebsiteID: "missing-script-url"},
		Matomo: &config.Matomo{URL: "https://matomo.example"},
	}
	markers := []string{
		"googletagmanager.com",
		"goatcounter.com/count",
		"plausible.io",
		"data-website-id",
		"_paq.push",
	}

	for _, themeName := range []string{"default", "polysleek"} {
		t.Run(themeName, func(t *testing.T) {
			html := generateTestSite(t, themeName, analytics)
			for _, marker := range markers {
				if strings.Contains(html, marker) {
					t.Errorf("unconfigured analytics emitted marker %q", marker)
				}
			}
		})
	}
}

func generateTestSite(t *testing.T, themeName string, analytics config.Analytics) string {
	t.Helper()

	cfg := &config.Config{
		Name:      "Test",
		Theme:     themeName,
		Analytics: analytics,
	}
	outputDir := t.TempDir()
	themesDir := filepath.Join("..", "..", "themes")
	gen, err := New(cfg, themesDir, outputDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := gen.Generate(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
