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

var organization = plugins_aladino.PluginBuiltIns().Functions["organization"].Code

func TestOrganization_WhenListMembersRequestFails(t *testing.T) {
	failMessage := "ListMembersRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetOrgsPublicMembersByOrg,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetOrgsMembersByOrg,
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

	args := []aladino.Value{}
	gotMembers, err := organization(mockedEnv, args)

	assert.Nil(t, gotMembers)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestOrganization(t *testing.T) {
	ghMembers := []*github.User{
		{Login: github.String("john")},
		{Login: github.String("jane")},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetOrgsPublicMembersByOrg,
				ghMembers,
			),
			mock.WithRequestMatch(
				mock.GetOrgsMembersByOrg,
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

	args := []aladino.Value{}
	gotMembers, err := organization(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMembers, gotMembers)
}
