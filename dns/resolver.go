package dns

import (
	"fmt"
	"github.com/mullerhx/dnsbuster/model"
	"net"
	"sync"
)

type Resolver struct {
	sem         chan struct{}
	wg          sync.WaitGroup
	UseDoH      bool
	Provider    string
	Result      *model.Result
	WildcardIPs map[string]bool
}

func NewResolver(threads int, useDoH bool, provider string, wildcardIPs map[string]bool) *Resolver {
	return &Resolver{
		sem:      make(chan struct{}, threads),
		UseDoH:   useDoH,
		Provider: provider,
		Result: &model.Result{
			Records:    make(map[string][]string),
			Subdomains: make(map[string][]string),
		},
		WildcardIPs: wildcardIPs,
	}
}

func (r *Resolver) ResolveDomain(domain string) {
	records := []string{"A", "AAAA", "MX", "NS", "TXT"}

	for _, t := range records {
		ips, err := net.LookupHost(domain)
		if err == nil && (t == "A" || t == "AAAA") {
			for _, ip := range ips {
				fmt.Printf("  [%s] %s\n", t, ip)
			}
		}
	}

	if mx, err := net.LookupMX(domain); err == nil {
		for _, m := range mx {
			fmt.Printf("  [MX] %s\n", m.Host)
		}
	}

	if ns, err := net.LookupNS(domain); err == nil {
		for _, n := range ns {
			fmt.Printf("  [NS] %s\n", n.Host)
		}
	}
}

func (r *Resolver) ResolveHost(host string) {
	r.sem <- struct{}{}
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		defer func() { <-r.sem }()

		var ips []string
		var err error

		if r.UseDoH {
			ips, err = DoHLookup(host, r.Provider)
		} else {
			ips, err = net.LookupHost(host)
		}

		if err != nil {
			return
		}

		// 🔥 Automatic wildcard filtering
		for _, ip := range ips {
			if r.WildcardIPs[ip] {
				continue
			}
			r.Result.Subdomains[host] = append(r.Result.Subdomains[host], ip)
		}
	}()
}

func (r *Resolver) Wait() {
	r.wg.Wait()
}
