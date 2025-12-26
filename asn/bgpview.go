package asn

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type bgpViewResponse struct {
	Data struct {
		IPv4Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"ipv4_prefixes"`
		IPv6Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"ipv6_prefixes"`
	} `json:"data"`
}

// Hardcoded IPs for api.bgpview.io (Cloudflare-backed)
var bgpViewIPs = []string{
	"104.26.12.93",
	"104.26.13.93",
	"172.67.74.234",
}

func LookupASN(asn string) ([]string, error) {
	asn = strings.ToUpper(strings.TrimSpace(asn))
	asn = strings.TrimPrefix(asn, "AS")

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{Timeout: 10 * time.Second}
			// connect directly to IP:443
			return dialer.DialContext(ctx, network, bgpViewIPs[0]+":443")
		},
		TLSClientConfig: &tls.Config{
			ServerName: "api.bgpview.io", // 🔥 keep SNI correct
		},
	}

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	url := fmt.Sprintf("https://api.bgpview.io/asn/%s/prefixes", asn)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bgpview API returned %s", resp.Status)
	}

	var r bgpViewResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	var prefixes []string
	for _, p := range r.Data.IPv4Prefixes {
		prefixes = append(prefixes, p.Prefix)
	}
	for _, p := range r.Data.IPv6Prefixes {
		prefixes = append(prefixes, p.Prefix)
	}

	return prefixes, nil
}
