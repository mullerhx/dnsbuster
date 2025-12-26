package dns

import "fmt"

func Bruteforce(domain string, subs []string, r *Resolver, wildcard bool) {
	for _, sub := range subs {
		host := sub + "." + domain
		r.ResolveHost(host)
	}
	r.Wait()

	if wildcard {
		fmt.Println("[!] Results may include wildcard responses")
	}
}
