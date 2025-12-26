package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mullerhx/dnsbuster/model"
)

func Write(result *model.Result, format string) {
	switch format {
	case "json":
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	case "xml":
		b, _ := xml.MarshalIndent(result, "", "  ")
		fmt.Println(xml.Header + string(b))
	default:
		printText(result)
	}
}

func printText(r *model.Result) {
	fmt.Println("\n[*] Records:")
	for t, v := range r.Records {
		for _, rec := range v {
			fmt.Printf("  [%s] %s\n", t, rec)
		}
	}

	fmt.Println("\n[*] Subdomains:")
	for host, ips := range r.Subdomains {
		for _, ip := range ips {
			fmt.Printf("  [+] %s -> %s\n", host, ip)
		}
	}
}
