package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var hasRequiredApprovals = plugins_aladino.PluginBuiltIns().Functions["hasRequiredApprovals"].Code

func TestHasRequiredApprovals_WhenErrorOccurs(t *testing.T) {
	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockedCodeReviewNumber := host.GetPullRequestNumber(mockedCodeReview)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)

	mockedLatestReviewsGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					latestReviews: latestOpinionatedReviews(last:100){
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
	}`, mockedCodeReviewNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		ghGraphQLHandler          func(http.ResponseWriter, *http.Request)
		inputTotalRequiredReviews lang.Value
		inputRequiredReviewsFrom  lang.Value
		wantErr                   error
	}{
		"when given total required approvals exceeds the size of the given list of required approvals": {
			inputTotalRequiredReviews: lang.BuildIntValue(2),
			inputRequiredReviewsFrom:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("john")}),
			wantErr:                   fmt.Errorf("hasRequiredApprovals: the number of required approvals exceeds the number of members from the given list of required approvals"),
		},
		"when get approved reviewers request fails": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					http.Error(w, "GetLatestReviewsRequestFail", http.StatusNotFound)
				}
			},
			inputTotalRequiredReviews: lang.BuildIntValue(1),
			inputRequiredReviewsFrom:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("john"), lang.BuildStringValue("test")}),
			wantErr:                   fmt.Errorf("non-200 OK status code: 404 Not Found body: \"GetLatestReviewsRequestFail\\n\""),
		},
	}

	for _, test := range tests {
		mockedEnv := aladino.MockDefaultEnv(t, nil, test.ghGraphQLHandler, aladino.MockBuiltIns(), nil)

		args := []lang.Value{test.inputTotalRequiredReviews, test.inputRequiredReviewsFrom}
		gotHasRequiredApprovals, gotErr := hasRequiredApprovals(mockedEnv, args)

		assert.Nil(t, gotHasRequiredApprovals)
		assert.Equal(t, test.wantErr, gotErr)
	}
}

func TestHasRequiredApprovals(t *testing.T) {
	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockedCodeReviewNumber := host.GetPullRequestNumber(mockedCodeReview)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)

	mockedLatestReviewsGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					latestReviews: latestOpinionatedReviews(last:100){
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
	}`, mockedCodeReviewNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		clientOptions             []mock.MockBackendOption
		ghGraphQLHandler          func(http.ResponseWriter, *http.Request)
		inputTotalRequiredReviews lang.Value
		inputRequiredReviewsFrom  lang.Value
		wantHasRequiredApprovals  lang.Value
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
			inputTotalRequiredReviews: lang.BuildIntValue(1),
			inputRequiredReviewsFrom:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("test")}),
			wantHasRequiredApprovals:  lang.BuildBoolValue(false),
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
			inputTotalRequiredReviews: lang.BuildIntValue(1),
			inputRequiredReviewsFrom:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("test")}),
			wantHasRequiredApprovals:  lang.BuildBoolValue(true),
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

		args := []lang.Value{test.inputTotalRequiredReviews, test.inputRequiredReviewsFrom}
		gotHasRequiredApprovals, gotErr := hasRequiredApprovals(mockedEnv, args)

		assert.Nil(t, gotErr)
		assert.Equal(t, test.wantHasRequiredApprovals, gotHasRequiredApprovals)
	}
}
