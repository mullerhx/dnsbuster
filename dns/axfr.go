package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

func TryAXFR(domain string) {
	nss, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return
	}

	for _, server := range nss.Servers {
		t := new(dns.Transfer)
		m := new(dns.Msg)
		m.SetAxfr(domain)

		ch, err := t.In(m, server+":53")
		if err != nil {
			continue
		}

		for env := range ch {
			if env.Error != nil {
				continue
			}
			fmt.Printf("[AXFR] %s\n", server)
			for _, rr := range env.RR {
				fmt.Println(rr.String())
			}
		}
	}
}
