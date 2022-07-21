// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v3/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var createdAt = plugins_aladino.PluginBuiltIns().Functions["createdAt"].Code

func TestCreatedAt(t *testing.T) {
	date := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		CreatedAt: &date,
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantCreatedAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", date.String())
	if err != nil {
		log.Fatalf("time.Parse failed: %v", err)
	}
	wantCreatedAt := aladino.BuildIntValue(int(wantCreatedAtTime.Unix()))

	args := []aladino.Value{}
	gotCreatedAt, err := createdAt(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCreatedAt, gotCreatedAt)
}
