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

var hasAnyReviewers = plugins_aladino.PluginBuiltIns().Functions["hasAnyReviewers"].Code

func TestHasAnyReviewers(t *testing.T) {
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
			wantErr: errors.New(`error getting reviewers. non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when there are no reviewers": {
			wantResult: lang.BuildBoolValue(false),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"nodes": []
								},
								"latestReviews": {
									"nodes": []
								}
							}
						}
					}
				}`)
			},
			wantErr: nil,
		},
		"when there are only requested reviewers": {
			wantResult: lang.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"nodes": [
										{
											"requestedReviewer": {
												"login": "test"
											}
										},
										{
											"requestedReviewer": {
												"slug": "test2"
											}
										}
									]
								},
								"latestReviews": {
									"nodes": []
								}
							}
						}
					}
				}`)
			},
			wantErr: nil,
		},
		"when there are only latest reviews": {
			wantResult: lang.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"nodes": []
								},
								"latestReviews": {
									"nodes": [
										{
											"author": {
												"login": "test"
											}
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
		"when there are both requested reviewers and latest reviews": {
			wantResult: lang.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"nodes": [
										{
											"requestedReviewer": {
												"login": "test"
											}
										},
										{
											"requestedReviewer": {
												"slug": "test2"
											}
										}
									]
								},
								"latestReviews": {
									"nodes": [
										{
											"author": {
												"login": "test"
											}
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

			res, err := hasAnyReviewers(env, []lang.Value{})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
