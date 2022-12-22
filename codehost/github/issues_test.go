// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github_test

import (
	"context"
	"net/http"
	"testing"

	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetLinkedProjects_WhenRequestFails(t *testing.T) {
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
	)

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)
	mockIssuerNumber := 10

	projects, err := mockedGithubClient.GetLinkedProjectsForIssue(context.Background(), mockOwner, mockRepo, mockIssuerNumber, 3)

	assert.NotNil(t, err)

	assert.Nil(t, projects)
}

func TestGetLinkedProjects_WhenProjectItemsNotFound(t *testing.T) {
	mockedGetLinkedProjectsQuery := `{
        "query": "query($issueNumber:Int! $projectItemsCursor:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                issue(number: $issueNumber) {
                    projectItems(first: 10, after: $projectItemsCursor) {
                        nodes {
                            id,
                            project {id,number,title}
                        },
                        pageInfo {endCursor,hasNextPage}
                    }
                }
            }
        }",
        "variables": {
            "issueNumber":10,
            "projectItemsCursor": "",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"foobar"
        }
    }`
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository":{
                "issue":{
                }
            }
        }
    }`
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			mockedQuery := utils.MinifyQuery(mockedGetLinkedProjectsQuery)
			switch query {
			case mockedQuery:
				utils.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
	)
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)
	mockIssuerNumber := 10

	projects, err := mockedGithubClient.GetLinkedProjectsForIssue(context.Background(), mockOwner, mockRepo, mockIssuerNumber, 3)

	assert.NotNil(t, err)

	assert.Nil(t, projects)
}

func TestGetLinkedProjects_WhenProjectFound(t *testing.T) {
	mockedGetLinkedProjectsQuery := `{
        "query": "query($issueNumber:Int! $projectItemsCursor:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                issue(number: $issueNumber) {
                    projectItems(first: 10, after: $projectItemsCursor) {
                        nodes {
                            id,
                            project {id,number,title}
                        },
                        pageInfo {endCursor,hasNextPage}
                    }
                }
            }
        }",
        "variables": {
            "issueNumber":10,
            "projectItemsCursor": "",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"foobar"
        }
    }`
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository":{
                "issue":{
                    "projectItems":{
                        "nodes":[
                            {
                                "id": "1",
                                "project":{
                                    "id": "1",
                                    "number": 1,
                                    "title": "Project 1"
                                }
                            }
                        ]
                    }
                }
            }
        }
    }`
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			mockedQuery := utils.MinifyQuery(mockedGetLinkedProjectsQuery)
			switch query {
			case mockedQuery:
				utils.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
	)
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)
	mockIssuerNumber := 10

	gotProjects, err := mockedGithubClient.GetLinkedProjectsForIssue(context.Background(), mockOwner, mockRepo, mockIssuerNumber, 3)

	wantProjects := []host.GQLProjectV2Item{
		{
			ID: "1",
			Project: host.ProjectV2{
				ID:     "1",
				Number: 1,
				Title:  "Project 1",
			},
		},
	}

	assert.Nil(t, err)

	assert.NotNil(t, gotProjects)
	assert.Equal(t, wantProjects, gotProjects)
}
