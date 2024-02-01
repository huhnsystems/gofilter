// Copyright (c) 2024, Julian Huhn
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

// Package network contains the methods and types for the network stack and for
// extracting data from network packets.
package network

import (
	"encoding/binary"
	"errors"
)

var (
	errIPv4TooShort    = errors.New("ipv4 header too short")
	errIPv4Nil         = errors.New("unexpected nil")
	errIPv4LenMismatch = errors.New("mismatched lengths")
)

const ipv4HeaderLen = 20

type ipv4Header struct {
	headerLength uint8
	totalLength  uint16
	protocol     uint8
	options      []byte
}

type ipv4Packet struct {
	ipv4Header
	ipv4Data []byte
}

func (h *ipv4Header) parse(b []byte) error {
	if h == nil || b == nil {
		return errIPv4Nil
	}

	if len(b) < ipv4HeaderLen {
		return errIPv4TooShort
	}

	// Header length is calculated by the last 4 bits multiplied with 32.
	h.headerLength = (b[0] & 15) << 2 // TODO: 15 or 31?
	if len(b) < int(h.headerLength) {
		return errIPv4LenMismatch
	}
	if h.headerLength > 20 {
		h.options = b[20:h.headerLength]
	}

	h.totalLength = binary.BigEndian.Uint16(b[2:4])
	if len(b) != int(h.totalLength) {
		return errIPv4LenMismatch
	}

	h.protocol = b[9]

	return nil
}

func (d *ipv4Packet) parse(b []byte) error {
	// Parse the header of the IPv4 packet.
	if err := d.ipv4Header.parse(b); err != nil {
		return err
	}

	// Extract the data from the packet.
	d.ipv4Data = b[d.headerLength:]

	return nil
}

// data returns the payload contained in transport layer segments.
func (d *ipv4Packet) data() ([]byte, error) {
	// TODO: Catch case if nor tcp or udp.
	switch d.protocol {
	case 6:
		tcp := new(tcpSegment)
		if err := tcp.parse(d.ipv4Data); err != nil {
			return nil, err
		}
		return tcp.data(), nil
	case 17:
		udp := new(udpSegment)
		if err := udp.parse(d.ipv4Data); err != nil {
			return nil, err
		}
		return udp.data(), nil
	}

	return nil, nil
}
