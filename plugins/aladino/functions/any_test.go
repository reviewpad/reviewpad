// Copyright (C) 2022 zola - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"strings"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var any = plugins_aladino.PluginBuiltIns().Functions["any"].Code

func TestAny(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	testCases := []struct {
		name    string
		args    []lang.Value
		res     lang.Value
		wantErr error
	}{
		{
			name: "matches one",
			args: []lang.Value{
				lang.BuildArrayValue(
					[]lang.Value{
						lang.BuildStringValue("a"),
						lang.BuildStringValue("e"),
						lang.BuildStringValue("f"),
					},
				),
				lang.BuildFunctionValue(func(args []lang.Value) lang.Value {
					val := args[0].(*lang.StringValue).Val
					return lang.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     lang.BuildBoolValue(true),
			wantErr: nil,
		},
		{
			name: "matches two",
			args: []lang.Value{
				lang.BuildArrayValue(
					[]lang.Value{
						lang.BuildStringValue("a"),
						lang.BuildStringValue("b"),
						lang.BuildStringValue("f"),
					},
				),
				lang.BuildFunctionValue(func(args []lang.Value) lang.Value {
					val := args[0].(*lang.StringValue).Val
					return lang.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     lang.BuildBoolValue(true),
			wantErr: nil,
		},
		{
			name: "matches all",
			args: []lang.Value{
				lang.BuildArrayValue(
					[]lang.Value{
						lang.BuildStringValue("a"),
						lang.BuildStringValue("b"),
						lang.BuildStringValue("c"),
					},
				),
				lang.BuildFunctionValue(func(args []lang.Value) lang.Value {
					val := args[0].(*lang.StringValue).Val
					return lang.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     lang.BuildBoolValue(true),
			wantErr: nil,
		},
		{
			name: "matches none",
			args: []lang.Value{
				lang.BuildArrayValue(
					[]lang.Value{
						lang.BuildStringValue("e"),
						lang.BuildStringValue("f"),
						lang.BuildStringValue("g"),
						lang.BuildStringValue("h"),
						lang.BuildStringValue("i"),
					},
				),
				lang.BuildFunctionValue(func(args []lang.Value) lang.Value {
					val := args[0].(*lang.StringValue).Val
					return lang.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     lang.BuildBoolValue(false),
			wantErr: nil,
		},
		{
			name: "empty list",
			args: []lang.Value{
				lang.BuildArrayValue(
					[]lang.Value{},
				),
				lang.BuildFunctionValue(func(args []lang.Value) lang.Value {
					val := args[0].(*lang.StringValue).Val
					return lang.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     lang.BuildBoolValue(false),
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := any(mockedEnv, tc.args)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.res, res)
		})
	}
}
