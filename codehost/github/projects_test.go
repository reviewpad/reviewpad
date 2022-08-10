// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github_test

import (
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino_actions "github.com/reviewpad/reviewpad/v3/plugins/aladino/actions"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectV2ByName_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
		aladino.MockBuiltIns(),
		nil,
	)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedEnv.GetPullRequest())
	mockRepo := host.GetPullRequestBaseRepoName(mockedEnv.GetPullRequest())
	mockProjectName := "reviewpad"

	project, err := mockedEnv.GetGithubClient().GetProjectV2ByName(mockedEnv.GetCtx(), mockOwner, mockRepo, mockProjectName)

	assert.NotNil(t, err)

	assert.Nil(t, project)
}

func TestGetProjectV2ByName_WhenProjectNotFound(t *testing.T) {
	mockedGetProjectByNameQuery := `{
        "query": "query($name:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectsV2(query: $name, first: 1, orderBy: {field: TITLE, direction: ASC}) {
                    nodes{
                        id,
                        number
                    }
                }
            }
        }",
        "variables": {
            "name": "reviewpad",
            "repositoryName": "default-mock-repo",
            "repositoryOwner": "john"
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
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{},
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(aladino.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				aladino.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedEnv.GetPullRequest())
	mockRepo := host.GetPullRequestBaseRepoName(mockedEnv.GetPullRequest())
	mockProjectName := "reviewpad"

	project, err := mockedEnv.GetGithubClient().GetProjectV2ByName(mockedEnv.GetCtx(), mockOwner, mockRepo, mockProjectName)

	assert.Equal(t, plugins_aladino_actions.ErrProjectNotFound, err)

	assert.Nil(t, project)
}

func TestGetProjectV2ByName_WhenProjectFound(t *testing.T) {
	mockedGetProjectByNameQuery := `{
        "query": "query($name:String! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectsV2(query: $name, first: 1, orderBy: {field: TITLE, direction: ASC}) {
                    nodes{id,number}
                }
            }
        }",
        "variables": {
            "name":"reviewpad",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"john"
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
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{},
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(aladino.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				aladino.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedEnv.GetPullRequest())
	mockRepo := host.GetPullRequestBaseRepoName(mockedEnv.GetPullRequest())
	mockProjectName := "reviewpad"

	project, err := mockedEnv.GetGithubClient().GetProjectV2ByName(mockedEnv.GetCtx(), mockOwner, mockRepo, mockProjectName)

	assert.Equal(t, nil, err)

	assert.NotNil(t, project)
}

func TestGetProjectFieldsByProjectNumber_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
		aladino.MockBuiltIns(),
		nil,
	)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedEnv.GetPullRequest())
	mockRepo := host.GetPullRequestBaseRepoName(mockedEnv.GetPullRequest())
	mockProjectNumber := 1
	mockRetryCount := 1

	project, err := mockedEnv.GetGithubClient().GetProjectFieldsByProjectNumber(mockedEnv.GetCtx(), mockOwner, mockRepo, uint64(mockProjectNumber), mockRetryCount)

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
                            ... on ProjectV2SingleSelectField {
                                id,
                                name,
                                options {
                                    id,
                                    name
                                }
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
            "repositoryOwner":"john"
        }
    }`
	mockedGetProjectByNameQueryBody := `{
        "data": {
            "repository":{
                "projectV2": null
            }
        }
    }`
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{},
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(aladino.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				aladino.MustWrite(
					res,
					mockedGetProjectByNameQueryBody,
				)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedEnv.GetPullRequest())
	mockRepo := host.GetPullRequestBaseRepoName(mockedEnv.GetPullRequest())
	mockProjectNumber := 1
	mockRetryCount := 1

	project, err := mockedEnv.GetGithubClient().GetProjectFieldsByProjectNumber(mockedEnv.GetCtx(), mockOwner, mockRepo, uint64(mockProjectNumber), mockRetryCount)

	assert.Equal(t, plugins_aladino_actions.ErrProjectNotFound, err)

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
                            ... on ProjectV2SingleSelectField {
                                id,
                                name,
                                options {
                                    id,
                                    name
                                }
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
            "repositoryOwner":"john"
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
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{},
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(aladino.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedGetProjectByNameQuery):
				if currentTry == 2 {
					aladino.MustWrite(
						res,
						mockedGetProjectFieldsQueryBody,
					)
					return
				}
				currentTry++
				aladino.MustWrite(res, "")
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedEnv.GetPullRequest())
	mockRepo := host.GetPullRequestBaseRepoName(mockedEnv.GetPullRequest())
	mockProjectNumber := 1
	mockRetryCount := 2

	fields, err := mockedEnv.GetGithubClient().GetProjectFieldsByProjectNumber(mockedEnv.GetCtx(), mockOwner, mockRepo, uint64(mockProjectNumber), mockRetryCount)

	assert.Equal(t, nil, err)

	assert.NotNil(t, fields)

	assert.Equal(t, 1, len(fields))
}
