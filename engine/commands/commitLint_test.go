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

func TestCommitLint(t *testing.T) {
	tests := map[string]struct {
		args       []string
		wantAction string
		wantErr    error
	}{
		"when invalid number of arguments": {
			args: []string{
				"test",
			},
			wantErr: errors.New("accepts no args, received 1"),
		},
		"when successful": {
			args:       []string{},
			wantAction: "$commitLint()",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := commands.CommitLintCmd()

			cmd.SetOut(out)
			cmd.SetArgs(test.args)

			err := cmd.Execute()

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantAction, out.String())
		})
	}
}
