// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var setProjectField = plugins_aladino.PluginBuiltIns().Actions["setProjectField"].Code

func TestSetProjectField_WhenGetLinkedProjectsRequestFails(t *testing.T) {
	failMessageForPullRequest := "GetLinkedProjectsForPullRequestRequestFail"
	failMessageForIssue := "GetLinkedProjectsForIssueRequestFail"

	mockedPRRepoName := aladino.DefaultMockPrRepoName
	mockedPRNum := aladino.DefaultMockPrNum
	mockedPROwner := aladino.DefaultMockPrOwner

	mockedGetLinkedProjectsForPullRequestGraphQLQuery := getMockedGetLinkedProjectsForPullRequestGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForIssueGraphQLQuery := getMockedGetLinkedProjectsForIssueGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	tests := map[string]struct {
		mockedEnv aladino.Env
		wantErr   string
	}{
		"when get linked projects for pull request fails": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForPullRequestGraphQLQuery):
						http.Error(w, failMessageForPullRequest, http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.PullRequest,
				},
			),
			wantErr: fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessageForPullRequest),
		},
		"when get linked projects for issue fails": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForIssueGraphQLQuery):
						http.Error(w, failMessageForIssue, http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.Issue,
				},
			),
			wantErr: fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessageForIssue),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []lang.Value{lang.BuildStringValue("reviewpad"), lang.BuildStringValue("size"), lang.BuildStringValue("50")}
			gotErr := setProjectField(test.mockedEnv, args)

			assert.EqualError(t, gotErr, test.wantErr)
		})
	}
}

func TestSetProjectField_WhenEntityIsInProject(t *testing.T) {
	mockedPRRepoName := aladino.DefaultMockPrRepoName
	mockedPRNum := aladino.DefaultMockPrNum
	mockedPROwner := aladino.DefaultMockPrOwner

	mockedGetLinkedProjectsForPullRequestGraphQLQuery := getMockedGetLinkedProjectsForPullRequestGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForPullRequestGraphQLQueryBody := `{
		"data": {
				"repository": {
				"pullRequest": {
					"projectItems": {
						"nodes": [
							{
								"id": "test",
								"project": {
									"id": "test",
									"number": 1,
									"title": "test"
								}
							}
						]
					}
				}
			}
		}
	}`

	mockedGetLinkedProjectsForIssueGraphQLQuery := getMockedGetLinkedProjectsForIssueGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForIssueGraphQLQueryBody := `{
		"data": {
				"repository": {
				"issue": {
					"projectItems": {
						"nodes": [
							{
								"id": "test",
								"project": {
									"id": "test",
									"number": 1,
									"title": "test"
								}
							}
						]
					}
				}
			}
		}
	}`

	mockedGetProjectFieldsByProjectNumberGraphQLQuery := getMockedProjectFieldsByProjectNumberGraphQLQuery(mockedPRRepoName, mockedPROwner)

	mockedSetProjectSingleSelectFieldMutation := getMockedUpdateProjectV2ItemFieldValueMutation()

	mockedGetProjectFieldsByProjectNumberGraphQLQueryBody := `{
		"data": {
			"repository": {
				"projectV2": {
					"fields": {
						"nodes": [
							{
								"__typename": "ProjectV2SingleSelectField",
								"id": "test",
								"name": "size",
								"options": [
									{
										"id": "test",
										"name": "50"
									}
								]
							}
						]
					}
				}
			}
		}
	}`

	mockedSetProjectSingleSelectFieldMutationBody := `{
		"data": {
			"updateProjectV2ItemFieldValue": {
				"clientMutationId": "test"
			}
		}
	}`

	tests := map[string]struct {
		mockedEnv aladino.Env
		wantErr   string
	}{
		"when pull request is already in the project": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					log.Printf("query: %s", query)
					log.Printf("mockedSetProjectSingleSelectFieldMutation: %s", utils.MinifyQuery(mockedSetProjectSingleSelectFieldMutation))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForPullRequestGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForPullRequestGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectFieldsByProjectNumberGraphQLQuery):
						utils.MustWrite(w, mockedGetProjectFieldsByProjectNumberGraphQLQueryBody)
					case utils.MinifyQuery(mockedSetProjectSingleSelectFieldMutation):
						utils.MustWrite(w, mockedSetProjectSingleSelectFieldMutationBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.PullRequest,
				},
			),
		},
		"when issue is already in the project": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForIssueGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForIssueGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectFieldsByProjectNumberGraphQLQuery):
						utils.MustWrite(w, mockedGetProjectFieldsByProjectNumberGraphQLQueryBody)
					case utils.MinifyQuery(mockedSetProjectSingleSelectFieldMutation):
						utils.MustWrite(w, mockedSetProjectSingleSelectFieldMutationBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.Issue,
				},
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []lang.Value{lang.BuildStringValue("test"), lang.BuildStringValue("size"), lang.BuildStringValue("50")}
			gotErr := setProjectField(test.mockedEnv, args)

			assert.Nil(t, gotErr)
		})
	}
}

func TestSetProjectField_WhenSetProjectFieldFails(t *testing.T) {
	mockedPRRepoName := aladino.DefaultMockPrRepoName
	mockedPRNum := aladino.DefaultMockPrNum
	mockedPROwner := aladino.DefaultMockPrOwner

	mockedGetLinkedProjectsForPullRequestGraphQLQuery := getMockedGetLinkedProjectsForPullRequestGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForPullRequestGraphQLQueryBody := `{
		"data": {
			"repository": {
				"pullRequest": {
					"projectItems": {
						"nodes": [
							{
								"id": "test",
								"project": {
									"id": "test",
									"number": 1,
									"title": "test"
								}
							}
						]
					}
				}
			}
		}
	}`

	mockedGetLinkedProjectsForIssueGraphQLQuery := getMockedGetLinkedProjectsForIssueGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForIssueGraphQLQueryBody := `{
		"data": {
			"repository": {
				"issue": {
					"projectItems": {
						"nodes": [
							{
								"id": "test",
								"project": {
									"id": "test",
									"number": 1,
									"title": "test"
								}
							}
						]
					}
				}
			}
		}
	}`

	mockedGetProjectFieldsByProjectNumberGraphQLQuery := getMockedProjectFieldsByProjectNumberGraphQLQuery(mockedPRRepoName, mockedPROwner)

	mockedGetProjectFieldsByProjectNumberGraphQLQueryBody := `{
		"data": {
			"repository": {
				"projectV2": {
					"fields": {
						"nodes": [
							{
								"__typename": "ProjectV2SingleSelectField",
								"id": "test",
								"name": "size",
								"options": [
									{
										"id": "test",
										"name": "50"
									}
								]
							}
						]
					}
				}
			}
		}
	}`

	mockedSetProjectSingleSelectFieldMutation := getMockedUpdateProjectV2ItemFieldValueMutation()

	tests := map[string]struct {
		mockedEnv aladino.Env
		wantErr   string
	}{
		"when set project field for pull request fails because the request to get the project fields by project number fails": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForPullRequestGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForPullRequestGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectFieldsByProjectNumberGraphQLQuery):
						http.Error(w, "GetProjectFieldsByProjectNumberRequestFail", http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.PullRequest,
				},
			),
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetProjectFieldsByProjectNumberRequestFail\\n\"",
		},
		"when set project field for pull request fails because the request to set the project single select field fails": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForPullRequestGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForPullRequestGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectFieldsByProjectNumberGraphQLQuery):
						utils.MustWrite(w, mockedGetProjectFieldsByProjectNumberGraphQLQueryBody)
					case utils.MinifyQuery(mockedSetProjectSingleSelectFieldMutation):
						http.Error(w, "SetProjectSingleSelectFieldRequestFail", http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.PullRequest,
				},
			),
			wantErr: "non-200 OK status code: 404 Not Found body: \"SetProjectSingleSelectFieldRequestFail\\n\"",
		},
		"when set project field for issue fails because the request to get the project fields by project number fails": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForIssueGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForIssueGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectFieldsByProjectNumberGraphQLQuery):
						http.Error(w, "GetProjectFieldsByProjectNumberRequestFail", http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.Issue,
				},
			),
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetProjectFieldsByProjectNumberRequestFail\\n\"",
		},
		"when set project field for issue fails because the request to set the project single select field fails": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForIssueGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForIssueGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectFieldsByProjectNumberGraphQLQuery):
						utils.MustWrite(w, mockedGetProjectFieldsByProjectNumberGraphQLQueryBody)
					case utils.MinifyQuery(mockedSetProjectSingleSelectFieldMutation):
						http.Error(w, "SetProjectSingleSelectFieldRequestFail", http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.Issue,
				},
			),
			wantErr: "non-200 OK status code: 404 Not Found body: \"SetProjectSingleSelectFieldRequestFail\\n\"",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []lang.Value{lang.BuildStringValue("test"), lang.BuildStringValue("size"), lang.BuildStringValue("50")}
			gotErr := setProjectField(test.mockedEnv, args)

			assert.EqualError(t, gotErr, test.wantErr)
		})
	}
}

func TestSetProjectField_WhenEntityIsNotInProjectAndAddToProjectFails(t *testing.T) {
	mockedPRRepoName := aladino.DefaultMockPrRepoName
	mockedPRNum := aladino.DefaultMockPrNum
	mockedPROwner := aladino.DefaultMockPrOwner

	mockedGetLinkedProjectsForPullRequestGraphQLQuery := getMockedGetLinkedProjectsForPullRequestGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForPullRequestGraphQLQueryBody := `{
		"data": {
				"repository": {
				"pullRequest": {
					"projectItems": {
						"nodes": []
					}
				}
			}
		}
	}`

	mockedGetLinkedProjectsForIssueGraphQLQuery := getMockedGetLinkedProjectsForIssueGraphQLQuery(mockedPRNum, mockedPRRepoName, mockedPROwner)

	mockedGetLinkedProjectsForIssueGraphQLQueryBody := `{
		"data": {
				"repository": {
				"issue": {
					"projectItems": {
						"nodes": []
					}
				}
			}
		}
	}`

	mockedGetProjectQuery := `{
        "query":"query($name: String! $repositoryName: String! $repositoryOwner: String!) {
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
        "variables":{
            "name":"test",
            "repositoryName":"default-mock-repo",
            "repositoryOwner":"foobar"
            }
        }`

	tests := map[string]struct {
		mockedEnv aladino.Env
		wantErr   string
	}{
		"when pull request is already in the project": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForPullRequestGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForPullRequestGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectQuery):
						http.Error(w, "GetProjectV2ByName", http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.PullRequest,
				},
			),
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetProjectV2ByName\\n\"",
		},
		"when issue is already in the project": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedGetLinkedProjectsForIssueGraphQLQuery):
						utils.MustWrite(w, mockedGetLinkedProjectsForIssueGraphQLQueryBody)
					case utils.MinifyQuery(mockedGetProjectQuery):
						http.Error(w, "GetProjectV2ByName", http.StatusNotFound)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&entities.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   entities.Issue,
				},
			),
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetProjectV2ByName\\n\"",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []lang.Value{lang.BuildStringValue("test"), lang.BuildStringValue("size"), lang.BuildStringValue("50")}
			gotErr := setProjectField(test.mockedEnv, args)

			assert.EqualError(t, gotErr, test.wantErr)
		})
	}
}

func getMockedGetLinkedProjectsForIssueGraphQLQuery(mockedPRNum int, mockedPRRepoName, mockedPROwner string) string {
	return fmt.Sprintf(`{
		"query": "query($issueNumber:Int!$projectItemsCursor:String!$repositoryName:String!$repositoryOwner:String!) {
			repository(owner: $repositoryOwner, name: $repositoryName) {
				issue(number: $issueNumber) {
					projectItems(first: 10, after: $projectItemsCursor) {
						nodes {
							id,
							project {
								id,
								number,
								title
							}
						},
						pageInfo {
							endCursor,
							hasNextPage
						}
					}
				}
			}
		}",
		"variables": {
			"issueNumber": %d,
			"projectItemsCursor": "",
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedPRNum, mockedPRRepoName, mockedPROwner)
}

func getMockedGetLinkedProjectsForPullRequestGraphQLQuery(mockedPRNum int, mockedPRRepoName, mockedPROwner string) string {
	return fmt.Sprintf(`{
		"query": "query($issueNumber:Int!$projectItemsCursor:String!$repositoryName:String!$repositoryOwner:String!) {
			repository(owner: $repositoryOwner, name: $repositoryName) {
				pullRequest(number: $issueNumber) {
					projectItems(first: 10, after: $projectItemsCursor) {
						nodes {
							id,
							project {
								id,
								number,
								title
							}
						},
						pageInfo {
							endCursor,
							hasNextPage
						}
					}
				}
			}
		}",
		"variables": {
			"issueNumber": %d,
			"projectItemsCursor": "",
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedPRNum, mockedPRRepoName, mockedPROwner)
}

func getMockedProjectFieldsByProjectNumberGraphQLQuery(mockedPRRepoName, mockedPROwner string) string {
	return fmt.Sprintf(`{
		"query": "query($afterCursor:String!$projectNumber:Int!$repositoryName:String!$repositoryOwner:String!) {
			repository(owner: $repositoryOwner, name: $repositoryName) {
				projectV2(number: $projectNumber) {
					fields(first: 50, after: $afterCursor, orderBy: {field: NAME, direction: ASC}) {
						pageInfo {
							hasNextPage,
							endCursor
						},
						nodes {
							__typename,
							...on ProjectV2SingleSelectField {
								id,
								name,
								options {
									id,
									name
								}
							},
							...on ProjectV2Field {
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
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedPRRepoName, mockedPROwner)
}

func getMockedUpdateProjectV2ItemFieldValueMutation() string {
	return fmt.Sprintf(`{
		"query": "mutation($input: UpdateProjectV2ItemFieldValueInput!) {
			updateProjectV2ItemFieldValue(input: $input) {
				clientMutationId
			}
		}",
		"variables": {
			"input": {
				"itemId": "test",
				"value": {
					"singleSelectOptionId":"test"
				},
				"projectId":"test",
				"fieldId":"test"
			}
		}
	}`)
}
