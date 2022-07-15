// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/collector"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/stretchr/testify/assert"
)

func TestGetCtx(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantCtx := context.Background()

	gotCtx := mockedEnv.GetCtx()

	assert.Equal(t, wantCtx, gotCtx)
}

func TestGetCollector(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	var wantCollector collector.Collector

	gotCollector := mockedEnv.GetCollector()

	assert.Equal(t, wantCollector, gotCollector)
}

func TestGetPullRequest(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantPullRequest, _, err := mockedEnv.GetClient().PullRequests.Get(
		mockedEnv.GetCtx(),
		"foobar",
		"default-mock-repo",
		6,
	)
	if err != nil {
		log.Fatal(err)
	}

	gotPullRequest := mockedEnv.GetPullRequest()

	assert.Equal(t, wantPullRequest, gotPullRequest)
}

func TestGetPatch(t *testing.T) {
	filename := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String(filename),
			Patch:    nil,
		},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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

	wantPatch := make(aladino.Patch)
	wantPatch[filename] = &aladino.File{
		Repr: &github.CommitFile{
			Filename: github.String(filename),
		},
		Diff: []*aladino.DiffBlock{},
	}

	gotPatch := mockedEnv.GetPatch()

	assert.Equal(t, wantPatch, gotPatch)
}

func TestGetRegisterMap(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantRegisterMap := make(aladino.RegisterMap)

	gotRegisterMap := mockedEnv.GetRegisterMap()

	assert.Equal(t, wantRegisterMap, gotRegisterMap)
}

func TestGetBuiltIns(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantBuiltIns := &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			"emptyFunction": {
				Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
				Code: nil,
			},
		},
		Actions: map[string]*aladino.BuiltInAction{
			"emptyAction": {
				Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
				Code: nil,
			},
		},
	}

	gotBuiltIns := mockedEnv.GetBuiltIns()

	assert.Equal(t, wantBuiltIns, gotBuiltIns)
}

func TestGetReport(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantReport := &aladino.Report{
		WorkflowDetails: make(map[string]aladino.ReportWorkflowDetails),
	}

	gotReport := mockedEnv.GetReport()

	assert.Equal(t, wantReport, gotReport)
}

func TestNewTypeEnv(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantTypeEnv := aladino.TypeEnv(map[string]aladino.Type{
		"emptyFunction": aladino.BuildFunctionType([]aladino.Type{}, nil),
		"emptyAction":   aladino.BuildFunctionType([]aladino.Type{}, nil),
	})

	gotTypeEnv := *aladino.NewTypeEnv(mockedEnv)

	assert.Equal(t, wantTypeEnv, gotTypeEnv)
}

// The MockDefaultEnv function to create a mocked default env, it uses the NewEvalEnv which is the function we wish to test
func TestNewEvalEnv_WhenGetPullRequestFilesFails(t *testing.T) {
	failMessage := "GetPullRequestFilesFail"
	_, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
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
		},
		nil,
	)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestNewEvalEnv_WhenNewFileFails(t *testing.T) {
	mockedPullRequestFileList := &[]*github.CommitFile{
		{Patch: github.String("@@a")},
	}
	_, err := mocks_aladino.MockDefaultEnv(
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

	assert.EqualError(t, err, "error in file patch : error in chunk lines parsing (1): missing lines info: @@a\npatch: @@a")
}

func TestNewEvalEnv(t *testing.T) {
	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String("default-mock-repo/file1.ts"),
			Patch:    nil,
		},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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

	wantEvalEnv := &aladino.BaseEnv{
		Ctx:         mockedEnv.GetCtx(),
		Client:      mockedEnv.GetClient(),
		ClientGQL:   mockedEnv.GetClientGQL(),
		Collector:   mockedEnv.GetCollector(),
		PullRequest: mockedEnv.GetPullRequest(),
		Patch:       mockedEnv.GetPatch(),
		RegisterMap: mockedEnv.GetRegisterMap(),
		BuiltIns:    mockedEnv.GetBuiltIns(),
		Report:      mockedEnv.GetReport(),
	}

	assert.Equal(t, wantEvalEnv, mockedEnv)
}
