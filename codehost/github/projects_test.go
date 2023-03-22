// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectV2ByName_WhenRequestFails(t *testing.T) {
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
	)

	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)
	mockProjectName := "reviewpad"

	project, err := mockedGithubClient.GetProjectV2ByName(context.Background(), mockOwner, mockRepo, mockProjectName)

	assert.NotNil(t, err)

	assert.Nil(t, project)
}

func TestGetProjectV2ByName_WhenProjectNotFound(t *testing.T) {
	mockedGetProjectByNameQuery := `{
        "query": "query($name:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectsV2(query: $name, first: 50, orderBy: {field: TITLE, direction: ASC}) {
                    nodes{
                        id,
                        number,
                        title
                    }
                }
            }
        }",
        "variables": {
            "name": "reviewpad",
            "repositoryName": "default-mock-repo",
            "repositoryOwner": "foobar"
        }
    }`
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository": {
                "projectsV2":{
                    "nodes":[]
                }
            }
        }
    }`
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				utils.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
	)
	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)
	mockProjectName := "reviewpad"

	project, err := mockedGithubClient.GetProjectV2ByName(context.Background(), mockOwner, mockRepo, mockProjectName)

	assert.Equal(t, host.ErrProjectNotFound, err)

	assert.Nil(t, project)
}

func TestGetProjectV2ByName_WhenProjectFound(t *testing.T) {
	mockedGetProjectByNameQuery := `{
        "query": "query($name:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectsV2(query: $name, first: 50, orderBy: {field: TITLE, direction: ASC}) {
                    nodes{id,number,title}
                }
            }
        }",
        "variables": {
            "name":"reviewpad",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"foobar"
        }
    }`
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository":{
                "projectsV2":{
                    "nodes":[
                        {
                            "id": "1",
                            "number": 1
                        }
                    ]
                }
            }
        }
    }`
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				utils.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
	)
	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)
	mockProjectName := "reviewpad"

	project, err := mockedGithubClient.GetProjectV2ByName(context.Background(), mockOwner, mockRepo, mockProjectName)

	assert.Nil(t, err)

	assert.NotNil(t, project)
}

func TestGetProjectFieldsByProjectNumber_WhenRequestFails(t *testing.T) {
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
	)
	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)
	mockProjectNumber := 1
	mockRetryCount := 1

	project, err := mockedGithubClient.GetProjectFieldsByProjectNumber(context.Background(), mockOwner, mockRepo, uint64(mockProjectNumber), mockRetryCount)

	assert.NotNil(t, err)

	assert.Nil(t, project)
}

func TestGetProjectFieldsByProjectNumber_WhenProjectNotFound(t *testing.T) {
	mockedGetProjectByNameQuery := `{
        "query": "query($afterCursor:String! $projectNumber:Int! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectV2(number: $projectNumber) {
                    fields(first: 50, after: $afterCursor, orderBy: {field: NAME, direction: ASC}) {
                        pageInfo {
                            hasNextPage,
                            endCursor
                        },
                        nodes{
                            __typename,
                            ... on ProjectV2SingleSelectField {
                                id,
                                name,
                                options {
                                    id,
                                    name
                                }
                            },
                            ... on ProjectV2Field {
                                id,
                                name,
                                dataType
                            }
                        }
                    }
                }
            }
        }",
        "variables": {
            "afterCursor": "",
            "projectNumber": 1,
            "repositoryName": "default-mock-repo",
            "repositoryOwner":"foobar"
        }
    }`
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository":{
                "projectV2": null
            }
        }
    }`
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			wantedQuery := utils.MinifyQuery(mockedGetProjectByNameQuery)
			switch query {
			case wantedQuery:
				utils.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
	)
	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)
	mockProjectNumber := 1
	mockRetryCount := 1

	project, err := mockedGithubClient.GetProjectFieldsByProjectNumber(context.Background(), mockOwner, mockRepo, uint64(mockProjectNumber), mockRetryCount)

	assert.Equal(t, host.ErrProjectNotFound, err)

	assert.Nil(t, project)

}

func TestGetProjectFieldsByProjectNumber_WhenRetrySuccessful(t *testing.T) {
	mockedGetProjectByNameQuery := `{
        "query":"query($afterCursor:String! $projectNumber:Int! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectV2(number: $projectNumber) {
                    fields(first: 50, after: $afterCursor, orderBy: {field: NAME, direction: ASC}) {
                        pageInfo{
                            hasNextPage,
                            endCursor
                        },
                        nodes{
                            __typename,
                            ... on ProjectV2SingleSelectField {
                                id,
                                name,
                                options {
                                    id,
                                    name
                                }
                            },
                            ... on ProjectV2Field {
                                id,
                                name,
                                dataType
                            }
                        }
                    }
                }
            }
        }",
        "variables":{
            "afterCursor":"",
            "projectNumber":1,
            "repositoryName":
            "default-mock-repo",
            "repositoryOwner":"foobar"
        }
    }`
	mockedGetProjectFieldsQueryBody := `{
        "data": {
            "repository":{
                "projectV2": {
                    "fields": {
                        "pageInfo": {
                            "hasNextPage": false,
                            "endCursor": ""
                        },
                        "nodes": [
                            {
                                "id": "1",
                                "name": "status",
                                "options": []
                            }
                        ]
                    }
                }
            }
        }
    }`
	currentTry := 1
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				if currentTry == 2 {
					utils.MustWrite(
						res,
						mockedGetProjectFieldsQueryBody,
					)
					return
				}
				currentTry++
				utils.MustWrite(res, "")
			}
		},
	)
	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)
	mockProjectNumber := 1
	mockRetryCount := 2

	fields, err := mockedGithubClient.GetProjectFieldsByProjectNumber(context.Background(), mockOwner, mockRepo, uint64(mockProjectNumber), mockRetryCount)

	assert.Equal(t, nil, err)

	assert.NotNil(t, fields)

	assert.Equal(t, 1, len(fields))
}

func TestGetProjectV2ByName_WhenSeveralProjectsFound(t *testing.T) {
	mockProjectName := "reviewpad"
	mockedGetProjectByNameQuery := fmt.Sprintf(`{
        "query": "query($name:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectsV2(query: $name, first: 50, orderBy: {field: TITLE, direction: ASC}) {
                    nodes{id,number,title}
                }
            }
        }",
        "variables": {
            "name":"%s",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"foobar"
        }
    }`, mockProjectName)
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository":{
                "projectsV2":{
                    "nodes":[
                        {
                            "id": "2",
                            "number": 2,
                            "title": "1eviewpad"
                        },
                        {
                            "id": "1",
                            "number": 1,
                            "title": "reviewpad"
                        }
                    ]
                }
            }
        }
    }`
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				utils.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
	)
	mockedCodeReview := aladino.GetDefaultMockReview()

	mockOwner := host.GetCodeReviewBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetCodeReviewBaseRepoName(mockedCodeReview)

	project, err := mockedGithubClient.GetProjectV2ByName(context.Background(), mockOwner, mockRepo, mockProjectName)

	assert.Nil(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, mockProjectName, project.Title)
}
