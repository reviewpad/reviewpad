package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var checkRunConclusion = plugins_aladino.PluginBuiltIns().Functions["checkRunConclusion"].Code

func TestCheckRunStatus_WhenRequestsFail(t *testing.T) {
	checkName := "test-check-run"

	tests := map[string]struct {
		clientOptions []mock.MockBackendOption
		wantErr       string
	}{
		"when list commits request fails": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"ListCommitsRequestFail",
						)
					}),
				),
			},
			wantErr: "ListCommitsRequestFail",
		},
		"when list check runs for ref request fails": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{SHA: github.String("1234abc")},
					},
				),
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"ListCheckRunsForRefRequestFail",
						)
					}),
				),
			},
			wantErr: "ListCheckRunsForRefRequestFail",
		},
	}

	for _, test := range tests {
		mockedEnv := aladino.MockDefaultEnv(
			t,
			test.clientOptions,
			nil,
			aladino.MockBuiltIns(),
			nil,
		)

		args := []aladino.Value{aladino.BuildStringValue(checkName)}
		gotStatus, err := checkRunConclusion(mockedEnv, args)

		assert.Nil(t, gotStatus)
		assert.Equal(t, test.wantErr, err.(*github.ErrorResponse).Message)
	}
}

func TestCheckRunStatus(t *testing.T) {
	checkName := "test-check-run"

	tests := map[string]struct {
		clientOptions []mock.MockBackendOption
		wantStatus    aladino.Value
	}{
		"when check run not found due to empty check runs": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write(mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue(""),
		},
		"when check run is missing in non empty check runs": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write(mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{
									{
										Name: github.String("dummy-check"),
									},
								},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue(""),
		},
		"when check run is found": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write(mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{
									{
										Name: github.String("dummy-check"),
									},
									{
										Name:   &checkName,
										Status: github.String("completed"),
									},
								},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue("completed"),
		},
	}

	for _, test := range tests {
		mockedEnv := aladino.MockDefaultEnv(
			t,
			append(
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						[]*github.RepositoryCommit{
							{SHA: github.String("1234abc")},
						},
					),
				},
				test.clientOptions...,
			),
			nil,
			aladino.MockBuiltIns(),
			nil,
		)

		args := []aladino.Value{aladino.BuildStringValue(checkName)}
		gotStatus, err := checkRunConclusion(mockedEnv, args)

		assert.Nil(t, err)
		assert.Equal(t, test.wantStatus, gotStatus)
	}
}
