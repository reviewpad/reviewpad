// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestGetInstallations(t *testing.T) {
	testCases := map[string]struct {
		mockBackendOptions []mock.MockBackendOption
		wantInstallations  []*github.Installation
		wantErrorMessage   string
	}{
		"when list installations request fails": {
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetAppInstallations,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"ListInstallationsRequestFail",
						)
					}),
				),
			},
			wantInstallations: nil,
			wantErrorMessage:  "ListInstallationsRequestFail",
		},
		"when list installations request succeeds": {
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatchPages(
					mock.GetAppInstallations,
					[]*github.Installation{
						{ID: github.Int64(1)},
					},
					[]*github.Installation{
						{ID: github.Int64(2)},
					},
				),
			},
			wantInstallations: []*github.Installation{
				{ID: github.Int64(1)},
				{ID: github.Int64(2)},
			},
			wantErrorMessage: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockedGithubClient := aladino.MockDefaultGithubAppClient(tc.mockBackendOptions)
			ctx := context.Background()

			gotInstallations, gotErr := mockedGithubClient.GetInstallations(ctx)

			assert.Equal(t, tc.wantInstallations, gotInstallations)
			if gotErr != nil {
				assert.Equal(t, tc.wantErrorMessage, gotErr.(*github.ErrorResponse).Message)
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}

func TestCreateInstallationToken(t *testing.T) {
	testCases := map[string]struct {
		mockBackendOptions  []mock.MockBackendOption
		installationID      int64
		installationOptions *github.InstallationTokenOptions
		wantToken           *github.InstallationToken
		wantErrorMessage    string
	}{
		"when create installation token request fails": {
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostAppInstallationsAccessTokensByInstallationId,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"CreateInstallationTokenRequestFail",
						)
					}),
				),
			},
			installationID:      1,
			installationOptions: &github.InstallationTokenOptions{},
			wantToken:           nil,
			wantErrorMessage:    "CreateInstallationTokenRequestFail",
		},
		"when create installation token request succeeds": {
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.PostAppInstallationsAccessTokensByInstallationId,
					&github.InstallationToken{Token: github.String("abc123")},
				),
			},
			installationID:      1,
			installationOptions: &github.InstallationTokenOptions{},
			wantToken: &github.InstallationToken{
				Token: github.String("abc123"),
			},
			wantErrorMessage: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockedGithubClient := aladino.MockDefaultGithubAppClient(tc.mockBackendOptions)
			ctx := context.Background()

			gotToken, gotErr := mockedGithubClient.CreateInstallationToken(ctx, tc.installationID, tc.installationOptions)

			assert.Equal(t, tc.wantToken, gotToken)
			if gotErr != nil {
				assert.Equal(t, tc.wantErrorMessage, gotErr.(*github.ErrorResponse).Message)
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}
