package caddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type Site struct {
	Domain      string `json:"domain"`
	ServicePort int    `json:"service_port"`
	TLSEnabled  bool   `json:"tls_enabled"`
	TLEmail     string `json:"tls_email,omitempty"`
	ExtraConfig string `json:"extra_config,omitempty"`
	Type        string `json:"type"`
	Root        string `json:"root,omitempty"`
}

type CaddyApp struct {
	Apps struct {
		HTTP struct {
			Servers map[string]CaddyServer `json:"servers"`
		} `json:"http"`
	} `json:"apps"`
}

type CaddyServer struct {
	Listen []string        `json:"listen"`
	Routes []CaddyRoute    `json:"routes"`
}

type CaddyRoute struct {
	Handle []CaddyHandler `json:"handle"`
	Match  []CaddyMatch   `json:"match,omitempty"`
}

type CaddyMatch struct {
	Host []string `json:"host"`
}

type CaddyHandler struct {
	Handler   string            `json:"handler"`
	Upstreams []CaddyUpstream   `json:"upstreams,omitempty"`
	Root      string            `json:"root,omitempty"`
}

type CaddyUpstream struct {
	Dial string `json:"dial"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{},
	}
}

func (c *Client) request(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

func (c *Client) GetConfig() (*CaddyApp, error) {
	resp, err := c.request("GET", "/config/", nil)
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}
	defer resp.Body.Close()
	var app CaddyApp
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &app, nil
}

func (c *Client) ListSites() ([]Site, error) {
	cfg, err := c.GetConfig()
	if err != nil {
		return nil, err
	}
	sites := make([]Site, 0)
	for _, srv := range cfg.Apps.HTTP.Servers {
		for _, route := range srv.Routes {
			for _, m := range route.Match {
				for _, host := range m.Host {
					site := Site{Domain: host, Type: "proxy"}
					for _, h := range route.Handle {
						if h.Handler == "reverse_proxy" && len(h.Upstreams) > 0 {
							fmt.Sscanf(h.Upstreams[0].Dial, "localhost:%d", &site.ServicePort)
						}
						if h.Handler == "file_server" && h.Root != "" {
							site.Type = "static"
							site.Root = h.Root
						}
					}
					sites = append(sites, site)
				}
			}
		}
	}
	return sites, nil
}

func (c *Client) serverName() (string, error) {
	resp, err := c.request("GET", "/config/apps/http/servers/", nil)
	if err != nil {
		return "", fmt.Errorf("get servers: %w", err)
	}
	defer resp.Body.Close()
	var servers map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return "", fmt.Errorf("decode servers: %w", err)
	}
	for name := range servers {
		return name, nil
	}
	return "", fmt.Errorf("no servers found in caddy config")
}

func (c *Client) AddSite(site Site) error {
	var handle []map[string]interface{}

	if site.Type == "static" {
		handle = []map[string]interface{}{
			{
				"handler": "file_server",
				"root":    site.Root,
			},
		}
	} else {
		handle = []map[string]interface{}{
			{
				"handler": "reverse_proxy",
				"upstreams": []map[string]interface{}{
					{"dial": fmt.Sprintf("localhost:%d", site.ServicePort)},
				},
			},
		}
	}

	config := map[string]interface{}{
		"@id":    site.Domain,
		"match":  []map[string]interface{}{{"host": []string{site.Domain}}},
		"handle": handle,
	}
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	srv, err := c.serverName()
	if err != nil {
		return fmt.Errorf("add site: %w", err)
	}
	resp, err := c.request("POST", "/config/apps/http/servers/"+srv+"/routes/", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("add site: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		// If error is due to missing 'routes' key (invalid traversal path), initialize routes as an empty array and retry
		if resp.StatusCode == 500 && strings.Contains(string(respBody), "invalid traversal path") && strings.Contains(string(respBody), "routes") {
			initResp, initErr := c.request("PUT", "/config/apps/http/servers/"+srv+"/routes", strings.NewReader("[]"))
			if initErr == nil {
				initResp.Body.Close()
				if initResp.StatusCode < 400 {
					// Retry original POST
					retryResp, retryErr := c.request("POST", "/config/apps/http/servers/"+srv+"/routes/", bytes.NewReader(body))
					if retryErr == nil {
						defer retryResp.Body.Close()
						if retryResp.StatusCode < 400 {
							return nil
						}
						respBody, _ = io.ReadAll(retryResp.Body)
						return fmt.Errorf("caddy api error %d (retry): %s", retryResp.StatusCode, string(respBody))
					}
				}
			}
		}
		return fmt.Errorf("caddy api error %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *Client) RemoveSite(domain string) error {
	path := fmt.Sprintf("/id/%s", domain)
	resp, err := c.request("DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("remove site: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("caddy api error %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *Client) Health() error {
	resp, err := c.request("GET", "/config/", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
