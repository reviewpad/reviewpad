// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var deleteHeadBranch = plugins_aladino.PluginBuiltIns().Actions["deleteHeadBranch"].Code

func TestDeleteHeadBranch_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					http.Error(w, "reference not found", http.StatusUnprocessableEntity)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	err := deleteHeadBranch(mockedEnv, []aladino.Value{})

	assert.NotNil(t, err)
}

func TestDeleteHeadBranch_Success(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	err := deleteHeadBranch(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
}
