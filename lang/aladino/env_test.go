// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestGetCtx_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantCtx := aladino.DefaultMockContext

	gotCtx := mockedEnv.GetCtx()

	assert.Equal(t, wantCtx, gotCtx)
}

func TestGetCollector_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantCollector := aladino.DefaultMockCollector

	gotCollector := mockedEnv.GetCollector()

	assert.Equal(t, wantCollector, gotCollector)
}

func TestGetPullRequest_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantPullRequest := aladino.GetDefaultMockPullRequestDetails()

	gotPullRequest := mockedEnv.GetPullRequest()

	assert.Equal(t, wantPullRequest, gotPullRequest)
}

func TestGetPatch_WithDefaultEnv(t *testing.T) {
	fileName := fmt.Sprintf("%v/file1.ts", aladino.DefaultMockPrRepoName)
	patch := "@@ -2,9 +2,11 @@ package main\n- func previous1() {\n+ func new1() {\n+\nreturn"
	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String(fileName),
			Patch:    github.String(patch),
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedFile1 := &aladino.File{
		Repr: &github.CommitFile{
			Filename: github.String(fileName),
			Patch:    github.String(patch),
		},
	}
	mockedFile1.AppendToDiff(false, 2, 2, 2, 3, " func previous1() {", " func new1() {\n")

	wantPatch := aladino.Patch{
		fileName: mockedFile1,
	}

	gotPatch := mockedEnv.GetPatch()

	assert.Equal(t, wantPatch, gotPatch)
}

func TestGetRegisterMap_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantRegisterMap := make(aladino.RegisterMap)

	gotRegisterMap := mockedEnv.GetRegisterMap()

	assert.Equal(t, wantRegisterMap, gotRegisterMap)
}

func TestGetBuiltIns_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantBuiltIns := aladino.MockBuiltIns()

	gotBuiltIns := mockedEnv.GetBuiltIns()

	assert.Equal(t, len(wantBuiltIns.Functions), len(gotBuiltIns.Functions))
	assert.Equal(t, len(wantBuiltIns.Actions), len(gotBuiltIns.Actions))

	// TODO: Search a way to compare functions
	for functionName, functionCode := range wantBuiltIns.Functions {
		assert.NotNil(t, gotBuiltIns.Functions[functionName])
		assert.Equal(t, functionCode.Type, gotBuiltIns.Functions[functionName].Type)
	}

	for actionName, actionCode := range wantBuiltIns.Actions {
		assert.NotNil(t, gotBuiltIns.Actions[actionName])
		assert.Equal(t, actionCode.Type, gotBuiltIns.Actions[actionName].Type)
	}
}

func TestGetReport_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantReport := &aladino.Report{
		WorkflowDetails: make(map[string]aladino.ReportWorkflowDetails),
	}

	gotReport := mockedEnv.GetReport()

	assert.Equal(t, wantReport, gotReport)
}

func TestNewTypeEnv_WithDefaultEnv(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantTypeEnv := aladino.TypeEnv(map[string]aladino.Type{
		"emptyFunction": aladino.BuildFunctionType([]aladino.Type{}, nil),
		"emptyAction":   aladino.BuildFunctionType([]aladino.Type{}, nil),
		"zeroConst":     aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		"returnStr":     aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
	})

	gotTypeEnv := aladino.NewTypeEnv(mockedEnv)

	assert.Equal(t, wantTypeEnv, gotTypeEnv)
}

func TestNewEvalEnv_WhenGetPullRequestFilesFails(t *testing.T) {
	failMessage := "GetPullRequestFilesFail"
	mockedGithubClient := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			aladino.GetDefaultMockPullRequestDetails(),
		),
		mock.WithRequestMatchHandler(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	))

	ctx := context.Background()
	mockedPullRequest, _, err := mockedGithubClient.PullRequests.Get(
		ctx,
		aladino.DefaultMockPrOwner,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrNum,
	)
	if err != nil {
		log.Fatalf("couldn't get pull request: %v", err)
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		false,
		mockedGithubClient,
		nil,
		aladino.DefaultMockCollector,
		mockedPullRequest,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestNewEvalEnv_WhenNewFileFails(t *testing.T) {
	mockedGithubClient := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			aladino.GetDefaultMockPullRequestDetails(),
		),
		mock.WithRequestMatch(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			&[]*github.CommitFile{
				{Patch: github.String("@@a")},
			},
		),
	))

	ctx := context.Background()
	mockedPullRequest, _, err := mockedGithubClient.PullRequests.Get(
		ctx,
		aladino.DefaultMockPrOwner,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrNum,
	)
	if err != nil {
		log.Fatalf("couldn't get pull request: %v", err)
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		false,
		mockedGithubClient,
		nil,
		aladino.DefaultMockCollector,
		mockedPullRequest,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.EqualError(t, err, "error in file patch : error in chunk lines parsing (1): missing lines info: @@a\npatch: @@a")
}

func TestNewEvalEnv(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	patch := "@@ -2,9 +2,11 @@ package main\n- func previous1() {\n+ func new1() {\n+\nreturn"
	mockedGithubClient := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			aladino.GetDefaultMockPullRequestDetails(),
		),
		mock.WithRequestMatch(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			&[]*github.CommitFile{
				{
					Filename: github.String(fileName),
					Patch:    github.String(patch),
				},
			},
		),
	))

	ctx := context.Background()
	mockedPullRequest, _, err := mockedGithubClient.PullRequests.Get(
		ctx,
		aladino.DefaultMockPrOwner,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrNum,
	)

	if err != nil {
		log.Fatalf("couldn't get pull request: %v", err)
	}

	gotEnv, err := aladino.NewEvalEnv(
		ctx,
		false,
		mockedGithubClient,
		nil,
		aladino.DefaultMockCollector,
		mockedPullRequest,
		nil,
		aladino.MockBuiltIns(),
	)

	mockedFile1 := &aladino.File{
		Repr: &github.CommitFile{
			Filename: github.String(fileName),
			Patch:    github.String(patch),
		},
	}
	mockedFile1.AppendToDiff(false, 2, 2, 2, 3, " func previous1() {", " func new1() {\n")

	mockedPatch := aladino.Patch{
		fileName: mockedFile1,
	}

	wantEnv := &aladino.BaseEnv{
		Ctx:    ctx,
		Client: mockedGithubClient,
		// TODO: Mock a ClientGQL
		ClientGQL:   nil,
		Collector:   aladino.DefaultMockCollector,
		PullRequest: mockedPullRequest,
		Patch:       mockedPatch,
		RegisterMap: aladino.RegisterMap(make(map[string]aladino.Value)),
		BuiltIns:    aladino.MockBuiltIns(),
		Report:      &aladino.Report{WorkflowDetails: make(map[string]aladino.ReportWorkflowDetails)},
		// TODO: Mock an event
		EventPayload: nil,
	}

	assert.Nil(t, err)
	assert.Equal(t, wantEnv.Ctx, gotEnv.GetCtx())
	assert.Equal(t, wantEnv.Client, gotEnv.GetClient())
	assert.Equal(t, wantEnv.ClientGQL, gotEnv.GetClientGQL())
	assert.Equal(t, wantEnv.Collector, gotEnv.GetCollector())
	assert.Equal(t, wantEnv.PullRequest, gotEnv.GetPullRequest())
	assert.Equal(t, wantEnv.Patch, gotEnv.GetPatch())
	assert.Equal(t, wantEnv.RegisterMap, gotEnv.GetRegisterMap())
	assert.Equal(t, wantEnv.EventPayload, gotEnv.GetEventPayload())

	assert.Equal(t, len(wantEnv.BuiltIns.Functions), len(gotEnv.GetBuiltIns().Functions))
	assert.Equal(t, len(wantEnv.BuiltIns.Actions), len(gotEnv.GetBuiltIns().Actions))
	for functionName, functionCode := range wantEnv.BuiltIns.Functions {
		assert.NotNil(t, gotEnv.GetBuiltIns().Functions[functionName])
		assert.Equal(t, functionCode.Type, gotEnv.GetBuiltIns().Functions[functionName].Type)
	}
	for actionName, actionCode := range wantEnv.BuiltIns.Actions {
		assert.NotNil(t, gotEnv.GetBuiltIns().Actions[actionName])
		assert.Equal(t, actionCode.Type, gotEnv.GetBuiltIns().Actions[actionName].Type)
	}

	assert.Equal(t, wantEnv.Report, gotEnv.GetReport())
}
