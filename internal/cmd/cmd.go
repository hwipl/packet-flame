package cmd

import (
	"flag"
	"fmt"
)

var (
	// pcap settings
	pcapDevice  string
	pcapFile    string
	pcapPromisc bool = true
	pcapSnaplen int  = 1024
	pcapTimeout int  = 0
	pcapMaxPkts int  = 0
	pcapMaxTime int  = 0
	pcapFilter  string

	// display settings
	showMACs       bool
	showIPs        bool

	// packet layer packet/byte counter
	counter = make(map[string]int)
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// pcap settings
	flag.StringVar(&pcapFile, "f", pcapFile,
		"read packets from a pcap file and set it to `file`")
	flag.StringVar(&pcapDevice, "i", pcapDevice, "read packets from "+
		"a network interface (default) and set it to `interface`")
	flag.BoolVar(&pcapPromisc, "pcap-promisc", pcapPromisc,
		"set network interface to promiscuous mode")
	flag.IntVar(&pcapSnaplen, "pcap-snaplen", pcapSnaplen,
		"set pcap timeout to `milliseconds`")
	flag.IntVar(&pcapTimeout, "pcap-timeout", pcapTimeout,
		"set pcap timeout to `milliseconds`")
	flag.IntVar(&pcapMaxPkts, "pcap-maxpkts", pcapMaxPkts,
		"set maximum packets to capture to `number` "+
			"(may require pcap-timeout argument)")
	flag.IntVar(&pcapMaxTime, "pcap-maxtime", pcapMaxTime,
		"set maximum capturing time to `seconds` "+
			"(may require pcap-timeout argument)")
	flag.StringVar(&pcapFilter, "pcap-filter", pcapFilter,
		"set pcap packet filter to `filter` (e.g.: \"not port 22\")")

	// display settings
	flag.BoolVar(&showMACs, "with-macs", showMACs,
		"show MAC addresses in output")
	flag.BoolVar(&showIPs, "with-ips", showIPs,
		"show IP addresses in output")

	// parse and overwrite default values of settings
	flag.Parse()
}

// printCounter prints all packet layers and their count
func printCounter() {
	for l, c := range counter {
		fmt.Println(l, c)
	}
}

// Run is the main entry point
func Run() {
	parseCommandLine()
	listen()
	printCounter()
}
