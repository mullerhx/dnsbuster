package model

type Result struct {
	Domain     string              `json:"domain,omitempty" xml:"domain,omitempty"`
	ASN        string              `json:"asn,omitempty" xml:"asn,omitempty"`
	Netblocks  []string            `json:"netblocks,omitempty" xml:"netblocks>cidr,omitempty"`
	Wildcard   bool                `json:"wildcard,omitempty" xml:"wildcard,omitempty"`
	Records    map[string][]string `json:"records,omitempty" xml:"records>record,omitempty"`
	Subdomains map[string][]string `json:"subdomains,omitempty" xml:"subdomains>host,omitempty"`
}
