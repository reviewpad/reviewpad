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
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/stretchr/testify/assert"
)

var addLabel = plugins_aladino.PluginBuiltIns().Actions["addLabel"].Code

func TestAddLabel_WhenArgumentIsInvalid(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildIntValue(1)}
	err = addLabel(mockedEnv, args)

	assert.EqualError(t, err, "addLabel: expecting string argument, got IntValue")
}

func TestAddLabel_WhenGetLabelRequestFails(t *testing.T) {
	failMessage := "GetLabelRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetReposLabelsByOwnerByRepoByName,
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

	args := []aladino.Value{aladino.BuildStringValue("test")}
	err = addLabel(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAddLabel_WhenAddLabelToIssueRequestFails(t *testing.T) {
	label := "bug"
	failMessage := "AddLabelsToIssueRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err = addLabel(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAddLabel(t *testing.T) {
	label := "bug"
	wantLabels := []string{
		label,
	}
	gotLabels := []string{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err = addLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabels, gotLabels)
}
