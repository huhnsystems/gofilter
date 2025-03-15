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

// gofilter is an OpenBSD string filter using divert(4) written in Golang.
//
// For more information try runinng:
//
//	gofilter -h
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/huhnsystems/gofilter/filter"
	"github.com/huhnsystems/gofilter/network"
)

func main() {
	portFlag := flag.Int("p", 700, "divert socket listening port")
	filterFlag := flag.String("f", "", "strings to filter, comma separated")
	flag.Parse()

	srv := network.NewSrv(context.Background(), *portFlag)

	if *filterFlag != "" {
		lst := strings.Split(*filterFlag, ",")
		log.Println("Filtering packets with the following strings:", lst)
		srv.Filter = filter.NewStringFilter(lst)
	}

	// Intercept sigint and shut down the program.
	shutdown := make(chan struct{})
	sigint := make(chan os.Signal)
	go func() {
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		srv.Shutdown()
		close(shutdown)
	}()

	// Start listening on IPv4 and IPv6.
	go func() {
		if err := srv.Listen(); err != nil {
			log.Fatalln(err)
		}
	}()
	log.Println("Listening on port", *portFlag)

	<-shutdown
	log.Println("Shutting down")
}
