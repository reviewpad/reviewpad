// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v52/github"
	"github.com/hasura/go-graphql-client"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	pbe "github.com/reviewpad/api/go/entities"
	pbs "github.com/reviewpad/api/go/services"
	api_mocks "github.com/reviewpad/api/go/services_mocks"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const DefaultMockPrID = 1234
const DefaultMockPrNum = 6
const DefaultMockPrOwner = "foobar"
const DefaultMockPrRepoName = "default-mock-repo"
const DefaultMockEventName = "pull_request"
const DefaultMockEventAction = "opened"
const DefaultMockEntityNodeID = "test"

var DefaultMockPrDate = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
var DefaultMockContext = context.Background()
var DefaultMockLogger = logrus.NewEntry(logrus.New()).WithField("prefix", "[aladino]")
var DefaultMockCollector, _ = collector.NewCollector("", "distinctId", "pull_request", "runnerName", nil)
var DefaultMockTargetEntity = &entities.TargetEntity{
	Owner:  DefaultMockPrOwner,
	Repo:   DefaultMockPrRepoName,
	Number: DefaultMockPrNum,
	Kind:   entities.PullRequest,
}
var DefaultMockEventDetails = &entities.EventDetails{
	EventName:   DefaultMockEventName,
	EventAction: DefaultMockEventAction,
}

func GetDefaultMockIssueDetails() *github.Issue {
	issueNum := DefaultMockPrNum
	issueId := int64(DefaultMockPrID)
	issueDate := DefaultMockPrDate

	return &github.Issue{
		ID:   &issueId,
		User: &github.User{Login: github.String("john")},
		Assignees: []*github.User{
			{Login: github.String("jane")},
		},
		Title:     github.String("Found a bug"),
		Body:      github.String("I'm having a problem with this"),
		URL:       github.String("https://foo.bar"),
		CreatedAt: &github.Timestamp{Time: issueDate},
		Comments:  github.Int(6),
		Number:    github.Int(issueNum),
		Milestone: &github.Milestone{
			Title: github.String("v1.0"),
		},
		Labels: []*github.Label{
			{
				ID:   github.Int64(1),
				Name: github.String("bug"),
			},
		},
	}
}

func MockBuiltIns() *BuiltIns {
	return &BuiltIns{
		Functions: map[string]*BuiltInFunction{
			"emptyFunction": {
				Type: BuildFunctionType([]Type{}, nil),
				Code: func(e Env, args []Value) (Value, error) {
					return nil, nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
			},
			"zeroConst": {
				Type: BuildFunctionType([]Type{}, BuildIntType()),
				Code: func(e Env, args []Value) (Value, error) {
					return BuildIntValue(0), nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
			},
			"returnStr": {
				Type: BuildFunctionType([]Type{BuildStringType()}, BuildStringType()),
				Code: func(e Env, args []Value) (Value, error) {
					return args[0].(*StringValue), nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
			},
		},
		Actions: map[string]*BuiltInAction{
			"emptyAction": {
				Type: BuildFunctionType([]Type{}, nil),
				Code: func(e Env, args []Value) error {
					return nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
			},
		},
	}
}

func mockTypes() map[string]Type {
	mockedBuiltIns := MockBuiltIns()
	builtInsType := make(map[string]Type)

	for builtInName, builtInFunction := range mockedBuiltIns.Functions {
		builtInsType[builtInName] = builtInFunction.Type
	}

	for builtInName, builtInAction := range mockedBuiltIns.Actions {
		builtInsType[builtInName] = builtInAction.Type
	}

	return builtInsType
}

func MockTypeEnv() TypeEnv {
	return TypeEnv(mockTypes())
}

func mockHttpClientWith(clientOptions ...mock.MockBackendOption) *http.Client {
	return mock.NewMockedHTTPClient(clientOptions...)
}

func mockEnvWith(prOwner string, prRepoName string, prNum int, githubClient *gh.GithubClient, codehostClient *codehost.CodeHostClient, eventPayload interface{}, builtIns *BuiltIns, targetEntity *entities.TargetEntity) (Env, error) {
	ctx := context.Background()

	env, err := NewEvalEnv(
		ctx,
		DefaultMockLogger,
		false,
		githubClient,
		codehostClient,
		DefaultMockCollector,
		targetEntity,
		eventPayload,
		builtIns,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return env, nil
}

// mockDefaultHttpClient mocks an HTTP client with default values and ready for Aladino.
// Being ready for Aladino means that at least the request to build Aladino Env need to be mocked.
// As for now, at least one github request need to be mocked in order to build Aladino Env, mainly:
// - Get issue details (i.e. /repos/{owner}/{repo}/issues/{issue_number})
func mockDefaultHttpClient(clientOptions []mock.MockBackendOption) *http.Client {
	defaultMocks := []mock.MockBackendOption{
		mock.WithRequestMatchHandler(
			// Mock request to get issue details
			mock.GetReposIssuesByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				utils.MustWriteBytes(w, mock.MustMarshal(GetDefaultMockIssueDetails()))
			}),
		),
	}

	mocks := append(clientOptions, defaultMocks...)

	return mockHttpClientWith(mocks...)
}

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using handler directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func MockDefaultGithubClient(ghApiClientOptions []mock.MockBackendOption, ghGraphQLHandler func(http.ResponseWriter, *http.Request)) *gh.GithubClient {
	client := github.NewClient(mockDefaultHttpClient(ghApiClientOptions))

	// Handle GraphQL
	var clientGQL *githubv4.Client
	var rawClientGQL *graphql.Client
	if ghGraphQLHandler != nil {
		mux := http.NewServeMux()
		mux.HandleFunc("/graphql", ghGraphQLHandler)
		clientGQL = githubv4.NewClient(&http.Client{Transport: localRoundTripper{handler: mux}})
		rawClientGQL = graphql.NewClient("https://api.github.com/graphql", &http.Client{Transport: localRoundTripper{handler: mux}})
	}

	return gh.NewGithubClient(client, clientGQL, rawClientGQL)
}

func MockDefaultGithubAppClient(ghApiClientOptions []mock.MockBackendOption) *gh.GithubAppClient {
	return &gh.GithubAppClient{Client: github.NewClient(mock.NewMockedHTTPClient(ghApiClientOptions...))}
}

// MockDefaultEnv mocks an Aladino Env with default values.
func MockDefaultEnv(
	t *testing.T,
	ghApiClientOptions []mock.MockBackendOption,
	ghGraphQLHandler func(http.ResponseWriter, *http.Request),
	builtIns *BuiltIns,
	eventPayload interface{},
) Env {
	prOwner := DefaultMockPrOwner
	prRepoName := DefaultMockPrRepoName
	prNum := DefaultMockPrNum
	githubClient := MockDefaultGithubClient(ghApiClientOptions, ghGraphQLHandler)
	codehostClient := GetDefaultCodeHostClient(t, GetDefaultPullRequestDetails(), GetDefaultPullRequestFileList(), nil, nil)

	mockedEnv, err := mockEnvWith(prOwner, prRepoName, prNum, githubClient, codehostClient, eventPayload, builtIns, DefaultMockTargetEntity)
	if err != nil {
		t.Fatalf("[MockDefaultEnv] failed to create mock env: %v", err)
	}

	return mockedEnv
}

// MockDefaultEnvWithTargetEntity mocks an Aladino Env with default values and a custom TargetEntity.
func MockDefaultEnvWithTargetEntity(
	t *testing.T,
	ghApiClientOptions []mock.MockBackendOption,
	ghGraphQLHandler func(http.ResponseWriter, *http.Request),
	builtIns *BuiltIns,
	eventPayload interface{},
	targetEntity *entities.TargetEntity,
) Env {
	prOwner := DefaultMockPrOwner
	prRepoName := DefaultMockPrRepoName
	prNum := DefaultMockPrNum
	githubClient := MockDefaultGithubClient(ghApiClientOptions, ghGraphQLHandler)
	codehostClient := GetDefaultCodeHostClient(t, GetDefaultPullRequestDetails(), GetDefaultPullRequestFileList(), nil, nil)

	mockedEnv, err := mockEnvWith(prOwner, prRepoName, prNum, githubClient, codehostClient, eventPayload, builtIns, targetEntity)
	if err != nil {
		t.Fatalf("[MockDefaultEnvWithTargetEntity] failed to create mock env: %v", err)
	}

	return mockedEnv
}

func MockDefaultEnvWithPullRequestAndFiles(
	t *testing.T,
	ghApiClientOptions []mock.MockBackendOption,
	ghGraphQLHandler func(http.ResponseWriter, *http.Request),
	pullRequest *pbc.PullRequest,
	files []*pbc.File,
	builtIns *BuiltIns,
	eventPayload interface{},
) Env {
	prOwner := DefaultMockPrOwner
	prRepoName := DefaultMockPrRepoName
	prNum := DefaultMockPrNum
	githubClient := MockDefaultGithubClient(ghApiClientOptions, ghGraphQLHandler)
	codehostClient := GetDefaultCodeHostClient(t, pullRequest, files, nil, nil)

	mockedEnv, err := mockEnvWith(prOwner, prRepoName, prNum, githubClient, codehostClient, eventPayload, builtIns, DefaultMockTargetEntity)
	if err != nil {
		t.Fatalf("[MockDefaultEnvWithCodeReview] failed to create mock env: %v", err)
	}

	return mockedEnv
}

func GetDefaultPullRequestDetails() *pbc.PullRequest {
	prNum := DefaultMockPrNum
	prOwner := DefaultMockPrOwner
	prRepoName := DefaultMockPrRepoName
	prUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/%v", prOwner, prRepoName, prNum)
	prDate := DefaultMockPrDate

	return &pbc.PullRequest{
		Id: "test",
		Author: &pbc.User{
			Login: "john",
		},
		Status: pbc.PullRequestStatus_OPEN,
		Assignees: []*pbc.User{
			{
				Login: "jane",
			},
		},
		Title:         "Amazing new feature",
		Description:   "Please pull these awesome changes in!",
		Url:           "https://foo.bar",
		CreatedAt:     timestamppb.New(prDate),
		CommentsCount: 6,
		CommitsCount:  5,
		Number:        int64(prNum),
		Milestone: &pbc.Milestone{
			Title: "v1.0",
		},
		Labels: []*pbc.Label{
			{
				Id:   "1",
				Name: "enhancement",
			},
			{
				Id:   "2",
				Name: "large",
			},
		},
		Head: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: prOwner,
				Uri:   prUrl,
				Name:  prRepoName,
			},
			Name: "new-topic",
		},
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: prOwner,
				Uri:   prUrl,
				Name:  prRepoName,
			},
			Name: "master",
		},
		RequestedReviewers: &pbc.RequestedReviewers{
			Users: []*pbc.User{
				{
					Login: "jane",
				},
			},
		},
		IsMerged:        true,
		RawRestResponse: `{"id":1234,"number":6,"state":"open","title":"Amazing new feature","body":"Please pull these awesome changes in!","created_at":"2009-11-17T20:34:58.651387237Z","labels":[{"id":1,"name":"enhancement"},{"id":2,"name":"large"}],"user":{"login":"john"},"merged":true,"comments":6,"commits":5,"url":"https://foo.bar","assignees":[{"login":"jane"}],"milestone":{"title":"v1.0"},"node_id":"test","requested_reviewers":[{"login":"jane"}],"head":{"ref":"new-topic","repo":{"owner":{"login":"foobar"},"name":"default-mock-repo","url":"https://api.github.com/repos/foobar/default-mock-repo/pulls/6"}},"base":{"ref":"master","repo":{"owner":{"login":"foobar"},"name":"default-mock-repo","url":"https://api.github.com/repos/foobar/default-mock-repo/pulls/6"}}}`,
	}
}

func GetDefaultMockPullRequestDetailsWith(pullRequest *pbc.PullRequest) *pbc.PullRequest {
	defaultCodeReview := GetDefaultPullRequestDetails()

	if pullRequest == nil {
		return defaultCodeReview
	}

	if pullRequest.Author != nil {
		defaultCodeReview.Author = pullRequest.Author
	}

	if pullRequest.Number != 0 {
		defaultCodeReview.Number = pullRequest.Number
	}

	if pullRequest.Head != nil {
		defaultCodeReview.Head = pullRequest.Head
	}

	if pullRequest.Base != nil {
		defaultCodeReview.Base = pullRequest.Base
	}

	if pullRequest.Assignees != nil {
		defaultCodeReview.Assignees = pullRequest.Assignees
	}

	if pullRequest.CommitsCount != 0 {
		defaultCodeReview.CommitsCount = pullRequest.CommitsCount
	}

	if pullRequest.Labels != nil {
		defaultCodeReview.Labels = pullRequest.Labels
	}

	if pullRequest.Milestone != nil {
		defaultCodeReview.Milestone = pullRequest.Milestone
	}

	if pullRequest.RequestedReviewers != nil {
		defaultCodeReview.RequestedReviewers = pullRequest.RequestedReviewers
	}

	if pullRequest.AdditionsCount != 0 {
		defaultCodeReview.AdditionsCount = pullRequest.AdditionsCount
	}

	if pullRequest.DeletionsCount != 0 {
		defaultCodeReview.DeletionsCount = pullRequest.DeletionsCount
	}

	if pullRequest.Title != "" {
		defaultCodeReview.Title = pullRequest.Title
	}

	if pullRequest.Description != "" {
		defaultCodeReview.Description = pullRequest.Description
	}

	if pullRequest.IsDraft != defaultCodeReview.IsDraft {
		defaultCodeReview.IsDraft = pullRequest.IsDraft
	}

	if pullRequest.Id != "" {
		defaultCodeReview.Id = pullRequest.Id
	}

	if pullRequest.UpdatedAt != nil {
		defaultCodeReview.UpdatedAt = pullRequest.UpdatedAt
	}

	if pullRequest.IsRebaseable != defaultCodeReview.IsRebaseable {
		defaultCodeReview.IsRebaseable = pullRequest.IsRebaseable
	}

	if pullRequest.IsMerged != defaultCodeReview.IsMerged {
		defaultCodeReview.IsMerged = pullRequest.IsMerged
	}

	if pullRequest.ClosedAt != nil {
		defaultCodeReview.ClosedAt = pullRequest.ClosedAt
	}

	if pullRequest.Status != defaultCodeReview.Status {
		defaultCodeReview.Status = pullRequest.Status
	}

	return defaultCodeReview
}

func GetDefaultPullRequestFileList() []*pbc.File {
	prRepoName := DefaultMockPrRepoName
	return []*pbc.File{
		{
			Filename: fmt.Sprintf("%v/file1.ts", prRepoName),
			Patch:    "@@ -2,9 +2,11 @@ package main\n- func previous1() {\n+ func new1() {\n+\nreturn",
		},
		{
			Filename: fmt.Sprintf("%v/file2.ts", prRepoName),
			Patch:    "@@ -2,9 +2,11 @@ package main\n- func previous2() {\n+ func new2() {\n+\nreturn",
		},
		{
			Filename: fmt.Sprintf("%v/file3.ts", prRepoName),
			Patch:    "@@ -2,9 +2,11 @@ package main\n- func previous3() {\n+ func new3() {\n+\nreturn",
		},
	}
}

func GetDefaultCodeHostClient(t *testing.T, pullRequest *pbc.PullRequest, files []*pbc.File, pullRequestErr error, fileErr error) *codehost.CodeHostClient {
	hostsClient := api_mocks.NewMockHostClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetPullRequestReply{
			PullRequest: pullRequest,
		}, pullRequestErr)

	hostsClient.EXPECT().
		GetPullRequestFiles(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetPullRequestFilesReply{
			Files: files,
		}, nil)

	return &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}
}

func GetDefaultCodeHostClientWithFiles(t *testing.T, files []*pbc.File, err error) *codehost.CodeHostClient {
	hostsClient := api_mocks.NewMockHostClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetPullRequestReply{
			PullRequest: GetDefaultPullRequestDetails(),
		}, nil)

	hostsClient.EXPECT().
		GetPullRequestFiles(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetPullRequestFilesReply{
			Files: files,
		}, err)

	return &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}
}
