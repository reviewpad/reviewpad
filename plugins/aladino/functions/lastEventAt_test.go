// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var lastEventAt = plugins_aladino.PluginBuiltIns().Functions["lastEventAt"].Code

func TestLastEvent_WhenIssueTimelineRequestFails(t *testing.T) {
	failMessage := "IssueTimelineRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposIssuesTimelineByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	gotVal, err := lastEventAt(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestLastEvent(t *testing.T) {
	nonLastEventDate := time.Date(2022, 04, 13, 20, 49, 13, 651387237, time.UTC)
	lastEventDate := time.Date(2022, 04, 16, 20, 49, 34, 0, time.UTC)

	tests := map[string]struct {
		mockedEnv aladino.Env
		wantVal   *aladino.IntValue
	}{
		"when last event is not a review": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesTimelineByOwnerByRepoByIssueNumber,
						[]*github.Timeline{
							{
								ID:        github.Int64(6430295168),
								Event:     github.String("locked"),
								CreatedAt: &nonLastEventDate,
							},
							{
								ID:        github.Int64(6430296748),
								Event:     github.String("labeled"),
								CreatedAt: &lastEventDate,
							},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantVal: aladino.BuildIntValue(int(lastEventDate.Unix())),
		},
		"when last event is a review": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesTimelineByOwnerByRepoByIssueNumber,
						[]*github.Timeline{
							{
								ID:        github.Int64(6430295168),
								Event:     github.String("locked"),
								CreatedAt: &nonLastEventDate,
							},
							{
								ID:          github.Int64(6430296748),
								Event:       github.String("reviewed"),
								SubmittedAt: &lastEventDate,
							},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantVal: aladino.BuildIntValue(int(lastEventDate.Unix())),
		},
		"when timeline has only one event": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesTimelineByOwnerByRepoByIssueNumber,
						[]*github.Timeline{
							{ID:        github.Int64(6430295168)},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantVal: aladino.BuildIntValue(int(aladino.DefaultMockPrDate.Unix())),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{}
			gotVal, err := lastEventAt(test.mockedEnv, args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}

}
