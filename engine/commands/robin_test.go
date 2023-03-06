// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v4/engine/commands"
	"github.com/stretchr/testify/assert"
)

func TestRobin(t *testing.T) {
	testCases := map[string]struct {
		args       []string
		wantAction string
		wantErr    error
	}{
		"when invalid number of arguments": {
			args: []string{
				"john",
				"jane",
			},
			wantErr: errors.New("accepts 1 arg(s), received 2"),
		},
		"when correct string": {
			args: []string{
				"\"john\"",
			},
			wantAction: `$robin("john")`,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := commands.RobinCmd()

			cmd.SetOut(out)
			cmd.SetArgs(test.args)

			err := cmd.Execute()

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantAction, out.String())
		})
	}
}
