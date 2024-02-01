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
	errUDPTooShort = errors.New("udp segment too short")
	errUDPNil      = errors.New("unexpected nil")
)

const (
	udpHeaderLen = 8
)

type udpHeader struct{}

type udpSegment struct {
	udpHeader
	udpData []byte
}

func (h *udpHeader) parse(b []byte) error {
	if h == nil || b == nil {
		return errUDPNil
	}
	if len(b) < udpHeaderLen {
		return errUDPTooShort
	}

	return nil
}

func (s *udpSegment) parse(b []byte) error {
	if err := s.udpHeader.parse(b); err != nil {
		return err
	}

	s.udpData = b[8:]

	return nil
}

// data returns the udpData from the corresponding udpSegment.
func (s *udpSegment) data() []byte {
	return s.udpData
}
