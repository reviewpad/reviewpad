// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAddRateLimitQuery(t *testing.T) {
	tests := map[string]struct {
		body    []byte
		want    []byte
		wantErr error
	}{
		"when error is request": {
			body: []byte(`{
				"query": "mutation { login }"
			`),
			want:    nil,
			wantErr: fmt.Errorf("failed to unmarshal GraphQL request: unexpected end of JSON input"),
		},
		"when error is query": {
			body: []byte(`{
				"query": "mutation { login "
			}`),
			want:    nil,
			wantErr: fmt.Errorf("failed to parse GraphQL request: Syntax Error GraphQL (1:18) Expected Name, found EOF\n\n1: mutation { login \n                    ^\n"),
		},
		"when request is mutation": {
			body: []byte(`{
				"query": "mutation { login }"
			}`),
			want: []byte(`{
				"query": "mutation { login }"
			}`),
		},
		"when request is query": {
			body: []byte(`{
				"query": "query { username }"
			}`),
			want: []byte(`{"query":"{\n  username\n  rateLimit {\n    limit\n    cost\n    remaining\n    resetAt\n  }\n}\n"}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := addRateLimitQuery(tt.body)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSetRateLimitFields(t *testing.T) {
	tests := map[string]struct {
		fields      logrus.Fields
		requestType string
		headers     http.Header
		body        []byte
		wantFields  logrus.Fields
		wantErr     error
	}{
		"when request type is rest with no ratelimit headers": {
			fields: logrus.Fields{
				"test": "test",
			},
			requestType: "rest",
			headers:     http.Header{},
			body:        nil,
			wantFields: logrus.Fields{
				"test": "test",
			},
			wantErr: errors.New("failed to convert x-ratelimit-limit header to int: strconv.Atoi: parsing \"\": invalid syntax"),
		},
		"when request type is rest with no ratelimit-remaining header": {
			fields: logrus.Fields{
				"test": "test",
			},
			requestType: "rest",
			headers: http.Header{
				textproto.CanonicalMIMEHeaderKey("x-ratelimit-limit"): []string{"1"},
			},
			body: nil,
			wantFields: logrus.Fields{
				"test": "test",
			},
			wantErr: errors.New("failed to convert x-ratelimit-remaining header to int: strconv.Atoi: parsing \"\": invalid syntax"),
		},
		"when request type is rest": {
			fields: logrus.Fields{
				"test": "test",
			},
			requestType: "rest",
			headers: http.Header{
				textproto.CanonicalMIMEHeaderKey("x-ratelimit-limit"):     []string{"5000"},
				textproto.CanonicalMIMEHeaderKey("x-ratelimit-remaining"): []string{"4999"},
				textproto.CanonicalMIMEHeaderKey("x-ratelimit-reset"):     []string{"1628601600"},
			},
			body: nil,
			wantFields: logrus.Fields{
				"test":                 "test",
				"rate_limit_per_hour":  5000,
				"rate_limit_remaining": 4999,
				"rate_limit_reset":     "1628601600",
				"rate_limit_cost":      1,
			},
		},
		"when request type is graphql with response error": {
			fields: logrus.Fields{
				"test": "test",
			},
			requestType: "graphql",
			body: []byte(`{
				"data": {
			}`),
			wantFields: logrus.Fields{
				"test": "test",
			},
			wantErr: errors.New("failed to unmarshal GraphQL response: unexpected end of JSON input"),
		},
		"when request type is graphql": {
			fields: logrus.Fields{
				"test": "test",
			},
			requestType: "graphql",
			body: []byte(`{
				"data": {
					"rateLimit": {
						"limit": 5000,
						"cost": 5,
						"remaining": 4995,
						"resetAt": "2021-08-10T00:00:00Z"
					}
				}
			}`),
			wantFields: logrus.Fields{
				"test":                 "test",
				"rate_limit_per_hour":  5000,
				"rate_limit_remaining": 4995,
				"rate_limit_reset":     "2021-08-10T00:00:00Z",
				"rate_limit_cost":      5,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fields, err := setRateLimitFields(tt.fields, tt.requestType, tt.headers, tt.body)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantFields, fields)
		})
	}
}
