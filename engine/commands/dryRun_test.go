// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine/commands"
	"github.com/stretchr/testify/assert"
)

func TestDryRun(t *testing.T) {
	testCases := map[string]struct {
		args       []string
		wantAction string
		wantErr    error
	}{
		"when invalid number of arguments": {
			args: []string{
				"john",
			},
			wantErr: errors.New("accepts no args, received 1"),
		},
		"when command is fully executed": {
			args:       []string{},
			wantAction: "$dryRun()",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := commands.DryRunCmd()

			cmd.SetOut(out)
			cmd.SetArgs(test.args)

			err := cmd.Execute()

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantAction, out.String())
		})
	}
}
