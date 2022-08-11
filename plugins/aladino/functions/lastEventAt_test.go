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
	mockedEnv := aladino.MockDefaultEnv(
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
	)

	wantLastEventAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", lastEventDate.String())
	if err != nil {
		assert.FailNow(t, "time.Parse failed", err)
	}
	wantLastEventAt := aladino.BuildIntValue(int(wantLastEventAtTime.Unix()))

	args := []aladino.Value{}
	gotLastEventAt, err := lastEventAt(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLastEventAt, gotLastEventAt)
}
