// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var merge = plugins_aladino.PluginBuiltIns(plugins_aladino.DefaultPluginConfig()).Actions["merge"].Code

type MergeRequestPostBody struct {
	MergeMethod string `json:"merge_method"`
}

func TestMerge_WhenMergeMethodIsUnsupported(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildStringValue("INVALID")}
	err := merge(mockedEnv, args)

	assert.EqualError(t, err, "merge: unsupported merge method INVALID")
}

func TestMerge_WhenNoMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "merge"
	var gotMergeMethod string
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}

func TestMerge_WhenMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "rebase"
	var gotMergeMethod string
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}
