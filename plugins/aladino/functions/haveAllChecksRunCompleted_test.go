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
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, test.mockBackendOptions, nil, nil, nil)

			res, err := haveAllChecksRunCompleted(env, []aladino.Value{test.checkRunsToIgnore})

			githubError := &github.ErrorResponse{}
			if errors.As(err, &githubError) {
				githubError.Response = nil
			}

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
