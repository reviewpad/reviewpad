package plugins_aladino_functions_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
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

		args := []lang.Value{lang.BuildStringValue(checkName)}
		gotStatus, err := checkRunConclusion(mockedEnv, args)

		assert.Nil(t, gotStatus)
		assert.Equal(t, test.wantErr, err.(*github.ErrorResponse).Message)
	}
}

func TestCheckRunStatus(t *testing.T) {
	checkName := "test-check-run"

	tests := map[string]struct {
		clientOptions []mock.MockBackendOption
		wantStatus    lang.Value
	}{
		"when there are no checks": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(
							&github.ListCheckRunsResults{
								CheckRuns: []*github.CheckRun{},
							},
						))
					}),
				),
			},
			wantStatus: lang.BuildStringValue(""),
		},
		"when check name is not found": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(
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
			wantStatus: lang.BuildStringValue(""),
		},
		"when check status is not completed": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(
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
			wantStatus: lang.BuildStringValue(""),
		},
		"when check conclusion is not set": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(
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
			wantStatus: lang.BuildStringValue(""),
		},
		"when check is eligible": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(
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
			wantStatus: lang.BuildStringValue("success"),
		},
		"when there are multiple checks": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(
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
			wantStatus: lang.BuildStringValue("failure"),
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

		args := []lang.Value{lang.BuildStringValue(checkName)}
		gotStatus, err := checkRunConclusion(mockedEnv, args)

		assert.Nil(t, err)
		assert.Equal(t, test.wantStatus, gotStatus)
	}
}
