package cmd

import (
	"fmt"

	"github.com/mullerhx/dnsbuster/dns"

	"github.com/spf13/cobra"
)

var axfrDomain string

var axfrCmd = &cobra.Command{
	Use:   "axfr",
	Short: "Attempt DNS zone transfer",
	Run: func(cmd *cobra.Command, args []string) {
		if axfrDomain == "" {
			fmt.Println("Domain is required")
			return
		}
		dns.TryAXFR(axfrDomain)
	},
}

func init() {
	rootCmd.AddCommand(axfrCmd)
	axfrCmd.Flags().StringVarP(&axfrDomain, "domain", "d", "", "Target domain")
	axfrCmd.MarkFlagRequired("domain")
}
