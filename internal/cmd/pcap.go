package cmd

import (
	"time"

	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

type handler struct{}

func getFlowAddresses(flow gopacket.Flow, prefix string) string {
	layers := ""
	src, dst := flow.Endpoints()
	if src.LessThan(dst) {
		layers += ";" + prefix + "_" + src.String()
		layers += ";" + prefix + "_" + dst.String()
	} else {
		layers += ";" + prefix + "_" + dst.String()
		layers += ";" + prefix + "_" + src.String()
	}

	// add flow direction information
	layers += ";" + prefix + "_" + src.String() + "->" + prefix + "_" +
		dst.String()

	return layers
}

func (h *handler) HandlePacket(packet gopacket.Packet) {
	layers := ""
	link := packet.LinkLayer()
	network := packet.NetworkLayer()
	transport := packet.TransportLayer()
	for _, l := range packet.Layers() {
		if layers != "" {
			layers += ";"
		}

		// add layer name
		layers += l.LayerType().String()

		// add address information for the link layer
		if showMACs && link != nil &&
			l.LayerType() == link.LayerType() {
			layers += getFlowAddresses(link.LinkFlow(), "MAC")
		}

		// add address information for the network layer
		if showIPs && network != nil &&
			l.LayerType() == network.LayerType() {
			layers += getFlowAddresses(network.NetworkFlow(), "IP")
		}

		// add port information for the transport layer
		if showPorts && transport != nil &&
			l.LayerType() == transport.LayerType() {
			layers += getFlowAddresses(transport.TransportFlow(),
				"Port")
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
