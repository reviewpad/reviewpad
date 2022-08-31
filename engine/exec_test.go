// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/engine/testutils"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestEval_WhenGitHubRequestsFail(t *testing.T) {
	tests := map[string]struct {
		inputReviewpadFilePath string
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
						w.Write(mock.MustMarshal(github.ErrorResponse{
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
						w.Write(mock.MustMarshal(github.ErrorResponse{
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

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			reviewpadFileData, err := utils.LoadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
			}

			reviewpadFile, err := testutils.ParseReviewpadFile(reviewpadFileData)
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
		clientOptions          []mock.MockBackendOption
		wantProgram            *engine.Program
		wantErr                string
	}{
		"when label has no name": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_unnamed_label.yml",
			clientOptions:          []mock.MockBackendOption{mockGetReposLabelsByOwnerByRepoByName("bug", "")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-unnamed-label")`),
				},
			),
		},
		"when label exists": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_without_description.yml",
			clientOptions:          []mock.MockBackendOption{mockGetReposLabelsByOwnerByRepoByName("bug", "")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-label-without-description")`),
				},
			),
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
			),
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
			),
		},
		"when label has no updates": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_label_with_repo_description.yml",
			clientOptions:          []mock.MockBackendOption{mockGetReposLabelsByOwnerByRepoByName("bug", "Something is wrong")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-label-with-no-updates")`),
				},
			),
		},
		"when group is valid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-valid-group")`),
				},
			),
		},
		"when group is invalid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_invalid_group.yml",
			wantErr:                "ProcessGroup:evalGroup expression is not a valid group",
		},
		"when workflow is invalid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_invalid_workflow.yml",
			wantErr:                "expression 1 is not a condition",
		},
		"when no workflow is activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_no_activated_workflows.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{},
			),
		},
		"when one workflow is activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_one_activated_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activate-one-workflow")`),
				},
			),
		},
		"when multiple workflows are activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_multiple_activated_workflows.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow-a")`),
					engine.BuildStatement(`$addLabel("activated-workflow-b")`),
				},
			),
		},
		"when workflow is activated and has extra actions": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_activated_workflow_with_extra_actions.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow")`),
					engine.BuildStatement(`$addLabel("workflow-with-extra-actions")`),
				},
			),
		},
		"when workflow is skipped": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_skipped_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow")`),
				},
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockAladinoInterpreter(mockedClient)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("mockAladinoInterpreter: %v", err))
			}

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			reviewpadFileData, err := utils.LoadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
			}

			reviewpadFile, err := testutils.ParseReviewpadFile(reviewpadFileData)
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
		engine.DefaultMockPullRequestTargetEntity,
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
