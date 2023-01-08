// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	clientREST *github.Client
	clientGQL  *githubv4.Client
	logger     *logrus.Entry
}

type GithubAppClient struct {
	*github.Client
}

func NewGithubClient(clientREST *github.Client, clientGQL *githubv4.Client) *GithubClient {
	return &GithubClient{
		clientREST: clientREST,
		clientGQL:  clientGQL,
	}
}

func NewGithubClientFromToken(ctx context.Context, logger *logrus.Entry, token string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	clientREST := github.NewClient(tc)
	clientGQL := githubv4.NewClient(tc)

	return &GithubClient{
		clientREST: clientREST,
		clientGQL:  clientGQL,
		logger:     logger,
	}
}

// FIXME: Remove these to hide the implementation details.
func (c *GithubClient) GetClientREST() *github.Client {
	return c.clientREST
}

func (c *GithubClient) GetClientGraphQL() *githubv4.Client {
	return c.clientGQL
}

func NewGithubAppClient(gitHubAppID int64, gitHubAppPrivateKey []byte) (*GithubAppClient, error) {
	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, gitHubAppID, gitHubAppPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub App client: %v", err)
	}

	return &GithubAppClient{github.NewClient(&http.Client{Transport: transport})}, nil
}
