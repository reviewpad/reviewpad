// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
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
	tests := map[string]struct {
		targetEntity *handler.TargetEntity
		wantErr      error
		wantRes      aladino.Value
	}{
		"when pull request": {
			targetEntity: aladino.DefaultMockTargetEntity,
			wantRes:      aladino.BuildStringValue(`{"id":1234,"number":6,"title":"Amazing new feature","body":"Please pull these awesome changes in!","created_at":"2009-11-17T20:34:58.651387237Z","labels":[{"id":1,"name":"enhancement"}],"user":{"login":"john"},"merged":true,"comments":6,"commits":5,"url":"https://foo.bar","assignees":[{"login":"jane"}],"milestone":{"title":"v1.0"},"requested_reviewers":[{"login":"jane"}],"head":{"ref":"new-topic","repo":{"owner":{"login":"foobar"},"name":"default-mock-repo","url":"https://api.github.com/repos/foobar/default-mock-repo/pulls/6"}},"base":{"ref":"master","repo":{"owner":{"login":"foobar"},"name":"default-mock-repo","url":"https://api.github.com/repos/foobar/default-mock-repo/pulls/6"}}}`),
		},
		"when issue": {
			targetEntity: &handler.TargetEntity{
				Kind:   handler.Issue,
				Owner:  "test",
				Repo:   "test",
				Number: 1,
			},
			wantRes: aladino.BuildStringValue(`{"id":1234,"number":6,"title":"Found a bug","body":"I'm having a problem with this","user":{"login":"john"},"labels":[{"id":1,"name":"bug"}],"comments":6,"created_at":"2009-11-17T20:34:58.651387237Z","url":"https://foo.bar","milestone":{"title":"v1.0"},"assignees":[{"login":"jane"}]}`),
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
