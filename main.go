package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func main() {
	domain := flag.String("d", "", "domain to check")
	flag.Parse()

	checkDomain(*domain)
}

func checkDomain(domain string) {
	has_mx := mx(domain)
	has_spf := txt(domain)
	has_dmarc := dmarc(domain)

	headerFmt := color.New(color.FgHiCyan, color.Underline).SprintfFunc()

	tbl := table.New("mx", "spf", "dmarc")
	tbl.WithHeaderFormatter(headerFmt)

	tbl.AddRow(has_mx, has_spf, has_dmarc)

	tbl.Print()
}

func mx(domain string) bool {
	mx, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Err: %v\n", err)
	}
	if len(mx) > 0 {
		return true
	} else {
		return false
	}

}

func txt(domain string) bool {
	txt, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Err: %v\n", err)
	}

	var has_mx bool
	for _, r := range txt {
		if strings.HasPrefix(r, "v=spf1") {
			has_mx = true
		} else {
			has_mx = false
		}
	}

	return has_mx
}

func dmarc(domain string) bool {
	dmarc, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Err: %v\n", err)
	}

	var has_dmarc bool
	for _, r := range dmarc {
		if strings.HasPrefix(r, "v=DMARC1") {
			has_dmarc = true
		} else {
			has_dmarc = false
		}
	}

	return has_dmarc
}
