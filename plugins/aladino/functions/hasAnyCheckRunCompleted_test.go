// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var hasAnyCheckRunCompleted = plugins_aladino.PluginBuiltIns().Functions["hasAnyCheckRunCompleted"].Code

func TestHasAnyCheckRunCompleted(t *testing.T) {
	tests := map[string]struct {
		checkRunsToIgnore  *lang.ArrayValue
		checkConclusions   *lang.ArrayValue
		mockBackendOptions []mock.MockBackendOption
		graphQLHandler     http.HandlerFunc
		wantResult         lang.Value
		wantErr            error
	}{
		"when getting last commit sha failed": {
			mockBackendOptions: []mock.MockBackendOption{},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			wantErr: errors.New(`failed to get last commit: non-200 OK status code: 400 Bad Request body: ""`),
		},
		"when there are no github commits": {
			mockBackendOptions: []mock.MockBackendOption{},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": []
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(false),
		},
		"when there are no check runs": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total:     github.Int(0),
						CheckRuns: []*github.CheckRun{},
					},
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(false),
		},
		"when no check runs are completed": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("in_progress"),
								Conclusion: github.String("pending"),
							},
							{
								Name:       github.String("test"),
								Status:     github.String("in_progress"),
								Conclusion: github.String("pending"),
							},
						},
					},
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(false),
		},
		"when all check runs are completed but not with a desired conclusion": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("failure")}),
			mockBackendOptions: []mock.MockBackendOption{
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
								Conclusion: github.String("success"),
							},
						},
					},
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(false),
		},
		"when one check run is completed but not with a desired conclusion": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("failure")}),
			mockBackendOptions: []mock.MockBackendOption{
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
								Status:     github.String("in_progress"),
								Conclusion: github.String("pending"),
							},
						},
					},
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(false),
		},
		"when all check runs are completed with a desired conclusion": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("failure")}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("completed"),
								Conclusion: github.String("failure"),
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
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(true),
		},
		"when one check runs is completed with a desired conclusion": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("success")}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("in_progress"),
								Conclusion: github.String("pending"),
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
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(true),
		},
		"when an ignore check run is completed with a desired conclusion": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("build")}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("success")}),
			mockBackendOptions: []mock.MockBackendOption{
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
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(false),
		},
		"when multiple conclusions": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("build")}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("success"), lang.BuildStringValue("failure")}),
			mockBackendOptions: []mock.MockBackendOption{
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
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(true),
		},
		"when no conclusions are provided": {
			checkRunsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			checkConclusions:  lang.BuildArrayValue([]lang.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{
						Total: github.Int(2),
						CheckRuns: []*github.CheckRun{
							{
								Name:       github.String("build"),
								Status:     github.String("in_progress"),
								Conclusion: github.String("pending"),
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
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"commits": {
									"nodes": [
										{
											"commit": {
												"oid": "b0b55a8a10139a324f3ccb1a6481862a4b5b5bcc"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			wantResult: lang.BuildBoolValue(true),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, test.mockBackendOptions, test.graphQLHandler, nil, nil)

			res, err := hasAnyCheckRunCompleted(env, []lang.Value{test.checkRunsToIgnore, test.checkConclusions})

			githubError := &github.ErrorResponse{}
			if errors.As(err, &githubError) {
				githubError.Response = nil
			}

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
