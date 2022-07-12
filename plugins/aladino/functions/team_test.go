// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var team = plugins_aladino.PluginBuiltIns().Functions["team"].Code

func TestTeam_WhenListTeamMembersBySlugRequestFails(t *testing.T) {
	teamSlug := "reviewpad-team"
	failMessage := "ListTeamMembersBySlugRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetOrgsTeamsMembersByOrgByTeamSlug,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(teamSlug)}
	gotMembers, err := team(mockedEnv, args)

	assert.Nil(t, gotMembers)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestTeam(t *testing.T) {
	teamSlug := "reviewpad-team"
    ghMembers := []*github.User{
        {Login: github.String("john")},
        {Login: github.String("jane")},
    }
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetOrgsTeamsMembersByOrgByTeamSlug,
			ghMembers,
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    members := make([]aladino.Value, len(ghMembers))
    for i, ghMember := range ghMembers {
		members[i] = aladino.BuildStringValue(ghMember.GetLogin())
	}

    wantMembers := aladino.BuildArrayValue(members)

	args := []aladino.Value{aladino.BuildStringValue(teamSlug)}
	gotMembers, err := team(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMembers, gotMembers)
}
