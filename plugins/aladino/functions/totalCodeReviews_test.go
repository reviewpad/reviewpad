// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var totalCodeReviews = plugins_aladino.PluginBuiltIns().Functions["totalCodeReviews"].Code

func TestTotalCodeReviews_WhenRequestFails(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()
	author := mockedPullRequest.GetUser().GetLogin()

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Total user code reviews request failed", http.StatusBadRequest)
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(author)}
	gotVal, gotErr := totalCodeReviews(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, gotErr, "non-200 OK status code: 400 Bad Request body: \"Total user code reviews request failed\\n\"")
}

func TestTotalCodeReviews(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()
	author := mockedPullRequest.GetUser().GetLogin()
	codeReviewsNum := 80

	mockedTotalPullRequestReviewContributionsQuery := fmt.Sprintf(`{
	    "query": "query($username:String!) {
			user(login:$username){
				contributionsCollection {
					totalPullRequestReviewContributions
				}
			}
	    }",
	    "variables": {
	        "username": "%s"
	    }
	}`, author)

	mockedTotalPullRequestReviewContributionsQueryBody := fmt.Sprintf(`{
		"data": {
			"user": {
				"contributionsCollection": {
					"totalPullRequestReviewContributions": %d
				}
			}
		}
	}`, codeReviewsNum)

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(res http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(aladino.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedTotalPullRequestReviewContributionsQuery):
				aladino.MustWrite(res, mockedTotalPullRequestReviewContributionsQueryBody)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	wantErr := ""
	wantVal := aladino.BuildIntValue(codeReviewsNum)

	args := []aladino.Value{aladino.BuildStringValue(author)}
	gotVal, gotErr := totalCodeReviews(mockedEnv, args)

	if gotErr != nil && gotErr.Error() != wantErr {
		assert.FailNow(t, "totalCodeReviews() error = %v, wantErr %v", gotErr, wantErr)
	}

	assert.Equal(t, wantVal, gotVal)
}
