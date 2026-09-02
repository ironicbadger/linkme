# linkme

A customizable link page generator built in Go.

![A screenshot of the application.](assets/screenshot.png)

## Usage

```bash
# Build the static site
go run ./cmd/linkme build

# Serve locally
go run ./cmd/linkme serve
```

## Docker

```bash
docker run -d -p 8080:80 ghcr.io/ironicbadger/linkme:latest
```

## Configuration

Edit `config/config.yml` to customize your links and appearance.

### Icons

Link and social icons use [Simple Icons](https://simpleicons.org/) by default. Set
`icon-provider: "lucide"` to use a generic [Lucide](https://lucide.dev/icons/)
icon instead:

```yaml
links:
  - title: "My blog"
    url: "https://example.com"
    icon: "rss"
    icon-provider: "lucide"
    color: "#404040"
```

An omitted or empty provider selects Simple Icons. Unknown providers and icon
names render a placeholder.

### Analytics

Supported: Google Analytics, GoatCounter, Plausible, Umami, and Matomo.

Example config:

```yaml
analytics:
  google:
    id: "G-XXXXXXX"
  goatcounter:
    id: "example"
    selfhosted: false
  plausible:
    domain: "example.com" # tracked site
    script_url: "" # optional; set to your self-hosted instance URL (e.g. https://plausible.example.com/js/script.js)
  umami:
    website_id: "00000000-0000-0000-0000-000000000000"
    script_url: "https://cloud.umami.is/script.js"
    host_url: "" # optional alternate collection endpoint
    respect_do_not_track: true
  matomo:
    url: "https://analytics.example.com"
    site_id: "1"
    disable_cookies: true
```

GoatCounter:
If `selfhosted: true`, `id` must be the full host (FQDN) of your instance. Otherwise `id` is used as the subdomain on `goatcounter.com`.

Plausible:
Set `domain` to the site you want to track. For self-hosted, set `script_url` to your instance’s script URL.

Umami:
Both `website_id` and `script_url` are required. `host_url` is optional and
overrides the collection endpoint. `respect_do_not_track` emits Umami's
Do Not Track setting.

Matomo:
Both `url` and `site_id` are required. Set `disable_cookies: true` for
cookieless tracking.

Analytics are disabled when their required values are empty. Configuring more
than one provider sends page-view data to each provider. Linkme does not add a
consent interface; choose settings appropriate for your privacy obligations.
Analytics IDs are public identifiers embedded in the generated HTML, not secrets.
