// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var dryRun = plugins_aladino.PluginBuiltIns().Actions["dryRun"].Code

func TestDryRun(t *testing.T) {
	var gotReport string
	tests := map[string]struct {
		clientOptions []mock.MockBackendOption
		reviewpadFile *engine.ReviewpadFile
		wantReport    string
		wantErr       error
	}{
		"when there is no reviewpad file": {
			reviewpadFile: &engine.ReviewpadFile{},
			wantErr:       nil,
		},
		"when program fails to build due to bad rule syntax": {
			reviewpadFile: &engine.ReviewpadFile{
				Rules: []engine.PadRule{
					{
						Name: "tautology",
						Spec: "1 == ",
					},
				},
				Workflows: []engine.PadWorkflow{
					{
						Name: "test",
						On: []handler.TargetEntityKind{
							"pull_request",
						},
						Rules: []engine.PadWorkflowRule{
							{
								Rule:         "tautology",
								ExtraActions: []string{"$addLabel(\"test\")"},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("parse error: failed to build AST on input 1 == "),
		},
		"when program fails to build due to bad pipeline trigger syntax": {
			reviewpadFile: &engine.ReviewpadFile{
				Rules: []engine.PadRule{
					{
						Name: "tautology",
						Spec: "1 == 1",
					},
				},
				Workflows: []engine.PadWorkflow{
					{
						Name: "test",
						On: []handler.TargetEntityKind{
							"pull_request",
						},
						Rules: []engine.PadWorkflowRule{
							{
								Rule:         "tautology",
								ExtraActions: []string{"$addLabel(\"test\")"},
							},
						},
					},
				},
				Pipelines: []engine.PadPipeline{
					{
						Name:    "test",
						Trigger: "1 == ",
						Stages: []engine.PadStage{
							{
								Actions: []string{
									"$addLabel(\"test\")",
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("parse error: failed to build AST on input 1 == "),
		},
		"when program fails to build due to bad pipeline stage actions syntax": {
			reviewpadFile: &engine.ReviewpadFile{
				Pipelines: []engine.PadPipeline{
					{
						Name:    "test",
						Trigger: "1 == 1",
						Stages: []engine.PadStage{
							{
								Actions: []string{
									"$dummy()",
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("no type for built-in dummy. Please check if the mode in the reviewpad.yml file supports it"),
		},
		"when program fails to build due to bad pipeline stage until syntax": {
			reviewpadFile: &engine.ReviewpadFile{
				Pipelines: []engine.PadPipeline{
					{
						Name:    "test",
						Trigger: "1 == 1",
						Stages: []engine.PadStage{
							{
								Actions: []string{
									"$addLabel(\"test\")",
								},
								Until: "1 ==",
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("parse error: failed to build AST on input 1 =="),
		},
		"when report fails to be posted": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"CreateCommentRequestFailed",
						)
					}),
				),
			},
			reviewpadFile: &engine.ReviewpadFile{
				Pipelines: []engine.PadPipeline{
					{
						Name:    "test",
						Trigger: "1 == 1",
						Stages: []engine.PadStage{
							{
								Actions: []string{
									"$addLabel(\"test\")",
								},
								Until: "1 == 1",
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("error on creating report comment "),
		},
		"when dry run is performed successfully": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						rawBody, _ := io.ReadAll(r.Body)
						body := github.IssueComment{}

						utils.MustUnmarshal(rawBody, &body)

						gotReport = *body.Body
					}),
				),
			},
			reviewpadFile: &engine.ReviewpadFile{
				Rules: []engine.PadRule{
					{
						Name: "tautology",
						Spec: "1 == 1",
					},
				},
				Workflows: []engine.PadWorkflow{
					{
						Name: "test",
						On: []handler.TargetEntityKind{
							"pull_request",
						},
						Rules: []engine.PadWorkflowRule{
							{
								Rule: "tautology",
							},
						},
						Actions: []string{"$addLabel(\"test-1\")"},
					},
				},
				Pipelines: []engine.PadPipeline{
					{
						Name:    "test",
						Trigger: "1 == 1",
						Stages: []engine.PadStage{
							{
								Actions: []string{
									"$addLabel(\"test-2\")",
								},
							},
						},
					},
				},
			},
			wantReport: "<!--@annotation-reviewpad-report-triggered-by-dry-run-command-->\n**Reviewpad Report** (Reviewpad ran in dry-run mode because the `$dryRun()` action was triggered)\n\n:scroll: **Executed actions**\n```yaml\n$addLabel(\"test-1\")\n$addLabel(\"test-2\")\n```\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithReviewpadFile(
				t,
				test.clientOptions,
				nil,
				plugins_aladino.PluginBuiltIns(),
				nil,
				test.reviewpadFile,
			)
			args := []aladino.Value{}
			gotErr := dryRun(mockedEnv, args)

			assert.Equal(t, test.wantReport, gotReport)
			assert.Equal(t, test.wantErr, gotErr)

			gotReport = ""
		})
	}
}
