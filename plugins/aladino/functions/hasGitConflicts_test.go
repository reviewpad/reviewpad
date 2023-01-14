// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var hasGitConflicts = plugins_aladino.PluginBuiltIns().Functions["hasGitConflicts"].Code

func TestHasGitConflicts_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Pull Request request failed", http.StatusBadRequest)
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	gotVal, gotErr := hasGitConflicts(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, gotErr, "non-200 OK status code: 400 Bad Request body: \"Pull Request request failed\\n\"")
}

func TestHasGitConflicts(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedPullRequestQuery := fmt.Sprintf(`{
	    "query": "query($pullRequestNumber:Int! $repositoryName:String! $repositoryOwner:String!) {
	        repository(owner: $repositoryOwner, name: $repositoryName) {
	            pullRequest(number: $pullRequestNumber) {
	                mergeable
	            }
	        }
	    }",
	    "variables": {
	        "pullRequestNumber": %d,
            "repositoryName": "%s",
	        "repositoryOwner": "%s"
	    }
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		mockedPullRequestQueryBody string
		wantVal                    aladino.Value
		wantErr                    string
	}{
		"when pull request has git conflicts": {
			mockedPullRequestQueryBody: pullRequestQueryBodyWith("CONFLICTING"),
			wantVal:                    aladino.BuildBoolValue(true),
		},
		"when pull request has no git conflicts and is mergeable": {
			mockedPullRequestQueryBody: pullRequestQueryBodyWith("MERGEABLE"),
			wantVal:                    aladino.BuildBoolValue(false),
		},
		"when pull request has no git conflicts and is not mergeable": {
			mockedPullRequestQueryBody: pullRequestQueryBodyWith("UNKNOWN"),
			wantVal:                    aladino.BuildBoolValue(false),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				nil,
				func(res http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedPullRequestQuery):
						utils.MustWrite(res, test.mockedPullRequestQueryBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			)

			args := []aladino.Value{}
			gotVal, gotErr := hasGitConflicts(mockedEnv, args)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "hasGitConflicts() error = %v, wantErr %v", gotErr, test.wantErr)
			}

			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}

func pullRequestQueryBodyWith(state string) string {
	return fmt.Sprintf(`{
        "data": {
            "repository":{
                "pullRequest":{
                    "mergeable": "%s"
                }
            }
        }
    }`, state)
}
