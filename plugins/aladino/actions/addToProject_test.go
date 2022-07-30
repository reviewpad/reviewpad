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

func TestAddToProject_WhenProjectNotFound(t *testing.T) {
	prNodeId := "PR_nodeId"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		NodeID: &prNodeId,
	})
	mockedRepositoryProjectQuery := "{\"query\":\"query($name:String!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){projectsV2(query: $name, first: 1, orderBy: {field: TITLE, direction: ASC}){nodes{id,fields(first: 50, orderBy: {field: NAME, direction: ASC}){nodes{... on ProjectV2SingleSelectField{id,name,options{id,name}}}}}}}}\",\"variables\":{\"name\":\"reviewpad\",\"repositoryName\":\"default-mock-repo\",\"repositoryOwner\":\"john\"}}\n"
	mockedAddProjectV2ItemByIdMutation := "{\"query\":\"mutation($input:AddProjectV2ItemByIdInput!){addProjectV2ItemById(input: $input){item{id}}}\",\"variables\":{\"input\":{\"projectId\":\"1\",\"contentId\":\"PR_nodeId\"}}}\n"
	mockedUpdateProjectV2ItemFieldValueMutation := "{\"query\":\"mutation($input:UpdateProjectV2ItemFieldValueInput!){updateProjectV2ItemFieldValue(input: $input){clientMutationId}}\",\"variables\":{\"input\":{\"itemId\":\"item_id\",\"value\":{\"singleSelectOptionId\":\"1\"},\"projectId\":\"1\",\"fieldId\":\"1\"}}}\n"

	gqlTestCases := []struct {
		name                              string
		query                             string
		body                              string
		addProjectV2ItemByIdBody          string
		updateProjectV2ItemFieldValueBody string
		expectedError                     error
	}{
		{
			name:          "when project not found",
			query:         mockedRepositoryProjectQuery,
			body:          `{"data": {"repository":{"projectsV2":{"nodes":[]}}}}`,
			expectedError: plugins_aladino_actions.ErrProjectNotFound,
		},
		{
			name:          "when project has no status field",
			query:         mockedRepositoryProjectQuery,
			body:          `{"data": {"repository":{"projectsV2":{"nodes":[{"id": "1", "fields": {}}]}}}}`,
			expectedError: plugins_aladino_actions.ErrProjectHasNoStatusField,
		},
		{
			name:          "when project status is not found",
			query:         mockedRepositoryProjectQuery,
			body:          `{"data":{"repository":{"projectsV2":{"nodes":[{"id":"1","fields":{"nodes":[{"id":"1","name":"status","options":[{"id":"1","name":"bug"},{"id":"2","name":"feature"}]}]}}]}}}}`,
			expectedError: plugins_aladino_actions.ErrProjectStatusNotFound,
		},
		{
			name:          "add project item error",
			query:         mockedRepositoryProjectQuery,
			body:          `{"data":{"repository":{"projectsV2":{"nodes":[{"id":"1","fields":{"nodes":[{"id":"1","name":"status","options":[{"id":"1","name":"to do"},{"id":"2","name":"in progress"}]}]}}]}}}}`,
			expectedError: io.EOF,
		},
		{
			name:                              "no error",
			query:                             mockedRepositoryProjectQuery,
			body:                              `{"data":{"repository":{"projectsV2":{"nodes":[{"id":"1","fields":{"nodes":[{"id":"1","name":"status","options":[{"id":"1","name":"to do"},{"id":"2","name":"in progress"}]}]}}]}}}}`,
			addProjectV2ItemByIdBody:          `{"data": {"addProjectV2ItemById": {"item": {"id": "item_id"}}}}`,
			updateProjectV2ItemFieldValueBody: `{"data": {"updateProjectV2ItemFieldValue": {"clientMutationId": "client_mutation_id"}}}}`,
			expectedError:                     nil,
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
					switch query {
					case testCase.query:
						aladino.MustWrite(
							res,
							testCase.body,
						)
					case mockedAddProjectV2ItemByIdMutation:
						aladino.MustWrite(
							res,
							testCase.addProjectV2ItemByIdBody,
						)
					case mockedUpdateProjectV2ItemFieldValueMutation:
						aladino.MustWrite(
							res,
							testCase.updateProjectV2ItemFieldValueBody,
						)
					}
				},
			)

			err := addToProject(mockedEnv, []aladino.Value{aladino.BuildStringValue("reviewpad"), aladino.BuildStringValue("to do")})

			assert.Equal(t, testCase.expectedError, err)
		})
	}

}
