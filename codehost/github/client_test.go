// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var dummyPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA0BUezcR7uycgZsfVLlAf4jXP7uFpVh4geSTY39RvYrAll0yh
q7uiQypP2hjQJ1eQXZvkAZx0v9lBYJmX7e0HiJckBr8+/O2kARL+GTCJDJZECpjy
97yylbzGBNl3s76fZ4CJ+4f11fCh7GJ3BJkMf9NFhe8g1TYS0BtSd/sauUQEuG/A
3fOJxKTNmICZr76xavOQ8agA4yW9V5hKcrbHzkfecg/sQsPMmrXixPNxMsqyOMmg
jdJ1aKr7ckEhd48ft4bPMO4DtVL/XFdK2wJZZ0gXJxWiT1Ny41LVql97Odm+OQyx
tcayMkGtMb1nwTcVVl+RG2U5E1lzOYpcQpyYFQIDAQABAoIBAAfUY55WgFlgdYWo
i0r81NZMNBDHBpGo/IvSaR6y/aX2/tMcnRC7NLXWR77rJBn234XGMeQloPb/E8iw
vtjDDH+FQGPImnQl9P/dWRZVjzKcDN9hNfNAdG/R9JmGHUz0JUddvNNsIEH2lgEx
C01u/Ntqdbk+cDvVlwuhm47MMgs6hJmZtS1KDPgYJu4IaB9oaZFN+pUyy8a1w0j9
RAhHpZrsulT5ThgCra4kKGDNnk2yfI91N9lkP5cnhgUmdZESDgrAJURLS8PgInM4
YPV9L68tJCO4g6k+hFiui4h/4cNXYkXnaZSBUoz28ICA6e7I3eJ6Y1ko4ou+Xf0V
csM8VFkCgYEA7y21JfECCfEsTHwwDg0fq2nld4o6FkIWAVQoIh6I6o6tYREmuZ/1
s81FPz/lvQpAvQUXGZlOPB9eW6bZZFytcuKYVNE/EVkuGQtpRXRT630CQiqvUYDZ
4FpqdBQUISt8KWpIofndrPSx6JzI80NSygShQsScWFw2wBIQAnV3TpsCgYEA3reL
L7AwlxCacsPvkazyYwyFfponblBX/OvrYUPPaEwGvSZmE5A/E4bdYTAixDdn4XvE
ChwpmRAWT/9C6jVJ/o1IK25dwnwg68gFDHlaOE+B5/9yNuDvVmg34PWngmpucFb/
6R/kIrF38lEfY0pRb05koW93uj1fj7Uiv+GWRw8CgYEAn1d3IIDQl+kJVydBKItL
tvoEur/m9N8wI9B6MEjhdEp7bXhssSvFF/VAFeQu3OMQwBy9B/vfaCSJy0t79uXb
U/dr/s2sU5VzJZI5nuDh67fLomMni4fpHxN9ajnaM0LyI/E/1FFPgqM+Rzb0lUQb
yqSM/ptXgXJls04VRl4VjtMCgYEAprO/bLx2QjxdPpXGFcXbz6OpsC92YC2nDlsP
3cfB0RFG4gGB2hbX/6eswHglLbVC/hWDkQWvZTATY2FvFps4fV4GrOt5Jn9+rL0U
elfC3e81Dw+2z7jhrE1ptepprUY4z8Fu33HNcuJfI3LxCYKxHZ0R2Xvzo+UYSBqO
ng0eTKUCgYEAxW9G4FjXQH0bjajntjoVQGLRVGWnteoOaQr/cy6oVii954yNMKSP
rezRkSNbJ8cqt9XQS+NNJ6Xwzl3EbuAt6r8f8VO1TIdRgFOgiUXRVNZ3ZyW8Hegd
kGTL0A6/0yAu9qQZlFbaD5bWhQo7eyx63u4hZGppBhkTSPikOYUPCH8=
-----END RSA PRIVATE KEY-----`)

func TestNewGithubAppClient(t *testing.T) {
	tests := map[string]struct {
		appID      int64
		privateKey []byte
		wantType   interface{}
		wantErr    string
	}{
		"fails when private key is invalid": {
			appID:      *github.Int64(1),
			privateKey: nil,
			wantType:   nil,
			wantErr:    "could not parse private key",
		},
		"succeeds when private key is valid": {
			appID:      *github.Int64(1),
			privateKey: dummyPrivateKey,
			wantType:   &host.GithubAppClient{},
			wantErr:    "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := host.NewGithubAppClient(tt.appID, tt.privateKey)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.IsType(t, tt.wantType, got)
			}
		})
	}
}

func TestGetAuthenticatedUserLogin(t *testing.T) {
	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`

	tests := map[string]struct {
		ghGraphQLHandler func(http.ResponseWriter, *http.Request)
		wantUser         string
		wantErr          string
	}{
		"when get authenticated user login request fails": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery) {
					http.Error(w, "GetAuthenticatedUserLoginRequestFail", http.StatusNotFound)
				}
			},
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetAuthenticatedUserLoginRequestFail\\n\"",
		},
		"when get authenticated user login request succeeds": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
					utils.MustWrite(w, `{
							"data": {
								"viewer": {
									"login": "test"
								}
							}
						}`)
				}
			},
			wantUser: "test",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				nil,
				test.ghGraphQLHandler,
				aladino.MockBuiltIns(),
				nil,
			)

			gotUser, gotErr := mockedEnv.GetGithubClient().GetAuthenticatedUserLogin()

			if gotErr != nil {
				assert.EqualError(t, gotErr, test.wantErr)
				assert.Equal(t, "", gotUser)
			} else {
				assert.Nil(t, gotErr)
				assert.Equal(t, test.wantUser, gotUser)
			}
		})
	}
}
