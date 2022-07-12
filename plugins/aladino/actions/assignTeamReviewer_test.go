// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignTeamReviewer = plugins_aladino.PluginBuiltIns().Actions["assignTeamReviewer"].Code

type TeamReviewersRequestPostBody struct {
	TeamReviewers []string `json:"team_reviewers"`
}

func TestAssignTeamReviewer_WhenNoTeamSlugsAreProvided(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})}
	err = assignTeamReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignTeamReviewer: requires at least 1 team to request for review")
}

func TestAssignTeamReviewer(t *testing.T) {
	teamA := "core"
	teamB := "reviewpad-project"
	wantTeamReviewers := []string{teamA, teamB}
	gotTeamReviewers := []string{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := TeamReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotTeamReviewers = body.TeamReviewers
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(teamA), aladino.BuildStringValue(teamB)})}
	err = assignTeamReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTeamReviewers, gotTeamReviewers)
}
