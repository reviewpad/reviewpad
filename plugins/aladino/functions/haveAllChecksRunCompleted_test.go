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

var haveAllChecksRunCompleted = plugins_aladino.PluginBuiltIns().Functions["haveAllChecksRunCompleted"].Code

func TestHaveAllChecksRunCompleted(t *testing.T) {
	tests := map[string]struct {
		checkRunsToIgnore   *lang.ArrayValue
		conclusion          *lang.StringValue
		conclusionsToIgnore *lang.ArrayValue
		mockBackendOptions  []mock.MockBackendOption
		graphQLHandler      http.HandlerFunc
		wantResult          lang.Value
		wantErr             error
	}{
		"when last commit sha failed": {
			mockBackendOptions: []mock.MockBackendOption{},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			wantErr: errors.New(`non-200 OK status code: 400 Bad Request body: ""`),
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
		"when listing check runs fails": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{}),
			conclusion:          lang.BuildStringValue(""),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						utils.MustWrite(w, `{"message": "mock error"}`)
					}),
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
			wantErr: &github.ErrorResponse{
				Message: "mock error",
			},
		},
		"when there are no check runs": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{}),
			conclusion:          lang.BuildStringValue(""),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
			wantResult: lang.BuildBoolValue(true),
		},
		"when all check runs are completed": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{}),
			conclusion:          lang.BuildStringValue(""),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
		"when all check runs are not completed with success conclusion": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{}),
			conclusion:          lang.BuildStringValue("success"),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
								Name:       github.String("run"),
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
		"when all check runs are completed with success conclusion": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{}),
			conclusion:          lang.BuildStringValue("success"),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
								Name:       github.String("run"),
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
		"when all check runs are completed with ignored": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("build")}),
			conclusion:          lang.BuildStringValue(""),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
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
		"when all check runs are completed with ignored and success conclusion": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("run")}),
			conclusion:          lang.BuildStringValue("success"),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
		"when all check runs are completed with ignored and failure conclusion": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("build")}),
			conclusion:          lang.BuildStringValue("failure"),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
								Name:       github.String("run"),
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
		"when reviewpad is one of the checks": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{}),
			conclusion:          lang.BuildStringValue(""),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{}),
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
								Name:   github.String("reviewpad"),
								Status: github.String("in-progress"),
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
		"when check runs are skipped with ignored conclusion and check run": {
			checkRunsToIgnore:   lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("build")}),
			conclusion:          lang.BuildStringValue("success"),
			conclusionsToIgnore: lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("skipped")}),
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
								Name:       github.String("run"),
								Status:     github.String("in_progress"),
								Conclusion: github.String("skipped"),
							},
							{
								Name:       github.String("test"),
								Status:     github.String("completed"),
								Conclusion: github.String("skipped"),
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

			res, err := haveAllChecksRunCompleted(env, []lang.Value{test.checkRunsToIgnore, test.conclusion, test.conclusionsToIgnore})

			githubError := &github.ErrorResponse{}
			if errors.As(err, &githubError) {
				githubError.Response = nil
			}

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
