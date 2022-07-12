// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasLinkedIssues = plugins_aladino.PluginBuiltIns().Functions["hasLinkedIssues"].Code

func TestHasLinkedIssues_WhenRequestFails(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	_, err = hasLinkedIssues(mockedEnv, []aladino.Value{})

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
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
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
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			query := mocks_aladino.MustRead(req.Body)
			switch query {
			case mockedGraphQLQuery:
				mocks_aladino.MustWrite(w, `{"data": {
                    "repository": {
                        "pullRequest": {
                            "closingIssuesReferences": {
                                "totalCount": 3
                            }
                        }
                    }
                }}`)
			}
		},
	)
	if err != nil {
		log.Fatalf("mockDefaultEvalEnvWithGQ failed: %v", err)
	}

	wantVal := aladino.BuildBoolValue(true)
	gotVal, err := hasLinkedIssues(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

// func TestHasLinkedIssues_WhenNoLinkedIssues(t *testing.T) {
// 	wantQuery := fmt.Sprintf(
// 		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){closingIssuesReferences{totalCount}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":%q,\"repositoryOwner\":%q}}\n",
// 		defaultMockPrNum,
// 		defaultMockPrRepoName,
// 		defaultMockPrOwner,
// 	)
// 	gotQuery := ""

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
// 		if got, want := req.Method, http.MethodPost; got != want {
// 			t.Errorf("got request method: %v, want: %v", got, want)
// 		}

// 		gotQuery = mustRead(req.Body)

// 		mustWrite(w, `{"data": {
// 			"repository": {
// 				"pullRequest": {
// 					"closingIssuesReferences": {
// 						"totalCount": 0
// 					}
// 				}
// 			}
// 		}}`)
// 	})

// 	authorLogin := "john"
// 	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
// 		Number: github.Int(defaultMockPrNum),
// 		User:   &github.User{Login: github.String(authorLogin)},
// 		Base: &github.PullRequestBranch{
// 			Repo: &github.Repository{
// 				Owner: &github.User{
// 					Login: github.String(defaultMockPrOwner),
// 				},
// 				Name: github.String(defaultMockPrRepoName),
// 			},
// 		},
// 	})
// 	mockedEnv, err := mocks_aladino.MockDefaultEvalEnvWithGQ(
// 		mux,
// 		mock.WithRequestMatchHandler(
// 			mock.GetReposPullsByOwnerByRepoByPullNumber,
// 			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
// 				w.Write(mock.MustMarshal(mockedPullRequest))
// 			}),
// 		),
// 	)
// 	if err != nil {
// 		log.Fatalf("mockDefaultEvalEnvWithGQ failed: %v", err)
// 	}

// 	wantVal := aladino.BuildBoolValue(false)
// 	gotVal, err := hasLinkedIssues(mockedEnv, []aladino.Value{})

// 	assert.Nil(t, err)
// 	assert.Equal(t, wantQuery, gotQuery)
// 	assert.Equal(t, wantVal, gotVal)
// }

// func mustRead(r io.Reader) string {
// 	b, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(b)
// }

// func mustWrite(w io.Writer, s string) {
// 	_, err := io.WriteString(w, s)
// 	if err != nil {
// 		panic(err)
// 	}
// }
