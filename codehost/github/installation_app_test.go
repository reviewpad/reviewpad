// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestGetInstallations_WhenListInstallationsRequestFails(t *testing.T) {
	failMessage := "ListInstallationsRequestFail"
	mockedGithubClient := aladino.MockDefaultGithubAppClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetAppInstallations,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
	)

	ctx := context.Background()

	gotInstallations, gotErr := mockedGithubClient.GetInstallations(ctx)

	assert.Nil(t, gotInstallations)
	assert.Equal(t, gotErr.(*github.ErrorResponse).Message, failMessage)
}

func TestGetInstallation(t *testing.T) {
	firstInstallationID := &github.Installation{ID: github.Int64(1)}
	secondInstallationID := &github.Installation{ID: github.Int64(2)}

	wantInstallations := []*github.Installation{
		firstInstallationID,
		secondInstallationID,
	}

	mockedGithubClient := aladino.MockDefaultGithubAppClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchPages(
				mock.GetAppInstallations,
				[]*github.Installation{
					firstInstallationID,
				},
				[]*github.Installation{
					secondInstallationID,
				},
			),
		},
	)

	ctx := context.Background()

	gotInstallations, gotErr := mockedGithubClient.GetInstallations(ctx)

	assert.Nil(t, gotErr)
	assert.Equal(t, wantInstallations, gotInstallations)
}

func TestCreateInstallationToken_WhenCreateInstallationTokenRequestFails(t *testing.T) {
	failMessage := "CreateInstallationTokenRequestFail"
	mockedGithubClient := aladino.MockDefaultGithubAppClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PostAppInstallationsAccessTokensByInstallationId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
	)

	ctx := context.Background()

	gotToken, gotErr := mockedGithubClient.CreateInstallationToken(ctx, *github.Int64(1), &github.InstallationTokenOptions{})

	assert.Nil(t, gotToken)
	assert.Equal(t, gotErr.(*github.ErrorResponse).Message, failMessage)
}

func TestCreateInstallationToken(t *testing.T) {
	wantToken := &github.InstallationToken{
		Token: github.String("abc123"),
	}
	mockedGithubClient := aladino.MockDefaultGithubAppClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.PostAppInstallationsAccessTokensByInstallationId,
				wantToken,
			),
		},
	)

	ctx := context.Background()

	gotToken, gotErr := mockedGithubClient.CreateInstallationToken(ctx, *github.Int64(1), &github.InstallationTokenOptions{})

	assert.Nil(t, gotErr)
	assert.Equal(t, wantToken, gotToken)
}
