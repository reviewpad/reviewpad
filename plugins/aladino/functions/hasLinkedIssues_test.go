// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var hasLinkedIssues = plugins_aladino.PluginBuiltIns().Functions["hasLinkedIssues"].Code

func TestHasLinkedIssues_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
		aladino.MockBuiltIns(),
		nil,
	)

	_, err := hasLinkedIssues(mockedEnv, []aladino.Value{})

	assert.NotNil(t, err)
}

func TestHasLinkedIssues_WhenHasLinkedIssues(t *testing.T) {
	mockedPrNum := 6
	mockedPrOwner := "foobar"
	mockedPrRepoName := "default-mock-repo"
	mockedAuthorLogin := "john"
	mockedGraphQLQuery := fmt.Sprintf(
		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){closingIssuesReferences{totalCount}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":%q,\"repositoryOwner\":%q}}\n",
		mockedPrNum,
		mockedPrRepoName,
		mockedPrOwner,
	)
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Number: github.Int(mockedPrNum),
		User:   &github.User{Login: github.String(mockedAuthorLogin)},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPrOwner),
				},
				Name: github.String(mockedPrRepoName),
			},
		},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MustRead(req.Body)
			switch query {
			case mockedGraphQLQuery:
				utils.MustWrite(
					w,
					`{"data": {
						"repository": {
							"pullRequest": {
								"closingIssuesReferences": {
									"totalCount": 3
								}
							}
						}
					}}`,
				)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	wantVal := aladino.BuildBoolValue(true)
	gotVal, err := hasLinkedIssues(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasLinkedIssues_WhenNoLinkedIssues(t *testing.T) {
	mockedPrNum := 6
	mockedPrOwner := "foobar"
	mockedPrRepoName := "default-mock-repo"
	mockedAuthorLogin := "john"
	mockedGraphQLQuery := fmt.Sprintf(
		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){closingIssuesReferences{totalCount}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":%q,\"repositoryOwner\":%q}}\n",
		mockedPrNum,
		mockedPrRepoName,
		mockedPrOwner,
	)
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Number: github.Int(mockedPrNum),
		User:   &github.User{Login: github.String(mockedAuthorLogin)},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPrOwner),
				},
				Name: github.String(mockedPrRepoName),
			},
		},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MustRead(req.Body)
			switch query {
			case mockedGraphQLQuery:
				utils.MustWrite(
					w,
					`{"data": {
						"repository": {
							"pullRequest": {
								"closingIssuesReferences": {
									"totalCount": 0
								}
							}
						}
					}}`,
				)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	wantVal := aladino.BuildBoolValue(false)
	gotVal, err := hasLinkedIssues(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
