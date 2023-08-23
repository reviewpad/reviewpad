// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v52/github"
	"github.com/hasura/go-graphql-client"
	"github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	base          http.RoundTripper
	clientREST    *github.Client
	clientGQL     *githubv4.Client
	rawClientGQL  *graphql.Client
	token         string
	totalRequests uint64
	logger        *logrus.Entry
}

type GithubAppClient struct {
	*github.Client
}

func (t *GithubClient) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var err error
	requestType := "rest"

	if strings.Contains(req.URL.Path, "graphql") {
		requestType = "graphql"
	}

	t.totalRequests++

	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	if t.logger != nil {
		t.logger.WithFields(logrus.Fields{
			"method":                req.Method,
			"url":                   req.URL.String(),
			"body":                  string(body),
			"current_request_count": t.totalRequests,
			"request_type":          requestType,
		}).Debug("github client request")
	}

	return t.base.RoundTrip(req)
}

func NewGithubClient(clientREST *github.Client, clientGQL *githubv4.Client, rawClientGQL *graphql.Client, logger *logrus.Entry) *GithubClient {
	return &GithubClient{
		clientREST:   clientREST,
		clientGQL:    clientGQL,
		rawClientGQL: rawClientGQL,
		logger:       logger,
	}
}

func NewGithubClientFromToken(ctx context.Context, token string, logger *logrus.Entry) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	clientREST := github.NewClient(tc)
	clientGQL := githubv4.NewClient(tc)
	rawClientGQL := graphql.NewClient("https://api.github.com/graphql", tc)

	client := &GithubClient{
		clientREST:   clientREST,
		clientGQL:    clientGQL,
		rawClientGQL: rawClientGQL,
		token:        token,
		logger:       logger,
		base:         tc.Transport,
	}

	// override the transport to count requests
	tc.Transport = client

	return client
}

// FIXME: Remove these to hide the implementation details.
func (c *GithubClient) GetClientREST() *github.Client {
	return c.clientREST
}

func (c *GithubClient) GetClientGraphQL() *githubv4.Client {
	return c.clientGQL
}

func (c *GithubClient) GetToken() string {
	return c.token
}

func (c *GithubClient) GetRawClientGraphQL() *graphql.Client {
	return c.rawClientGQL
}

func (c *GithubClient) GetAuthenticatedUserLogin() (string, error) {
	clientGQL := c.GetClientGraphQL()

	var userLogin struct {
		Viewer struct {
			Login string
		}
	}

	err := clientGQL.Query(context.Background(), &userLogin, nil)
	if err != nil {
		return "", err
	}

	return userLogin.Viewer.Login, nil
}

func (c *GithubClient) GetTotalRequests() uint64 {
	return c.totalRequests
}

func (c *GithubClient) WithLogger(ctx context.Context, logger *logrus.Entry) *GithubClient {
	return NewGithubClientFromToken(ctx, c.GetToken(), logger)
}

func NewGithubAppClient(gitHubAppID int64, gitHubAppPrivateKey []byte) (*GithubAppClient, error) {
	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, gitHubAppID, gitHubAppPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub App client: %v", err)
	}

	return &GithubAppClient{github.NewClient(&http.Client{Transport: transport})}, nil
}
