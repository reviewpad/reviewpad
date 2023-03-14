package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var hasRequiredApprovals = plugins_aladino.PluginBuiltIns().Functions["hasRequiredApprovals"].Code

func TestHasRequiredApprovals_WhenErrorOccurs(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedLatestReviewsGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					latestReviews(last:100){
						nodes{
							author{login},
							state
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber":%d,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		ghGraphQLHandler          func(http.ResponseWriter, *http.Request)
		inputTotalRequiredReviews aladino.Value
		inputRequiredReviewsFrom  aladino.Value
		wantErr                   error
	}{
		"when given total required approvals exceeds the size of the given list of required approvals": {
			inputTotalRequiredReviews: aladino.BuildIntValue(2),
			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john")}),
			wantErr:                   fmt.Errorf("hasRequiredApprovals: the number of required approvals exceeds the number of members from the given list of required approvals"),
		},
		"when get approved reviewers request fails": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					http.Error(w, "GetLatestReviewsRequestFail", http.StatusNotFound)
				}
			},
			inputTotalRequiredReviews: aladino.BuildIntValue(1),
			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john"), aladino.BuildStringValue("test")}),
			wantErr:                   fmt.Errorf("non-200 OK status code: 404 Not Found body: \"GetLatestReviewsRequestFail\\n\""),
		},
	}

	for _, test := range tests {
		mockedEnv := aladino.MockDefaultEnv(t, nil, test.ghGraphQLHandler, aladino.MockBuiltIns(), nil)

		args := []aladino.Value{test.inputTotalRequiredReviews, test.inputRequiredReviewsFrom}
		gotHasRequiredApprovals, gotErr := hasRequiredApprovals(mockedEnv, args)

		assert.Nil(t, gotHasRequiredApprovals)
		assert.Equal(t, test.wantErr, gotErr)
	}
}

func TestHasRequiredApprovals(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedLatestReviewsGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					latestReviews(last:100){
						nodes{
							author{login},
							state
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber":%d,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		clientOptions             []mock.MockBackendOption
		ghGraphQLHandler          func(http.ResponseWriter, *http.Request)
		inputTotalRequiredReviews aladino.Value
		inputRequiredReviewsFrom  aladino.Value
		wantHasRequiredApprovals  aladino.Value
		wantErr                   string
	}{
		"when there is not enough required approvals": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"latestReviews": {
										"nodes": [
											{
												"state": "CHANGES_REQUESTED",
												"author": {
													"login": "test"
												}
											}
										]
									}
								}
							}
						}
					}`)
				}
			},
			inputTotalRequiredReviews: aladino.BuildIntValue(1),
			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("test")}),
			wantHasRequiredApprovals:  aladino.BuildBoolValue(false),
			wantErr:                   "",
		},
		"when there is enough required approvals": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"latestReviews": {
										"nodes": [
											{
												"state": "APPROVED",
												"author": {
													"login": "test"
												}
											}
										]
									}
								}
							}
						}
					}`)
				}
			},
			inputTotalRequiredReviews: aladino.BuildIntValue(1),
			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("test")}),
			wantHasRequiredApprovals:  aladino.BuildBoolValue(true),
			wantErr:                   "",
		},
	}

	for _, test := range tests {
		mockedEnv := aladino.MockDefaultEnv(
			t,
			nil,
			test.ghGraphQLHandler,
			aladino.MockBuiltIns(),
			nil,
		)

		args := []aladino.Value{test.inputTotalRequiredReviews, test.inputRequiredReviewsFrom}
		gotHasRequiredApprovals, gotErr := hasRequiredApprovals(mockedEnv, args)

		assert.Nil(t, gotErr)
		assert.Equal(t, test.wantHasRequiredApprovals, gotHasRequiredApprovals)
	}
}
