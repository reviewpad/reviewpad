package plugins_aladino_functions_test

import (
	"net/http"
	"testing"
	"time"

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
		"when there are no checks": {
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
		"when check name is not found": {
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
		"when check status is not completed": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write(mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{
									{
										Name:   github.String("test-check-run"),
										Status: github.String("in_progress"),
									},
								},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue(""),
		},
		"when check conclusion is not set": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write(mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{
									{
										Name:       github.String("test-check-run"),
										Status:     github.String("completed"),
										Conclusion: nil,
									},
								},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue(""),
		},
		"when check is eligible": {
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
										Name:       &checkName,
										Status:     github.String("completed"),
										Conclusion: github.String("success"),
									},
								},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue("success"),
		},
		"when there are multiple checks": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write(mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{
									{
										Name:        github.String("test-check-run"),
										Status:      github.String("completed"),
										CompletedAt: &github.Timestamp{Time: time.Now()},
										Conclusion:  github.String("success"),
									},
									{
										Name:        github.String("test-check-run"),
										Status:      github.String("completed"),
										CompletedAt: &github.Timestamp{Time: time.Now().Add(1 * time.Hour)},
										Conclusion:  github.String("failure"),
									},
								},
							},
						))
					}),
				),
			},
			wantStatus: aladino.BuildStringValue("failure"),
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
