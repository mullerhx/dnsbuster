package dns

import (
	"math/rand"
	"net"
	"time"
)

func HasWildcard(domain string) bool {
	rand.Seed(time.Now().UnixNano())
	test := randomString(12) + "." + domain
	_, err := net.LookupHost(test)
	return err == nil
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func DetectWildcardIPs(domain string) map[string]bool {
	test := randomString(12) + "." + domain
	ips, err := net.LookupHost(test)

	m := make(map[string]bool)
	if err != nil {
		return m
	}

	for _, ip := range ips {
		m[ip] = true
	}
	return m
}
