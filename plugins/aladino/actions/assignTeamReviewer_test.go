// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var assignTeamReviewer = plugins_aladino.PluginBuiltIns().Actions["assignTeamReviewer"].Code

type TeamReviewersRequestPostBody struct {
	TeamReviewers []string `json:"team_reviewers"`
}

func TestAssignTeamReviewer_WhenNoTeamSlugsAreProvided(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})}
	err := assignTeamReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignTeamReviewer: requires at least 1 team to request for review")
}

func TestAssignTeamReviewer(t *testing.T) {
	teamA := "core"
	teamB := "reviewpad-project"
	wantTeamReviewers := []string{teamA, teamB}
	gotTeamReviewers := []string{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := TeamReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotTeamReviewers = body.TeamReviewers
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(teamA), aladino.BuildStringValue(teamB)})}
	err := assignTeamReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTeamReviewers, gotTeamReviewers)
}
