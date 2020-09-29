package cmd

import (
	"flag"
	"fmt"
	"strings"
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
	showPorts      bool
	showDirections bool

	// packet layer packet/byte counter
	counter = newCountMap()

	// http settings
	httpListen string

	// flamegraph tool settings
	flamegraphExe  string
	flamegraphArgs []string
)

// parseFlamegraphArgs parses arguments for the flamegraph tool and constructs
// an argument list containing options beginning with "-" and other texts like
// the following example:
// ["--opt1", "text", "--opt2", "-opt3", "this is another text"]
func parseFlamegraphArgs(fgArgs string) []string {
	args := []string{}
	text := ""
	for _, s := range strings.Fields(fgArgs) {
		if strings.HasPrefix(s, "-") {
			// add concatenated text as separate argument
			if text != "" {
				args = append(args, text)
				text = ""
			}

			// add option beginning with "-" as separate argument
			args = append(args, s)
			continue
		}

		// concatenate text until new option beginning with "-"
		// is found (see above) or we reach the end (see below)
		if text != "" {
			text += " "
		}
		text += s
	}

	// add remaining text as separate argument
	if text != "" {
		args = append(args, text)
	}

	return args
}

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
	flag.BoolVar(&showPorts, "with-ports", showPorts,
		"show ports in output")
	flag.BoolVar(&showDirections, "with-directions", showDirections,
		"show flow directions in MAC, IP, port output")

	// http settings
	flag.StringVar(&httpListen, "http", httpListen,
		"use http server and set the listen `address` (e.g.: :8000)")

	// flamegraph tool settings
	flag.StringVar(&flamegraphExe, "fg-exe", flamegraphExe,
		"set the flamegraph `executable` for graph generation "+
			"(e.g.: flamegraph.pl)")
	fgArgs := flag.String("fg-args", "",
		"set `arguments` for the flamegraph executable "+
			"(e.g.: \"--title test graph\")")

	// parse and overwrite default values of settings
	flag.Parse()

	// parse flamegraph tool arguments
	if *fgArgs != "" {
		flamegraphArgs = parseFlamegraphArgs(*fgArgs)
	}
}

// printCounter prints all packet layers and their count
func printCounter() {
	for l, c := range counter.getAll() {
		fmt.Println(l, c)
	}
}

// Run is the main entry point
func Run() {
	parseCommandLine()
	if httpListen != "" {
		go runHTTPServer(httpListen)
	}
	listen()
	printCounter()
}
