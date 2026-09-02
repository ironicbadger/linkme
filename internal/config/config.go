package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigVersion string     `yaml:"configVersion"`
	Name          string     `yaml:"name"`
	Subtitle      string     `yaml:"subtitle"`
	Description   string     `yaml:"description"`
	Avatar        string     `yaml:"avatar"`
	Favicon       string     `yaml:"favicon"`
	Theme         string     `yaml:"theme"`
	Background    Background `yaml:"background"`
	Links         []Link     `yaml:"links"`
	Sections      []Section  `yaml:"sections"`
	Socials       []Social   `yaml:"socials"`
	Analytics     Analytics  `yaml:"analytics"`
}

type Analytics struct {
	Google      *GoogleAnalytics `yaml:"google"`
	GoatCounter *GoatCounter     `yaml:"goatcounter"`
	Plausible   *Plausible       `yaml:"plausible"`
	Umami       *Umami           `yaml:"umami"`
	Matomo      *Matomo          `yaml:"matomo"`
}

type GoogleAnalytics struct {
	ID string `yaml:"id"`
}

type GoatCounter struct {
	ID         string `yaml:"id"`
	Selfhosted bool   `yaml:"selfhosted"`
}

type Plausible struct {
	Domain    string `yaml:"domain"`
	ScriptURL string `yaml:"script_url"`
}

type Umami struct {
	WebsiteID         string `yaml:"website_id"`
	ScriptURL         string `yaml:"script_url"`
	HostURL           string `yaml:"host_url"`
	RespectDoNotTrack bool   `yaml:"respect_do_not_track"`
}

type Matomo struct {
	URL            string `yaml:"url"`
	SiteID         string `yaml:"site_id"`
	DisableCookies bool   `yaml:"disable_cookies"`
}

func (m *Matomo) BaseURL() string {
	return strings.TrimRight(strings.TrimSpace(m.URL), "/") + "/"
}

type Background struct {
	Type    string  `yaml:"type"`    // color, image, gradient, particles
	Value   string  `yaml:"value"`   // color code, image path, or gradient definition
	Blur    int     `yaml:"blur"`    // blur amount for images
	Opacity float64 `yaml:"opacity"` // overlay opacity
}

type Link struct {
	Title        string `yaml:"title"`
	URL          string `yaml:"url"`
	Icon         string `yaml:"icon"`          // Simple Icons slug
	IconURL      string `yaml:"iconUrl"`       // Custom icon URL (overrides icon)
	IconProvider string `yaml:"icon-provider"` // simpleicon (default) or lucide
	Color        string `yaml:"color"`         // hex color for button
}

type Section struct {
	Title string `yaml:"title"`
	Links []Link `yaml:"links"`
}

type Social struct {
	Icon         string `yaml:"icon"`          // Simple Icons slug
	IconProvider string `yaml:"icon-provider"` // simpleicon (default) or lucide
	URL          string `yaml:"url"`
	Color        string `yaml:"color"` // hex color for icon
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.ConfigVersion == "" {
		cfg.ConfigVersion = "1.0"
	}
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}
	if cfg.Background.Type == "" {
		cfg.Background.Type = "color"
		cfg.Background.Value = "#1e1f26"
	}
	if cfg.Background.Opacity == 0 {
		cfg.Background.Opacity = 1.0
	}

	return &cfg, nil
}
