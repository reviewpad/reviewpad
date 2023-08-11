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

var getReviewers = plugins_aladino.PluginBuiltIns().Functions["getReviewers"].Code

func TestGetReviewers(t *testing.T) {
	tests := map[string]struct {
		graphqlHandler func(http.ResponseWriter, *http.Request)
		wantResult     lang.Value
		wantErr        error
		state          string
	}{
		"when graphql query errors": {
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: errors.New(`error getting all reviewers. non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when there are no reviewers": {
			wantResult: lang.BuildArrayValue([]lang.Value{}),
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
		"when there is no filter": {
			wantResult: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("test1"),
				lang.BuildStringValue("test2"),
				lang.BuildStringValue("test3"),
			}),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"nodes": [
										{
											"requestedReviewer": {
												"login": "test1"
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
											"state": "APPROVED",
											"author": {
												"login": "test3"
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
		"when approved only": {
			state: "approved",
			wantResult: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("test3"),
			}),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"reviewRequests": {
									"nodes": [
										{
											"requestedReviewer": {
												"login": "test1"
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
											"state": "APPROVED",
											"author": {
												"login": "test3"
											}
										},
										{
											"state": "CHANGES_REQUESTED",
											"author": {
												"login": "test4"
											}
										},
										{
											"state": "COMMENTED",
											"author": {
												"login": "test5"
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

			res, err := getReviewers(env, []lang.Value{lang.BuildStringValue(test.state)})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
