// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var merge = plugins_aladino.PluginBuiltIns().Actions["merge"].Code

type MergeRequestPostBody struct {
	MergeMethod string `json:"merge_method"`
}

func TestMerge_WhenMergeMethodIsUnsupported(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue("INVALID")}
	err := merge(mockedEnv, args)

	assert.EqualError(t, err, "merge: unsupported merge method INVALID")
}

func TestMerge_WhenNoMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "merge"
	var gotMergeMethod string

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedIsGitHubMergeQueueEnabledGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": null
			}
		}
	}`

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PutReposPullsMergeByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := MergeRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotMergeMethod = body.MergeMethod
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				utils.MustWrite(w, mockedIsGitHubMergeQueueEnabledGQLQueryBody)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}

func TestMerge_WhenMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "rebase"
	var gotMergeMethod string

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedIsGitHubMergeQueueEnabledGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": null
			}
		}
	}`

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PutReposPullsMergeByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := MergeRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotMergeMethod = body.MergeMethod
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				utils.MustWrite(w, mockedIsGitHubMergeQueueEnabledGQLQueryBody)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}

func TestMerge_WhenMergeIsOnDraftPullRequest(t *testing.T) {
	wantMergeMethod := "rebase"

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Status:  pbc.PullRequestStatus_OPEN,
		IsDraft: true,
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}

func TestMerge_WhenMergeIsClosedPullRequest(t *testing.T) {
	wantMergeMethod := "rebase"

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Status: pbc.PullRequestStatus_CLOSED,
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}

func TestMerge_WhenCheckGitHubMergeQueueIsEnabledRequestFails(t *testing.T) {
	wantMergeMethod := "merge"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				http.Error(w, "IsGitHubMergeQueueEnabledRequestFail", http.StatusNotFound)
			}
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.EqualError(t, err, "non-200 OK status code: 404 Not Found body: \"IsGitHubMergeQueueEnabledRequestFail\\n\"")
}

func TestMerge_WhenGitHubMergeQueueEntriesRequestFails(t *testing.T) {
	wantMergeMethod := "merge"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedIsGitHubMergeQueueEnabledGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test"
				}
			}
		}
	}`

	mockedGitHubMergeQueueEntriesGQLQuery := getMockedGitHubMergeQueueEntriesGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			log.Printf("query: %v\n", query)
			log.Printf("mockedIsGitHubMergeQueueEnabledGQLQuery: %v\n", utils.MinifyQuery(mockedGitHubMergeQueueEntriesGQLQuery))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				utils.MustWrite(w, mockedIsGitHubMergeQueueEnabledGQLQueryBody)
			case utils.MinifyQuery(mockedGitHubMergeQueueEntriesGQLQuery):
				http.Error(w, "GetGitHubMergeQueueEntriesRequestFail", http.StatusNotFound)
			}
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.EqualError(t, err, "non-200 OK status code: 404 Not Found body: \"GetGitHubMergeQueueEntriesRequestFail\\n\"")
}

func TestMerge_WhenGitHubMergeQueueIsONAndPullRequestIsNotOnTheQueue(t *testing.T) {
	wantMergeMethod := "merge"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()

	entityNodeID := aladino.DefaultMockEntityNodeID

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedIsGitHubMergeQueueEnabledGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test"
				}
			}
		}
	}`

	mockedGitHubMergeQueueEntriesGQLQuery := getMockedGitHubMergeQueueEntriesGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedGitHubMergeQueueEntriesGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test",
					"entries": {
						"nodes": []
					}
				}
			}
		}
	}`

	mockedAddPullRequestToGitHubMergeQueueMutation := fmt.Sprintf(`{
        "query": "mutation($input:EnqueuePullRequestInput!) {
            enqueuePullRequest(input: $input) {
                clientMutationId
	        }
	    }",
	    "variables":{
	        "input":{
	            "pullRequestId": "%v"
	        }
	    }
    }`, entityNodeID)

	mockedAddPullRequestToGitHubMergeQueueBody := `{"data": {"enqueuePullRequest": {"clientMutationId": "client_mutation_id"}}}}`

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				utils.MustWrite(w, mockedIsGitHubMergeQueueEnabledGQLQueryBody)
			case utils.MinifyQuery(mockedGitHubMergeQueueEntriesGQLQuery):
				utils.MustWrite(w, mockedGitHubMergeQueueEntriesGQLQueryBody)
			case utils.MinifyQuery(mockedAddPullRequestToGitHubMergeQueueMutation):
				utils.MustWrite(w, mockedAddPullRequestToGitHubMergeQueueBody)
			}
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}

func TestMerge_WhenGitHubMergeQueueIsONAndPullRequestIsNotOnTheQueueAndEnqueueFails(t *testing.T) {
	wantMergeMethod := "merge"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()

	entityNodeID := aladino.DefaultMockEntityNodeID

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedIsGitHubMergeQueueEnabledGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test"
				}
			}
		}
	}`

	mockedGitHubMergeQueueEntriesGQLQuery := getMockedGitHubMergeQueueEntriesGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedGitHubMergeQueueEntriesGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test",
					"entries": {
						"nodes": []
					}
				}
			}
		}
	}`

	mockedAddPullRequestToGitHubMergeQueueMutation := fmt.Sprintf(`{
        "query": "mutation($input:EnqueuePullRequestInput!) {
            enqueuePullRequest(input: $input) {
                clientMutationId
	        }
	    }",
	    "variables":{
	        "input":{
	            "pullRequestId": "%v"
	        }
	    }
    }`, entityNodeID)

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				utils.MustWrite(w, mockedIsGitHubMergeQueueEnabledGQLQueryBody)
			case utils.MinifyQuery(mockedGitHubMergeQueueEntriesGQLQuery):
				utils.MustWrite(w, mockedGitHubMergeQueueEntriesGQLQueryBody)
			case utils.MinifyQuery(mockedAddPullRequestToGitHubMergeQueueMutation):
				http.Error(w, "EnqueuePullRequestFail", http.StatusNotFound)
			}
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}

func TestMerge_WhenPullRequestIsOnGitHubMergeQueue(t *testing.T) {
	wantMergeMethod := "merge"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockBranch := mockedCodeReview.GetBase().GetName()
	mockPrNum := mockedCodeReview.GetNumber()

	mockedIsGitHubMergeQueueEnabledGQLQuery := getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedIsGitHubMergeQueueEnabledGQLQueryBody := `{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test"
				}
			}
		}
	}`

	mockedGitHubMergeQueueEntriesGQLQuery := getMockedGitHubMergeQueueEntriesGQLQuery(mockBranch, mockRepo, mockOwner)

	mockedGitHubMergeQueueEntriesGQLQueryBody := fmt.Sprintf(`{
		"data": {
			"repository": {
				"mergeQueue": {
					"id": "test",
					"entries": {
						"nodes": [
							{
								"pullRequest": {
									"number": %d
								}
							}
						]
					}
				}
			}
		}
	}`, mockPrNum)

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedIsGitHubMergeQueueEnabledGQLQuery):
				utils.MustWrite(w, mockedIsGitHubMergeQueueEnabledGQLQueryBody)
			case utils.MinifyQuery(mockedGitHubMergeQueueEntriesGQLQuery):
				utils.MustWrite(w, mockedGitHubMergeQueueEntriesGQLQueryBody)
			}
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}

func getMockedIsGitHubMergeQueueGQLQuery(mockBranch, mockRepo, mockOwner string) string {
	return fmt.Sprintf(`{
		"query":"query($branchName:String!$repositoryName:String!$repositoryOwner:String!) {
			repository(owner: $repositoryOwner,name: $repositoryName) {
				mergeQueue(branch:$branchName) {
					id
				}
			}
		}",
		"variables":{
			"branchName":"%s",
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockBranch, mockRepo, mockOwner)
}

func getMockedGitHubMergeQueueEntriesGQLQuery(mockBranch, mockRepo, mockOwner string) string {
	return fmt.Sprintf(`{
		"query":"query($branchName:String!$cursor:String$repositoryName:String!$repositoryOwner:String!) {
			repository(owner:$repositoryOwner,name:$repositoryName) {
				mergeQueue(branch:$branchName) {
					id,
					entries(first:100,after:$cursor) {
						nodes {
							pullRequest {
								number
							}
						},
						pageInfo{
							endCursor,
							hasNextPage
						}
					}
				}
			}
		}",
		"variables":{
			"branchName":"%s",
			"cursor":null,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockBranch, mockRepo, mockOwner)
}
