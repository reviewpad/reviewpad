package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var haveAllChecksRunCompleted = plugins_aladino.PluginBuiltIns().Functions["haveAllChecksRunCompleted"].Code

func TestHaveAllChecksRunCompleted(t *testing.T) {
	tests := map[string]struct {
		checkRunsToIgnore  *aladino.ArrayValue
		conclusion         *aladino.StringValue
		mockBackendOptions []mock.MockBackendOption
		wantResult         aladino.Value
		wantErr            error
	}{
		"when get pull request commits failed": {
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
						aladino.MustWrite(w, `{"message": "mock error"}`)
					})),
			},
			wantErr: &github.ErrorResponse{
				Message: "mock error",
			},
		},
		"when there are no github commits": {
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]github.Commit{},
				),
			},
			wantResult: aladino.BuildBoolValue(false),
		},
		"when listing check runs fails": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue(""),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						aladino.MustWrite(w, `{"message": "mock error"}`)
					}),
				),
			},
			wantErr: &github.ErrorResponse{
				Message: "mock error",
			},
		},
		"when there are no check runs": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue(""),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total:     github.Int(0),
						CheckRuns: []*github.CheckRun{},
					},
				),
			},
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue(""),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("completed"),
								Conclusion: github.String("success"),
							},
							{
								Name:       github.String("test"),
								Status:     github.String("completed"),
								Conclusion: github.String("failure"),
							},
						},
					},
				),
			},
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are not completed with conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue("success"),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("completed"),
								Conclusion: github.String("success"),
							},
							{
								Name:       github.String("run"),
								Status:     github.String("completed"),
								Conclusion: github.String("failure"),
							},
						},
					},
				),
			},
			wantResult: aladino.BuildBoolValue(false),
		},
		"when all check runs are completed with conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue("success"),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("completed"),
								Conclusion: github.String("success"),
							},
							{
								Name:       github.String("run"),
								Status:     github.String("completed"),
								Conclusion: github.String("success"),
							},
						},
					},
				),
			},
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed with ignored": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("build")}),
			conclusion:        aladino.BuildStringValue(""),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:   github.String("build"),
								Status: github.String("in_progress"),
							},
							{
								Name:       github.String("run"),
								Status:     github.String("completed"),
								Conclusion: github.String("failure"),
							},
						},
					},
				),
			},
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed with ignored and conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("run")}),
			conclusion:        aladino.BuildStringValue("success"),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
					[]*github.RepositoryCommit{
						{
							SHA: github.String("2a4d4e32a88ee6cb6520ee7232ce6217"),
						},
						{
							SHA: github.String("1992f2a6442859ff07f452282e1cd5b0"),
						},
					},
				),
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("completed"),
								Conclusion: github.String("success"),
							},
							{
								Name:       github.String("run"),
								Status:     github.String("completed"),
								Conclusion: github.String("failure"),
							},
							{
								Name:       github.String("test"),
								Status:     github.String("completed"),
								Conclusion: github.String("success"),
							},
						},
					},
				),
			},
			wantResult: aladino.BuildBoolValue(true),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, test.mockBackendOptions, nil, nil, nil)

			res, err := haveAllChecksRunCompleted(env, []aladino.Value{test.checkRunsToIgnore, test.conclusion})

			githubError := &github.ErrorResponse{}
			if errors.As(err, &githubError) {
				githubError.Response = nil
			}

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
