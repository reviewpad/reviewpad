// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/gorilla/mux"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var removeLabels = plugins_aladino.PluginBuiltIns().Actions["removeLabels"].Code

func TestRemoveLabels_WhenRequestFails(t *testing.T) {
	defaultLabel := "default-label"

	tests := map[string]struct {
		mockedEnv aladino.Env
		wantErr   string
	}{
		"when request to remove a label fails": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
								Labels: []*github.Label{
									{
										ID:   github.Int64(1234),
										Name: github.String(defaultLabel),
									},
								},
							})))
						}),
					),
					mock.WithRequestMatch(
						mock.GetReposLabelsByOwnerByRepo,
						[]*github.Label{
							{Name: github.String(defaultLabel)},
						},
					),
					mock.WithRequestMatchHandler(
						mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"DeleteLabelsFail",
							)
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantErr: "DeleteLabelsFail",
		},
	}

	inputLabels := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(defaultLabel)})}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := removeLabels(test.mockedEnv, inputLabels)

			assert.Equal(t, test.wantErr, gotErr.(*github.ErrorResponse).Message)
		})
	}
}

func TestRemoveLabels(t *testing.T) {
	gotRemovedLabels := []string{}

	labelAInIssue := "label-a-in-issue"
	labelBNotInIssue := "label-b-not-in-issue"

	tests := map[string]struct {
		inputLabels                []string
		labelExistsInLabelsSection bool
		wantRemovedLabels          []string
		wantErr                    string
	}{
		"when no labels are provided": {
			inputLabels:                []string{},
			labelExistsInLabelsSection: false,
			wantRemovedLabels:          []string{},
			wantErr:                    "removeLabels: no labels provided",
		},
		"when label exists in the labels section of reviewpad.yml": {
			inputLabels: []string{
				labelAInIssue,
			},
			labelExistsInLabelsSection: true,
			wantRemovedLabels: []string{
				labelAInIssue,
			},
		},
		"when a label is already applied to the issue": {
			inputLabels: []string{
				labelAInIssue,
			},
			labelExistsInLabelsSection: false,
			wantRemovedLabels: []string{
				labelAInIssue,
			},
		},
		"when a label is not yet applied to the issue": {
			inputLabels: []string{
				labelBNotInIssue,
			},
			labelExistsInLabelsSection: false,
			wantRemovedLabels:          []string{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
								Labels: []*github.Label{
									{
										ID:   github.Int64(1234),
										Name: github.String(labelAInIssue),
									},
								},
							})))
						}),
					),
					mock.WithRequestMatch(
						mock.GetReposLabelsByOwnerByRepo,
						[]*github.Label{
							{
								ID:   github.Int64(1234),
								Name: github.String(labelAInIssue),
							},
							{
								ID:   github.Int64(5678),
								Name: github.String(labelBNotInIssue),
							},
						},
					),
					mock.WithRequestMatchHandler(
						mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							vars := mux.Vars(r)
							removedLabelName := vars["name"]
							gotRemovedLabels = append(gotRemovedLabels, removedLabelName)
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			if test.labelExistsInLabelsSection {
				for _, inputLabel := range test.inputLabels {
					mockedEnv.GetRegisterMap()[aladino.BuildInternalLabelID(inputLabel)] = aladino.BuildStringValue(inputLabel)
				}
			}

			labels := make([]aladino.Value, len(test.inputLabels))

			for i, inputLabel := range test.inputLabels {
				labels[i] = aladino.BuildStringValue(inputLabel)
			}

			args := []aladino.Value{aladino.BuildArrayValue(labels)}
			gotErr := removeLabels(mockedEnv, args)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "RemoveLabels() error = %v, wantErr %v", gotErr, test.wantErr)
			}
			assert.Equal(t, test.wantRemovedLabels, gotRemovedLabels)

			// Since this is a variable common to all tests we need to reset its value at the end of each test
			gotRemovedLabels = []string{}
		})
	}
}
