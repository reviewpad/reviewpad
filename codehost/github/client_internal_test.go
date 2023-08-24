// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"fmt"
	"testing"

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
