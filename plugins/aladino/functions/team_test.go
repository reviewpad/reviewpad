// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var team = plugins_aladino.PluginBuiltIns().Functions["team"].Code

func TestTeam_WhenListTeamMembersBySlugRequestFails(t *testing.T) {
	teamSlug := "reviewpad-team"
	failMessage := "ListTeamMembersBySlugRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(teamSlug)}
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
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetOrgsTeamsMembersByOrgByTeamSlug,
				ghMembers,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantMembers := lang.BuildArrayValue([]lang.Value{
		lang.BuildStringValue("john"),
		lang.BuildStringValue("jane"),
	})

	args := []lang.Value{lang.BuildStringValue(teamSlug)}
	gotMembers, err := team(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMembers, gotMembers)
}
