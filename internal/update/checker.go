package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const AppVersion = "1.0.0"

const repoOwner = "aevynt"
const repoName = "goPanel"

type Asset struct {
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Assets      []Asset   `json:"assets"`
}

type CheckResult struct {
	CurrentVersion string   `json:"current_version"`
	LatestVersion  string   `json:"latest_version"`
	HasUpdate      bool     `json:"has_update"`
	Release        *Release `json:"release,omitempty"`
	CheckedAt      string   `json:"checked_at"`
	Error          string   `json:"error,omitempty"`
}

type Checker struct {
	mu       sync.Mutex
	cache    *CheckResult
	cacheTTL time.Duration
}

func NewChecker() *Checker {
	return &Checker{cacheTTL: 1 * time.Hour}
}

func (c *Checker) Check() *CheckResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if c.cache != nil {
		checked, err := time.Parse(time.RFC3339, c.cache.CheckedAt)
		if err == nil && now.Sub(checked) < c.cacheTTL {
			return c.cache
		}
	}
	result := c.fetchLatest()
	result.CheckedAt = now.Format(time.RFC3339)
	c.cache = result
	return result
}

func (c *Checker) fetchLatest() *CheckResult {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &CheckResult{CurrentVersion: AppVersion, Error: err.Error(), CheckedAt: time.Now().Format(time.RFC3339)}
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "goPanel")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &CheckResult{CurrentVersion: AppVersion, Error: err.Error(), CheckedAt: time.Now().Format(time.RFC3339)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &CheckResult{CurrentVersion: AppVersion, Error: fmt.Sprintf("GitHub API returned %d", resp.StatusCode), CheckedAt: time.Now().Format(time.RFC3339)}
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return &CheckResult{CurrentVersion: AppVersion, Error: err.Error(), CheckedAt: time.Now().Format(time.RFC3339)}
	}

	current := strings.TrimPrefix(AppVersion, "v")
	latest := strings.TrimPrefix(release.TagName, "v")

	return &CheckResult{
		CurrentVersion: AppVersion,
		LatestVersion:  release.TagName,
		HasUpdate:      compareVersions(latest, current) > 0,
		Release:        &release,
	}
}

func compareVersions(a, b string) int {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		var ai, bi int
		n, err := fmt.Sscanf(as[i], "%d", &ai)
		if n != 1 || err != nil {
			return strings.Compare(a, b)
		}
		n, err = fmt.Sscanf(bs[i], "%d", &bi)
		if n != 1 || err != nil {
			return strings.Compare(a, b)
		}
		if ai != bi {
			return ai - bi
		}
	}
	return len(as) - len(bs)
}
