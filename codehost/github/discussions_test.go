// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github_test

import (
	"context"
	"net/http"
	"testing"

	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetDiscussionComments(t *testing.T) {
	mockedGetDiscussionCommentsQuery := `{
        "query": "query($afterCursor:String $name:String! $number:Int! $owner:String!) {
            repository(owner: $owner, name: $name) {
                discussion(number: $number) {
                    id, 
                    body,
                    author {
                        login
                    },
                    comments(first: 50, after: $afterCursor) {
                        pageInfo {
                            hasNextPage,
                            endCursor
                        },
                        nodes {
                            id,
                            body,
                            author {
                                login
                            }
                        }
                    }
                }
            }
        }",
        "variables": {
            "afterCursor": null,
            "name":"default-mock-repo",
            "number":82,
            "owner":"foobar"
        }
    }`
	mockedGetDiscussionCommentsQueryBody := `{
        "data": {
            "repository": {
              "discussion": {
                "id": "D_kwDOJRdjrs4AT7xD",
                "body": "This is a test discussion",
                "author": {
                  "login": "marcelosousa"
                },
                "comments": {
                  "nodes": [
                    {
                      "id": "DC_kwDOJRdjrs4AWz6g",
                      "body": "Reply to a discussion",
                      "author": {
                        "login": "marcelosousa"
                      }
                    },
                    {
                      "id": "DC_kwDOJRdjrs4AW0Hs",
                      "body": "calling from GraphiQL",
                      "author": {
                        "login": "reviewpad-bot"
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
			wantQuery := utils.MinifyQuery(mockedGetDiscussionCommentsQuery)
			switch query {
			case wantQuery:
				utils.MustWrite(
					res,
					mockedGetDiscussionCommentsQueryBody,
				)
			}
		},
	)
	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockDiscussionNum := 82

	project, err := mockedGithubClient.GetDiscussionComments(context.Background(), mockOwner, mockRepo, mockDiscussionNum)

	assert.Nil(t, err)

	assert.NotNil(t, project)
}
