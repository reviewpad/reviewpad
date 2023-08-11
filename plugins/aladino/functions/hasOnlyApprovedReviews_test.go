package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var hasOnlyApprovedReviews = plugins_aladino.PluginBuiltIns().Functions["hasOnlyApprovedReviews"].Code

func TestHasOnlyApprovedReviews(t *testing.T) {
	tests := map[string]struct {
		graphqlHandler func(http.ResponseWriter, *http.Request)
		wantResult     lang.Value
		wantErr        error
	}{
		"when graphql query errors": {
			wantResult: (lang.Value)(nil),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: errors.New(`error getting approved reviews. non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when there are review requests with no response": {
			wantResult: lang.BuildBoolValue(false),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"totalCount": 2
								},
								"latestReviews": {
									"nodes": [
										{
											"state": "APPROVED"
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantErr: nil,
		},
		"when all reviewers have responded but one with requested changes": {
			wantResult: lang.BuildBoolValue(false),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"totalCount": 0
								},
								"latestReviews": {
									"nodes": [
										{
											"state": "APPROVED"
										},
										{
											"state": "CHANGES_REQUESTED"
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantErr: nil,
		},
		"when all reviewers have responded with approval": {
			wantResult: lang.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"totalCount": 0
								},
								"latestReviews": {
									"nodes": [
										{
											"state": "APPROVED"
										},
										{
											"state": "APPROVED"
										},
										{
											"state": "APPROVED"
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, []mock.MockBackendOption{}, test.graphqlHandler, aladino.MockBuiltIns(), nil)

			res, err := hasOnlyApprovedReviews(env, []lang.Value{})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
