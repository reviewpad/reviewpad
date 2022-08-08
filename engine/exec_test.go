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
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils/file"
	"github.com/stretchr/testify/assert"
)

func TestEval_WhenRequestFails(t *testing.T) {
	tests := map[string]struct {
		inputReviewpadFilePath string
		clientOptions          []mock.MockBackendOption
		wantErr                string
	}{
		"when get label request fails": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_label.yml",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(mock.MustMarshal(github.ErrorResponse{
							// An error response may also consist of a 404 status code.
							// However, in this context, such response is not an error since it only means a label does not exist.
							// So we force an error by setting it to a status code different than 404.
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
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_label.yml",
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

			mockedAladinoInterpreter, err := mockDefaultAladinoInterpreterWith(mockedClient)
			if err != nil {
				assert.FailNow(t, "mockDefaultAladinoInterpreterWith: %v", err)
			}

			mockedEnv, err := engine.MockDefaultEnvWith(mockedClient, mockedAladinoInterpreter)
			if err != nil {
				assert.FailNow(t, "engine MockDefaultEnvWith: %v", err)
			}

			reviewpadFileData, err := file.LoadReviewpadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, "Error reading reviewpad file: %v", err)
			}

			reviewpadFile, err := file.ParseReviewpadFile(reviewpadFileData)
			if err != nil {
				assert.FailNow(t, "Error parsing reviewpad file: %v", err)
			}

			gotProgram, err := engine.Eval(reviewpadFile, mockedEnv)

			assert.Nil(t, gotProgram)
			assert.Equal(t, err.(*github.ErrorResponse).Message, test.wantErr)
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
		// Label use cases
		"when label has no name": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_unnamed_label.yml",
			clientOptions:          []mock.MockBackendOption{getLabelRequest("bug")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"test-unnamed-label\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "test-workflow",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"test-unnamed-label\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
		"when label exists": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_label.yml",
			clientOptions:          []mock.MockBackendOption{getLabelRequest("bug")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"test-valid-label\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "test-workflow",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"test-valid-label\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
		"when label does not exist": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_label.yml",
			clientOptions: []mock.MockBackendOption{
				getLabelRequest(""),
				createLabelRequest("bug", "f29513", ""),
			},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"test-valid-label\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "test-workflow",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"test-valid-label\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
		// Group use cases
		"when group is valid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			clientOptions:          []mock.MockBackendOption{getLabelRequest("test-valid-group")},
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"test-valid-group\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "test-workflow",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"test-valid-group\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
		"when group is invalid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_invalid_group.yml",
			clientOptions:          []mock.MockBackendOption{getLabelRequest("test-invalid-group")},
			wantErr:                "ProcessGroup:evalGroup expression is not a valid group",
		},
		// Workflow use cases
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
					engine.BuildStatement(
						"$addLabel(\"activate-one-workflow\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "activated-workflow",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"activate-one-workflow\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
		"when multiple workflows are activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_multiple_activated_workflows.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"activated-workflow-a\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "activated-workflow-a",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"activated-workflow-a\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
					engine.BuildStatement(
						"$addLabel(\"activated-workflow-b\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:      "activated-workflow-b",
								AlwaysRun: true,
								Rules:     []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions:   []string{"$addLabel(\"activated-workflow-b\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
		"when workflow is activated and has extra actions": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_activated_workflow_with_extra_actions.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"activated-workflow\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name: "activated-workflow",
								Rules: []engine.PadWorkflowRule{
									{
										Rule:         "tautology",
										ExtraActions: []string{"$addLabel(\"workflow-with-extra-actions\")"},
									},
								},
								Actions: []string{"$addLabel(\"activated-workflow\")"},
							},
							[]engine.PadWorkflowRule{
								{
									Rule:         "tautology",
									ExtraActions: []string{"$addLabel(\"workflow-with-extra-actions\")"},
								},
							},
						),
					),
					engine.BuildStatement(
						"$addLabel(\"workflow-with-extra-actions\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name: "activated-workflow",
								Rules: []engine.PadWorkflowRule{
									{
										Rule:         "tautology",
										ExtraActions: []string{"$addLabel(\"workflow-with-extra-actions\")"},
									},
								},
								Actions: []string{"$addLabel(\"activated-workflow\")"},
							},
							[]engine.PadWorkflowRule{
								{
									Rule:         "tautology",
									ExtraActions: []string{"$addLabel(\"workflow-with-extra-actions\")"},
								},
							},
						),
					),
				},
			),
		},
		"when workflow is skipped": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_skipped_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(
						"$addLabel(\"activated-workflow\")",
						engine.BuildMetadata(
							engine.PadWorkflow{
								Name:    "activated-workflow",
								Rules:   []engine.PadWorkflowRule{{Rule: "tautology"}},
								Actions: []string{"$addLabel(\"activated-workflow\")"},
							},
							[]engine.PadWorkflowRule{{Rule: "tautology"}},
						),
					),
				},
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockDefaultAladinoInterpreterWith(mockedClient)
			if err != nil {
				assert.FailNow(t, "mockDefaultAladinoInterpreterWith: %v", err)
			}

			mockedEnv, err := engine.MockDefaultEnvWith(mockedClient, mockedAladinoInterpreter)
			if err != nil {
				assert.FailNow(t, "engine MockDefaultEnvWith: %v", err)
			}

			reviewpadFileData, err := file.LoadReviewpadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, "Error reading reviewpad file: %v", err)
			}

			reviewpadFile, err := file.ParseReviewpadFile(reviewpadFileData)
			if err != nil {
				assert.FailNow(t, "Error parsing reviewpad file: %v", err)
			}

			gotProgram, gotErr := engine.Eval(reviewpadFile, mockedEnv)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "Load() error = %v, wantErr %v", gotErr, test.wantErr)
			}
			assert.Equal(t, test.wantProgram, gotProgram)
		})
	}
}

func mockDefaultAladinoInterpreterWith(client *github.Client) (engine.Interpreter, error) {
	dryRun := false
	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		engine.DefaultMockCtx,
		dryRun,
		client,
		// TODO: add mocked github GQL client
		nil,
		engine.DefaultMockCollector,
		engine.GetDefaultMockPullRequestDetails(),
		engine.DefaultMockEventPayload,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		return nil, fmt.Errorf("aladino NewInterpreter returned unexpected error: %v", err)
	}

	return mockedAladinoInterpreter, nil
}

func getLabelRequest(label string) mock.MockBackendOption {
	var labelsMockedResponse *github.Label

	if label != "" {
		labelsMockedResponse = &github.Label{
			Name: github.String(label),
		}
	} else {
		labelsMockedResponse = &github.Label{}
	}

	return mock.WithRequestMatch(
		mock.GetReposLabelsByOwnerByRepoByName,
		labelsMockedResponse,
	)
}

func createLabelRequest(name, color, description string) mock.MockBackendOption {
	labelsMockedResponse := &github.Label{
		Name:        github.String(name),
		Color:       github.String(color),
		Description: github.String(description),
	}

	return mock.WithRequestMatch(
		mock.PostReposLabelsByOwnerByRepo,
		labelsMockedResponse,
	)
}
