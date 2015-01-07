// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sprintly

import (
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

func appendArgs(urlString string, args interface{}) (string, error) {
	v := reflect.ValueOf(args)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return urlString, nil
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	values, err := query.Values(args)
	if err != nil {
		return "", err
	}

	u.RawQuery = values.Encode()
	return u.String(), nil
}

func encodeArgs(args interface{}) (io.Reader, error) {
	v := reflect.ValueOf(args)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, nil
	}

	values, err := query.Values(args)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(values.Encode()), nil
}
