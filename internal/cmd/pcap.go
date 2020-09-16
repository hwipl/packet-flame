package cmd

import (
	"time"

	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

type handler struct{}

func (h *handler) HandlePacket(packet gopacket.Packet) {
	layers := ""
	transport := packet.TransportLayer()
	for _, l := range packet.Layers() {
		if layers != "" {
			layers += ";"
		}

		// add layer name
		layers += l.LayerType().String()

		// add port information for the transport layer
		if transport != nil && l.LayerType() == transport.LayerType() {
			flow := transport.TransportFlow()
			src, dst := flow.Endpoints()
			if src.LessThan(dst) {
				layers += ";Port_" + src.String()
				layers += ";Port_" + dst.String()
			} else {
				layers += ";Port_" + dst.String()
				layers += ";Port_" + src.String()
			}
		}

	}
	counter[layers]++
}

// listen captures packets on the network interface or reads them from a pcap
// file and parses them
func listen() {
	// create handler
	var handler handler

	// create listener
	listener := pcap.Listener{
		PacketHandler: &handler,
		File:          pcapFile,
		Device:        pcapDevice,
		Promisc:       pcapPromisc,
		Snaplen:       pcapSnaplen,
		Timeout:       time.Duration(pcapTimeout) * time.Millisecond,
		Filter:        pcapFilter,
		MaxPkts:       pcapMaxPkts,
		MaxTime:       time.Duration(pcapMaxTime) * time.Second,
	}

	// start listen loop
	listener.Prepare()
	listener.Loop()
}
