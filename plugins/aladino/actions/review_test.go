// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var review = plugins_aladino.PluginBuiltIns().Actions["review"].Code

type ReviewRequestPostBody struct {
	Event string `json:"event"`
	Body  string `json:"body"`
}

func TestReview_WhenReviewMethodIsUnsupported(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildStringValue("INVALID"), aladino.BuildStringValue("some comment")}
	err := review(mockedEnv, args)

	assert.EqualError(t, err, "review: unsupported review event INVALID")
}

func TestReview_WhenBodyIsNotEmpty(t *testing.T) {
	reviewBody := "test comment"
	tests := map[string]struct {
		inputReviewEvent string
	}{
		"when review event is APPROVE": {
			inputReviewEvent: "APPROVE",
		},
		"when review event is REQUEST_CHANGES": {
			inputReviewEvent: "REQUEST_CHANGES",
		},
		"when review event is COMMENT": {
			inputReviewEvent: "COMMENT",
		},
	}
	for _, test := range tests {
		var (
			gotReviewEvent string
			gotReviewBody  string
		)
		mockedEnv := aladino.MockDefaultEnv(
			t,
			[]mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						rawBody, _ := ioutil.ReadAll(r.Body)
						body := ReviewRequestPostBody{}
						json.Unmarshal(rawBody, &body)
						gotReviewEvent = body.Event
						gotReviewBody = body.Body
					}),
				),
			},
			nil,
			aladino.MockBuiltIns(),
			nil,
		)
		args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue(reviewBody)}
		err := review(mockedEnv, args)
		assert.Nil(t, err)
		assert.Equal(t, test.inputReviewEvent, gotReviewEvent)
		assert.Equal(t, reviewBody, gotReviewBody)
	}
}

func TestReview_WhenBodyIsEmpty(t *testing.T) {
	tests := map[string]struct {
		inputReviewEvent string
		wantErr          string
	}{
		"when review event is APPROVE": {
			inputReviewEvent: "APPROVE",
			wantErr:          "",
		},
		"when review event is REQUEST_CHANGES": {
			inputReviewEvent: "REQUEST_CHANGES",
			wantErr:          "review: comment required in REQUEST_CHANGES event",
		},
		"when review event is COMMENT": {
			inputReviewEvent: "COMMENT",
			wantErr:          "review: comment required in COMMENT event",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
						&ReviewRequestPostBody{},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)
			args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue("")}
			gotErr := review(mockedEnv, args)
			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "review() error = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}
