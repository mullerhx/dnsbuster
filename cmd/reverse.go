package cmd

import (
	"fmt"
	"github.com/mullerhx/dnsbuster/dns"
	"github.com/spf13/cobra"
)

var cidr string
var rthreads int

var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "Reverse DNS enumeration",
	Run: func(cmd *cobra.Command, args []string) {
		if cidr == "" {
			fmt.Println("CIDR is required")
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
		dns.ReverseLookup(cidr, resolver)
	},
}

func init() {
	rootCmd.AddCommand(reverseCmd)

	reverseCmd.Flags().StringVarP(&cidr, "range", "r", "", "CIDR range (e.g. 8.8.8.0/24)")
	reverseCmd.Flags().IntVarP(&rthreads, "threads", "t", 20, "Concurrent threads")
	reverseCmd.MarkFlagRequired("range")
}
