package codehost_test

import (
	"testing"
	"time"

	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/stretchr/testify/assert"
)

func TestHasReview(t *testing.T) {
	userLogin := "john"

	tests := map[string]struct {
		inputReviews      []*codehost.Review
		expectedHasReview bool
	}{
		"when there are no reviews": {
			inputReviews:      []*codehost.Review{},
			expectedHasReview: false,
		},
		"when user has done a review": {
			inputReviews: []*codehost.Review{
				{
					User: &codehost.User{
						Login: userLogin,
					},
				},
			},
			expectedHasReview: true,
		},
		"when user has done no review": {
			inputReviews:      []*codehost.Review{},
			expectedHasReview: false,
		},
	}

	for _, test := range tests {
		gotHasReview := codehost.HasReview(test.inputReviews, userLogin)

		assert.Equal(t, test.expectedHasReview, gotHasReview)
	}
}

func TestLastReview(t *testing.T) {
	userLogin := "john"

	reviewATime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	reviewA := &codehost.Review{
		User: &codehost.User{
			Login: userLogin,
		},
		SubmittedAt: &reviewATime,
	}

	reviewBTime := time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)
	reviewB := &codehost.Review{
		User: &codehost.User{
			Login: userLogin,
		},
		SubmittedAt: &reviewBTime,
	}

	tests := map[string]struct {
		inputReviews       []*codehost.Review
		expectedLastReview *codehost.Review
	}{
		"when there are no reviews": {
			inputReviews:       []*codehost.Review{},
			expectedLastReview: nil,
		},
		"when user has done one review": {
			inputReviews: []*codehost.Review{
				reviewA,
			},
			expectedLastReview: reviewA,
		},
		"when user has done more than one review": {
			inputReviews: []*codehost.Review{
				reviewA,
				reviewB,
			},
			expectedLastReview: reviewB,
		},
	}

	for _, test := range tests {
		gotLastReview := codehost.LastReview(test.inputReviews, userLogin)

		assert.Equal(t, test.expectedLastReview, gotLastReview)
	}
}
