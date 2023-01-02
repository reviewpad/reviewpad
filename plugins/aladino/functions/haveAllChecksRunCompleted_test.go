package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var haveAllChecksRunCompleted = plugins_aladino.PluginBuiltIns().Functions["haveAllChecksRunCompleted"].Code

func TestHaveAllChecksRunCompleted(t *testing.T) {
	tests := map[string]struct {
		checkRunsToIgnore  *aladino.ArrayValue
		conclusion         *aladino.StringValue
		mockBackendOptions []mock.MockBackendOption
		graphQLHandler     http.HandlerFunc
		wantResult         aladino.Value
		wantErr            error
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
			wantResult: aladino.BuildBoolValue(false),
		},
		"when listing check runs fails": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue(""),
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
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue(""),
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
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue(""),
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
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are not completed with success conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue("success"),
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
			wantResult: aladino.BuildBoolValue(false),
		},
		"when all check runs are completed with success conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{}),
			conclusion:        aladino.BuildStringValue("success"),
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
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed with ignored": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("build")}),
			conclusion:        aladino.BuildStringValue(""),
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
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed with ignored and success conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("run")}),
			conclusion:        aladino.BuildStringValue("success"),
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
			wantResult: aladino.BuildBoolValue(true),
		},
		"when all check runs are completed with ignored and failure conclusion": {
			checkRunsToIgnore: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("build")}),
			conclusion:        aladino.BuildStringValue("failure"),
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
			wantResult: aladino.BuildBoolValue(true),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, test.mockBackendOptions, test.graphQLHandler, nil, nil)

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
