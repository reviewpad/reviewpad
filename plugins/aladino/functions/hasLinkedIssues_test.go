// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
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

	_, err := hasLinkedIssues(mockedEnv, []lang.Value{})

	assert.NotNil(t, err)
}

func TestHasLinkedIssues_WhenHasLinkedIssues(t *testing.T) {
	mockedPrNum := int64(6)
	mockedPrOwner := "foobar"
	mockedPrRepoName := "default-mock-repo"
	mockedAuthorLogin := "john"
	mockedGraphQLQuery := fmt.Sprintf(
		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){closingIssuesReferences{totalCount}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":%q,\"repositoryOwner\":%q}}\n",
		mockedPrNum,
		mockedPrRepoName,
		mockedPrOwner,
	)
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Number: mockedPrNum,
		Author: &pbc.User{Login: mockedAuthorLogin},
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: mockedPrOwner,
				Name:  mockedPrRepoName,
			},
		},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
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
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	wantVal := lang.BuildBoolValue(true)
	gotVal, err := hasLinkedIssues(mockedEnv, []lang.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasLinkedIssues_WhenNoLinkedIssues(t *testing.T) {
	mockedPrNum := int64(6)
	mockedPrOwner := "foobar"
	mockedPrRepoName := "default-mock-repo"
	mockedAuthorLogin := "john"
	mockedGraphQLQuery := fmt.Sprintf(
		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){closingIssuesReferences{totalCount}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":%q,\"repositoryOwner\":%q}}\n",
		mockedPrNum,
		mockedPrRepoName,
		mockedPrOwner,
	)
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Number: mockedPrNum,
		Author: &pbc.User{Login: mockedAuthorLogin},
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: mockedPrOwner,
				Name:  mockedPrRepoName,
			},
		},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
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
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	wantVal := lang.BuildBoolValue(false)
	gotVal, err := hasLinkedIssues(mockedEnv, []lang.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
