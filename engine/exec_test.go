// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbe "github.com/reviewpad/api/go/entities"
	api_mocks "github.com/reviewpad/api/go/mocks"
	pbs "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/engine"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEval_WhenGitHubRequestsFail(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
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
	hostsClient := api_mocks.NewMockHostsClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetCodeReview(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetCodeReviewReply{
			Review: aladino.GetDefaultMockReview(),
		}, nil)

	codehostClient := &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockAladinoInterpreter(mockedClient, codehostClient)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("mockAladinoInterpreter: %v", err))
			}

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter, engine.DefaultMockTargetEntity, engine.DefaultMockEventDetails)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			reviewpadFileData, err := utils.ReadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
			}

			reviewpadFile, err := engine.Load(test.inputContext, logger, test.inputGitHubClient, reviewpadFileData)
			if err != nil {
				assert.FailNow(t, "Error parsing reviewpad file: %v", err)
			}

			gotProgram, err := engine.EvalConfigurationFile(reviewpadFile, mockedEnv)

			assert.Nil(t, gotProgram)
			assert.Equal(t, test.wantErr, err.(*github.ErrorResponse).Message)
		})
	}
}

func TestEvalCommand(t *testing.T) {
	tests := map[string]struct {
		inputContext      context.Context
		inputGitHubClient *gh.GithubClient
		clientOptions     []mock.MockBackendOption
		targetEntity      *handler.TargetEntity
		eventDetails      *handler.EventDetails
		wantProgram       *engine.Program
		wantErr           string
	}{
		"when invalid assign reviewer command": {
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "created",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						ID:   github.Int64(1),
						Body: github.String("/reviewpad assign-reviewers john, jane, 1 random"),
					},
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
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["john"], 1, "reviewpad")`),
				},
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "created",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						ID:   github.Int64(1),
						Body: github.String("/reviewpad assign-reviewers john -t 1"),
					},
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
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["jane-12","john01"], 1, "random")`),
				},
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "created",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						ID:   github.Int64(1),
						Body: github.String("/reviewpad assign-reviewers jane-12,john01 --total-reviewers 1 -p random"),
					},
				},
			},
		},
		"when assign reviewer has reviewpad policy": {
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["john","jane"], 1, "reviewpad")`),
				},
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "created",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						ID:   github.Int64(1),
						Body: github.String("/reviewpad assign-reviewers john,jane -t 1 --review-policy reviewpad"),
					},
				},
			},
		},
		"when assign reviewer has round-robin policy": {
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$assignReviewer(["john","johnny"], 1, "round-robin")`),
				},
			),
			targetEntity: &handler.TargetEntity{
				Kind:   engine.DefaultMockTargetEntity.Kind,
				Number: engine.DefaultMockTargetEntity.Number,
				Owner:  engine.DefaultMockTargetEntity.Owner,
				Repo:   engine.DefaultMockTargetEntity.Repo,
			},
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "created",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						ID:   github.Int64(1),
						Body: github.String("/reviewpad assign-reviewers john,johnny -t 1 -p round-robin"),
					},
				},
			},
		},
	}

	hostsClient := api_mocks.NewMockHostsClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetCodeReview(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetCodeReviewReply{
			Review: aladino.GetDefaultMockReview(),
		}, nil)

	codehostClient := &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockAladinoInterpreter(mockedClient, codehostClient)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("mockAladinoInterpreter: %v", err))
			}

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter, test.targetEntity, test.eventDetails)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			gotProgram, gotErr := engine.EvalCommand(test.eventDetails.Payload.(*github.IssueCommentEvent).GetComment().GetBody(), mockedEnv)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, fmt.Sprintf("Eval() error = %v, wantErr %v", gotErr, test.wantErr))
			}
			assert.Equal(t, test.wantProgram, gotProgram)
		})
	}
}

func TestEvalConfigurationFile(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	tests := map[string]struct {
		inputReviewpadFilePath string
		inputContext           context.Context
		inputGitHubClient      *gh.GithubClient
		clientOptions          []mock.MockBackendOption
		targetEntity           *handler.TargetEntity
		eventDetails           *handler.EventDetails
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
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when group is valid": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_valid_group.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("test-valid-group")`),
				},
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
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when one workflow is activated": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_one_activated_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activate-one-workflow")`),
				},
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
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
		"when workflow is skipped": {
			inputReviewpadFilePath: "testdata/exec/reviewpad_with_skipped_workflow.yml",
			wantProgram: engine.BuildProgram(
				[]*engine.Statement{
					engine.BuildStatement(`$addLabel("activated-workflow")`),
				},
			),
			targetEntity: engine.DefaultMockTargetEntity,
		},
	}

	hostsClient := api_mocks.NewMockHostsClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetCodeReview(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetCodeReviewReply{
			Review: aladino.GetDefaultMockReview(),
		}, nil)

	codehostClient := &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := engine.MockGithubClient(test.clientOptions)

			mockedAladinoInterpreter, err := mockAladinoInterpreter(mockedClient, codehostClient)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("mockAladinoInterpreter: %v", err))
			}

			mockedEnv, err := engine.MockEnvWith(mockedClient, mockedAladinoInterpreter, test.targetEntity, test.eventDetails)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("engine MockEnvWith: %v", err))
			}

			reviewpadFileData, err := utils.ReadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
			}

			reviewpadFile, err := engine.Load(test.inputContext, logger, test.inputGitHubClient, reviewpadFileData)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Error parsing reviewpad file: %v", err))
			}

			gotProgram, gotErr := engine.EvalConfigurationFile(reviewpadFile, mockedEnv)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, fmt.Sprintf("Eval() error = %v, wantErr %v", gotErr, test.wantErr))
			}
			assert.Equal(t, test.wantProgram, gotProgram)
		})
	}
}

func mockAladinoInterpreter(githubClient *gh.GithubClient, codehostClient *codehost.CodeHostClient) (engine.Interpreter, error) {
	dryRun := false
	logger := logrus.NewEntry(logrus.New())
	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		engine.DefaultMockCtx,
		logger,
		dryRun,
		githubClient,
		codehostClient,
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
