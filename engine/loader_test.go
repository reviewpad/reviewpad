// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/jarcoal/httpmock"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

type httpMockResponder struct {
	url       string
	responder httpmock.Responder
}

func TestLoadWithAST(t *testing.T) {
	inputContext := context.Background()
	inputGitHubClient := aladino.MockDefaultGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(
						[]github.RepositoryContent{
							{
								Name:        github.String("reviewpad_with_no_extends_a.yml"),
								DownloadURL: github.String("https://raw.githubusercontent.com/reviewpad_with_no_extends_a.yml"),
							},
							{
								Name:        github.String("reviewpad_with_no_extends_b.yml"),
								DownloadURL: github.String("https://raw.githubusercontent.com/reviewpad_with_no_extends_b.yml"),
							},
						},
					))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: "/reviewpad_with_no_extends_a.yml",
					Method:  "GET",
				},
				httpmock.File("testdata/loader/reviewpad_with_no_extends_a.yml").Bytes(),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: "/reviewpad_with_no_extends_b.yml",
					Method:  "GET",
				},
				httpmock.File("testdata/loader/reviewpad_with_no_extends_b.yml").Bytes(),
			),
		},
		nil,
	)

	inputReviewpadFilePath := "testdata/loader/reviewpad_with_extends.yml"
	reviewpadFileData, err := utils.ReadFile(inputReviewpadFilePath)
	if err != nil {
		assert.FailNow(t, "Error reading reviewpad file: %v", err)
	}

	gotReviewpadFile, err := engine.Load(inputContext, inputGitHubClient, reviewpadFileData)
	if err != nil {
		assert.FailNow(t, "Error parsing reviewpad file: %v", err)
	}

	noIgnoreErrors := true
	noMetricsOnMerge := true
	wantReviewpadFile := &engine.ReviewpadFile{
		Version:        "reviewpad.com/v3.x",
		Edition:        "enterprise",
		Mode:           "verbose",
		Recipes:        map[string]*bool{},
		IgnoreErrors:   &noIgnoreErrors,
		MetricsOnMerge: &noMetricsOnMerge,
		Extends:        []string{},
		Labels: map[string]engine.PadLabel{
			"small": {
				Color: "#aa12ab",
			},
			"medium": {
				Color: "#a8c3f7",
			},
		},
		Groups: []engine.PadGroup{
			{
				Name: "owners",
				Kind: "developers",
				Spec: "[\"anonymous\"]",
			},
		},
		Rules: []engine.PadRule{
			{
				Name: "is-medium",
				Kind: "patch",
				Spec: "$size([]) > 30 && $size([]) <= 100",
			},
			{
				Name: "is-small",
				Kind: "patch",
				Spec: "$size([]) < 30",
			},
			{
				Name: "$isElementOf($author(), $group(\"owners\"))",
				Kind: "patch",
				Spec: "$isElementOf($author(), $group(\"owners\"))",
			},
			{
				Name: "true",
				Kind: "patch",
				Spec: "true",
			},
		},
		Workflows: []engine.PadWorkflow{
			{
				Name:      "add-label-with-medium-size",
				AlwaysRun: false,
				On: []handler.TargetEntityKind{
					"pull_request",
				},
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "is-medium",
					},
				},
				Actions: []string{
					"$addLabel(\"medium\")",
				},
			},
			{
				Name:      "info-owners",
				AlwaysRun: false,
				On: []handler.TargetEntityKind{
					"pull_request",
				},
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "$isElementOf($author(), $group(\"owners\"))",
					},
				},
				Actions: []string{
					"$info(\"bob has authored a PR\")",
				},
			},
			{
				Name:      "check-title",
				AlwaysRun: false,
				On: []handler.TargetEntityKind{
					"pull_request",
				},
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "true",
					},
				},
				Actions: []string{
					"$titleLint()",
				},
			},
			{
				Name:      "add-label-with-small-size",
				AlwaysRun: false,
				On: []handler.TargetEntityKind{
					"pull_request",
				},
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "is-small",
					},
				},
				Actions: []string{
					"$addLabel(\"small\")",
				},
			},
		},
	}

	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestLoad(t *testing.T) {
	tests := map[string]struct {
		httpMockResponders     []httpMockResponder
		inputReviewpadFilePath string
		inputContext           context.Context
		inputGitHubClient      *gh.GithubClient
		wantReviewpadFilePath  string
		wantErr                string
	}{
		"when the file has a parsing error": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_parse_error.yml",
			wantErr:                "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file imports a nonexistent file": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_import_of_nonexistent_file.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/nonexistent_file",
					responder: httpmock.NewErrorResponder(fmt.Errorf("file doesn't exist")),
				},
			},
			wantErr: "Get \"https://foo.bar/nonexistent_file\": file doesn't exist",
		},
		"when the file imports a file that has a parsing error": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_import_file_with_parse_error.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_parse_error.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_parse_error.yml").Bytes()),
				},
			},
			wantErr: "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file has cyclic imports": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_cyclic_dependency_a.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_cyclic_dependency_b.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_cyclic_dependency_b.yml").Bytes()),
				},
				{
					url:       "https://foo.bar/reviewpad_with_cyclic_dependency_a.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_cyclic_dependency_a.yml").Bytes()),
				},
			},
			wantErr: "loader: cyclic dependency",
		},
		"when the file has import chains": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_imports_chain.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_no_imports.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_no_imports.yml").Bytes()),
				},
				{
					url:       "https://foo.bar/reviewpad_with_one_import.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_one_import.yml").Bytes()),
				},
			},
			wantReviewpadFilePath: "testdata/loader/reviewpad_appended.yml",
		},
		"when the file has no issues": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_no_imports.yml",
			wantReviewpadFilePath:  "testdata/loader/reviewpad_with_no_imports.yml",
		},
		"when the file requires action transformation": {
			inputReviewpadFilePath: "testdata/loader/transform/reviewpad_before_action_transform.yml",
			wantReviewpadFilePath:  "testdata/loader/transform/reviewpad_after_action_transform.yml",
		},
		"when the file requires extra action transformation": {
			inputReviewpadFilePath: "testdata/loader/transform/reviewpad_before_extra_action_transform.yml",
			wantReviewpadFilePath:  "testdata/loader/transform/reviewpad_after_extra_action_transform.yml",
		},
		"when the file has no on field": {
			inputReviewpadFilePath: "testdata/loader/transform/reviewpad_before_on_transform.yml",
			wantReviewpadFilePath:  "testdata/loader/transform/reviewpad_after_on_transform.yml",
		},
		"when the file has a rule with no kind": {
			inputReviewpadFilePath: "testdata/loader/transform/reviewpad_before_kind_transform.yml",
			wantReviewpadFilePath:  "testdata/loader/transform/reviewpad_after_kind_transform.yml",
		},
		"when the file has inline rules": {
			inputReviewpadFilePath: "testdata/loader/process/reviewpad_with_inline_rules.yml",
			wantReviewpadFilePath:  "testdata/loader/process/reviewpad_with_inline_rules_after_processing.yml",
		},
		"when the file has invalid inline rule": {
			inputReviewpadFilePath: "testdata/loader/process/reviewpad_with_invalid_inline_rule.yml",
			wantErr:                "unknown rule type int",
		},
		"when the file has an inline rule with extra actions": {
			inputReviewpadFilePath: "testdata/loader/process/reviewpad_with_inline_rules_with_extra_actions.yml",
			wantReviewpadFilePath:  "testdata/loader/process/reviewpad_with_inline_rules_with_extra_actions_after_processing.yml",
		},
		"when the file has multiple inline rules": {
			inputReviewpadFilePath: "testdata/loader/process/reviewpad_with_multiple_inline_rules.yml",
			wantReviewpadFilePath:  "testdata/loader/process/reviewpad_with_multiple_inline_rules_after_processing.yml",
		},
		"when the file has non-existing extends": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_extends.yml",
			inputContext:           context.Background(),
			inputGitHubClient: aladino.MockDefaultGithubClient(
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposContentsByOwnerByRepoByPath,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"loader: extends file not found",
							)
						}),
					),
				},
				nil,
			),
			wantErr: "loader: extends file not found",
		},
		"when the file has invalid extends": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_invalid_extends.yml",
			inputContext:           context.Background(),
			wantErr:                "fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go",
		},
		"when the file has an extends": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_extends.yml",
			wantReviewpadFilePath:  "testdata/loader/process/reviewpad_with_extends_after_processing.yml",
			inputContext:           context.Background(),
			inputGitHubClient: aladino.MockDefaultGithubClient(
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposContentsByOwnerByRepoByPath,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(
								[]github.RepositoryContent{
									{
										Name:        github.String("reviewpad_with_no_extends_a.yml"),
										DownloadURL: github.String("https://raw.githubusercontent.com/reviewpad_with_no_extends_a.yml"),
									},
									{
										Name:        github.String("reviewpad_with_no_extends_b.yml"),
										DownloadURL: github.String("https://raw.githubusercontent.com/reviewpad_with_no_extends_b.yml"),
									},
								},
							))
						}),
					),
					mock.WithRequestMatch(
						mock.EndpointPattern{
							Pattern: "/reviewpad_with_no_extends_a.yml",
							Method:  "GET",
						},
						httpmock.File("testdata/loader/reviewpad_with_no_extends_a.yml").Bytes(),
					),
					mock.WithRequestMatch(
						mock.EndpointPattern{
							Pattern: "/reviewpad_with_no_extends_b.yml",
							Method:  "GET",
						},
						httpmock.File("testdata/loader/reviewpad_with_no_extends_b.yml").Bytes(),
					),
				},
				nil,
			),
		},
		"when the file has cyclic extends": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_cyclic_extends_dependency_a.yml",
			inputContext:           context.Background(),
			inputGitHubClient: aladino.MockDefaultGithubClient(
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposContentsByOwnerByRepoByPath,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(
								[]github.RepositoryContent{
									{
										Name:        github.String("reviewpad_with_cyclic_extends_dependency_a.yml"),
										DownloadURL: github.String("https://raw.githubusercontent.com/reviewpad_with_cyclic_extends_dependency_a.yml"),
									},
									{
										Name:        github.String("reviewpad_with_cyclic_extends_dependency_b.yml"),
										DownloadURL: github.String("https://raw.githubusercontent.com/reviewpad_with_cyclic_extends_dependency_b.yml"),
									},
								},
							))
						}),
					),
					mock.WithRequestMatch(
						mock.EndpointPattern{
							Pattern: "/reviewpad_with_cyclic_extends_dependency_a.yml",
							Method:  "GET",
						},
						httpmock.File("testdata/loader/reviewpad_with_cyclic_extends_dependency_a.yml").Bytes(),
					),
					mock.WithRequestMatch(
						mock.EndpointPattern{
							Pattern: "/reviewpad_with_cyclic_extends_dependency_b.yml",
							Method:  "GET",
						},
						httpmock.File("testdata/loader/reviewpad_with_cyclic_extends_dependency_b.yml").Bytes(),
					),
				},
				nil,
			),
			wantErr: "loader: cyclic extends dependency",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.httpMockResponders != nil {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				registerHttpResponders(test.httpMockResponders)
			}

			var wantReviewpadFile *engine.ReviewpadFile
			if test.wantReviewpadFilePath != "" {
				wantReviewpadFileData, err := utils.ReadFile(test.wantReviewpadFilePath)
				if err != nil {
					assert.FailNow(t, "Error reading reviewpad file: %v", err)
				}

				wantReviewpadFile, err = engine.Load(test.inputContext, test.inputGitHubClient, wantReviewpadFileData)
				if err != nil {
					assert.FailNow(t, "Error parsing reviewpad file: %v", err)
				}
			}

			reviewpadFileData, err := utils.ReadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, "Error reading reviewpad file: %v", err)
			}

			gotReviewpadFile, gotErr := engine.Load(test.inputContext, test.inputGitHubClient, reviewpadFileData)

			if gotErr != nil {
				githubError := &github.ErrorResponse{}
				if errors.As(gotErr, &githubError) {
					if githubError.Message != test.wantErr {
						assert.FailNow(t, "Load() error = %v, wantErr %v", gotErr, test.wantErr)
					}
				} else {
					if gotErr.Error() != test.wantErr {
						assert.FailNow(t, "Load() error = %v, wantErr %v", gotErr, test.wantErr)
					}
				}
			}

			assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
		})
	}
}

func registerHttpResponders(httpMockResponders []httpMockResponder) {
	for _, httpMockResponder := range httpMockResponders {
		httpmock.RegisterResponder("GET", httpMockResponder.url, httpMockResponder.responder)
	}
}
