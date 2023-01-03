// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestEval_WhenGitHubRequestsFail(t *testing.T) {
	tests := map[string]struct {
		inputReviewpadFilePath string
		inputContext           context.Context
		inputGitHubClient      *gh.GithubClient
		clientOptions          []mock.MockBackendOption
		wantErr                string
	}{
		"when get label request fails": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_without_description.yml",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						engine.MustWriteBytes(w, mock.MustMarshal(github.ErrorResponse{
							// An error response may also consist of a 404 status code.
							// However, in this context, such response means a label does not exist.
							Response: &http.Response{
								StatusCode: 500,
							},
							Message: "GetLabelRequestFailed",
						}))
					}),
				),
			},
			wantErr: "GetLabelRequestFailed",
		},
		"when create label request fails": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_without_description.yml",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						engine.MustWriteBytes(w, mock.MustMarshal(github.ErrorResponse{
							Response: &http.Response{
								StatusCode: 404,
							},
							Message: "Resource not found",
						}))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.PostReposLabelsByOwnerByRepo,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"CreateLabelRequestFailed",
						)
					}),
				),
			},
			wantErr: "CreateLabelRequestFailed",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockAladinoInterpreter(mockedClient)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("mockAladinoInterpreter: %v", err))
			}

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter, engine.DefaultMockTargetEntity, engine.DefaultMockEventData)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			reviewpadFileData, err := utils.ReadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
			}

			reviewpadFile, err := engine.Load(test.inputContext, test.inputGitHubClient, reviewpadFileData)
			if err != nil {
				assert.FailNow(t, "Error parsing reviewpad file: %v", err)
			}

			gotProgram, err := engine.Eval(reviewpadFile, mockedEnv)

			assert.Nil(t, gotProgram)
			assert.Equal(t, test.wantErr, err.(*github.ErrorResponse).Message)
		})
	}
}

func TestEval(t *testing.T) {
	tests := map[string]struct {
		inputReviewpadFilePath string
		inputContext           context.Context
		inputGitHubClient      *gh.GithubClient
		clientOptions          []mock.MockBackendOption
		targetEntity           *handler.TargetEntity
		eventData              *handler.EventData
		wantProgram            *engine.Program
		wantErr                string
	}{
		"when label has no name": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_unnamed_label.yml",
			clientOptions: []mock.MockBackendOption{
				mockGetReposLabelsByOwnerByRepoByName("bug", ""),
				mockPatchReposLabelsByOwnerByRepo("bug", ""),
			},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-unnamed-label")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when label exists": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_without_description.yml",
			clientOptions:          []mock.MockBackendOption{mockGetReposLabelsByOwnerByRepoByName("bug", "")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-label-without-description")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when label does not exist": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_without_description.yml",
			clientOptions: []mock.MockBackendOption{
				mockGetReposLabelsByOwnerByRepoByName("bug", ""),
				mockPostReposLabelsByOwnerByRepo("bug", ""),
			},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-label-without-description")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when label has updates": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_with_diff_repo_description.yml",
			clientOptions: []mock.MockBackendOption{
				mockGetReposLabelsByOwnerByRepoByName("bug", "Something is wrong"),
				mockPatchReposLabelsByOwnerByRepo("bug", "Something is definitely wrong"),
			},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-label-update")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when label has no updates": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_with_repo_description.yml",
			clientOptions:          []mock.MockBackendOption{mockGetReposLabelsByOwnerByRepoByName("bug", "Something is wrong")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-label-with-no-updates")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when group is valid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-valid-group")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when group is invalid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_invalid_group.yml",
			wantErr:                "ProcessGroup:evalGroup expression is not a valid group",
			targetEntity:           engine.DefaultMockTargetEntity,
		},
		"when workflow is invalid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_invalid_workflow.yml",
			wantErr:                "expression 1 is not a condition",
			targetEntity:           engine.DefaultMockTargetEntity,
		},
		"when no workflow is activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_no_activated_workflows.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when one workflow is activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_one_activated_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activate-one-workflow")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when multiple workflows are activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_multiple_activated_workflows.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow-a")`),
					engine.BuildStatement(`$addLabel("activated-workflow-b")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when workflow is activated and has extra actions": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_activated_workflow_with_extra_actions.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow")`),
					engine.BuildStatement(`$addLabel("workflow-with-extra-actions")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when workflow is skipped": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_skipped_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow")`),
				},
				false,
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when invalid assign reviewer command": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventData: &handler.EventData{
				EventName:   "issue_comment",
				EventAction: "created",
				Comment: &github.IssueComment{
					ID:   github.Int64(1),
					Body: github.String("/reviewpad assign-reviewers john, jane, 1 random"),
				},
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
					&github.IssueComment{},
				),
			},
			wantErr: "accepts 1 arg(s), received 4",
		},
		"when missing policy in assign reviewer command args": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["john"], 1, "reviewpad")`),
				},
				true,
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventData: &handler.EventData{
				EventName:   "issue_comment",
				EventAction: "created",
				Comment: &github.IssueComment{
					ID:   github.Int64(1),
					Body: github.String("/reviewpad assign-reviewers john -t 1"),
				},
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
					&github.IssueComment{},
				),
			},
		},
		"when assign reviewer has random policy": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["jane-12","john01"], 1, "random")`),
				},
				true,
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventData: &handler.EventData{
				EventName:   "issue_comment",
				EventAction: "created",
				Comment: &github.IssueComment{
					ID:   github.Int64(1),
					Body: github.String("/reviewpad assign-reviewers jane-12,john01 --total-reviewers 1 -p random"),
				},
			},
		},
		"when assign reviewer has reviewpad policy": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["john","jane"], 1, "reviewpad")`),
				},
				true,
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventData: &handler.EventData{
				EventName:   "issue_comment",
				EventAction: "created",
				Comment: &github.IssueComment{
					ID:   github.Int64(1),
					Body: github.String("/reviewpad assign-reviewers john,jane -t 1 --review-policy reviewpad"),
				},
			},
		},
		"when assign reviewer has round-robin policy": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["john","johnny"], 1, "round-robin")`),
				},
				true,
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventData: &handler.EventData{
				EventName:   "issue_comment",
				EventAction: "created",
				Comment: &github.IssueComment{
					ID:   github.Int64(1),
					Body: github.String("/reviewpad assign-reviewers john,johnny -t 1 -p round-robin"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockAladinoInterpreter(mockedClient)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("mockAladinoInterpreter: %v", err))
			}

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter, test.targetEntity, test.eventData)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			reviewpadFileData, err := utils.ReadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
			}

			reviewpadFile, err := engine.Load(test.inputContext, test.inputGitHubClient, reviewpadFileData)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error parsing reviewpad file: %v", err))
			}

			gotProgram, gotErr := engine.Eval(reviewpadFile, mockedEnv)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, fmt.Sprintf("Eval() error = %v, wantErr %v", gotErr, test.wantErr))
			}
			assert.Equal(t, test.wantProgram, gotProgram)
		})
	}
}

func mockAladinoInterpreter(githubClient *gh.GithubClient) (engine.Interpreter, error) {
	dryRun := false
	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		engine.DefaultMockCtx,
		dryRun,
		githubClient,
		engine.DefaultMockCollector,
		engine.DefaultMockTargetEntity,
		engine.DefaultMockEventPayload,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		return nil, fmt.Errorf("aladino NewInterpreter returned unexpected error: %v", err)
	}

	return mockedAladinoInterpreter, nil
}

func mockGetReposLabelsByOwnerByRepoByName(name, description string) mock.MockBackendOption {
	labelsMockedResponse := &github.Label{}

	if name != "" {
		labelsMockedResponse = &github.Label{
			Name:        github.String(name),
			Description: github.String(description),
		}
	}

	return mock.WithRequestMatch(
		mock.GetReposLabelsByOwnerByRepoByName,
		labelsMockedResponse,
	)
}

func mockPatchReposLabelsByOwnerByRepo(name, description string) mock.MockBackendOption {
	labelsMockedResponse := &github.Label{
		Name:        github.String(name),
		Description: github.String(description),
	}

	return mock.WithRequestMatch(
		mock.PatchReposLabelsByOwnerByRepoByName,
		labelsMockedResponse,
	)
}

func mockPostReposLabelsByOwnerByRepo(name, description string) mock.MockBackendOption {
	labelsMockedResponse := &github.Label{
		Name:        github.String(name),
		Description: github.String(description),
	}

	return mock.WithRequestMatch(
		mock.PostReposLabelsByOwnerByRepo,
		labelsMockedResponse,
	)
}
