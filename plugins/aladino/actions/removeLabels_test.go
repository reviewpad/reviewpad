// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/gorilla/mux"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var removeLabels = plugins_aladino.PluginBuiltIns().Actions["removeLabels"].Code

func TestRemoveLabels_WhenDeleteLabelRequestFails(t *testing.T) {
	defaultLabel := "default-label"
	failMessage := "DeleteLabelsFail"

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					utils.MustWriteBytes(w, mock.MustMarshal(github.ErrorResponse{
						// An error response may also consist of a 404 status code.
						// However, in this context, such response means a label does not exist.
						Response: &http.Response{
							StatusCode: 500,
						},
						Message: failMessage,
					}))
				}),
			),
		},
		nil,
		aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
			Labels: []*pbc.Label{
				{
					Id:   "1234",
					Name: "defaultLabel",
				},
			},
		}),
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	wantErr := failMessage

	inputLabels := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(defaultLabel)})}
	gotErr := removeLabels(mockedEnv, inputLabels)

	assert.Equal(t, wantErr, gotErr.(*github.ErrorResponse).Message)
}

func TestRemoveLabels_WhenLabelDoesNotExist(t *testing.T) {
	defaultLabel := "default-label"

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesLabelsByOwnerByRepoByIssueNumberByName,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					utils.MustWriteBytes(w, mock.MustMarshal(github.ErrorResponse{
						Response: &http.Response{
							StatusCode: 404,
						},
						Message: "Resource not found",
					}))
				}),
			),
		},
		nil,
		aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
			Labels: []*pbc.Label{
				{
					Id:   "1234",
					Name: defaultLabel,
				},
			},
		}),
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	inputLabels := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(defaultLabel)})}

	gotErr := removeLabels(mockedEnv, inputLabels)

	assert.Nil(t, gotErr)
}

func TestRemoveLabels(t *testing.T) {
	gotRemovedLabels := []string{}

	labelAInIssue := "label-a-in-issue"

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
		"when label does not exist in the labels section of reviewpad.yml": {
			inputLabels: []string{
				labelAInIssue,
			},
			labelExistsInLabelsSection: false,
			wantRemovedLabels: []string{
				labelAInIssue,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				[]mock.MockBackendOption{
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
				aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
					Labels: []*pbc.Label{
						{
							Id:   "1234",
							Name: labelAInIssue,
						},
					},
				}),
				aladino.GetDefaultPullRequestFileList(),
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
