// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v52/github"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
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

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName *string                `json:"operationName,omitempty"`
}

type GraphQLRateLimitResponse struct {
	Data struct {
		RateLimit struct {
			Limit     int64  `json:"limit"`
			Cost      int64  `json:"cost"`
			Remaining int64  `json:"remaining"`
			ResetAt   string `json:"resetAt"`
		} `json:"rateLimit"`
	} `json:"data"`
}

func (t *GithubClient) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var resErr error
	requestType := "rest"

	if strings.Contains(req.URL.Path, "graphql") {
		requestType = "graphql"
	}

	t.totalRequests++

	if req.Body != nil {
		body, resErr = io.ReadAll(req.Body)
		if resErr != nil {
			return nil, resErr
		}

		if requestType == "graphql" {
			body, resErr = addRateLimitQuery(body)
			if resErr != nil {
				return nil, resErr
			}
		}

		req.ContentLength = int64(len(body))
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	fields := logrus.Fields{
		"method":                req.Method,
		"url":                   req.URL.String(),
		"body":                  string(body),
		"current_request_count": t.totalRequests,
		"request_type":          requestType,
	}

	res, resErr := t.base.RoundTrip(req)
	if resErr != nil {
		return nil, resErr
	}

	resBody, resErr := io.ReadAll(res.Body)
	if resErr != nil {
		return nil, resErr
	}

	fields, err := setRateLimitFields(fields, requestType, res.Header, resBody)
	if err != nil {
		if t.logger != nil {
			t.logger.WithError(err).Error("failed to set rate limit fields")
		}
	}

	if requestType == "graphql" {
		// now we have to remove the rateLimit query from the response
		// so that the client can unmarshal the response into the expected struct
		resBody, err = removeRateLimitFromBody(resBody)
		if err != nil {
			return nil, err
		}
	}

	res.Body = io.NopCloser(bytes.NewReader(resBody))

	if t.logger != nil {
		t.logger.WithFields(fields).Debug("github client request")
	}

	return res, resErr
}

func setRateLimitFields(fields logrus.Fields, requestType string, headers http.Header, body []byte) (logrus.Fields, error) {
	if requestType == "rest" {
		limit, err := strconv.Atoi(headers.Get("x-ratelimit-limit"))
		if err != nil {
			return fields, err
		}

		remaining, err := strconv.Atoi(headers.Get("x-ratelimit-remaining"))
		if err != nil {
			return fields, err
		}

		fields["rate_limit_per_hour"] = limit
		fields["rate_limit_remaining"] = remaining
		fields["rate_limit_reset"] = headers.Get("x-ratelimit-reset")
		fields["rate_limit_cost"] = 1

		return fields, nil
	}

	rateLimitResponse := &GraphQLRateLimitResponse{}

	if err := json.Unmarshal(body, rateLimitResponse); err != nil {
		return fields, err
	}

	fields["rate_limit_per_hour"] = rateLimitResponse.Data.RateLimit.Limit
	fields["rate_limit_remaining"] = rateLimitResponse.Data.RateLimit.Remaining
	fields["rate_limit_reset"] = rateLimitResponse.Data.RateLimit.ResetAt
	fields["rate_limit_cost"] = rateLimitResponse.Data.RateLimit.Cost

	return fields, nil
}

// removeRateLimitFromBody removes the rate limit field from the GraphQL response
// so that the client can unmarshal the response into the expected struct
func removeRateLimitFromBody(body []byte) ([]byte, error) {
	r := map[string]interface{}{}

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	if _, ok := r["data"].(map[string]interface{}); ok {
		delete(r["data"].(map[string]interface{}), "rateLimit")
	}

	return json.Marshal(r)
}

// addRateLimitQuery adds a rate limit query to the GraphQL request
func addRateLimitQuery(body []byte) ([]byte, error) {
	graphQLRequest := &GraphQLRequest{}

	err := json.Unmarshal(body, graphQLRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL request: %s", err.Error())
	}

	doc, err := parser.Parse(parser.ParseParams{
		Source: string(graphQLRequest.Query),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse GraphQL request: %s", err.Error())
	}

	// if for some reason the query is empty, just return the body
	if len(doc.Definitions) == 0 {
		return body, nil
	}

	def := doc.Definitions[0].(*ast.OperationDefinition)

	if def.GetKind() != "OperationDefinition" || def.Operation != "query" {
		return body, nil
	}

	def.SelectionSet.Selections = append(def.SelectionSet.Selections, &ast.Field{
		Kind: "Field",
		Name: &ast.Name{
			Kind:  "Name",
			Value: "rateLimit",
		},
		SelectionSet: &ast.SelectionSet{
			Kind: "SelectionSet",
			Selections: []ast.Selection{
				&ast.Field{
					Kind: "Field",
					Name: &ast.Name{
						Kind:  "Name",
						Value: "limit",
					},
				},
				&ast.Field{
					Kind: "Field",
					Name: &ast.Name{
						Kind:  "Name",
						Value: "cost",
					},
				},
				&ast.Field{
					Kind: "Field",
					Name: &ast.Name{
						Kind:  "Name",
						Value: "remaining",
					},
				},
				&ast.Field{
					Kind: "Field",
					Name: &ast.Name{
						Kind:  "Name",
						Value: "resetAt",
					},
				},
			},
		},
	})

	graphQLRequest.Query = printer.Print(doc).(string)

	return json.Marshal(graphQLRequest)
}
