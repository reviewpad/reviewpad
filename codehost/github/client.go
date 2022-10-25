// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	clientREST *github.Client
	clientGQL  *githubv4.Client
}

func NewGithubClient(clientREST *github.Client, clientGQL *githubv4.Client) *GithubClient {
	return &GithubClient{
		clientREST: clientREST,
		clientGQL:  clientGQL,
	}
}

func NewGithubClientFromToken(ctx context.Context, token string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	clientREST := github.NewClient(tc)
	clientGQL := githubv4.NewClient(tc)

	return &GithubClient{
		clientREST: clientREST,
		clientGQL:  clientGQL,
	}
}

// FIXME: Remove these to hide the implementation details.
func (c *GithubClient) GetClientREST() *github.Client {
	return c.clientREST
}

func (c *GithubClient) GetClientGraphQL() *githubv4.Client {
	return c.clientGQL
}
