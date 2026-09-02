package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ironicbadger/linkme/internal/config"
	"github.com/ironicbadger/linkme/internal/theme"
)

type Generator struct {
	cfg       *config.Config
	themePath string
	outputDir string
	theme     *theme.ThemeManifest
}

type TemplateData struct {
	Config          *config.Config
	Theme           *theme.ThemeManifest
	Links           []LinkData
	Sections        []SectionData
	Socials         []SocialData
	BackgroundCSS   template.CSS
	DescriptionHTML template.HTML
	ThemeStyles     []string // CSS files to include
	ThemeScripts    []string // JS files to include
	BuildDate       string   // Date when the site was generated
	AssetVersion    string   // Cache-busting version for theme assets
	Analytics       *config.Analytics
}

type SectionData struct {
	Title string
	Links []LinkData
}

type LinkData struct {
	Title   string
	URL     string
	IconSVG template.HTML
	IconURL string
	Color   string
}

type SocialData struct {
	URL     string
	IconSVG template.HTML
	Color   string
}

func New(cfg *config.Config, themesDir, outputDir string) (*Generator, error) {
	themePath := filepath.Join(themesDir, cfg.Theme)

	// Load theme manifest
	themeManifest, err := theme.Load(themePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load theme manifest: %w", err)
	}

	return &Generator{
		cfg:       cfg,
		themePath: themePath,
		outputDir: outputDir,
		theme:     themeManifest,
	}, nil
}

func (g *Generator) Generate() error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Load and execute template
	tmplPath := filepath.Join(g.themePath, "template.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	data := g.prepareTemplateData()

	// Render HTML
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write index.html
	indexPath := filepath.Join(g.outputDir, "index.html")
	if err := os.WriteFile(indexPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %w", err)
	}

	// Copy theme assets (CSS, JS, etc.)
	if err := g.copyThemeAssets(); err != nil {
		return fmt.Errorf("failed to copy theme assets: %w", err)
	}

	// Copy user assets (avatar, etc.)
	if err := g.copyUserAssets(); err != nil {
		return fmt.Errorf("failed to copy user assets: %w", err)
	}

	return nil
}

func (g *Generator) prepareTemplateData() *TemplateData {
	// Convert newlines in description to <br> tags
	descHTML := template.HTML(strings.ReplaceAll(template.HTMLEscapeString(g.cfg.Description), "\n", "<br>"))
	builtAt := time.Now()

	data := &TemplateData{
		Config:          g.cfg,
		Theme:           g.theme,
		BackgroundCSS:   g.generateBackgroundCSS(),
		DescriptionHTML: descHTML,
		ThemeStyles:     g.theme.Styles,
		ThemeScripts:    g.theme.Scripts,
		BuildDate:       builtAt.Format("Jan 2, 2006"),
		AssetVersion:    builtAt.UTC().Format("20060102150405"),
		Analytics:       &g.cfg.Analytics,
	}

	// Prepare links with SVG icons
	for _, link := range g.cfg.Links {
		data.Links = append(data.Links, LinkData{
			Title:   link.Title,
			URL:     link.URL,
			IconSVG: template.HTML(GetIconSVG(link.Icon, link.IconProvider)),
			IconURL: link.IconURL,
			Color:   link.Color,
		})
	}

	// Prepare sections with their links
	for _, section := range g.cfg.Sections {
		sectionData := SectionData{
			Title: section.Title,
		}
		for _, link := range section.Links {
			sectionData.Links = append(sectionData.Links, LinkData{
				Title:   link.Title,
				URL:     link.URL,
				IconSVG: template.HTML(GetIconSVG(link.Icon, link.IconProvider)),
				IconURL: link.IconURL,
				Color:   link.Color,
			})
		}
		data.Sections = append(data.Sections, sectionData)
	}

	// Prepare social links with SVG icons
	for _, social := range g.cfg.Socials {
		data.Socials = append(data.Socials, SocialData{
			URL:     social.URL,
			IconSVG: template.HTML(GetIconSVG(social.Icon, social.IconProvider)),
			Color:   social.Color,
		})
	}

	return data
}

func (g *Generator) generateBackgroundCSS() template.CSS {
	bg := g.cfg.Background
	var css string

	switch bg.Type {
	case "color":
		css = fmt.Sprintf("background-color: %s;", bg.Value)
	case "image":
		css = fmt.Sprintf(`background-image: url('%s');
			background-size: cover;
			background-position: center;
			background-repeat: no-repeat;`, bg.Value)
		if bg.Blur > 0 {
			css += fmt.Sprintf(" filter: blur(%dpx);", bg.Blur)
		}
	case "gradient":
		css = fmt.Sprintf("background: %s;", bg.Value)
	case "particles":
		css = fmt.Sprintf("background-color: %s;", bg.Value)
	default:
		css = "background-color: #1e1f26;"
	}

	return template.CSS(css)
}

func (g *Generator) copyThemeAssets() error {
	// Copy styles directory if it exists
	stylesDir := filepath.Join(g.themePath, "styles")
	if _, err := os.Stat(stylesDir); err == nil {
		dstStyles := filepath.Join(g.outputDir, "styles")
		if err := copyDir(stylesDir, dstStyles); err != nil {
			return err
		}
	}

	// Copy scripts directory if it exists
	scriptsDir := filepath.Join(g.themePath, "scripts")
	if _, err := os.Stat(scriptsDir); err == nil {
		dstScripts := filepath.Join(g.outputDir, "scripts")
		if err := copyDir(scriptsDir, dstScripts); err != nil {
			return err
		}
	}

	// Copy assets directory if it exists
	assetsDir := filepath.Join(g.themePath, "assets")
	if _, err := os.Stat(assetsDir); err == nil {
		dstAssets := filepath.Join(g.outputDir, "assets")
		if err := copyDir(assetsDir, dstAssets); err != nil {
			return err
		}
	}

	// Backward compatibility: copy root-level styles.css
	srcCSS := filepath.Join(g.themePath, "styles.css")
	if _, err := os.Stat(srcCSS); err == nil {
		dstCSS := filepath.Join(g.outputDir, "styles.css")
		if err := copyFile(srcCSS, dstCSS); err != nil {
			return err
		}
	}

	// Backward compatibility: copy any root-level JS files
	entries, err := os.ReadDir(g.themePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".js") {
			src := filepath.Join(g.themePath, entry.Name())
			dst := filepath.Join(g.outputDir, entry.Name())
			if err := copyFile(src, dst); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) copyUserAssets() error {
	// Copy avatar if specified and it's a local file (not a URL)
	if g.cfg.Avatar != "" && !strings.HasPrefix(g.cfg.Avatar, "http://") && !strings.HasPrefix(g.cfg.Avatar, "https://") {
		srcAvatar := filepath.Join("assets", g.cfg.Avatar)
		if _, err := os.Stat(srcAvatar); err == nil {
			dstAvatar := filepath.Join(g.outputDir, g.cfg.Avatar)
			if err := copyFile(srcAvatar, dstAvatar); err != nil {
				return err
			}
		}
	}

	if g.cfg.Favicon != "" && !strings.HasPrefix(g.cfg.Favicon, "http://") && !strings.HasPrefix(g.cfg.Favicon, "https://") {
		srcFavicon := filepath.Join("assets", g.cfg.Favicon)
		if _, err := os.Stat(srcFavicon); err == nil {
			dstFavicon := filepath.Join(g.outputDir, g.cfg.Favicon)
			if err := copyFile(srcFavicon, dstFavicon); err != nil {
				return err
			}
		}
	}

	// Copy user assets directory if it exists
	userAssetsDir := "assets"
	if _, err := os.Stat(userAssetsDir); err == nil {
		dstAssets := filepath.Join(g.outputDir, "assets")
		if err := copyDir(userAssetsDir, dstAssets); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
