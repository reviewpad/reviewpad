// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/api/go/services"
	"github.com/reviewpad/api/go/services_mocks"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v3/plugins/aladino/services"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestCommentHasAnnotation(t *testing.T) {
	tests := map[string]struct {
		comment    string
		annotation string
		wantVal    bool
	}{
		"single annotation": {
			comment:    "reviewpad-an: foo",
			annotation: "foo",
			wantVal:    true,
		},
		"single annotation with comment": {
			comment:    "reviewpad-an: foo\n hello world",
			annotation: "foo",
			wantVal:    true,
		},
		"multiple annotation empty": {
			comment:    "reviewpad-an: foo     bar",
			annotation: "",
			wantVal:    false,
		},
		"multiple annotation single spaced first": {
			comment:    "reviewpad-an: foo bar",
			annotation: "foo",
			wantVal:    true,
		},
		"multiple annotation single spaced": {
			comment:    "reviewpad-an: foo bar",
			annotation: "bar",
			wantVal:    true,
		},
		"multiple annotation multi spaced": {
			comment:    "reviewpad-an: foo   bar",
			annotation: "bar",
			wantVal:    true,
		},
		"starting with empty line": {
			comment:    "\n  reviewpad-an: foo",
			annotation: "foo",
			wantVal:    true,
		},
		"starting with spaces": {
			comment:    "  reviewpad-an: foo",
			annotation: "foo",
			wantVal:    true,
		},
		"annotation not found": {
			comment:    "reviewpad-an: foo",
			annotation: "bar",
			wantVal:    false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotVal := commentHasAnnotation(test.comment, test.annotation)

			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}

func TestHasAnnotationCode_WhenGetSymbolsFromPatchFails(t *testing.T) {
	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPullRequest := &github.PullRequest{
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPRRepoOwner),
				},
				URL:  github.String(mockedPRUrl),
				Name: github.String(mockedPRRepoName),
			},
			Ref: github.String("new-topic"),
			SHA: github.String(mockedHeadSHA),
		},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPRRepoOwner),
				},
				URL:  github.String(mockedPRUrl),
				Name: github.String(mockedPRRepoName),
			},
			Ref: github.String("master"),
		},
	}

	mockedPatchFilePath := "test"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)

	// Since the patch is simply passed around, it can be an empty string
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedPullRequestFileList := []*github.CommitFile{
		{
			SHA:      github.String(mockedBlobId),
			Filename: github.String(mockedPatchFileRelativeName),
			Patch:    github.String(mockedPatch),
		},
	}

	failMessage := "DownloadContents"

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
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
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("critical")}
	gotVal, gotErr := hasAnnotationCode(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.Equal(t, gotErr.(*github.ErrorResponse).Message, failMessage)
}

func TestHasAnnotationCode(t *testing.T) {
	tests := map[string]struct {
		args    []aladino.Value
		wantVal aladino.Value
	}{
		"has annotation": {
			args:    []aladino.Value{aladino.BuildStringValue("critical")},
			wantVal: aladino.BuildTrueValue(),
		},
		"does not have annotation": {
			args:    []aladino.Value{aladino.BuildStringValue("foo")},
			wantVal: aladino.BuildFalseValue(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockedSemanticClient := services_mocks.NewMockSemanticClient(controller)

			mockedPRRepoOwner := "mock-reviewpad"
			mockedPRRepoName := "test"
			mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
			mockedHeadSHA := "abc123"

			mockedPullRequest := &github.PullRequest{
				Head: &github.PullRequestBranch{
					Repo: &github.Repository{
						Owner: &github.User{
							Login: github.String(mockedPRRepoOwner),
						},
						URL:  github.String(mockedPRUrl),
						Name: github.String(mockedPRRepoName),
					},
					Ref: github.String("new-topic"),
					SHA: github.String(mockedHeadSHA),
				},
				Base: &github.PullRequestBranch{
					Repo: &github.Repository{
						Owner: &github.User{
							Login: github.String(mockedPRRepoOwner),
						},
						URL:  github.String(mockedPRUrl),
						Name: github.String(mockedPRRepoName),
					},
					Ref: github.String("master"),
				},
			}

			mockedPatchFilePath := "test"
			mockedPatchFileName := "crawler.go"
			mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
			mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

			// Since the blob and the patch are simply passed around, they can be an empty string
			mockedBlob := ""
			mockedPatch := ""
			mockedBlobId := "1234"

			mockedPullRequestFileList := []*github.CommitFile{
				{
					SHA:      github.String(mockedBlobId),
					Filename: github.String(mockedPatchFileRelativeName),
					Patch:    github.String(mockedPatch),
				},
			}

			mockedSymbols := &entities.Symbols{
				Files: map[string]*entities.File{},
				Symbols: map[string]*entities.Symbol{
					"Crawl": {
						Id:   "1",
						Name: "Crawl",
						Type: "Function",
						CodeComments: []*entities.SymbolDocumentation{
							{
								Code: "reviewpad-an: critical",
							},
						},
					},
				},
			}

			mockedSemanticClient.EXPECT().GetSymbols(
				gomock.Any(),
				&services.GetSymbolsRequest{
					Uri:      mockedPRUrl,
					CommitId: mockedHeadSHA,
					Filepath: mockedPatchFileRelativeName,
					Blob:     []byte(fmt.Sprintf("%#v", mockedBlob)),
					BlobId:   mockedBlobId,
					Diff:     &entities.ResolveFileDiff{Blocks: []*entities.ResolveBlockDiff{}}}).Return(
				&services.GetSymbolsReply{
					Symbols: mockedSymbols,
				}, nil,
			)

			mockBuiltIns := &aladino.BuiltIns{
				Services: map[string]interface{}{
					plugins_aladino_services.SEMANTIC_SERVICE_KEY: mockedSemanticClient,
				},
			}

			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposContentsByOwnerByRepoByPath,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
								{
									Name:        github.String(mockedPatchFileName),
									Path:        github.String(mockedPatchFilePath),
									DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
								},
							}))
						}),
					),
					mock.WithRequestMatch(
						mock.EndpointPattern{
							Pattern: mockedPatchLocation,
							Method:  "GET",
						},
						mockedBlob,
					),
				},
				nil,
				mockBuiltIns,
				nil,
			)

			gotVal, gotErr := hasAnnotationCode(mockedEnv, test.args)

			assert.Nil(t, gotErr)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
