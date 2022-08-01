// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package plugins_aladino_actions_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	plugins_aladino_actions "github.com/reviewpad/reviewpad/v3/plugins/aladino/actions"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var addToProject = plugins_aladino.PluginBuiltIns().Actions["addToProject"].Code

func TestAddToProject_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
	)

	err := addToProject(mockedEnv, []aladino.Value{aladino.BuildStringValue("reviewpad"), aladino.BuildStringValue("to do")})

	assert.NotNil(t, err)
}

func TestAddToProject(t *testing.T) {
	prNodeId := "PR_nodeId"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		NodeID: &prNodeId,
	})
	mockedGetProjectQuery := `{
        "query":"query($name: String! $repositoryName: String! $repositoryOwner: String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectsV2(query: $name, first: 1, orderBy: {field: TITLE, direction: ASC}) {
                    nodes{
                        id,
                        number
                    }
                }
            }
        }",
        "variables":{
            "name":"reviewpad",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"john"
            }
        }`
	mockedGetProjectFieldsQuery := `{
        "query":"query($afterCursor:String !$projectNumber:Int! $repositoryName:String! $repositoryOwner:String!) {
            repository(owner: $repositoryOwner, name: $repositoryName) {
                projectV2(number: $projectNumber) {
                    fields(first: 50, after: $afterCursor, orderBy: {field: NAME, direction: ASC}) {
                        pageInfo {
                            hasNextPage,
                            endCursor
                        },
                        nodes {
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
            "afterCursor":"",
            "projectNumber":1,
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"john"
        }
    }`
	mockedAddProjectV2ItemByIdMutation := `{
        "query": "mutation($input:AddProjectV2ItemByIdInput!) {
            addProjectV2ItemById(input: $input) {
                item {
                    id
                }
            }
        }",
        "variables":{
            "input":{
                "projectId": "1",
                "contentId": "PR_nodeId"
            }
        }
    }`
	mockedUpdateProjectV2ItemFieldValueMutation := `{
        "query": "mutation($input:UpdateProjectV2ItemFieldValueInput!) {
            updateProjectV2ItemFieldValue(input: $input) {
                clientMutationId
            }
        }",
        "variables": {
            "input":{
                "itemId": "item_id",
                "value": {
                    "singleSelectOptionId": "1"
                },
                "projectId": "1",
                "fieldId": "1"
            }
        }
    }`

	gqlTestCases := []struct {
		name                                  string
		getProjectQuery                       string
		getProjectQueryBody                   string
		getProjectFieldsQuery                 string
		getProjectFieldsBody                  string
		addProjectV2ItemByIdMutation          string
		updateProjectV2ItemFieldValueMutation string
		body                                  string
		addProjectV2ItemByIdBody              string
		updateProjectV2ItemFieldValueBody     string
		expectedError                         error
	}{
		{
			name:                "when project not found",
			getProjectQuery:     mockedGetProjectQuery,
			getProjectQueryBody: `{"data": {"repository":{"projectsV2":{"nodes":[]}}}}`,
			expectedError:       utils.ErrProjectNotFound,
		},
		{
			name:                  "error getting project fields",
			getProjectQuery:       mockedGetProjectQuery,
			getProjectQueryBody:   `{"data": {"repository":{"projectsV2":{"nodes":[{"id": "1", "number": 1}]}}}}`,
			getProjectFieldsQuery: mockedGetProjectFieldsQuery,
			getProjectFieldsBody:  ``,
			expectedError:         io.EOF,
		},
		{
			name:                  "when project has no status field",
			getProjectQuery:       mockedGetProjectQuery,
			getProjectQueryBody:   `{"data": {"repository":{"projectsV2":{"nodes":[{"id": "1", "number": 1}]}}}}`,
			getProjectFieldsQuery: mockedGetProjectFieldsQuery,
			getProjectFieldsBody:  `{"data": {"repository":{"projectV2": {"fields": {"pageInfo": {"hasNextPage": false, "endCursor": null}, "nodes": []}}}}}`,
			expectedError:         plugins_aladino_actions.ErrProjectHasNoStatusField,
		},
		{
			name:                  "when project status option is not found",
			getProjectQuery:       mockedGetProjectQuery,
			getProjectQueryBody:   `{"data": {"repository":{"projectsV2":{"nodes":[{"id": "1", "number": 1}]}}}}`,
			getProjectFieldsQuery: mockedGetProjectFieldsQuery,
			getProjectFieldsBody:  `{"data": {"repository":{"projectV2": {"fields": {"pageInfo": {"hasNextPage": false, "endCursor": null}, "nodes": [{"id": "1", "name": "status", "options": []}]}}}}}`,
			expectedError:         plugins_aladino_actions.ErrProjectStatusNotFound,
		},
		{
			name:                         "add project item error",
			getProjectQuery:              mockedGetProjectQuery,
			getProjectQueryBody:          `{"data": {"repository":{"projectsV2":{"nodes":[{"id": "1", "number": 1}]}}}}`,
			getProjectFieldsQuery:        mockedGetProjectFieldsQuery,
			getProjectFieldsBody:         `{"data": {"repository":{"projectV2": {"fields": {"pageInfo": {"hasNextPage": false, "endCursor": null}, "nodes": [{"id": "1", "name": "status", "options": [{"id": "1", "name": "to do"}]}]}}}}}`,
			addProjectV2ItemByIdMutation: mockedAddProjectV2ItemByIdMutation,
			expectedError:                io.EOF,
		},
		{
			name:                                  "no error",
			getProjectQuery:                       mockedGetProjectQuery,
			getProjectQueryBody:                   `{"data": {"repository":{"projectsV2":{"nodes":[{"id": "1", "number": 1}]}}}}`,
			getProjectFieldsQuery:                 mockedGetProjectFieldsQuery,
			getProjectFieldsBody:                  `{"data": {"repository":{"projectV2": {"fields": {"pageInfo": {"hasNextPage": false, "endCursor": null}, "nodes": [{"id": "1", "name": "status", "options": [{"id": "1", "name": "to do"}]}]}}}}}`,
			addProjectV2ItemByIdMutation:          mockedAddProjectV2ItemByIdMutation,
			addProjectV2ItemByIdBody:              `{"data": {"addProjectV2ItemById": {"item": {"id": "item_id"}}}}`,
			updateProjectV2ItemFieldValueMutation: mockedUpdateProjectV2ItemFieldValueMutation,
			updateProjectV2ItemFieldValueBody:     `{"data": {"updateProjectV2ItemFieldValue": {"clientMutationId": "client_mutation_id"}}}}`,
			expectedError:                         nil,
		},
	}

	for _, testCase := range gqlTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
							w.Write(mock.MustMarshal(mockedPullRequest))
						}),
					),
				},
				func(res http.ResponseWriter, req *http.Request) {
					query := aladino.MustRead(req.Body)
					query = utils.MinifyQuery(query)
					switch query {
					case utils.MinifyQuery(testCase.getProjectQuery):
						aladino.MustWrite(res, testCase.getProjectQueryBody)
					case utils.MinifyQuery(testCase.getProjectFieldsQuery):
						aladino.MustWrite(res, testCase.getProjectFieldsBody)
					case utils.MinifyQuery(testCase.addProjectV2ItemByIdMutation):
						aladino.MustWrite(res, testCase.addProjectV2ItemByIdBody)
					case utils.MinifyQuery(testCase.updateProjectV2ItemFieldValueMutation):
						aladino.MustWrite(res, testCase.updateProjectV2ItemFieldValueBody)
					}
				},
			)

			err := addToProject(mockedEnv, []aladino.Value{aladino.BuildStringValue("reviewpad"), aladino.BuildStringValue("to do")})

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
