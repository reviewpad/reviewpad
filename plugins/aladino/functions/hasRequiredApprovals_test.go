package plugins_aladino_functions_test

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/google/go-github/v49/github"
// 	"github.com/migueleliasweb/go-github-mock/src/mock"
// 	"github.com/reviewpad/reviewpad/v4/lang/aladino"
// 	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
// 	"github.com/stretchr/testify/assert"
// )

// var hasRequiredApprovals = plugins_aladino.PluginBuiltIns().Functions["hasRequiredApprovals"].Code

// func TestHasRequiredApprovals_WhenErrorOccurs(t *testing.T) {
// 	failMessage := "ListReviewsRequestFail"

// 	tests := map[string]struct {
// 		clientOptions             []mock.MockBackendOption
// 		inputTotalRequiredReviews aladino.Value
// 		inputRequiredReviewsFrom  aladino.Value
// 		wantErr                   string
// 	}{
// 		"when given total required approvals exceeds the size of the given list of required approvals": {
// 			clientOptions:             []mock.MockBackendOption{},
// 			inputTotalRequiredReviews: aladino.BuildIntValue(2),
// 			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john")}),
// 			wantErr:                   "hasRequiredApprovals: the number of required approvals exceeds the number of members from the given list of required approvals",
// 		},
// 		"when list reviews request fails": {
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatchHandler(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 						mock.WriteError(
// 							w,
// 							http.StatusInternalServerError,
// 							failMessage,
// 						)
// 					}),
// 				),
// 			},
// 			inputTotalRequiredReviews: aladino.BuildIntValue(1),
// 			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john"), aladino.BuildStringValue("mary")}),
// 			wantErr:                   failMessage,
// 		},
// 	}

// 	for _, test := range tests {
// 		mockedEnv := aladino.MockDefaultEnv(t, test.clientOptions, nil, aladino.MockBuiltIns(), nil)

// 		args := []aladino.Value{test.inputTotalRequiredReviews, test.inputRequiredReviewsFrom}
// 		gotHasRequiredApprovals, gotErr := hasRequiredApprovals(mockedEnv, args)

// 		assert.Nil(t, gotHasRequiredApprovals)
// 		assert.NotNil(t, gotErr)

// 		var gotErrMessage string
// 		if err, ok := gotErr.(*github.ErrorResponse); ok {
// 			gotErrMessage = err.Message
// 		} else {
// 			gotErrMessage = gotErr.Error()
// 		}

// 		assert.Equal(t, test.wantErr, gotErrMessage)
// 	}
// }

// func TestHasRequiredApprovals(t *testing.T) {
// 	tests := map[string]struct {
// 		clientOptions             []mock.MockBackendOption
// 		inputTotalRequiredReviews aladino.Value
// 		inputRequiredReviewsFrom  aladino.Value
// 		wantHasRequiredApprovals  aladino.Value
// 		wantErr                   string
// 	}{
// 		"when there is not enough required approvals": {
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{},
// 				),
// 			},
// 			inputTotalRequiredReviews: aladino.BuildIntValue(1),
// 			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john"), aladino.BuildStringValue("mary")}),
// 			wantHasRequiredApprovals:  aladino.BuildBoolValue(false),
// 			wantErr:                   "",
// 		},
// 		"when there is enough required approvals": {
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{
// 						{
// 							ID:    github.Int64(1),
// 							Body:  github.String("Here is the body for the review."),
// 							State: github.String("APPROVED"),
// 							User: &github.User{
// 								Login: github.String("mary"),
// 							},
// 						},
// 					},
// 				),
// 			},
// 			inputTotalRequiredReviews: aladino.BuildIntValue(1),
// 			inputRequiredReviewsFrom:  aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john"), aladino.BuildStringValue("mary")}),
// 			wantHasRequiredApprovals:  aladino.BuildBoolValue(true),
// 			wantErr:                   "",
// 		},
// 	}

// 	for _, test := range tests {
// 		mockedEnv := aladino.MockDefaultEnv(t, test.clientOptions, nil, aladino.MockBuiltIns(), nil)

// 		args := []aladino.Value{test.inputTotalRequiredReviews, test.inputRequiredReviewsFrom}
// 		gotHasRequiredApprovals, gotErr := hasRequiredApprovals(mockedEnv, args)

// 		assert.Nil(t, gotErr)
// 		assert.Equal(t, test.wantHasRequiredApprovals, gotHasRequiredApprovals)
// 	}
// }
