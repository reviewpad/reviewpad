// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var addLabel = plugins_aladino.PluginBuiltIns().Actions["addLabel"].Code

func TestAddLabel_WhenAddLabelToIssueRequestFails(t *testing.T) {
	label := "bug"
	failMessage := "AddLabelsToIssueRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{
					Name: github.String(label),
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
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

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err := addLabel(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAddLabel_WhenLabelIsInEnvironment(t *testing.T) {
	label := "bug"
	wantLabels := []string{
		label,
	}
	gotLabels := []string{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := []string{}

					json.Unmarshal(rawBody, &body)

					gotLabels = body
				}),
			),
		},
		nil,
	)
	internalLabelID := aladino.BuildInternalLabelID(label)
	mockedEnv.GetRegisterMap()[internalLabelID] = aladino.BuildStringValue(label)

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err := addLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantLabels, gotLabels)
}

func TestAddLabel_WhenLabelIsNotInEnvironment(t *testing.T) {
	label := "bug"
	wantLabels := []string{
		label,
	}
	gotLabels := []string{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := []string{}

					json.Unmarshal(rawBody, &body)

					gotLabels = body
				}),
			),
		},
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err := addLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantLabels, gotLabels)
}
