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

var merge = plugins_aladino.PluginBuiltIns().Actions["merge"].Code

type MergeRequestPostBody struct {
	MergeMethod string `json:"merge_method"`
}

func TestMerge_WhenMergeMethodIsUnsupported(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue("INVALID")}
	err = merge(mockedEnv, args)

	assert.EqualError(t, err, "merge: unsupported merge method INVALID")
}

func TestMerge_WhenNoMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "merge"
	var gotMergeMethod string
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PutReposPullsMergeByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := MergeRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotMergeMethod = body.MergeMethod
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}

func TestMerge_WhenMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "rebase"
	var gotMergeMethod string
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PutReposPullsMergeByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := MergeRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotMergeMethod = body.MergeMethod
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(wantMergeMethod)}
	err = merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}
