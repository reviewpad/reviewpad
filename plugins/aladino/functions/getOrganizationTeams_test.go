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

var getOrganizationTeams = plugins_aladino.PluginBuiltIns().Functions["getOrganizationTeams"].Code

func TestGetOrganizationTeams(t *testing.T) {
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
			wantErr: errors.New(`error getting organization teams. non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when there are no teams": {
			wantResult: lang.BuildArrayValue([]lang.Value{}),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repositoryOwner": {
							"teams": {
								"nodes": []
							}
						}
					}
				}`)
			},
			wantErr: nil,
		},
		"when there are teams": {
			wantResult: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("developers"),
				lang.BuildStringValue("team"),
			}),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repositoryOwner": {
							"teams": {
								"nodes": [
									{
										"slug": "developers"
									},
									{
										"slug": "team"
									}
								]
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

			res, err := getOrganizationTeams(env, []lang.Value{})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
