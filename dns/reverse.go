package dns

import (
	"fmt"
	"net"
)

func ReverseLookup(cidr string, r *Resolver) {
	ips, err := expandCIDR(cidr)
	if err != nil {
		return
	}

	for _, ip := range ips {
		r.sem <- struct{}{}
		r.wg.Add(1)

		go func(ip string) {
			defer r.wg.Done()
			defer func() { <-r.sem }()
			names, err := net.LookupAddr(ip)
			if err == nil {
				for _, n := range names {
					fmt.Printf("[PTR] %s -> %s\n", ip, n)
				}
			}
		}(ip)
	}

	r.Wait()
}

func expandCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
