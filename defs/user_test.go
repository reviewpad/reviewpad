// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs_test

import (
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/defs"
	"github.com/stretchr/testify/assert"
)

func TestFromGithubUsers(t *testing.T) {
	tests := map[string]struct {
		githubUsers []*github.User
		users       defs.Users
	}{
		"when single user": {
			githubUsers: []*github.User{
				{
					NodeID: github.String("1"),
					Login:  github.String("login"),
				},
			},
			users: []defs.User{
				{
					ID:    "1",
					Login: "login",
				},
			},
		},
		"when multiple users": {
			githubUsers: []*github.User{
				{
					NodeID: github.String("1"),
					Login:  github.String("login"),
				},
				{
					NodeID: github.String("2"),
					Login:  github.String("login"),
				},
			},
			users: []defs.User{
				{
					ID:    "1",
					Login: "login",
				},
				{
					ID:    "2",
					Login: "login",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			users := defs.FromGithubUsers(test.githubUsers)
			assert.Equal(t, test.users, users)
		})
	}
}
