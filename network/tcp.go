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
	"errors"
)

var (
	errTCPTooShort = errors.New("tcp segment too short")
	errTCPNil      = errors.New("unexpected nil")
)

const (
	tcpHeaderLen = 20
)

type tcpHeader struct {
	offset  uint8
	options []byte
}

type tcpSegment struct {
	tcpHeader
	tcpData []byte
}

func (h *tcpHeader) parse(b []byte) error {
	if h == nil || b == nil {
		return errTCPNil
	}
	if len(b) < tcpHeaderLen {
		return errTCPTooShort
	}

	h.offset = b[12] >> 2
	if h.offset > 20 {
		h.options = b[20:h.offset]
	}

	return nil
}

func (s *tcpSegment) parse(b []byte) error {
	if err := s.tcpHeader.parse(b); err != nil {
		return err
	}

	s.tcpData = b[s.offset:]

	return nil
}

// data returns the tcpData from the corresponding tcpSegment.
func (s *tcpSegment) data() []byte {
	return s.tcpData
}
