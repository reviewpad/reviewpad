// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestValidateLabelColor(t *testing.T) {
	tests := map[string]struct {
		arg     *PadLabel
		wantErr error
	}{
		"hex value with #": {
			arg:     &PadLabel{Color: "#91cf60"},
			wantErr: nil,
		},
		"hex value without #": {
			arg:     &PadLabel{Color: "91cf60"},
			wantErr: nil,
		},
		"invalid hex value with #": {
			arg:     &PadLabel{Color: "#91cg60"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"invalid hex value without #": {
			arg:     &PadLabel{Color: "91cg60"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"invalid hex value because of size with #": {
			arg:     &PadLabel{Color: "#91cg6"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"invalid hex value because of size without #": {
			arg:     &PadLabel{Color: "91cg6"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"english color": {
			arg:     &PadLabel{Color: "red"},
			wantErr: execError("evalLabel: color code not valid"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := validateLabelColor(test.arg)
			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}

func TestCreateLabel(t *testing.T) {
	tests := map[string]struct {
		labelName     string
		label         *PadLabel
		clientOptions []mock.MockBackendOption
		wantErr       error
	}{
		"when validate label color fails": {
			labelName: "test",
			label: &PadLabel{
				Name: "test",
				// invalid hex value because of size
				Color:       "#91cf6",
				Description: "test",
			},
			clientOptions: []mock.MockBackendOption{},
			wantErr:       execError("evalLabel: color code not valid"),
		},
		"when label color has leading #": {
			labelName: "test-name",
			label: &PadLabel{
				Color:       "#91cf60",
				Description: "test description",
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostReposLabelsByOwnerByRepo,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var gotLabel github.Label
						json.NewDecoder(r.Body).Decode(&gotLabel)

						assert.Equal(t, "test-name", gotLabel.GetName())
						assert.Equal(t, "91cf60", gotLabel.GetColor())
						assert.Equal(t, "test description", gotLabel.GetDescription())

						w.WriteHeader(http.StatusOK)
					}),
				),
			},
			wantErr: nil,
		},
		"when label color does not have leading #": {
			labelName: "test-name",
			label: &PadLabel{
				Color:       "91cf60",
				Description: "test description",
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostReposLabelsByOwnerByRepo,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var gotLabel github.Label
						json.NewDecoder(r.Body).Decode(&gotLabel)

						assert.Equal(t, "test-name", gotLabel.GetName())
						assert.Equal(t, "91cf60", gotLabel.GetColor())
						assert.Equal(t, "test description", gotLabel.GetDescription())

						w.WriteHeader(http.StatusOK)
					}),
				),
			},
			wantErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "engine MockEnvWith: %v", err)
			}

			gotErr := createLabel(mockedEnv, &test.labelName, test.label)

			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}

func TestCheckLabelExists_WhenGetLabelFails(t *testing.T) {
	tests := map[string]struct {
		labelName     string
		clientOptions []mock.MockBackendOption
		wantVal       *github.Label
		wantErr       string
	}{
		"when get label request 500": {
			labelName: "test",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(mock.MustMarshal(github.ErrorResponse{
							Response: &http.Response{
								StatusCode: 500,
							},
							Message: "GetLabelRequestFailed",
						}))
					}),
				),
			},
			wantErr: "GetLabelRequestFailed",
			wantVal: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "MockEnvWith: %v", err)
			}

			gotVal, gotErr := checkLabelExists(mockedEnv, test.labelName)

			assert.Equal(t, test.wantVal, gotVal)
			assert.Equal(t, test.wantErr, gotErr.(*github.ErrorResponse).Message)
		})
	}
}

func TestCheckLabelExists(t *testing.T) {
	labelName := "test"
	label := &github.Label{
		Name: github.String(labelName),
	}

	tests := map[string]struct {
		labelName     string
		clientOptions []mock.MockBackendOption
		wantVal       *github.Label
		wantErr       error
	}{
		"when get label request 404": {
			labelName: labelName,
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(mock.MustMarshal(github.ErrorResponse{
							Response: &http.Response{
								StatusCode: 404,
							},
						}))
					}),
				),
			},
			wantErr: nil,
			wantVal: nil,
		},
		"when get label request 200": {
			labelName: labelName,
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write(mock.MustMarshal(label))
					}),
				),
			},
			wantErr: nil,
			wantVal: label,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "MockEnvWith: %v", err)
			}

			gotVal, gotErr := checkLabelExists(mockedEnv, test.labelName)

			assert.Equal(t, test.wantVal, gotVal)
			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}

func TestCheckLabelHasUpdates(t *testing.T) {
	labelName := "test"
	description := "Lorem Ipsum"

	tests := map[string]struct {
		label   *PadLabel
		ghLabel *github.Label
		wantVal bool
		wantErr error
	}{
		"label does not exist in the repository": {
			label: &PadLabel{
				Name: labelName,
			},
			ghLabel: nil,
			wantVal: false,
			wantErr: fmt.Errorf("checkLabelHasUpdates: impossible to check updates on a empty github label"),
		},
		"label has an empty description": {
			label: &PadLabel{
				Name: labelName,
			},
			ghLabel: &github.Label{
				Name: github.String(description),
			},
			wantVal: false,
			wantErr: nil,
		},
		"label does not have any updates": {
			label: &PadLabel{
				Name:        labelName,
				Description: description,
			},
			ghLabel: &github.Label{
				Name:        github.String(labelName),
				Description: github.String(description),
			},
			wantVal: false,
			wantErr: nil,
		},
		"label has updates": {
			label: &PadLabel{
				Name:        labelName,
				Description: description,
			},
			ghLabel: &github.Label{
				Name:        github.String(labelName),
				Description: github.String("Updated label description"),
			},
			wantVal: true,
			wantErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(nil)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "MockEnvWith: %v", err)
			}

			gotVal, gotErr := checkLabelHasUpdates(mockedEnv, test.label, test.ghLabel)

			assert.Equal(t, test.wantVal, gotVal)
			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}

func TestUpdateLabel_WhenEditLabelRequestFails(t *testing.T) {
	labelName := "test-name"
	tests := map[string]struct {
		labelName     string
		label         *PadLabel
		clientOptions []mock.MockBackendOption
		wantErr       string
	}{
		"when label update is not successful": {
			labelName: labelName,
			label: &PadLabel{
				Color:       "#91cf60",
				Description: "test description",
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PatchReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(mock.MustMarshal(github.ErrorResponse{
							Response: &http.Response{
								StatusCode: 500,
							},
							Message: "EditLabelRequestFailed",
						}))
					}),
				),
			},
			wantErr: "EditLabelRequestFailed",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "engine MockEnvWith: %v", err)
			}

			gotErr := updateLabel(mockedEnv, &test.labelName, test.label)

			assert.Equal(t, test.wantErr, gotErr.(*github.ErrorResponse).Message)
		})
	}
}

func TestUpdateLabel(t *testing.T) {
	labelName := "test-name"
	tests := map[string]struct {
		labelName     string
		label         *PadLabel
		clientOptions []mock.MockBackendOption
		wantErr       error
	}{
		"when label update is successful": {
			labelName: labelName,
			label: &PadLabel{
				Color:       "#91cf60",
				Description: "test description",
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PatchReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var gotLabel github.Label
						json.NewDecoder(r.Body).Decode(&gotLabel)

						assert.Equal(t, "test description", gotLabel.GetDescription())

						w.WriteHeader(http.StatusOK)
					}),
				),
			},
			wantErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "engine MockEnvWith: %v", err)
			}

			gotErr := updateLabel(mockedEnv, &test.labelName, test.label)

			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}
