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
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var team = plugins_aladino.PluginBuiltIns().Functions["team"].Code

func TestTeam_WhenListTeamMembersBySlugRequestFails(t *testing.T) {
	teamSlug := "reviewpad-team"
	failMessage := "ListTeamMembersBySlugRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
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
		},
		nil,
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
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetOrgsTeamsMembersByOrgByTeamSlug,
				ghMembers,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantMembers := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("john"),
		aladino.BuildStringValue("jane"),
	})

	args := []aladino.Value{aladino.BuildStringValue(teamSlug)}
	gotMembers, err := team(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMembers, gotMembers)
}
