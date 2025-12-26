package dns

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

var dohProviders = map[string]string{
	"cloudflare": "https://cloudflare-dns.com/dns-query",
	"google":     "https://dns.google/dns-query",
	"quad9":      "https://dns.quad9.net/dns-query",
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// DoHLookup performs a DNS A lookup over HTTPS (RFC 8484)
func DoHLookup(host, provider string) ([]string, error) {
	url, ok := dohProviders[provider]
	if !ok {
		return nil, errors.New("unknown DoH provider")
	}

	// Build DNS query
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), dns.TypeA)
	msg.RecursionDesired = true

	wire, err := msg.Pack()
	if err != nil {
		return nil, err
	}

	// Use POST (more widely supported)
	req, err := http.NewRequest("POST", url, bytes.NewReader(wire))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("DoH query failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reply := new(dns.Msg)
	if err := reply.Unpack(body); err != nil {
		return nil, err
	}

	var ips []string
	for _, ans := range reply.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}

	return ips, nil
}
