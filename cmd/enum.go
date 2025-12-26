package cmd

import (
	"bufio"
	"fmt"
	"github.com/mullerhx/dnsbuster/output"
	"os"
	"strings"

	"github.com/mullerhx/dnsbuster/dns"

	"github.com/spf13/cobra"
)

var (
	outputFmt string
	useDoH    bool
	dohProv   string
)

var (
	domain   string
	wordlist string
	threads  int
	axfr     bool
)

var enumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Enumerate DNS records and subdomains",
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" {
			fmt.Println("Domain is required")
			return
		}

		wildcardIPs := dns.DetectWildcardIPs(domain)

		resolver := dns.NewResolver(
			threads,
			useDoH,
			dohProv,
			wildcardIPs,
		)

		resolver.Result.Domain = domain
		resolver.Result.Wildcard = len(wildcardIPs) > 0

		fmt.Println("[*] Resolving base records")
		resolver.ResolveDomain(domain)

		fmt.Println("[*] Checking wildcard DNS")
		wildcard := dns.HasWildcard(domain)
		fmt.Printf("    Wildcard DNS: %v\n", wildcard)

		if axfr {
			fmt.Println("[*] Testing AXFR")
			dns.TryAXFR(domain)
		}

		if wordlist != "" {
			fmt.Println("[*] Brute forcing subdomains")
			file, err := os.Open(wordlist)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			var subs []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				sub := strings.TrimSpace(scanner.Text())
				if sub != "" {
					subs = append(subs, sub)
				}
			}

			dns.Bruteforce(domain, subs, resolver, wildcard)
			output.Write(resolver.Result, outputFmt)
		}
	},
}

func init() {
	rootCmd.AddCommand(enumCmd)

	enumCmd.Flags().StringVarP(&domain, "domain", "d", "", "Target domain")
	enumCmd.Flags().StringVarP(&wordlist, "wordlist", "w", "", "Subdomain wordlist")
	enumCmd.Flags().IntVarP(&threads, "threads", "t", 20, "Number of concurrent threads")
	enumCmd.Flags().BoolVar(&axfr, "axfr", false, "Attempt zone transfer")
	enumCmd.Flags().StringVar(&outputFmt, "output", "text", "Output format: text|json|xml")
	enumCmd.Flags().BoolVar(&useDoH, "doh", false, "Use DNS over HTTPS")
	enumCmd.Flags().StringVar(&dohProv, "doh-provider", "cloudflare", "DoH provider")

	enumCmd.MarkFlagRequired("domain")
}
