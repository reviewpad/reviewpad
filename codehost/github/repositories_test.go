// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultRepositoryBranch_WhenGetRepositoryRequestFails(t *testing.T) {
	failMessage := "GetRepositoryRequestFail"
	mockedGithubClient := aladino.MockDefaultGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposByOwnerByRepo,
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

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()
	defaultBranch, err := mockedGithubClient.GetDefaultRepositoryBranch(
		context.Background(),
		mockedCodeReview.Base.Repo.Owner,
		mockedCodeReview.Base.Repo.GetName(),
	)

	assert.Equal(t, "", defaultBranch)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetDefaultRepositoryBranch(t *testing.T) {
	mockedCodeReview := aladino.GetDefaultPullRequestDetails()
	expectedDefaultBranch := mockedCodeReview.Base.Name

	mockedGithubClient := aladino.MockDefaultGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposByOwnerByRepo,
				&github.Repository{
					DefaultBranch: github.String(expectedDefaultBranch),
				},
			),
		},
		nil,
	)

	defaultBranch, err := mockedGithubClient.GetDefaultRepositoryBranch(
		context.Background(),
		mockedCodeReview.Base.Repo.Owner,
		mockedCodeReview.Base.Repo.Name,
	)

	assert.Equal(t, expectedDefaultBranch, defaultBranch)
	assert.Nil(t, err)
}

func TestDownloadContents_WhenInvalidDownloadMethodIsProvided(t *testing.T) {
	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		nil,
	)

	contents, err := mockedGithubClient.DownloadContents(
		context.Background(),
		"test/crawler.go",
		&pbc.Branch{
			Repo: &pbc.Repository{
				Owner: mockedCodeReview.Base.Repo.Owner,
				Name:  mockedCodeReview.Base.Repo.Name,
			},
		},
		&gh.DownloadContentsOptions{
			Method: "invalid",
		},
	)

	assert.Nil(t, contents)
	assert.Equal(t, err.Error(), "invalid download method specified")
}

func TestDownloadContents_WhenDownloadContentsRequestFails(t *testing.T) {
	failMessage := "DownloadContentsRequestFail"

	mockedPRRepoOwner := aladino.DefaultMockPrOwner
	mockedPRRepoName := aladino.DefaultMockPrRepoName
	mockedPRNumber := aladino.DefaultMockPrNum
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", mockedPRRepoOwner, mockedPRRepoName, mockedPRNumber)
	mockedPRBaseSHA := "abc123"
	mockedPRBaseBranchName := "master"

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: mockedPRBaseBranchName,
			Sha:  mockedPRBaseSHA,
		},
	})

	mockedPatchFilePath := "test/crawler.go"

	mockedFiles := []*github.CommitFile{
		{Filename: github.String(mockedPatchFilePath)},
	}

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
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
		mockedCodeReview,
		mockedFiles,
		aladino.MockBuiltIns(),
		nil,
	)

	contents, err := mockedEnv.GetGithubClient().DownloadContents(
		context.Background(),
		mockedPatchFilePath,
		mockedCodeReview.Base,
		&gh.DownloadContentsOptions{
			Method: gh.DownloadMethodSHA,
		},
	)

	assert.Nil(t, contents)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestDownloadContents(t *testing.T) {
	mockedPRRepoOwner := aladino.DefaultMockPrOwner
	mockedPRRepoName := aladino.DefaultMockPrRepoName
	mockedPRNumber := aladino.DefaultMockPrNum
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", mockedPRRepoOwner, mockedPRRepoName, mockedPRNumber)
	mockedPRBaseSHA := "abc123"
	mockedPRBaseBranchName := "master"

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: mockedPRBaseBranchName,
			Sha:  mockedPRBaseSHA,
		},
	})

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%s/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%s/%s/%s", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	mockedFiles := []*github.CommitFile{
		{Filename: github.String(mockedPatchFileRelativeName)},
	}

	mockedBlob := "test-blob"
	expectedContents := []byte(string("\"" + mockedBlob + "\""))

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%s", mockedPatchLocation)),
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
		mockedCodeReview,
		mockedFiles,
		aladino.MockBuiltIns(),
		nil,
	)

	contents, err := mockedEnv.GetGithubClient().DownloadContents(
		context.Background(),
		mockedPatchFileRelativeName,
		mockedCodeReview.Base,
		&gh.DownloadContentsOptions{
			Method: gh.DownloadMethodBranchName,
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, expectedContents, contents)
}
