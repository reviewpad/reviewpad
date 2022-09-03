// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"encoding/json"
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
				assert.FailNow(t, "engine MockDefaultEnvWith: %v", err)
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
		wantVal       bool
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
			wantVal: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "MockDefaultEnvWith: %v", err)
			}

			gotVal, gotErr := checkLabelExists(mockedEnv, test.labelName)

			assert.Equal(t, test.wantVal, gotVal)
			assert.Equal(t, test.wantErr, gotErr.(*github.ErrorResponse).Message)
		})
	}
}

func TestCheckLabelExists(t *testing.T) {
	tests := map[string]struct {
		labelName     string
		clientOptions []mock.MockBackendOption
		wantVal       bool
		wantErr       error
	}{
		"when get label request 404": {
			labelName: "test",
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
			wantVal: false,
		},
		"when get label request 200": {
			labelName: "test",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposLabelsByOwnerByRepoByName,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					}),
				),
			},
			wantErr: nil,
			wantVal: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedClient := MockGithubClient(test.clientOptions)

			mockedEnv, err := MockEnvWith(mockedClient, nil)
			if err != nil {
				assert.FailNow(t, "MockDefaultEnvWith: %v", err)
			}

			gotVal, gotErr := checkLabelExists(mockedEnv, test.labelName)

			assert.Equal(t, test.wantVal, gotVal)
			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}
