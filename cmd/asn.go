package cmd

import (
	"fmt"

	"github.com/mullerhx/dnsbuster/asn"
	"github.com/mullerhx/dnsbuster/dns"
	"github.com/mullerhx/dnsbuster/model"
	"github.com/mullerhx/dnsbuster/output"

	"github.com/spf13/cobra"
)

var (
	asnValue   string
	doReverse  bool
	asnThreads int
	asnOutput  string
)

var asnCmd = &cobra.Command{
	Use:   "asn",
	Short: "Expand ASN into IP netblocks",
	Run: func(cmd *cobra.Command, args []string) {
		if asnValue == "" {
			fmt.Println("ASN is required (e.g. AS15169)")
			return
		}

		cidrs, err := asn.LookupASN(asnValue)
		if err != nil {
			fmt.Println("ASN lookup failed:", err)
			return
		}

		result := &model.Result{
			ASN:       asnValue,
			Netblocks: cidrs,
		}

		fmt.Printf("[*] ASN %s → %d netblocks\n", asnValue, len(cidrs))

		if doReverse {
			fmt.Println("[*] Performing reverse DNS on netblocks")
			resolver := dns.NewResolver(asnThreads, false, "", nil)

			for _, cidr := range cidrs {
				dns.ReverseLookup(cidr, resolver)
			}
		}

		output.Write(result, asnOutput)
	},
}

func init() {
	rootCmd.AddCommand(asnCmd)

	asnCmd.Flags().StringVarP(&asnValue, "asn", "a", "", "ASN number (e.g. AS15169)")
	asnCmd.Flags().BoolVar(&doReverse, "reverse", false, "Perform reverse DNS on netblocks")
	asnCmd.Flags().IntVarP(&asnThreads, "threads", "t", 50, "Concurrent threads")
	asnCmd.Flags().StringVar(&asnOutput, "output", "text", "Output format: text|json|xml")

	asnCmd.MarkFlagRequired("asn")
}
