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

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var addLabel = plugins_aladino.PluginBuiltIns().Actions["addLabel"].Code

func TestAddLabel_WhenAddLabelToIssueRequestFails(t *testing.T) {
	label := "bug"
	failMessage := "AddLabelsToIssueRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err = addLabel(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAddLabel_WhenLabelIsInEnvironment(t *testing.T) {
	label := "bug"
	wantLabels := []string{
		label,
	}
	gotLabels := []string{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}
	internalLabelID := aladino.BuildInternalLabelID(label)
	mockedEnv.GetRegisterMap()[internalLabelID] = aladino.BuildStringValue(label)

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err = addLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantLabels, gotLabels)
}

func TestAddLabel_WhenLabelIsNotInEnvironment(t *testing.T) {
	label := "bug"
	wantLabels := []string{
		label,
	}
	gotLabels := []string{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err = addLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantLabels, gotLabels)
}
