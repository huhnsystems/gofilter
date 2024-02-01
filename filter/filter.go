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

// Package filter contains efficient methods and types to test for the existence
// of given strings in a payload.
package filter

import (
	"bytes"
	"slices"
)

type Filter interface {
	Filter([]byte) bool
}

// StringFilter contains a set of strings which should be filtered.
type StringFilter [][]byte

// NewStringFilter creates a new set of strings which should be filtered.
func NewStringFilter(strs []string) StringFilter {
	filter := StringFilter(make([][]byte, len(strs)))

	for i, v := range strs {
		filter[i] = []byte(v)
	}

	return filter
}

// Filter tests a payload to determine if it contains a string of StringFilter.
// TODO: Simplify the Filter() function and make it more efficient.
func (flt StringFilter) Filter(payload []byte) bool {
	return slices.ContainsFunc(flt, func(needle []byte) bool {
		return bytes.Contains(payload, needle)
	})
}
