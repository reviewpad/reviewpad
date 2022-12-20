// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var content = plugins_aladino.PluginBuiltIns().Functions["content"].Code

func TestContent(t *testing.T) {
	mockPRJSON, err := json.Marshal(aladino.GetDefaultMockPullRequestDetails())
	assert.Nil(t, err)

	mockIssueJSON, err := json.Marshal(aladino.GetDefaultMockIssueDetails())
	assert.Nil(t, err)

	tests := map[string]struct {
		targetEntity *handler.TargetEntity
		wantErr      error
		wantRes      aladino.Value
	}{
		"when pull request": {
			targetEntity: aladino.DefaultMockTargetEntity,
			wantRes:      aladino.BuildStringValue(string(mockPRJSON)),
		},
		"when issue": {
			targetEntity: &handler.TargetEntity{
				Kind:   handler.Issue,
				Owner:  "test",
				Repo:   "test",
				Number: 1,
			},
			wantRes: aladino.BuildStringValue(string(mockIssueJSON)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnvWithTargetEntity(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesByOwnerByRepoByIssueNumber,
						aladino.GetDefaultMockIssueDetails(),
					),
				},
				func(w http.ResponseWriter, r *http.Request) {},
				aladino.MockBuiltIns(),
				nil,
				test.targetEntity,
			)

			res, err := content(env, []aladino.Value{})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantRes, res)
		})
	}
}
