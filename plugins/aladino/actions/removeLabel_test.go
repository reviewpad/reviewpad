// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/gorilla/mux"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var removeLabel = plugins_aladino.PluginBuiltIns().Actions["removeLabel"].Code

func TestRemoveLabel_WhenLabelIsNotAppliedToPullRequest(t *testing.T) {
	wantLabel := "bug"
	var isLabelRemoved bool
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposLabelsByOwnerByRepoByName,
			github.Label{
				Name: github.String(wantLabel),
			},
		),
		mock.WithRequestMatchHandler(
			mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If the remove label request was performed then the label was removed
				isLabelRemoved = true
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(wantLabel)}
	err = removeLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, isLabelRemoved, "The label should not be removed")
}

func TestRemoveLabel_WhenLabelIsAppliedToPullRequestAndLabelIsInEnvironment(t *testing.T) {
	wantLabel := "enhancement"
	var gotLabel string
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposLabelsByOwnerByRepoByName,
			github.Label{
				Name: github.String(wantLabel),
			},
		),
		mock.WithRequestMatchHandler(
			mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				gotLabel = vars["name"]
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	internalLabelID := aladino.BuildInternalLabelID(wantLabel)
	mockedEnv.GetRegisterMap()[internalLabelID] = aladino.BuildStringValue(wantLabel)

	args := []aladino.Value{aladino.BuildStringValue(wantLabel)}
	err = removeLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabel, gotLabel)
}

func TestRemoveLabel_WhenLabelIsAppliedToPullRequestAndLabelIsNotInEnvironment(t *testing.T) {
	wantLabel := "enhancement"
	var gotLabel string
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposLabelsByOwnerByRepoByName,
			github.Label{
				Name: github.String(wantLabel),
			},
		),
		mock.WithRequestMatchHandler(
			mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				gotLabel = vars["name"]
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(wantLabel)}
	err = removeLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabel, gotLabel)
}
