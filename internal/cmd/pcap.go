package cmd

import (
	"time"

	"github.com/google/gopacket"
	"github.com/hwipl/packet-go/pkg/pcap"
)

type handler struct{}

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
		if link != nil && l.LayerType() == link.LayerType() {
			flow := link.LinkFlow()
			src, dst := flow.Endpoints()
			if src.LessThan(dst) {
				layers += ";MAC_" + src.String()
				layers += ";MAC_" + dst.String()
			} else {
				layers += ";MAC_" + dst.String()
				layers += ";MAC_" + src.String()
			}

			// add flow direction information
			layers += ";MAC_" + src.String() + "->MAC_" +
				dst.String()
		}

		// add address information for the network layer
		if network != nil && l.LayerType() == network.LayerType() {
			flow := network.NetworkFlow()
			src, dst := flow.Endpoints()
			if src.LessThan(dst) {
				layers += ";IP_" + src.String()
				layers += ";IP_" + dst.String()
			} else {
				layers += ";IP_" + dst.String()
				layers += ";IP_" + src.String()
			}

			// add flow direction information
			layers += ";IP_" + src.String() + "->IP_" +
				dst.String()
		}

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

			// add flow direction information
			layers += ";Port_" + src.String() + "->Port_" +
				dst.String()
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
