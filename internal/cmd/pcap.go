package cmd

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

type handler struct{}

func (h *handler) HandlePacket(packet gopacket.Packet) {
	layers := ""
	for _, l := range packet.Layers() {
		if layers != "" {
			layers += ","
		}
		layers += l.LayerType().String()

	}
	fmt.Println(layers)
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
