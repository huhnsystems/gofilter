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
	"context"
	"syscall"

	"github.com/huhnsystems/gofilter/filter"
)

type Srv struct {
	// filter is the string filter to be applied to the packets.
	Filter filter.Filter

	fd4  int
	port int

	// buffer contains the next package to be processed.
	buffer [syscall.IP_MAXPACKET]byte

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSrv creates a Srv.
func NewSrv(ctx context.Context, port int) *Srv {
	srv := new(Srv)
	srv.port = port
	srv.ctx, srv.cancel = context.WithCancel(ctx)
	return srv
}

// Listen handles the network packets to be processed.
// TODO: Make IPv6 work.
func (srv *Srv) Listen() error {
	var err error
	srv.fd4, err = syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_RAW,
		syscall.IPPROTO_DIVERT,
	)
	if err != nil {
		return err
	}

	// Bind to a divert(4) socket.
	addr := syscall.SockaddrInet4{Port: srv.port}
	if err := syscall.Bind(srv.fd4, &addr); err != nil {
		return err
	}

	for {
		select {
		case <-srv.ctx.Done():
			return srv.ctx.Err()
		default:
			plen, peer, err := syscall.Recvfrom(
				srv.fd4,
				srv.buffer[:],
				0,
			)
			if err != nil {
				return err
			}
			// Skip iteration, if there is no packet with content.
			if plen == 0 {
				continue
			}

			// Read the packet from the buffer.
			packet := new(ipv4Packet)
			if err := packet.parse(srv.buffer[:plen]); err != nil {
				return err
			}

			// Extract the data from the packet.
			data, err := packet.data()
			if err != nil {
				return err
			}

			// Skip iteration, if the filter matches the data.
			if srv.Filter != nil && srv.Filter.Filter(data) {
				continue
			}

			// Return the packet to the kernel for further
			// processing.
			if err := syscall.Sendto(
				srv.fd4,
				srv.buffer[:plen],
				0,
				peer,
			); err != nil {
				return err
			}
		}
	}
}

// Shutdown cancels the context of the Srv.
func (srv *Srv) Shutdown() {
	srv.cancel()
}
